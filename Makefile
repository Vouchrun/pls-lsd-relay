VERSION := $(shell git describe --tags --always)
COMMIT  := $(shell git log -1 --format='%H')

all: build

LD_FLAGS = -X github.com/stafiprotocol/eth-lsd-relay/cmd.Version=$(VERSION) \
	-X github.com/stafiprotocol/eth-lsd-relay/cmd.Commit=$(COMMIT) \

BUILD_FLAGS := -ldflags '$(LD_FLAGS)'

get:
	@echo "  >  \033[32mDownloading & Installing all the modules...\033[0m "
	go mod tidy && go mod download

build:
	@echo " > \033[32mBuilding relay...\033[0m "
	go build -mod readonly $(BUILD_FLAGS) -o build/lsd-relay main.go

build-linux:
	@GOOS=linux GOARCH=amd64 go build --mod readonly $(BUILD_FLAGS) -o ./build/lsd-relay main.go

abi:
	@echo " > \033[32mGenabi...\033[0m "
	abigen --abi ./bindings/Erc20/erc20_abi.json --pkg erc20 --type Erc20 --out ./bindings/Erc20/Erc20.go
	abigen --abi ./bindings/DepositContract/depositcontract_abi.json --pkg deposit_contract --type DepositContract --out ./bindings/DepositContract/DepositContract.go
	abigen --abi ./bindings/LsdNetworkFactory/lsdnetworkfactory_abi.json --pkg lsd_network_factory --type LsdNetworkFactory --out ./bindings/LsdNetworkFactory/LsdNetworkFactory.go
	abigen --abi ./bindings/NetworkWithdrawal/networkwithdrawal_abi.json --pkg network_withdrawal --type NetworkWithdrawal --out ./bindings/NetworkWithdrawal/NetworkWithdrawal.go
	abigen --abi ./bindings/NodeDeposit/nodedeposit_abi.json --pkg node_deposit --type NodeDeposit --out ./bindings/NodeDeposit/NodeDeposit.go
	abigen --abi ./bindings/NetworkProposal/networkproposal_abi.json --pkg network_proposal --type NetworkProposal --out ./bindings/NetworkProposal/NetworkProposal.go
	abigen --abi ./bindings/NetworkBalances/networkbalances_abi.json --pkg network_balances --type NetworkBalances --out ./bindings/NetworkBalances/NetworkBalances.go
	abigen --abi ./bindings/UserDeposit/userdeposit_abi.json --pkg user_deposit --type UserDeposit --out ./bindings/UserDeposit/UserDeposit.go
	abigen --abi ./bindings/FeePool/feepool_abi.json --pkg fee_pool --type FeePool --out ./bindings/FeePool/FeePool.go

# make abi_json abi contracts_repo_path=../eth-lsd-contracts
abi_json:
	cat $(contracts_repo_path)/artifacts/contracts/LsdNetworkFactory.sol/LsdNetworkFactory.json | jq '.abi' > ./bindings/LsdNetworkFactory/lsdnetworkfactory_abi.json
	cat $(contracts_repo_path)/artifacts/contracts/NetworkWithdrawal.sol/NetworkWithdrawal.json | jq '.abi' > ./bindings/NetworkWithdrawal/networkwithdrawal_abi.json
	cat $(contracts_repo_path)/artifacts/contracts/NodeDeposit.sol/NodeDeposit.json         | jq '.abi' > ./bindings/NodeDeposit/nodedeposit_abi.json
	cat $(contracts_repo_path)/artifacts/contracts/NetworkProposal.sol/NetworkProposal.json | jq '.abi' > ./bindings/NetworkProposal/networkproposal_abi.json
	cat $(contracts_repo_path)/artifacts/contracts/NetworkBalances.sol/NetworkBalances.json | jq '.abi' > ./bindings/NetworkBalances/networkbalances_abi.json
	cat $(contracts_repo_path)/artifacts/contracts/UserDeposit.sol/UserDeposit.json         | jq '.abi' > ./bindings/UserDeposit/userdeposit_abi.json
	cat $(contracts_repo_path)/artifacts/contracts/FeePool.sol/FeePool.json                 | jq '.abi' > ./bindings/FeePool/feepool_abi.json

clean:
	@echo " > \033[32mCleanning build files ...\033[0m "
	rm -rf build
fmt :
	@echo " > \033[32mFormatting go files ...\033[0m "
	go fmt ./...

swagger:
	@echo "  >  \033[32mBuilding swagger docs...\033[0m "
	swag init --parseDependency

get-lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest

lint:
	golangci-lint run ./... --skip-files ".+_test.go"

.PHONY: all lint test race msan tools clean build abi_json abi
