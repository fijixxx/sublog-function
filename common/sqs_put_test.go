package common

import (
	"testing"
)

type TestQue struct {
}

func (q *TestQue) SQSPut(p *SQSParams) string {
	s := p.SQSUrl + p.Message.ID
	return s
}

func TestSQSPut(t *testing.T) {
	m := &MsgSQS{
		ID: "ID",
	}
	p := &SQSParams{
		SQSUrl:  "SQSUrl",
		Message: *m,
	}
	q := TestQue{}
	t.Run("SQSPut test", func(t *testing.T) {
		q.SQSPut(p)
	})
}
