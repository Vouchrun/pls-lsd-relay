package task

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"reflect"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/config"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

var (
	minAmountNeedStake   = decimal.NewFromBigInt(big.NewInt(31), 18)
	minAmountNeedDeposit = decimal.NewFromBigInt(big.NewInt(32), 18)

	superNodeDepositAmount = decimal.NewFromBigInt(big.NewInt(1), 18)
	superNodeStakeAmount   = decimal.NewFromBigInt(big.NewInt(31), 18)

	blocksOfOneYear = decimal.NewFromInt(2629800)
)

// only support stafi super node account now !!!
// 0. find next key index and cache validator status on start
// 1. update validator status(on execution/ssv/beacon) periodically
// 2. check stakepool balance periodically, call stake/deposit if match
// 3. register validator on ssv, if status is staked on stafi contract
// 4. remove validator on ssv, if status is exited on beacon
type Task struct {
	stop            chan struct{}
	eth1StartHeight uint64
	eth1Endpoint    string
	eth2Endpoint    string

	keyPair *secp256k1.Keypair

	gasLimit    *big.Int
	maxGasPrice *big.Int

	// --- need init on start
	dev bool

	connection *connection.Connection
	eth1Client *ethclient.Client
	eth2Config beacon.Eth2Config

	eth1WithdrawalAdress common.Address

	handlers     []func() error
	handlersName []string

	dealedEth1Block uint64
}

func NewTask(cfg *config.Config, superNodeKeyPair *secp256k1.Keypair) (*Task, error) {
	if !common.IsHexAddress(cfg.Contracts.LsdTokenAddress) {
		return nil, fmt.Errorf("SsvTokenAddress contract address fmt err")
	}

	gasLimitDeci, err := decimal.NewFromString(cfg.GasLimit)
	if err != nil {
		return nil, err
	}

	if gasLimitDeci.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("gas limit is zero")
	}
	maxGasPriceDeci, err := decimal.NewFromString(cfg.MaxGasPrice)
	if err != nil {
		return nil, err
	}
	if maxGasPriceDeci.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("max gas price is zero")
	}

	eth1client, err := ethclient.Dial(cfg.Eth1Endpoint)
	if err != nil {
		return nil, err
	}
	s := &Task{
		stop:            make(chan struct{}),
		eth1Endpoint:    cfg.Eth1Endpoint,
		eth2Endpoint:    cfg.Eth2Endpoint,
		eth1Client:      eth1client,
		keyPair:         superNodeKeyPair,
		gasLimit:        gasLimitDeci.BigInt(),
		maxGasPrice:     maxGasPriceDeci.BigInt(),
		eth1StartHeight: utils.TheMergeBlockNumber,
	}

	return s, nil
}

func (task *Task) Start() error {
	logrus.Info("start...")
	var err error
	task.connection, err = connection.NewConnection(task.eth1Endpoint, task.eth2Endpoint, task.keyPair,
		task.gasLimit, task.maxGasPrice)
	if err != nil {
		return err
	}

	chainId, err := task.eth1Client.ChainID(context.Background())
	if err != nil {
		return err
	}

	task.eth2Config, err = task.connection.Eth2Client().GetEth2Config()
	if err != nil {
		return err
	}

	switch chainId.Uint64() {
	case 1: //mainnet
		task.dev = false
		if !bytes.Equal(task.eth2Config.GenesisForkVersion, params.MainnetConfig().GenesisForkVersion) {
			return fmt.Errorf("endpoint network not match")
		}
		task.dealedEth1Block = 17705353

	case 11155111: // sepolia
		task.dev = true
		if !bytes.Equal(task.eth2Config.GenesisForkVersion, params.SepoliaConfig().GenesisForkVersion) {
			return fmt.Errorf("endpoint network not match")
		}
		task.dealedEth1Block = 9354882
	case 5: // goerli
		task.dev = true
		if !bytes.Equal(task.eth2Config.GenesisForkVersion, params.PraterConfig().GenesisForkVersion) {
			return fmt.Errorf("endpoint network not match")
		}
		task.dealedEth1Block = 9403883

	default:
		return fmt.Errorf("unsupport chainId: %d", chainId.Int64())
	}
	if err != nil {
		return err
	}

	err = task.initContract()
	if err != nil {
		return err
	}

	return nil
}

func (task *Task) Stop() {
	close(task.stop)
}

func (task *Task) initContract() error {

	return nil
}

func (task *Task) appendHandlers(handlers ...func() error) {
	for _, handler := range handlers {

		funcNameRaw := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()

		splits := strings.Split(funcNameRaw, "/")
		funcName := splits[len(splits)-1]

		task.handlersName = append(task.handlersName, funcName)
		task.handlers = append(task.handlers, handler)
	}
}

func (task *Task) waitTxOk(txHash common.Hash) error {
	_, err := utils.WaitTxOkCommon(task.eth1Client, txHash)
	if err != nil {
		return err
	}
	return nil
}
