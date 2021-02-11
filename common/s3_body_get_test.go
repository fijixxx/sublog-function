package common

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type fakeS3Client struct {
	S3Client
	FakeSetBody func() []byte
}

func (sc *fakeS3Client) SetBody(i *s3.GetObjectInput) []byte {
	return sc.FakeSetBody()
}

func TestS3BodyGet(t *testing.T) {
	bp := &S3Params{
		BucketName: "BucketName",
		ObjectKey:  "ObjectKey",
	}

	sess := session.Must(session.NewSession())
	s3c := NewS3Client(s3.New(sess, aws.NewConfig().WithRegion(Region)))
	fc := &fakeS3Client{
		S3Client: s3c,
		FakeSetBody: func() []byte {
			return []byte("Test")
		}}

	t.Run("S3BodyGet test", func(t *testing.T) {
		si := fc.S3Client.SetParams(bp)
		ss := fc.SetBody(si)
		sv := fc.S3Client.GetBody(ss)
		_ = sv
	})
}
