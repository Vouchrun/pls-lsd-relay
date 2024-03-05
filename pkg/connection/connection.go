// Copyright 2024 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package connection

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/forta-network/go-multicall"
	"github.com/samber/lo"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/config"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon/client"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/types"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/gomicrobee"
)

var Gwei5 = big.NewInt(5e9)
var Gwei10 = big.NewInt(10e9)
var Gwei20 = big.NewInt(20e9)

type Connection struct {
	endpoints   []config.Endpoint
	kp          *secp256k1.Keypair
	gasLimit    *big.Int
	maxGasPrice *big.Int

	eth1Client  ContractBackend
	eth2Clients []*client.StandardHttpClient

	txOpts      *bind.TransactOpts
	callOpts    bind.CallOpts
	optsLock    sync.Mutex
	multiCaller *multicall.Caller

	latestMultiCallMicrobeeSystem gomicrobee.System[*multicall.Call, *MultiCall]
}

// NewConnection returns an uninitialized connection, must call Connection.Connect() before using.
func NewConnection(endpoints []config.Endpoint, kp *secp256k1.Keypair, gasLimit, maxGasPrice *big.Int) (*Connection, error) {
	if kp != nil {
		if maxGasPrice.Cmp(big.NewInt(0)) <= 0 {
			return nil, fmt.Errorf("max gas price empty")
		}
		if gasLimit.Cmp(big.NewInt(0)) <= 0 {
			return nil, fmt.Errorf("gas limit empty")
		}
	}
	c := &Connection{
		endpoints:   endpoints,
		kp:          kp,
		gasLimit:    gasLimit,
		maxGasPrice: maxGasPrice,
	}

	err := retry.Do(c.connect, retry.Delay(time.Second), retry.Attempts(3))
	if err != nil {
		return nil, err
	}

	if err = c.initMulticall(); err != nil {
		return nil, err
	}

	return c, nil
}

// Connect starts the ethereum WS connection
func (c *Connection) connect() error {
	if err := c.connectEth1(); err != nil {
		return err
	}

	chainId, err := c.eth1Client.ChainID(context.Background())
	if err != nil {
		return err
	}

	if err = c.connectEth2(chainId); err != nil {
		return err
	}

	if c.kp != nil {
		// Construct tx opts, call opts, and nonce mechanism
		opts, err := c.newTransactOpts(big.NewInt(0), c.gasLimit)
		if err != nil {
			return err
		}
		c.txOpts = opts
		c.callOpts = bind.CallOpts{Pending: false, From: c.kp.CommonAddress(), BlockNumber: nil, Context: context.Background()}
	} else {
		c.callOpts = bind.CallOpts{Pending: false, From: common.Address{}, BlockNumber: nil, Context: context.Background()}
	}
	return nil
}

func (c *Connection) connectEth1() (err error) {
	c.eth1Client, err = NewEth1Client(lo.Map(c.endpoints, func(e config.Endpoint, i int) string { return e.Eth1 }))
	return
}

func (c *Connection) connectEth2(chainId *big.Int) error {
	c.eth2Clients = make([]*client.StandardHttpClient, 0, len(c.endpoints))
	for _, e := range c.endpoints {
		client, err := client.NewStandardHttpClient(e.Eth2, chainId)
		if err != nil {
			return err
		}
		c.eth2Clients = append(c.eth2Clients, client)
	}
	return nil
}

// newTransactOpts builds the TransactOpts for the connection's keypair.
func (c *Connection) newTransactOpts(value, gasLimit *big.Int) (*bind.TransactOpts, error) {
	privateKey := c.kp.PrivateKey()
	address := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := c.eth1Client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, err
	}
	chainId, err := c.eth1Client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = uint64(gasLimit.Int64())
	auth.Context = context.Background()

	return auth, nil
}

func (c *Connection) Keypair() *secp256k1.Keypair {
	return c.kp
}

func (c *Connection) Eth1Client() ContractBackend {
	return c.eth1Client
}

