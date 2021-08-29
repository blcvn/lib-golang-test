#!/bin/bash

UNIQUE_ACCOUNT_ID=$(date +"%Y%m%d%H%m%s")
UNIQUE_ACCOUNT_ID="0${UNIQUE_ACCOUNT_ID}"

SCRIPT="go run client-app/main.go create_account_network --accountId=${UNIQUE_ACCOUNT_ID}"

bash -c "$SCRIPT"

SCRIPT2="go run client-app/main.go update_account --accountId=${UNIQUE_ACCOUNT_ID}"
bash -c "$SCRIPT2"