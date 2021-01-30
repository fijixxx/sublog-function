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

// S3GetBody は バケット名とオブジェクトキーからオブジェクトを取得して String の Body を返す
func S3GetBody(s3c *s3.S3, sp *S3Params) (string, error) {
	obj, err := s3c.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(sp.BucketName),
		Key:    aws.String(sp.ObjectKey),
	})

	objrc := obj.Body
	defer objrc.Close()
	bb, err := ioutil.ReadAll(objrc)

	sb := string(bb)
	logger.Logger(1, "Body: "+sb)

	return sb, err
}