func (c *Connection) TxOpts() *bind.TransactOpts {
	return c.txOpts
}

func (c *Connection) MultiCaller() *multicall.Caller {
	return c.multiCaller
}

func (c *Connection) CallOpts(blocknumber *big.Int) *bind.CallOpts {
	newCallOpts := c.callOpts
	newCallOpts.BlockNumber = blocknumber
	return &newCallOpts
}

func (c *Connection) CallOptsOn(targetBlockNumber uint64) *bind.CallOpts {
	newCallOpts := c.callOpts
	newCallOpts.BlockNumber = big.NewInt(int64(targetBlockNumber))
	return &newCallOpts
}

// return suggest gastipcap gasfeecap
func (c *Connection) SafeEstimateFee(ctx context.Context) (*big.Int, *big.Int, error) {
	gasTipCap, err := c.eth1Client.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, nil, err
	}
	gasFeeCap, err := c.eth1Client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, err
	}

	if gasFeeCap.Cmp(Gwei20) < 0 {
		gasFeeCap = new(big.Int).Add(gasFeeCap, Gwei5)
	} else {
		gasFeeCap = new(big.Int).Add(gasFeeCap, Gwei10)
	}

	if gasFeeCap.Cmp(c.maxGasPrice) > 0 {
		gasFeeCap = c.maxGasPrice
	}

	return gasTipCap, gasFeeCap, nil
}

// LockAndUpdateOpts acquires a lock on the opts before updating the nonce
// and gas price.
func (c *Connection) LockAndUpdateTxOpts() error {
	c.optsLock.Lock()

	gasTipCap, gasFeeCap, err := c.SafeEstimateFee(context.Background())
	if err != nil {
		c.optsLock.Unlock()
		return err
	}
	c.txOpts.GasTipCap = gasTipCap
	c.txOpts.GasFeeCap = gasFeeCap

	nonce, err := c.eth1Client.NonceAt(context.Background(), c.txOpts.From, nil)
	if err != nil {
		c.optsLock.Unlock()
		return err
	}
	c.txOpts.Nonce.SetUint64(nonce)
	return nil
}

func (c *Connection) UnlockTxOpts() {
	c.optsLock.Unlock()
}

// LatestBlock returns the latest block from the current chain
func (c *Connection) Eth1LatestBlock() (uint64, error) {
	header, err := c.eth1Client.BlockNumber(context.Background())
	if err != nil {
		return 0, err
	}
	return header, nil
}

func (c *Connection) GetValidatorStatus(ctx context.Context, pubkey types.ValidatorPubkey, opts *beacon.ValidatorStatusOptions) (validatorStatus beacon.ValidatorStatus, err error) {
	for _, client := range c.eth2Clients {
		validatorStatus, err = client.GetValidatorStatus(ctx, pubkey, opts)
		if err == nil {
			return
		}
	}
	return
}

func (c *Connection) GetValidatorStatuses(ctx context.Context, pubkeys []types.ValidatorPubkey, opts *beacon.ValidatorStatusOptions) (validatorStatus map[types.ValidatorPubkey]beacon.ValidatorStatus, err error) {
	for _, client := range c.eth2Clients {
		validatorStatus, err = client.GetValidatorStatuses(ctx, pubkeys, opts)
		if err == nil {
			return
		}
	}
	return
}

func (c *Connection) GetBeaconBlock(blockId uint64) (block beacon.BeaconBlock, exist bool, err error) {
	for _, client := range c.eth2Clients {
		block, exist, err = client.GetBeaconBlock(blockId)
		if exist {
			return
		}
	}
	return
}

func (c *Connection) GetEth2Config() (cfg beacon.Eth2Config, err error) {
	for _, client := range c.eth2Clients {
		cfg, err = client.GetEth2Config()
		if err == nil {
			return
		}
	}
	return
}

func (c *Connection) GetBeaconHead() (head beacon.BeaconHead, err error) {
	for _, client := range c.eth2Clients {
		head, err = client.GetBeaconHead()
		if err == nil {
			return
		}
	}
	return
}
