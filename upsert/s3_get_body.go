package upsert

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3GetBody バケット名とオブジェクトキーからオブジェクトを取得して、 String の Body を返す
func S3GetBody(s3c *s3.S3, bn string, ok string) (string, error) {
	obj, err := s3c.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bn),
		Key:    aws.String(ok),
	})
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
	fmt.Printf("[INFO] Object: %v\n", obj)

	objrc := obj.Body
	defer objrc.Close()
	bb, err := ioutil.ReadAll(objrc)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	sb := string(bb)

	return sb, err
}
