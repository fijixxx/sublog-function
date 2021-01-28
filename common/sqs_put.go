package common

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
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

// SQSPut SQS へメッセージ送信
func SQSPut(qc *sqs.SQS, sp *SQSParams) error {
	qm, err := json.Marshal(&sp.Message.ID)
	sqm := string(qm)
	setSQS := &sqs.SendMessageInput{
		MessageBody: aws.String(sqm),
		QueueUrl:    aws.String(sp.SQSUrl),
	}

	sqsRes, err := qc.SendMessage(setSQS)
	Logger(0, "SQSMessageID: "+*sqsRes.MessageId)

	return err
}
