// Copyright 2024 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package connection

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"sync"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/forta-network/go-multicall"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon/client"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/types"
	gomicrobee "github.com/stephennancekivell/go-micro-bee"
)

var Gwei5 = big.NewInt(5e9)
var Gwei10 = big.NewInt(10e9)
var Gwei20 = big.NewInt(20e9)

type Connection struct {
	eth1Endpoint string
	eth2Endpoint string
	kp           *secp256k1.Keypair
	gasLimit     *big.Int
	maxGasPrice  *big.Int
	eth1Client   *ethclient.Client
	eth1Rpc      *rpc.Client
	eth2Client   *client.StandardHttpClient
	txOpts       *bind.TransactOpts
	callOpts     bind.CallOpts
	optsLock     sync.Mutex
	multiCaller  *multicall.Caller

	latestMultiCallMicrobeeSystem       gomicrobee.System[*multicall.Call, *MultiCall]
	latestValicatorStatusMicrobeeSystem gomicrobee.System[types.ValidatorPubkey, BeaconValidatorStatus]
}

// NewConnection returns an uninitialized connection, must call Connection.Connect() before using.
func NewConnection(eth1Endpoint, eth2Endpoint string, kp *secp256k1.Keypair, gasLimit, maxGasPrice *big.Int) (*Connection, error) {
	if kp != nil {
		if maxGasPrice.Cmp(big.NewInt(0)) <= 0 {
			return nil, fmt.Errorf("max gas price empty")
		}
		if gasLimit.Cmp(big.NewInt(0)) <= 0 {
			return nil, fmt.Errorf("gas limit empty")
		}
	}
	c := &Connection{
		eth1Endpoint: eth1Endpoint,
		eth2Endpoint: eth2Endpoint,
		kp:           kp,
		gasLimit:     gasLimit,
		maxGasPrice:  maxGasPrice,
	}

	err := retry.Do(c.connect, retry.Delay(time.Second), retry.Attempts(3))
	if err != nil {
		return nil, err
	}

	if err = c.initMulticall(); err != nil {
		return nil, err
	}
	c.initBatchBeaconChain()

	return c, nil
}

// Connect starts the ethereum WS connection
func (c *Connection) connect() error {
	var rpcClient *rpc.Client
	var err error
	// Start http or ws client
	u, err := url.Parse(c.eth1Endpoint)
	if err != nil {
		return err
	}
	switch u.Scheme {
	case "http", "https":
		rpcClient, err = rpc.DialHTTP(c.eth1Endpoint)
	case "ws", "wss":
		rpcClient, err = rpc.DialWebsocket(context.Background(), c.eth1Endpoint, fmt.Sprintf("/%s", u.Scheme))
	default:
		err = fmt.Errorf("unsupport scheme: %s", u.Scheme)
	}
	if err != nil {
		return err
	}

	c.eth1Client = ethclient.NewClient(rpcClient)

	c.eth1Rpc = rpcClient

	chainId, err := c.eth1Client.ChainID(context.Background())
	if err != nil {
		return err
	}

	// eth2 client
	if len(c.eth2Endpoint) != 0 {
		c.eth2Client, err = client.NewStandardHttpClient(c.eth2Endpoint, chainId)
		if err != nil {
			return err
		}
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

func (c *Connection) Eth1Client() *ethclient.Client {
	return c.eth1Client
}

func (c *Connection) Eth2Client() *client.StandardHttpClient {
	return c.eth2Client
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

// EnsureHasBytecode asserts if contract code exists at the specified address
func (c *Connection) EnsureHasBytecode(addr common.Address) error {
	code, err := c.eth1Client.CodeAt(context.Background(), addr, nil)
	if err != nil {
		return err
	}

	if len(code) == 0 {
		return fmt.Errorf("no bytecode found at %s", addr.Hex())
	}
	return nil
}

func (c *Connection) Eth2BeaconHead() (beacon.BeaconHead, error) {
	return c.eth2Client.GetBeaconHead()
}

func (c *Connection) GetValidatorStatus(pubkey types.ValidatorPubkey, opts *beacon.ValidatorStatusOptions) (beacon.ValidatorStatus, error) {
	return c.eth2Client.GetValidatorStatus(pubkey, opts)
}

func (c *Connection) GetValidatorStatuses(pubkeys []types.ValidatorPubkey, opts *beacon.ValidatorStatusOptions) (map[types.ValidatorPubkey]beacon.ValidatorStatus, error) {
	return c.eth2Client.GetValidatorStatuses(pubkeys, opts)
}

func (c *Connection) GetBeaconBlock(blockId uint64) (beacon.BeaconBlock, bool, error) {
	return c.eth2Client.GetBeaconBlock(blockId)
}

func (c *Connection) GetValidatorStatusByIndex(index string, opts *beacon.ValidatorStatusOptions) (beacon.ValidatorStatus, error) {
	return c.eth2Client.GetValidatorStatusByIndex(index, opts)
}
