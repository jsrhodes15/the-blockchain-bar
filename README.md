#The blockchain bar
Learning some blockchain and Go things

> A block chain built from scratch in Go
 - Fully functional, accepts transactions and syncs all nodes connected to network
 - Accounts in the form of 'wallets'
 - HTTP API server to host full blockchain
 - CLI built using Cobra package

##Usage
clone the repository & cd into project

####Install the CLI
```bash
go install ./cmd/...
```

Show all commands
```bash
tbb help
```

Start a local server
```bash
tbb run --dataDir=[/absolute/path/to/dir]
# 'dataDir' sets you want config stored, defaults to $HOME/.tbb
```

Get balances

*CLI*
```bash
tbb balances list
```
*API*
```
curl http://localhost:8080/balances/list | jq
# 'jq' is for formatting, if you don't have it, can omit
```

Add a Transaction  
*_the first time you do this, you will want to use "jrhodes" as from, as that is the only account with "coins" to transfer_
*CLI*
```bash
tbb tx add --from=[from acct] --to=[to acct] --value=[value]
```
*API*
```bash
curl --location --request POST --header "Content-Type: application/json" --data '{"from":"[someAccount]","to":"[someAccount]","value":[someNumber]}' http://localhost:8080/tx/add  
```