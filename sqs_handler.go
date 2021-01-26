package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSHandler SQS 関連処理のハンドラー
func SQSHandler(sc *secretsmanager.SecretsManager, sess *session.Session, region string, iid string) error {
	// SQS 送信先取得
	sn := "sublogHighlighterSQS"
	sk := "url"
	qu, err := SecretGet(sc, sn, sk)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	// SQS に記事 ID を送信
	qc := sqs.New(sess, aws.NewConfig().WithRegion(region))
	if err := SQSPut(qc, iid, qu); err != nil {
		log.Printf("[ERROR] %v", err)
	}

	return err
}
