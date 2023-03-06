#!/bin/sh

cmd="aws"
while getopts l flag
do
    case "${flag}" in
        l) cmd="aws --endpoint-url=http://localhost:4566";;
    esac
done

export $(cat ./src/api/.env | grep -v '#' | sed 's/\r$//' | awk '/=/ {print $1}' )

${cmd} ses verify-email-identity --email-address ${EMAIL_TO}