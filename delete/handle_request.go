package delete

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/fijixxx/sublog-function/common"
	"github.com/fijixxx/sublog-function/logger"
	"github.com/guregu/dynamo"
)

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
HandleRequest は S3の DELETE イベントにトリガーされ、
EventNotification で渡ってきた fileName で DynamoDB をクエリし、
対象の fileName を持つレコードを DELETE する
*/
func HandleRequest(ctx context.Context, event events.S3Event) error {
	je, err := json.Marshal(event)
	if err != nil {
		logger.Logger(2, err.Error())
	}
	logger.Logger(1, "Event: "+string(je))

	// AWS SDK セッション作成
	sess := session.Must(session.NewSession())

	// SecretManager クライアントセットアップ
	sc := secretsmanager.New(sess, aws.NewConfig().WithRegion(common.Region))

	// S3 EventNotification から objKey を取得
	ok := event.Records[0].S3.Object.Key

	// fileName を切り出し
	fn := strings.Replace(ok, "meta/", "", 1)
	fn = strings.Replace(fn, ".toml", "", 1)

	// fileName を元に既存レコードの有無をチェック
	db := dynamo.New(sess, aws.NewConfig().WithRegion(common.Region))

	var rc Sublog
	tn := "sublog"
	ixn := "fileName-index"
	hk := "fileName"

	table := db.Table(tn)
	if err := table.Get(hk, fn).Index(ixn).One(&rc); err != nil {
		logger.Logger(2, err.Error())
	}
	logger.Logger(1, "fileName: "+rc.FileName)

	jrc, err := json.Marshal(&rc)
	if err != nil {
		logger.Logger(2, err.Error())
	}
	logger.Logger(1, "Item: "+string(jrc))

	// DELETE 処理
	dhk := "id"
	dv := rc.ID
	if err := table.Delete(dhk, dv).Run(); err != nil {
		logger.Logger(2, err.Error())
	}

	// Slack への通知処理
	if err := common.Notificator(sc, rc.Title, rc.FileName); err != nil {
		logger.Logger(2, err.Error())
	}

	return err
}
