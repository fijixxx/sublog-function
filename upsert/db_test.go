package upsert

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type fakeDB struct {
	DB
	FakeDynamoPut func() []byte
}

func (dnm *fakeDB) DynamoPut() []byte {
	return dnm.FakeDynamoPut()
}

func TestDB(t *testing.T) {
	dfr := &DBFrom{
		TableName: "sublog",
		IndexName: "fileName-index",
		HashKey:   "fileName",
	}

	sess := session.Must(session.NewSession())
	dbc := NewDB(dynamo.New(sess, aws.NewConfig().WithRegion("ap-northeast-1")))
	fc := &fakeDB{
		DB: dbc,
		FakeDynamoPut: func() []byte {
			return []byte("Test")
		}}

	t.Run("DB test", func(t *testing.T) {
		_, err := fc.DB.DynamoCheckRecord(dfr, "202012070108")
		if err != nil {
			t.Error(err)
		}
		fc.DynamoPut()
	})
}
