package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// SQSPutter SQS 関連処理のハンドラー
func SQSPutter(sc *secretsmanager.SecretsManager, sess *session.Session, iid string) error {
	// SQS 送信先取得
	sp := &SMParams{
		SName:     "sublogHighlighterSQS",
		SecretKey: "url",
	}
	svc := NewSecretClient(sc)
	si := svc.SetParams(sp)
	ss, err := svc.SetSecretString(si)
	if err != nil {
		return err
	}
	qu, err := svc.GetSecret(ss, sp)
	if err != nil {
		return err
	}

	// SQS に記事 ID を送信
	msg := &MsgSQS{
		ID: iid,
	}
	qp := &SQSParams{
		SQSUrl:  qu,
		Message: *msg,
	}
	q := NewQue(sqs.New(sess, aws.NewConfig().WithRegion(Region)))
	err = q.SQSPut(qp)

	return err
}
