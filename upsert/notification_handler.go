package upsert

import (
	"log"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// NotificationHandler Notification 関連のハンドラー
func NotificationHandler(sc *secretsmanager.SecretsManager, it string, fn string) error {
	// Slack 送信先取得
	sn := "sublog_slack_url"
	sk := "url"
	su, err := SecretGet(sc, sn, sk)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	// Slack 通知処理
	if err := Notification(su, it, fn, err); err != nil {
		log.Printf("[ERROR] %v", err)
	}

	return err
}
