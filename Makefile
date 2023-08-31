VERSION := $(shell git describe --tags)
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
	go build -mod readonly $(BUILD_FLAGS) -o build/relay main.go

build-linux:
	@GOOS=linux GOARCH=amd64 go build --mod readonly $(BUILD_FLAGS) -o ./build/relay main.go

install:
	@echo " > \033[32mInstalling relay...\033[0m "
	go install -mod readonly $(BUILD_FLAGS) ./...

abi:
	@echo " > \033[32mGenabi...\033[0m "
	abigen --abi ./bindings/Erc20/erc20_abi.json --pkg erc20 --type Erc20 --out ./bindings/Erc20/Erc20.go
	abigen --abi ./bindings/DepositContract/depositcontract_abi.json --pkg deposit_contract --type DepositContract --out ./bindings/DepositContract/DepositContract.go
	abigen --abi ./bindings/LsdNetworkFactory/lsdnetworkfactory_abi.json --pkg lsd_network_factory --type LsdNetworkFactory --out ./bindings/LsdNetworkFactory/LsdNetworkFactory.go
	abigen --abi ./bindings/NetworkWithdraw/networkwithdraw_abi.json --pkg network_withdraw --type NetworkWithdraw --out ./bindings/NetworkWithdraw/NetworkWithdraw.go
	abigen --abi ./bindings/NodeDeposit/nodedeposit_abi.json --pkg node_deposit --type NodeDeposit --out ./bindings/NodeDeposit/NodeDeposit.go
	abigen --abi ./bindings/NetworkProposal/networkproposal_abi.json --pkg network_proposal --type NetworkProposal --out ./bindings/NetworkProposal/NetworkProposal.go
	abigen --abi ./bindings/NetworkBalances/networkbalances_abi.json --pkg network_balances --type NetworkBalances --out ./bindings/NetworkBalances/NetworkBalances.go
	abigen --abi ./bindings/UserDeposit/userdeposit_abi.json --pkg user_deposit --type UserDeposit --out ./bindings/UserDeposit/UserDeposit.go


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

.PHONY: all lint test race msan tools clean build
