package common

import "os"

// Fcn は 呼び出し元の Lambda 名
var Fcn = os.Getenv("AWS_LAMBDA_FUNCTION_NAME")

// Region は 呼び出し元の Lambda を実行するリージョン
var Region = os.Getenv("AWS_REGION")
