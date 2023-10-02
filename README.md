# Lib Golang Test


# Folder structure: 
## Connect:
## Consensus:
### Blocks 
### peer:
   
### raft: 
### raft-grpc:
### raft3: 
This project collected all test cases to test lib which will be used in products

## Config git with private account:
git config --global url."git@github.com:".insteadOf "https://github.com/"
export GOPRIVATE=github.com/binhnt-teko/*


## Test client app

### Prepare

1. Set all `ServiceTest` to `false` at `loyalty-client-app/config/config.yaml`
2. Make sure we have `network_docker` folder at `loyalty-client-app/network_docker`
same as config in `loyalty-client-app/config/fabric-config.json`
3. Make sure our network running have same chaincode identify at `loyalty-client-app/config/config.yaml`
```text
channel1:
   chaincode.name: "loyalty_cc_3"
   chaincode.lang: "golang"
   channel.id: "vnpay-channel-1"
channel2:
   chaincode.name: "loyalty_cc_3"
   chaincode.lang: "golang"
   channel.id: "vnpay-channel-2"
channel3:
   chaincode.name: "loyalty_cc_3"
   chaincode.lang: "golang"
   channel.id: "vnpay-channel-3"
```
4. Please remember i using `03002` as branchId which mapping to `channel2`
This is old mapping, please change to branchId which is mapping to `channel2` at your environment.

5. Start the generator `make start id=3`

### Testing account
1. make scripts directory +x `chmod +x scripts/*`
2. Run first test with `go run client-app/main.go create_account`
3. Start with case to test as below
4. Create new account `make create-account-new` which create a new account, should be success
5. Create 2 account with same accountId `make create-account-existed`, 2nd response must be empty.
6. Create new account and get info `make create-account-and-get-info`, 2nd response must contain new account information which created.
7. Create new account and get info balance list `make create-account-and-get-info-balance`, 2nd response must contain new account information.
8. Create new account and update its state to `inactive` : `make update-account`, 2nd response must contain success message.
9. Create new account by branchId `101`, get channel info, and move account to other channel and get response channel info again : `make move-account`
### Testing transaction

1. Credit to network Accounts `make credit-trans`
2. Create account network, create account member, credit to account network, transfer to account member and get balance both `make transfer-trans`
3. Create account network, create account member, move account member to other channel, credit to account network, transfer to account member and get balance both `make not-same-transfer-trans`


# Conflict protobuf namespace 

GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn


# Export to gen 

export PATH=$PATH:~/go/bin/