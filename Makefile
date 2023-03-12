IMAGE_TAG_SERVERLESS=falconim_serverless
IMAGE_TAG_GO=falconim_go
S3_BUCKET_NAME=falconim
AWS_DEFAULT_REGION=eu-west-2
AWS_LOCAL_COMMAND=aws --endpoint-url=http://localhost:4566 --region=${AWS_DEFAULT_REGION}
SERVERLESS_STAGE=local

# DEV ENV specific

init_dev_env:
	make localstack_up
	make build_image_go
	make init_s3_bucket
	make create_local_secrets
	make verify_local_email
	make build_api
	make deploy_api

build_image_go:
	docker build --target go --tag ${IMAGE_TAG_GO} .

build_image_serverless:
	docker build --target serverless --tag ${IMAGE_TAG_SERVERLESS} .

run_cmd_go:
	docker run -v ${PWD}/src/api:/go/src/app ${IMAGE_TAG_GO} ${cmd}

localstack_up:
	DEBUG=true docker-compose up -d

delete_s3_bucket:
	${AWS_LOCAL_COMMAND} s3 rm s3://${S3_BUCKET_NAME} --recursive

init_s3_bucket:
	${AWS_LOCAL_COMMAND} s3api create-bucket --bucket ${S3_BUCKET_NAME} --create-bucket-configuration LocationConstraint=us-west-2
	${AWS_LOCAL_COMMAND} s3api put-bucket-policy --bucket ${S3_BUCKET_NAME} --policy file://bucket_policy.json

sync_s3_bucket:
	cd src/frontend && ${AWS_LOCAL_COMMAND} s3 sync ./ s3://${S3_BUCKET_NAME}
	${AWS_LOCAL_COMMAND} s3 website s3://${S3_BUCKET_NAME}/ --index-document index.html

build_api:
	docker run -v ${PWD}/src/api:/go/src/app:rw -w /go/src/app/cmd/entrypoint ${IMAGE_TAG_GO} sh -c 'GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o ../../bin/entrypoint/bootstrap -ldflags="-s -w" .'

deploy_api:
	cd src/api && serverless deploy --stage ${SERVERLESS_STAGE} --region eu-west-2

deploy_api_prod:
	cd src/api && serverless deploy --stage prod --region ${AWS_DEFAULT_REGION} --force --verbose

create_local_secrets:
	sh bin/create_secrets.sh -l

verify_local_email:
	sh bin/verify_email.sh -l

# ci specific

docker_build_tag_push:
	make build_docker_image_go
	make build_docker_image_serverless
	docker tag ${IMAGE_TAG_SERVERLESS}:latest ghcr.io/adzfaulkner/${IMAGE_TAG_SERVERLESS}:latest
	docker tag ${IMAGE_TAG_GO}:latest ghcr.io/adzfaulkner/${IMAGE_TAG_GO}:latest
	docker push ghcr.io/adzfaulkner/${IMAGE_TAG_SERVERLESS}:latest
	docker push ghcr.io/adzfaulkner/${IMAGE_TAG_GO}:latest