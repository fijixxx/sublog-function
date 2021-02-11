package common

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/fijixxx/sublog-function/logger"
)

// MsgSQS SQS に送るメッセージ型
type MsgSQS struct {
	ID string `json:"id"`
}

// SQSParams は SQS の URL とメッセージ内容を含む
type SQSParams struct {
	SQSUrl  string
	Message MsgSQS
}

// Que は SQSへの操作関数を持つ
type Que interface {
	SQSPut(p *SQSParams) error
}

// NewQue は SQS 操作用の que 構造体を返す
func NewQue(qc *sqs.SQS) Que {
	return &que{
		client: qc,
	}
}

type que struct {
	client *sqs.SQS
}

// SQSPut SQS へメッセージ送信
func (q *que) SQSPut(p *SQSParams) error {
	qm, err := json.Marshal(&p.Message)
	sqm := string(qm)
	setSQS := &sqs.SendMessageInput{
		MessageBody: aws.String(sqm),
		QueueUrl:    aws.String(p.SQSUrl),
	}

	sqsRes, err := q.client.SendMessage(setSQS)
	logger.Logger(1, "SQSMessageID: "+*sqsRes.MessageId)

	return err
}
