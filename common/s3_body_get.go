package common

import (
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fijixxx/sublog-function/logger"
)

// S3Params は S3 クライアントに渡すバケット名とオブジェクトキーを含む
type S3Params struct {
	BucketName string
	ObjectKey  string
}

// S3Client は S3 への取得処理を行う
type S3Client interface {
	SetParams(p *S3Params) *s3.GetObjectInput
	SetBody(i *s3.GetObjectInput) ([]byte, error)
	GetBody(b []byte) string
}

// NewS3Client は 新しい S3 クライアントを返す
func NewS3Client(c *s3.S3) S3Client {
	return &s3client{
		client: c,
	}
}

type s3client struct {
	client *s3.S3
}

func (c *s3client) SetParams(p *S3Params) *s3.GetObjectInput {
	return &s3.GetObjectInput{
		Bucket: aws.String(p.BucketName),
		Key:    aws.String(p.ObjectKey),
	}
}

func (c *s3client) SetBody(i *s3.GetObjectInput) ([]byte, error) {
	obj, err := c.client.GetObject(i)
	if err != nil {
		return nil, err
	}
	objrc := obj.Body
	defer objrc.Close()

	bb, err := ioutil.ReadAll(objrc)
	return bb, err
}

func (c *s3client) GetBody(b []byte) string {
	sb := string(b)
	logger.Logger(1, "Body: "+sb)
	return sb
}
