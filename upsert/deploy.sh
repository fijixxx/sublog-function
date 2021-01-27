aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin $AWS_ACCOUNT_ID.dkr.ecr.ap-northeast-1.amazonaws.com
docker build -t upsert -f upsert/Dockerfile .
docker tag upsert:latest $AWS_ACCOUNT_ID.dkr.ecr.ap-northeast-1.amazonaws.com/upsert:latest
docker push $AWS_ACCOUNT_ID.dkr.ecr.ap-northeast-1.amazonaws.com/upsert:latest