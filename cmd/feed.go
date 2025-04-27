package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/ohnishi/feed/common"
	"github.com/pkg/errors"
)

type subscription struct {
	Title   string    `json:"title"`
	URL     string    `json:"url"`
	Fetched time.Time `json:"fetched"`
}

func feedToMail(dest string, retry int) error {
	filePath := filepath.Join(dest, "fetchinfo.jsonl")
	subscriptions, err := readSubscriptions(filePath)
	if err != nil {
		return err
	}

	var bodys []string
	fp := gofeed.NewParser()
	for i, subscription := range subscriptions {
		var latest time.Time
		var feed *gofeed.Feed
		for cnt := 1; cnt <= retry; cnt++ {
			feed, err = fp.ParseURL(subscription.URL)
			if err == nil {
				break
			} else if cnt == retry {
				return err
			}
			time.Sleep(3 * time.Second)
		}

		siteBodys := []string{subscription.Title}
		for _, item := range feed.Items {
			published, err := parseLocal(item.Published)
			if err != nil {
				return err
			}

			if latest.Before(published) {
				latest = published
			}

			if subscription.Fetched.Before(published) {
				siteBodys = append(siteBodys, "  - "+item.Title)
				siteBodys = append(siteBodys, "    - "+item.Link)
			}
		}
		if len(siteBodys) > 1 {
			if len(bodys) > 0 {
				bodys = append(bodys, "")
			}
			bodys = append(bodys, siteBodys...)
		}

		if subscription.Fetched.Before(latest) {
			subscriptions[i].Fetched = latest
		}
	}

	// err = common.ChatworkNotify(strings.Join(bodys, "\n"), "", 1)
	// if err != nil {
	// 	return err
	// }

	// err = common.MailNotify(strings.Join(bodys, "\n"), "")
	// if err != nil {
	// 	return err
	// }

	err = common.MailNotifyByResend(strings.Join(bodys, "\n"), "")
	if err != nil {
		return err
	}

	return writeSubscription(filepath.Join(dest, "fetchinfo.jsonl"), subscriptions)
}

func writeSubscription(path string, subscriptions []subscription) error {
	if len(subscriptions) == 0 {
		return nil
	}
	f, err := createOutFile(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, subscription := range subscriptions {

		err = appendOutFile(f, subscription)
		if err != nil {
			return err
		}
	}
	if err := f.Sync(); err != nil {
		return errors.Wrap(err, "failed to sync file")
	}
	return nil
}

func createOutFile(path string) (*os.File, error) {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create output directory: %s", dir)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", path)
	}
	return f, nil
}

func appendOutFile(f *os.File, v interface{}) error {
	jsonl, err := toJSON(v)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(jsonl)); err != nil {
		return errors.Wrapf(err, "failed to write line: %s", jsonl)
	}
	return nil
}
func toJSON(r interface{}) (string, error) {
	jsonStr, err := json.Marshal(r)
	if err != nil {
		return "", errors.Wrapf(err, "could not marshal: %v", r)
	}
	return fmt.Sprintf("%s\n", jsonStr), nil
}

func readSubscriptions(path string) ([]subscription, error) {
	if !FileExists(path) {
		return nil, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file: %s", path)
	}
	defer f.Close()

	var subscriptions []subscription
	d := json.NewDecoder(f)
	for d.More() {
		var s subscription
		if err := d.Decode(&s); err != nil {
			return nil, errors.Wrap(err, "could not unmarshal subscription")
		}
		subscriptions = append(subscriptions, s)
	}
	return subscriptions, nil
}

func parseLocal(value string) (t time.Time, err error) {
	for _, layout := range []string{time.RFC1123, time.RFC3339, time.RFC1123Z} {
		t, err = time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			break
		}
	}
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "cannot parse as %q", time.Local)
	}
	return t, nil
}

// FileExists はファイルが存在するかどうかを確認します。
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}
