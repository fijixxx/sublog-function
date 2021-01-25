package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
)

// Config meta.toml 自体を定義
type Config struct {
	Meta Meta
}

// Meta meta.toml の中身
type Meta struct {
	Category string   `toml:"category"`
	FileName string   `toml:"fileName"`
	Tag      []string `toml:"tag"`
	Title    string   `toml:"title"`
}

// Sublog DynamoDB のレコード型
type Sublog struct {
	ID          string   `dynamo:"id,hash"`
	CreatedAt   string   `dynamo:"createdAt,range"`
	FileName    string   `dynamo:"fileName" index:"fileName-index"`
	Category    string   `dynamo:"category"`
	Media       string   `dynamo:"media"`
	Title       string   `dynamo:"title"`
	EyeCatchURL string   `dynamo:"eyeCatchURL"`
	Tag         []string `dynamo:"tag"`
	UpdatedAt   string   `dynamo:"updatedAt"`
	Body        string   `dynamo:"body"`
}

// MsgSQS SQS に送るメッセージ型
type MsgSQS struct {
	ID string `json:"id"`
}

/*
HandleRequest S3の PUT イベントにトリガーされ、
PUT された meta レコードを変換して
DynamoDB へ UPSERT する
*/
func HandleRequest(ctx context.Context, event events.S3Event) (string, error) {
	fmt.Printf("[INFO] event: ")

	// AWS SDK セッション作成
	region := "ap-northeast-1"
	sess := session.Must(session.NewSession())

	// tomlファイル格納バケットを取得
	secC := secretsmanager.New(sess, aws.NewConfig().WithRegion(region))
	inputS3 := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("sublog_assets_bucket_name"),
		VersionStage: aws.String("AWSCURRENT"),
	}

	rSecS3, err := secC.GetSecretValue(inputS3)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	}

	sSecS3 := aws.StringValue(rSecS3.SecretString)
	umSecS3 := make(map[string]interface{})
	if err := json.Unmarshal([]byte(sSecS3), &umSecS3); err != nil {
		fmt.Printf("[ERROR] %v", err)
	}
	s3Sec := umSecS3["name"].(string)

	// fileName を取得
	bucketName := s3Sec
	objectKey := event.Records[0].S3.Object.Key
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	}
	_objectKey := strings.Replace(objectKey, "meta/", "", 1)
	fileName := strings.Replace(_objectKey, ".toml", "", 1)

	// S3 clientを作成
	svc := s3.New(sess)

	obj, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	}
	fmt.Printf("%v", obj)

	bd := obj.Body
	defer bd.Close()
	bb, err := ioutil.ReadAll(bd)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	}

	sb := string(bb)

	var config Config

	if _, err := toml.Decode(sb, &config); err != nil {
		fmt.Printf("[ERROR] %v", err)
	}
	s := config
	fmt.Printf("[INFO] toml: %#v", s)

	// ID, 作成日などのメタ情報を作成
	u, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	}
	uu := u.String()

	t := time.Now()
	ts := t.String()

	slColors := []string{
		"Indianred", "Lightcoral", "Salmon", "Darksalmon", "Lightsalmon", "Crimson", "Red", "Firebrick", "Darkred", "Pink", "Lightpink", "Hotpink", "Deeppink", "Mediumvioletred", "Palevioletred", "Lightsalmon", "Coral", "Tomato", "Orangered", "Darkorange", "Orange", "Gold", "Yellow", "Lightyellow", "Lemonchiffon", "Lightgoldenrodyellow", "Papayawhip", "Moccasin", "Peachpuff", "Palegoldenrod", "Khaki", "Darkkhaki", "Greenyellow", "Chartreuse", "Lawngreen", "Lime", "Limegreen", "Palegreen", "Lightgreen", "Mediumspringgreen", "Springgreen", "Mediumseagreen", "Seagreen", "Forestgreen", "Green", "Darkgreen", "Yellowgreen", "Olivedrab", "Olive", "Darkolivegreen", "Mediumaquamarine", "Darkseagreen", "Lightseagreen", "Darkcyan", "Teal", "Aqua", "Cyan", "Lightcyan", "Paleturquoise", "Aquamarine", "Turquoise", "Mediumturquoise", "Darkturquoise", "Cadetblue", "Steelblue", "Lightsteelblue", "Powderblue", "Lightblue", "Skyblue", "Lightskyblue", "Deepskyblue", "Dodgerblue", "Cornflowerblue", "Mediumslateblue", "Royalblue", "Blue", "Mediumblue", "Darkblue", "Navy", "Midnightblue", "Lavender", "Thistle", "Plum", "Violet", "Orchid", "Fuchsia", "Magenta", "Mediumorchid", "Mediumpurple", "Rebeccapurple", "Blueviolet", "Darkviolet", "Darkorchid", "Darkmagenta", "Purple", "Indigo", "Slateblue", "Darkslateblue", "Mediumslateblue", "Cornsilk", "Blanchedalmond", "Bisque", "Navajowhite", "Wheat", "Burlywood", "Tan", "Rosybrown", "Sandybrown", "Goldenrod", "Darkgoldenrod", "Peru", "Chocolate", "Saddlebrown", "Sienna", "Brown", "Maroon", "Snow", "Honeydew", "Mintcream", "Azure", "Aliceblue", "Ghostwhite", "Whitesmoke", "Seashell", "Beige", "Oldlace", "Floralwhite", "Ivory", "Antiquewhite", "Linen", "Lavenderblush", "Mistyrose", "Gainsboro", "Lightgray", "Silver", "Darkgray", "Gray", "Dimgray", "Lightslategray", "Slategray", "Darkslategray", "Black"}

	rand.Seed(time.Now().UnixNano())
	cix := rand.Intn(len(slColors) - 1)

	// 構造体に toml 変換データと作成したメタデータをマッピング
	item := Sublog{
		ID:          uu,
		CreatedAt:   ts,
		FileName:    s.Meta.FileName,
		Category:    s.Meta.Category,
		Media:       "sublog",
		Title:       s.Meta.Title,
		EyeCatchURL: slColors[cix],
		Tag:         s.Meta.Tag,
		UpdatedAt:   ts,
	}
	fmt.Printf("[INFO] item: %#v", item)

	// DynamoDB の設定
	db := dynamo.New(sess, aws.NewConfig().WithRegion(region))
	table := db.Table("sublog")

	// fileName を元に既存レコードの有無をチェック
	var recC = Sublog{}
	ixName := "fileName-index"
	if err := table.Get("fileName", fileName).Index(ixName).One(&recC); err != nil {
		fmt.Printf("[ERROR] %v", err)
	}
	fmt.Printf("[INFO] isUpdateQueryResponseFileName: %#v", recC.FileName)

	// 既存レコード存在した場合、一部項目を上書き（作成日など）
	if recC.ID != "" {
		p := &item
		p.ID = recC.ID
		p.CreatedAt = recC.CreatedAt
		p.EyeCatchURL = recC.EyeCatchURL
	}

	// PUT処理
	if err := table.Put(item).Run(); err != nil {
		fmt.Printf("[ERROR] %v", err)
	}

	// SQS 送信先取得
	inputSQS := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("sublogHighlighterSQS"),
		VersionStage: aws.String("AWSCURRENT"),
	}
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	}

	rSecSQS, err := secC.GetSecretValue(inputSQS)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	}
	sSecSQS := aws.StringValue(rSecSQS.SecretString)
	umSecSQS := make(map[string]interface{})
	if err := json.Unmarshal([]byte(sSecSQS), &umSecSQS); err != nil {
		fmt.Printf("[ERROR] %v", err)
	}
	sqsSec := umSecSQS["url"].(string)

	// SQS に記事 ID を送信
	sqsC := sqs.New(sess, aws.NewConfig().WithRegion(region))

	msgSQS := new(MsgSQS)
	msgSQS.ID = item.ID

	msgJSQS, err := json.Marshal(msgSQS)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	}
	msgBSQS := string(msgJSQS)

	setSQS := &sqs.SendMessageInput{
		MessageBody: aws.String(msgBSQS),
		QueueUrl:    aws.String(sqsSec),
	}

	sqsRes, err := sqsC.SendMessage(setSQS)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	} else {
		fmt.Print("[INFO] SQSMessageID", *sqsRes.MessageId)
	}

	// Slack 送信先取得
	inputSlk := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String("sublog_slack_url"),
		VersionStage: aws.String("AWSCURRENT"),
	}
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	}

	rSecSlk, err := secC.GetSecretValue(inputSlk)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	}
	sSecSlk := aws.StringValue(rSecSlk.SecretString)
	umSecSlk := make(map[string]interface{})
	if err := json.Unmarshal([]byte(sSecSlk), &umSecSlk); err != nil {
		fmt.Printf("[ERROR] %v", err)
	}
	slkSec := umSecSlk["url"].(string)

	// Slack 通知処理
	urlSlk := slkSec
	jsonStr := `{"channel":"#notification",
					"username":"webhookbot",
					"text":"title: ` + item.Title + "\n" + "\n" + `fileName: ` + item.FileName + "\n" + "\n" + `記事を作成しました。",
					"icon_emoji":"ghost"}`
	fmt.Printf(jsonStr)
	req, err := http.NewRequest(
		"POST",
		urlSlk,
		bytes.NewBuffer([]byte(jsonStr)),
	)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	}
	defer resp.Body.Close()

	return "ok", nil
}

func main() {
	lambda.Start(HandleRequest)
}
