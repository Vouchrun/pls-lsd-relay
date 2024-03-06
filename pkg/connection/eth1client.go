package connection

import (
	"context"
	"fmt"
	"math/big"
	"net/url"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
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
}

var _ ContractBackend = &Eth1Client{}

type Eth1Client struct {
	clients []*ethclient.Client
}

func NewEth1Client(endpoints []string) (*Eth1Client, error) {
	if len(endpoints) == 0 {
		return nil, fmt.Errorf("endpoints can not be empty")
	}

	clients := make([]*ethclient.Client, len(endpoints))
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

		clients[i] = ethclient.NewClient(rpcClient)
	}

	return &Eth1Client{
		clients,
	}, nil
}

func (c *Eth1Client) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (balance *big.Int, err error) {
	for _, client := range c.clients {
		balance, err = client.BalanceAt(ctx, account, blockNumber)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) BlockByNumber(ctx context.Context, number *big.Int) (block *types.Block, err error) {
	for _, client := range c.clients {
		block, err = client.BlockByNumber(ctx, number)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) ChainID(ctx context.Context) (id *big.Int, err error) {
	for _, client := range c.clients {
		id, err = client.ChainID(ctx)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (nonce uint64, err error) {
	for _, client := range c.clients {
		nonce, err = client.NonceAt(ctx, account, blockNumber)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) BlockNumber(ctx context.Context) (number uint64, err error) {
	for _, client := range c.clients {
		number, err = client.BlockNumber(ctx)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) (bytes []byte, err error) {
	for _, client := range c.clients {
		bytes, err = client.CallContract(ctx, call, blockNumber)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (bytes []byte, err error) {
	for _, client := range c.clients {
		bytes, err = client.CodeAt(ctx, contract, blockNumber)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	for _, client := range c.clients {
		gas, err = client.EstimateGas(ctx, call)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) FilterLogs(ctx context.Context, query ethereum.FilterQuery) (logs []types.Log, err error) {
	for _, client := range c.clients {
		logs, err = client.FilterLogs(ctx, query)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) HeaderByNumber(ctx context.Context, number *big.Int) (header *types.Header, err error) {
	for _, client := range c.clients {
		header, err = client.HeaderByNumber(ctx, number)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) PendingCodeAt(ctx context.Context, account common.Address) (bytes []byte, err error) {
	for _, client := range c.clients {
		bytes, err = client.PendingCodeAt(ctx, account)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) PendingNonceAt(ctx context.Context, account common.Address) (nonce uint64, err error) {
	for _, client := range c.clients {
		nonce, err = client.PendingNonceAt(ctx, account)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) SendTransaction(ctx context.Context, tx *types.Transaction) (err error) {
	for _, client := range c.clients {
		err = client.SendTransaction(ctx, tx)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (sub ethereum.Subscription, err error) {
	for _, client := range c.clients {
		sub, err = client.SubscribeFilterLogs(ctx, query, ch)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) SuggestGasPrice(ctx context.Context) (price *big.Int, err error) {
	for _, client := range c.clients {
		price, err = client.SuggestGasPrice(ctx)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) SuggestGasTipCap(ctx context.Context) (cap *big.Int, err error) {
	for _, client := range c.clients {
		cap, err = client.SuggestGasTipCap(ctx)
		if err == nil {
			return
		}
	}
	return
}

func (c *Eth1Client) WaitTxOkCommon(txHash common.Hash) (blockNumber uint64, err error) {
	for _, client := range c.clients {
		blockNumber, err = waitTxOkCommon(client, txHash)
		if err == nil {
			return
		}
	}
	if err != nil {
		logrus.Errorf("find err: %s, will shutdown.", err.Error())
		utils.ShutdownRequestChannel <- struct{}{}
	}
	return
}

func waitTxOkCommon(client *ethclient.Client, txHash common.Hash) (blockNumber uint64, err error) {
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
