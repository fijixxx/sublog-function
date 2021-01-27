package upsert

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// MsgSQS SQS に送るメッセージ型
type MsgSQS struct {
	ID string `json:"id"`
}

// SQSPut SQS へメッセージ送信
func SQSPut(qc *sqs.SQS, mb string, qu string) error {
	var msg MsgSQS
	msg.ID = mb

	jmsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
	smsg := string(jmsg)

	setSQS := &sqs.SendMessageInput{
		MessageBody: aws.String(smsg),
		QueueUrl:    aws.String(qu),
	}

	sqsRes, err := qc.SendMessage(setSQS)
	if err != nil {
		fmt.Printf("[ERROR] %v", err)
	} else {
		fmt.Print("[INFO] SQSMessageID", *sqsRes.MessageId)
	}
	fmt.Printf("[INFO] sqsRes: %v\n", sqsRes)

	return err
}
