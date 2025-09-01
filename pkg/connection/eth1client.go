package connection

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/avast/retry-go/v4"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

type ContractBackend interface {
	bind.ContractBackend
	NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error)
	BlockNumber(ctx context.Context) (uint64, error)
	ChainID(ctx context.Context) (*big.Int, error)
	BlockByNumber(ctx context.Context, number *big.Int) (*types.Block, error)
	BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error)
	WaitTxOkCommon(txHash common.Hash) (blockNumber uint64, err error)
	Debug_TraceBlockByNumber(ctx context.Context, number *big.Int, tracer Tracer) ([]TxResult, error)
}

var _ ContractBackend = &Eth1Client{}

type underlyingEth1Client struct {
	*ethclient.Client
	endpoint string

	latestBlock      *types.Block
	outOfSync        bool
	healthCheckError error
	lastCheckedAt    time.Time
}

type Eth1Client struct {
	clients []*underlyingEth1Client
}

func NewEth1Client(endpoints []string) (*Eth1Client, error) {
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("endpoints can not be empty")
	}

	clients := make([]*underlyingEth1Client, len(endpoints))
	for i, e := range endpoints {
		var rpcClient *rpc.Client
		var err error
		// Start http or ws client
		u, err := url.Parse(e)
		if err != nil {
			return nil, err
		}
		switch u.Scheme {
		case "http", "https":
			rpcClient, err = rpc.DialHTTP(e)
		case "ws", "wss":
			rpcClient, err = rpc.DialWebsocket(context.Background(), e, fmt.Sprintf("/%s", u.Scheme))
		default:
			err = fmt.Errorf("unsupported scheme: %s", u.Scheme)
		}
		if err != nil {
			return nil, err
		}

		client := &underlyingEth1Client{
			Client:   ethclient.NewClient(rpcClient),
			endpoint: e,
		}
		checkHealth(client)
		clients[i] = client
	}

	utils.SafeGoWithRestart(func() {
		for {
			time.Sleep(time.Minute)
			for i := range clients {
				checkHealth(clients[i])
			}
		}
	})

	return &Eth1Client{
		clients,
	}, nil
}

func (c *Eth1Client) getHealthyClients() ([]*underlyingEth1Client, error) {
	clients := make([]*underlyingEth1Client, 0, len(c.clients))
	errMsgs := make([]string, 0, len(c.clients))
	for _, client := range c.clients {
		if client.healthCheckError != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("endpoint: %s checked at %s health check err: %s", client.endpoint, client.lastCheckedAt, client.healthCheckError.Error()))
			continue
		}
		if client.outOfSync {
			errMsgs = append(errMsgs, fmt.Sprintf("endpoint: %s checked at %s latest block number: %d", client.endpoint, client.lastCheckedAt, client.latestBlock.NumberU64()))
			continue
		}

		clients = append(clients, client)
	}
	if len(clients) == 0 {
		return nil, fmt.Errorf("all eth1 endpoints are out of sync: %s", strings.Join(errMsgs, ";"))
	}
	return clients, nil
}

func checkHealth(client *underlyingEth1Client) {
	block, err := retry.DoWithData(
		func() (*types.Block, error) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
			defer cancel()
			return client.BlockByNumber(ctx, nil)
		},
		retry.Delay(time.Second),
		retry.Attempts(5),
	)
	client.healthCheckError = nil
	client.latestBlock = block
	if err != nil {
		client.healthCheckError = err
	} else if client.latestBlock == nil {
		client.healthCheckError = fmt.Errorf("failed to get latest block")
	} else {
		client.outOfSync = block.Time() < uint64(time.Now().Add(-time.Minute*5).Unix())
	}
	client.lastCheckedAt = time.Now()
}

func (c *Eth1Client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (balance *big.Int, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		balance, err = client.BalanceAt(ctx, account, blockNumber)
		if err == nil {
			return
		}
	}
	return
}

type TxTrace struct {
	From    string      `json:"from"`
	Gas     string      `json:"gas"`
	GasUsed string      `json:"gasUsed"`
	To      string      `json:"to"`
	Input   string      `json:"input"`
	Output  string      `json:"output"`
	Calls   []TxTrace   `json:"calls"`
	Value   hexutil.Big `json:"value"`
	// Value   hexutil.Uint64 `json:"value"`
	Type string `json:"type"`
}

