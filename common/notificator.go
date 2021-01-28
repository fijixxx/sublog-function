package common

import (
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// Notificator は Slack へ更新通知を送る
func Notificator(sc *secretsmanager.SecretsManager, it string, fn string) error {
	// Slack 送信先取得
	sp := &SMParams{
		SName:     "sublog_slack_url",
		SecretKey: "url",
	}
	su, err := SecretGet(sc, sp)

	// Slack 通知処理
	err = Notify(su, it, fn, err)

	return err
}
