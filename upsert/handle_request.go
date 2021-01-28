package upsert

import (
	"context"
	"encoding/json"
	"math/rand"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/fijixxx/sublog-function/common"
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

/*
HandleRequest S3の PUT イベントにトリガーされ、
PUT された meta.toml ファイルを変換して
DynamoDB へ UPSERT する
*/
func HandleRequest(ctx context.Context, event events.S3Event) error {
	je, err := json.Marshal(event)
	if err != nil {
		common.Logger(1, err.Error())
	}
	common.Logger(0, "Event: "+string(je))

	// AWS SDK セッション作成
	sess := session.Must(session.NewSession())

	// SecretManager クライアントセットアップ
	sc := secretsmanager.New(sess, aws.NewConfig().WithRegion(region))

	// S3 EventNotification から objKey を取得
	ok := event.Records[0].S3.Object.Key

	// fileName を切り出し
	fn := strings.Replace(ok, "meta/", "", 1)
	fn = strings.Replace(fn, ".toml", "", 1)

	// バケットから toml データを取得
	tb, err := common.S3BodyGetter(sc, sess, ok)
	if err != nil {
		common.Logger(1, err.Error())
	}

	// toml ファイルを config Config にマッピング
	var config Config
	if _, err := toml.Decode(tb, &config); err != nil {
		common.Logger(1, err.Error())
	}

	// ID, 作成日などのメタ情報を作成
	u, err := uuid.NewRandom()
	if err != nil {
		common.Logger(1, err.Error())
	}
	uu := u.String()

	t := time.Now()
	ts := t.String()

	rand.Seed(time.Now().UnixNano())
	cix := rand.Intn(len(colorList) - 1)

	// item に toml 変換データと作成したメタデータをマッピング
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

	// fileName を元に既存レコードの有無をチェック
	db := dynamo.New(sess, aws.NewConfig().WithRegion(region))

	var orc Sublog
	tn := "sublog"
	ixn := "fileName-index"
	hk := "fileName"

	table := db.Table(tn)
	if err := table.Get(hk, fn).Index(ixn).One(&orc); err != nil {
		common.Logger(1, err.Error())
	}
	common.Logger(0, "fileName: "+orc.FileName)

	// 既存レコードが存在した場合、一部項目を上書き（作成日など）
	if orc.ID != "" {
		p := &item
		p.ID = orc.ID
		p.CreatedAt = orc.CreatedAt
		p.EyeCatchURL = orc.EyeCatchURL
	}

	jitem, err := json.Marshal(&item)
	if err != nil {
		common.Logger(1, err.Error())
	}
	common.Logger(0, "Item: "+string(jitem))

	// PUT処理
	if err := table.Put(item).Run(); err != nil {
		common.Logger(1, err.Error())
	}

	// 後続処理用に SQS へメッセージを送信
	if err := common.SQSPutter(sc, sess, item.ID); err != nil {
		common.Logger(1, err.Error())
	}

	// Slack への通知処理
	if err := common.Notificator(sc, item.Title, item.FileName); err != nil {
		common.Logger(1, err.Error())
	}

	return err
}
