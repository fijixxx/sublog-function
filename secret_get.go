package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// SecretGet シークレットマネージャーからシークレットを取得
func SecretGet(sc *secretsmanager.SecretsManager, sn string, sk string) (string, error) {
	isc := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(sn),
		VersionStage: aws.String("AWSCURRENT"),
	}

	rsv, err := sc.GetSecretValue(isc)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	ssv := aws.StringValue(rsv.SecretString)
	umsv := make(map[string]interface{})
	if err := json.Unmarshal([]byte(ssv), &umsv); err != nil {
		log.Printf("[ERROR] %v", err)
	}
	sv := umsv[sk].(string)
	return sv, err
}
