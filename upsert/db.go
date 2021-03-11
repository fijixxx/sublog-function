package upsert

import (
	"github.com/fijixxx/sublog-function/logger"
	"github.com/guregu/dynamo"
)

// DBFrom は DynamoDB の取得先情報を含む
type DBFrom struct {
	TableName string
	IndexName string
	HashKey   string
}

// DB は DynamoDB に対する処理を定義する
type DB interface {
	DynamoCheckRecord(f *DBFrom, fn string) (Sublog, error)
	DynamoPut(f *DBFrom, i *Sublog) error
}

// NewDB は 新しい DynamoDB クライアントを返す
func NewDB(d *dynamo.DB) DB {
	return &db{
		db: d,
	}
}

type db struct {
	db *dynamo.DB
}

// DynamoCheckRecord は DB 情報とファイル名を受け取って、クエリ結果のレコードを 1 件返す
func (dnm *db) DynamoCheckRecord(f *DBFrom, fn string) (Sublog, error) {
	table := dnm.db.Table(f.TableName)

	var orc Sublog
	if err := table.Get(f.HashKey, fn).Index(f.IndexName).One(&orc); err != nil {
		return Sublog{}, err
	}
	logger.Logger(1, "fileName: "+orc.FileName)

	return orc, nil
}

// DynamoPut は DB 情報と Sublog オブジェクトを受け取って、レコードを挿入する
func (dnm *db) DynamoPut(f *DBFrom, i *Sublog) error {
	table := dnm.db.Table(f.TableName)

	if err := table.Put(i).Run(); err != nil {
		return err
	}

	return nil
}
