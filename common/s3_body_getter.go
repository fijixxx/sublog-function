package common

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// S3BodyGetter はバケット内のオブジェクトの Body を String で返す
func S3BodyGetter(sc *secretsmanager.SecretsManager, sess *session.Session, ok string) (string, error) {
	// toml ファイルを格納しているバケット実名を SM から取得
	sp := &SMParams{
		SName:     "sublog_assets_bucket_name",
		SecretKey: "name",
	}
	svc := NewSecretClient(sc)
	si := svc.SetParams(sp)
	ss, err := svc.SetSecretString(si)
	if err != nil {
		return "", err
	}
	bn, err := svc.GetSecret(ss, sp)
	if err != nil {
		return "", err
	}

	// S3 から toml ファイルの Body を取得
	s3c := s3.New(sess, aws.NewConfig().WithRegion(Region))

	bp := &S3Params{
		BucketName: bn,
		ObjectKey:  ok,
	}

	c := NewS3Client(s3c)
	sip := c.SetParams(bp)
	sb, err := c.SetBody(sip)
	if err != nil {
		return "", err
	}
	tb := c.GetBody(sb)

	return tb, nil
}
