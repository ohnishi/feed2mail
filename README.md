# feed

購読中の RSS/Atom フィードを巡回し、新着記事を Resend 経由でメール通知します。
GitHub Actions で 1 日 3 回（JST 8:00 / 12:00 / 18:00）自動実行されます。

## 実行

```sh
export MAIL_APIKEY='re_xxxxxxxx'      # Resend の API キー
export MAIL_TO='you@example.com'      # 宛先（カンマ区切りで複数可）
go run .
```

| 環境変数 | 必須 | 説明 |
| --- | --- | --- |
| `MAIL_APIKEY` | ○ | Resend の API キー |
| `MAIL_TO` | ○ | 通知の宛先。カンマ区切りで複数指定可 |
| `MAIL_FROM` | | 送信元。未設定なら `config.yaml` の `mail.from` |

設定ファイルは既定で `config.yaml` を読みます。`-config` で変更できます。

```sh
go run . -config path/to/config.yaml
```

## 設定

購読リスト・除外ルール・保持期間などは [config.yaml](config.yaml) にあります。
フィードを増やすときは `subscriptions` に追記するだけで、再ビルドは不要です。

宛先と API キーはリポジトリに残さないため、設定ファイルではなく環境変数から読みます。

## 新着の判定

同じ記事を二度通知しないよう、通知済みのリンクを `cache/seen.jsonl` に記録します。
保持期間は `seen_retention`（既定 30 日）で、これを過ぎた記録は捨てます。

公開時刻だけで判定すると、後から遡って追加された記事を取りこぼします。
リンク単位で既読を持つことでこれを防いでいます。

まだ既読記録の無いフィード（新規購読、またはこの方式への移行直後）は、
過去記事を一斉に送りつけないよう `cache/fetched.jsonl` の取得済み時刻より
新しいものだけを通知し、残りは既読として取り込むだけにとどめます。

## 状態ファイル

いずれも 1 レコード 1 行の JSONL です。GitHub Actions が実行のたびにコミットします。

| ファイル | 内容 |
| --- | --- |
| `cache/fetched.jsonl` | 購読フィードごとの最終取得時刻 |
| `cache/seen.jsonl` | 通知済みリンク（保持期間内のみ） |

## 構成

```
main.go              設定の読み込みと依存の組み立て
internal/config      設定ファイルと環境変数の読み込み・検証
internal/feed        フィード取得（gofeed への依存をここに閉じ込める）
internal/store       購読状態と既読リンクの永続化
internal/render      通知本文の整形
internal/notify      通知の送信
internal/app         上記を束ねるオーケストレーション
```

## テスト

```sh
go test ./...
```

実ネットワークと実時間には依存しません。待機処理と現在時刻は注入で差し替えています。
