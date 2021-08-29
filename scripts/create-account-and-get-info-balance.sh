#!/bin/bash

UNIQUE_ACCOUNT_ID=$(date +"%Y%m%d%H%m%s")
UNIQUE_ACCOUNT_ID="2${UNIQUE_ACCOUNT_ID}"

SCRIPT="go run client-app/main.go create_account_member --accountId=${UNIQUE_ACCOUNT_ID}"

bash -c "$SCRIPT"

SCRIPT2="go run client-app/main.go get_account_info_balance --accountId=${UNIQUE_ACCOUNT_ID}"

bash -c "$SCRIPT2"