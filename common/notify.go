package common

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// SlackPutContent は Slack に送信する内容を含む
type SlackPutContent struct {
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	IconEmoji string `json:"icon_emoji"`
}

// Notify は Slack に通知をする
func Notify(su string, ti string, fn string, err error) error {
	v := ""
	if err != nil {
		v = `"` + fcn + ` の処理でエラーが発生しました",`
	} else {
		v = `"title: ` + ti + "\n" + "\n" + `fileName: ` + fn + "\n" + "\n" + fcn + ` の処理を完了しました。",`
	}

	b := &SlackPutContent{
		Channel:   "#notification",
		Username:  "sublogUpdateBot",
		Text:      v,
		IconEmoji: "ghost",
	}

	jb, err := json.Marshal(b)

	req, err := http.NewRequest(
		"POST",
		su,
		bytes.NewBuffer([]byte(jb)),
	)

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Logger(1, err.Error())
	}
	defer resp.Body.Close()

	return err
}
