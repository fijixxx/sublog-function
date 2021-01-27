package upsert

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// Notification Slack通知処理
func Notification(su string, ti string, fn string, err error) error {
	v := ""
	if err != nil {
		v = `"` + fcn + ` 関数でエラーが発生しました",`
	} else {
		v = `"title: ` + ti + "\n" + "\n" + `fileName: ` + fn + "\n" + "\n" + fcn + ` の処理を完了しました。",`
	}

	bj := `
		{"channel":"#notification",
		"username":"webhookbot",
		"text":` + v + `
		"icon_emoji":"ghost"}`
	req, err := http.NewRequest(
		"POST",
		su,
		bytes.NewBuffer([]byte(bj)),
	)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
	defer resp.Body.Close()

	return err
}
