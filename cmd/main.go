package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const MAX_RETRY = 3

var (
	dest string
)

func newFeedToMailCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feed",
		Short: "Fetch rss feed",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := feedToMail(dest, MAX_RETRY)
			if err != nil {
				return err
			}
			return nil
		},
	}
	gopath := os.Getenv("GOPATH")
	cmd.PersistentFlags().StringVar(&dest, "dest", filepath.Join(gopath, "cache"), "dest dir path")

	return cmd
}

func newResetSubscriptionsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "reset subscriptions",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := resetSubscriptions(dest)
			if err != nil {
				return err
			}
			return nil
		},
	}
	gopath := os.Getenv("GOPATH")
	cmd.PersistentFlags().StringVar(&dest, "dest", filepath.Join(gopath, "cache"), "dest dir path")

	return cmd
}

func main() {
	rootCmd := &cobra.Command{Use: "fetch"}
	rootCmd.AddCommand(
		newFeedToMailCommand(),
		newResetSubscriptionsCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
