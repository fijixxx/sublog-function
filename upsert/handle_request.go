package upsert

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/google/uuid"
	"github.com/guregu/dynamo"
)

const (
	fcn    = "upsert"
	region = "ap-northeast-1"
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

/*
HandleRequest S3の PUT イベントにトリガーされ、
PUT された meta レコードを変換して
DynamoDB へ UPSERT する
*/
func HandleRequest(ctx context.Context, event events.S3Event) error {
	fmt.Printf("[INFO] event: %v\n", event)

	// AWS SDK セッション作成
	sess := session.Must(session.NewSession())

	// SecretManager クライアントセットアップ
	sc := secretsmanager.New(sess, aws.NewConfig().WithRegion(region))

	// S3 EventNotification から objKey を取得
	ok := event.Records[0].S3.Object.Key

	// fileName を切り出し
	fn := strings.Replace(ok, "meta/", "", 1)
	fn = strings.Replace(fn, ".toml", "", 1)

	// バケットから Toml データを取得
	tb, _ := S3Handler(sc, sess, ok)

	// toml ファイルを Config 構造体にマッピング
	var config Config
	if _, err := toml.Decode(tb, &config); err != nil {
		log.Printf("[ERROR] %v", err)
	}

	// ID, 作成日などのメタ情報を作成
	u, err := uuid.NewRandom()
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
	uu := u.String()

	t := time.Now()
	ts := t.String()

	rand.Seed(time.Now().UnixNano())
	cix := rand.Intn(len(colorList) - 1)

	// 構造体に toml 変換データと作成したメタデータをマッピング
	item := Sublog{
		ID:          uu,
		CreatedAt:   ts,
		FileName:    config.Meta.FileName,
		Category:    config.Meta.Category,
		Media:       "sublog",
		Title:       config.Meta.Title,
		EyeCatchURL: colorList[cix],
		Tag:         config.Meta.Tag,
		UpdatedAt:   ts,
	}
	fmt.Printf("[INFO] item: %#v\n", item)

	// fileName を元に既存レコードの有無をチェック
	db := dynamo.New(sess, aws.NewConfig().WithRegion(region))

	var orc Sublog
	tn := "sublog"
	ixn := "fileName-index"
	hk := "fileName"

	table := db.Table(tn)
	if err := table.Get(hk, fn).Index(ixn).One(&orc); err != nil {
		log.Printf("[ERROR] %v", err)
	}
	fmt.Printf("[INFO] fileName: %v\n", orc.FileName)

	// 既存レコードが存在した場合、一部項目を上書き（作成日など）
	if orc.ID != "" {
		p := &item
		p.ID = orc.ID
		p.CreatedAt = orc.CreatedAt
		p.EyeCatchURL = orc.EyeCatchURL
	}
	// PUT処理
	if err := table.Put(item).Run(); err != nil {
		log.Printf("[ERROR] %v", err)
	}

	// SQS 関連処理
	SQSHandler(sc, sess, item.ID)

	// Slack 通知関連処理
	NotificationHandler(sc, item.Title, item.FileName)

	return err
}
