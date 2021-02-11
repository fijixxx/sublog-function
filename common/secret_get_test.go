package common

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type fakeSecretClient struct {
	SecretClient
	FakeSetSecretString func() string
}

func (sc *fakeSecretClient) SetSecretString(i *secretsmanager.GetSecretValueInput) string {
	return sc.FakeSetSecretString()
}

func TestSecretGet(t *testing.T) {
	sp := &SMParams{
		SName:     "TestSecretName",
		SecretKey: "name",
	}

	sess := session.Must(session.NewSession())
	tc := NewSecretClient(secretsmanager.New(sess, aws.NewConfig().WithRegion(os.Getenv("AWS_REGION"))))
	fc := &fakeSecretClient{
		SecretClient: tc,
		FakeSetSecretString: func() string {
			return `{"name": "testvalue"}`
		},
	}

	t.Run("SecretGet test", func(t *testing.T) {
		si := fc.SecretClient.SetParams(sp)
		ss := fc.SetSecretString(si)
		sv, err := fc.SecretClient.GetSecret(ss, sp)
		if err != nil {
			t.Error(err)
		}
		_ = sv
	})
}
