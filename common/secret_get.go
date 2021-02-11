package common

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// SMParams は シークレットマネージャー名とシークレットのキーを含む
type SMParams struct {
	SName     string
	SecretKey string
}

// SecretClient は シークレットマネージャーからシークレットを取得する
type SecretClient interface {
	SetParams(sp *SMParams) *secretsmanager.GetSecretValueInput
	SetSecretString(i *secretsmanager.GetSecretValueInput) (string, error)
	GetSecret(s string, sp *SMParams) (string, error)
}

// NewSecretClient は 新しいシークレットマネージャークライアントを生成する
func NewSecretClient(sc *secretsmanager.SecretsManager) SecretClient {
	return &secretclient{
		client: sc,
	}
}

type secretclient struct {
	client *secretsmanager.SecretsManager
}

// SetParams は シークレット取得時に必要な GetSecretValueInput を設定して返す
func (sc *secretclient) SetParams(p *SMParams) *secretsmanager.GetSecretValueInput {
	return &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(p.SName),
		VersionStage: aws.String("AWSCURRENT"),
	}
}

// SetSecretString は GetSecretValueInput をもらって実際にシークレットを取得して返す
func (sc *secretclient) SetSecretString(i *secretsmanager.GetSecretValueInput) (string, error) {
	rsv, err := sc.client.GetSecretValue(i)
	return aws.StringValue(rsv.SecretString), err
}

// GetSecret は 取得したシークレットの中から必要な文字列を SMParams のキーを使って取り出して返す
func (sc *secretclient) GetSecret(s string, sp *SMParams) (string, error) {
	umsv := make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &umsv)
	sv := umsv[sp.SecretKey].(string)

	return sv, err
}
