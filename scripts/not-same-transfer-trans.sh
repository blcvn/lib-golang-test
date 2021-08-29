#!/bin/bash

UNIQUE_ACCOUNT_ID=$(date +"%Y%m%d%H%m%s")
UNIQUE_ACCOUNT_ID="2${UNIQUE_ACCOUNT_ID}"

UNIQUE_ACCOUNT_MB_ID="2${UNIQUE_ACCOUNT_ID}"

echo $UNIQUE_ACCOUNT_ID
echo $UNIQUE_ACCOUNT_MB_ID


SCRIPT="go run client-app/main.go create_account_network --accountId=${UNIQUE_ACCOUNT_ID}"

bash -c "$SCRIPT"

SCRIPT="go run client-app/main.go create_account_member --accountId=${UNIQUE_ACCOUNT_MB_ID}"

bash -c "$SCRIPT"


SCRIPT="go run client-app/main.go move_account --accountId=${UNIQUE_ACCOUNT_MB_ID}"

bash -c "$SCRIPT"


SCRIPT2="go run client-app/main.go trans_credit --accountId=${UNIQUE_ACCOUNT_ID}"
bash -c "$SCRIPT2"


SCRIPT3="go run client-app/main.go trans_transfer --accountId=${UNIQUE_ACCOUNT_ID} --receiveId=${UNIQUE_ACCOUNT_MB_ID}"
bash -c "$SCRIPT3"

SCRIPT3="go run client-app/main.go get_account_info_balance --accountId=${UNIQUE_ACCOUNT_ID}"
bash -c "$SCRIPT3"

SCRIPT3="go run client-app/main.go get_account_info_balance --accountId=${UNIQUE_ACCOUNT_MB_ID}"
bash -c "$SCRIPT3"