type Tracer struct {
	Tracer       string         `json:"tracer"`
	TracerConfig map[string]any `json:"tracerConfig,omitempty"`
}

type TxResult struct {
	TxHash common.Hash
	Result TxTrace
}

func (c *Eth1Client) Debug_TraceBlockByNumber(ctx context.Context, number *big.Int, tracer Tracer) (result []TxResult, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		rpcClient := client.Client.Client()
		if err = rpcClient.CallContext(ctx, &result, "debug_traceBlockByNumber", number, tracer); err == nil {
			return
		}
	}

	return
}

func (c *Eth1Client) BlockByNumber(ctx context.Context, number *big.Int) (block *types.Block, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		block, err = client.BlockByNumber(ctx, number)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) ChainID(ctx context.Context) (id *big.Int, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		id, err = client.ChainID(ctx)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (nonce uint64, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		nonce, err = client.NonceAt(ctx, account, blockNumber)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) BlockNumber(ctx context.Context) (number uint64, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		number, err = client.BlockNumber(ctx)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) (bytes []byte, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		bytes, err = client.CallContract(ctx, call, blockNumber)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (bytes []byte, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		bytes, err = client.CodeAt(ctx, contract, blockNumber)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		gas, err = client.EstimateGas(ctx, call)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) FilterLogs(ctx context.Context, query ethereum.FilterQuery) (logs []types.Log, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		logs, err = client.FilterLogs(ctx, query)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) HeaderByNumber(ctx context.Context, number *big.Int) (header *types.Header, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		header, err = client.HeaderByNumber(ctx, number)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) PendingCodeAt(ctx context.Context, account common.Address) (bytes []byte, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		bytes, err = client.PendingCodeAt(ctx, account)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) PendingNonceAt(ctx context.Context, account common.Address) (nonce uint64, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		nonce, err = client.PendingNonceAt(ctx, account)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) SendTransaction(ctx context.Context, tx *types.Transaction) (err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		err = client.SendTransaction(ctx, tx)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (sub ethereum.Subscription, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		sub, err = client.SubscribeFilterLogs(ctx, query, ch)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) SuggestGasPrice(ctx context.Context) (price *big.Int, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		price, err = client.SuggestGasPrice(ctx)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) SuggestGasTipCap(ctx context.Context) (cap *big.Int, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		cap, err = client.SuggestGasTipCap(ctx)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) WaitTxOkCommon(txHash common.Hash) (blockNumber uint64, err error) {
	var clients []*underlyingEth1Client
	clients, err = c.getHealthyClients()
	if err != nil {
		return
	}

	for _, client := range clients {
		blockNumber, err = waitTxOkCommon(client, txHash)
		if err == nil {
			return
		}
	}
	return
}

func waitTxOkCommon(client *underlyingEth1Client, txHash common.Hash) (blockNumber uint64, err error) {
	retry := 0
	for {
		if retry > utils.RetryLimit {
			return 0, fmt.Errorf("waitTx %s reach retry limit", txHash.String())
		}
		_, pending, err := client.TransactionByHash(context.Background(), txHash)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"hash": txHash.String(),
				"err":  err.Error(),
			}).Warn("TransactionByHash")

			time.Sleep(utils.RetryInterval)
			retry++
			continue
		} else {
			if pending {
				logrus.WithFields(logrus.Fields{
					"hash":    txHash.String(),
					"pending": pending,
				}).Warn("TransactionByHash")

				time.Sleep(utils.RetryInterval)
				retry++
				continue
			} else {
				// check status
				var receipt *types.Receipt
				subRetry := 0
				for {
					if subRetry > utils.RetryLimit {
						return 0, fmt.Errorf("TransactionReceipt %s reach retry limit", txHash.String())
					}

					receipt, err = client.TransactionReceipt(context.Background(), txHash)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"hash": txHash.String(),
							"err":  err.Error(),
						}).Warn("tx TransactionReceipt")

						time.Sleep(utils.RetryInterval)
						subRetry++
						continue
					}
					break
				}

				if receipt.Status == 1 { //success
					blockNumber = receipt.BlockNumber.Uint64()
					break
				} else { //failed
					return 0, fmt.Errorf("tx %s failed", txHash.String())
				}
			}
		}
	}

	logrus.WithFields(logrus.Fields{
		"tx": txHash.String(),
	}).Info("tx send ok")

	return blockNumber, nil
}
