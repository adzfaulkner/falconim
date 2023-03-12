#!/bin/sh

cmd="aws"
while getopts l flag
do
    case "${flag}" in
        l) cmd="aws --endpoint-url=http://localhost:4566";;
    esac
done

export $(cat ./src/api/.env | grep -v '#' | sed 's/\r$//' | awk '/=/ {print $1}' )

${cmd} \
  ssm put-parameter \
  --name "/FalconIM/EMAIL_TO" \
  --type "SecureString" \
  --value ${EMAIL_TO} \
  --overwrite

${cmd} \
  ssm put-parameter \
  --name "/FalconIM/EMAIL_FROM" \
  --type "SecureString" \
  --value ${EMAIL_FROM} \
  --overwrite

${cmd} \
  ssm put-parameter \
  --name "/FalconIM/RECAPTCHA_SECRET" \
  --type "SecureString" \
  --value ${RECAPTCHA_SECRET} \
  --overwrite