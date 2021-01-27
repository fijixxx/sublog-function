package upsert

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// S3Handler S3 関連処理のハンドラー
func S3Handler(sc *secretsmanager.SecretsManager, sess *session.Session, ok string) (string, error) {
	// toml ファイルを格納しているバケット名を取得
	sn := "sublog_assets_bucket_name"
	sk := "name"
	bn, err := SecretGet(sc, sn, sk)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
	fmt.Printf("[INFO] BucketName: %v\n", bn)

	// S3 から toml ファイルの Body を取得
	s3c := s3.New(sess, aws.NewConfig().WithRegion(region))

	// S3 EventNotification から objKey を取得
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
	tb, err := S3GetBody(s3c, bn, ok)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
	fmt.Printf("[INFO] Toml: %v\n", tb)

	return tb, err
}
