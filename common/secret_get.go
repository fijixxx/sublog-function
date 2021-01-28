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

// SecretGet は シークレットマネージャーからシークレットを取得して String で返す
func SecretGet(sc *secretsmanager.SecretsManager, sp *SMParams) (string, error) {
	isc := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(sp.SName),
		VersionStage: aws.String("AWSCURRENT"),
	}

	rsv, err := sc.GetSecretValue(isc)

	ssv := aws.StringValue(rsv.SecretString)
	umsv := make(map[string]interface{})
	err = json.Unmarshal([]byte(ssv), &umsv)
	sv := umsv[sp.SecretKey].(string)

	return sv, err
}
