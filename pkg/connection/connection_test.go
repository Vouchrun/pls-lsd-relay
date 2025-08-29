package connection_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strings"
	"testing"
	"time"

	// "github.com/ethereum/go-ethereum/common"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	network_withdraw "github.com/stafiprotocol/eth-lsd-relay/bindings/NetworkWithdraw"
	node_deposit "github.com/stafiprotocol/eth-lsd-relay/bindings/NodeDeposit"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/config"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
	"github.com/stafiprotocol/eth-lsd-relay/service"
	"github.com/stretchr/testify/assert"
)

func TestCallOpts(t *testing.T) {
	endpoints := []config.Endpoint{
		{Eth1: os.Getenv("ETH1_ENDPOINT"), Eth2: os.Getenv("ETH2_ENDPOINT")},
	}
	c, err := connection.NewConnection(endpoints, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	oldopts := c.CallOpts(nil)
	t.Log(oldopts)
	newopts := c.CallOpts(big.NewInt(5))
	t.Log(oldopts)
	t.Log(newopts)

	newopts2 := c.CallOpts(big.NewInt(7))
	t.Log(oldopts)
	t.Log(newopts)
	t.Log(newopts2)

	gasPrice, err := c.Eth1Client().SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	gasTip, err := c.Eth1Client().SuggestGasTipCap(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(gasPrice.String(), gasTip.String())
}

func TestSafeEstimateFee(t *testing.T) {
	endpoints := []config.Endpoint{
		{Eth1: os.Getenv("ETH1_ENDPOINT"), Eth2: os.Getenv("ETH2_ENDPOINT")},
	}

	maxGasPriceDeci := decimal.RequireFromString("1100").Mul(utils.GweiDeci)
	c, err := connection.NewConnection(endpoints, nil, nil, maxGasPriceDeci.BigInt(), new(big.Float).SetFloat64(2))
	if err != nil {
		t.Fatal(err)
	}

	gasTipCap, gasFeeCap, err := c.SafeEstimateFee(context.Background())
	if err != nil {
		newErr := fmt.Errorf("LockAndUpdateTxOpts err: %w", err)
		var gasErr *connection.GasPriceError
		if errors.As(newErr, &gasErr) {
			t.Fatal("gas error:", gasErr)
		} else {
			t.Fatal("not gas err", err)
		}
	}
	fmt.Println(gasTipCap)
	fmt.Println(gasFeeCap)
}

func TestBlockReward(t *testing.T) {
	endpoints := []config.Endpoint{
		{Eth1: os.Getenv("ETH1_ENDPOINT"), Eth2: os.Getenv("ETH2_ENDPOINT")},
	}
	c, err := connection.NewConnection(endpoints, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	eth1Block, err := c.Eth1Client().BlockByNumber(context.Background(), big.NewInt(859542))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v", eth1Block.Coinbase())

	eth1Block, err = c.Eth1Client().BlockByNumber(context.Background(), big.NewInt(859543))
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v", eth1Block.Coinbase())
}

func TestEth2Config(t *testing.T) {
	s := make([]int64, 0)
	sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })

	logrus.SetLevel(logrus.DebugLevel)
	endpoints := []config.Endpoint{
		{Eth1: os.Getenv("ETH1_ENDPOINT"), Eth2: os.Getenv("ETH2_ENDPOINT")},
	}
	c, err := connection.NewConnection(endpoints, nil, nil, nil, nil)
	assert.Nil(t, err)
	config, err := c.GetEth2Config()
	assert.Nil(t, err)
	cfgBytes, err := json.MarshalIndent(config, "", "  ")
	assert.Nil(t, err)
	t.Log(string(cfgBytes))
	timestamp := utils.StartTimestampOfEpoch(config, 10383)
	t.Log(timestamp)
}

func TestBlockDetail(t *testing.T) {
	s := make([]int64, 0)
	sort.SliceStable(s, func(i, j int) bool { return s[i] < s[j] })

	logrus.SetLevel(logrus.DebugLevel)
	endpoints := []config.Endpoint{
		{Eth1: os.Getenv("ETH1_ENDPOINT"), Eth2: os.Getenv("ETH2_ENDPOINT")},
	}
	c, err := connection.NewConnection(endpoints, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = c.GetBeaconBlock(7312423)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBalance(t *testing.T) {
	cc, err := ethclient.Dial(os.Getenv("ETH1_ENDPOINT"))
	if err != nil {
		t.Fatal(err)
	}
	blockNumber, err := cc.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(blockNumber)
	tx, err := cc.TransactionReceipt(context.Background(), common.HexToHash("0x7e1bd5879335a0bc5d088f7709d76ba257de6b00473bc441c65fa9eedd552e57"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tx.Logs)
}

func TestGettingFirstNodeStakeEvent(t *testing.T) {
	endpoints := []config.Endpoint{
		{Eth1: os.Getenv("ETH1_ENDPOINT"), Eth2: os.Getenv("ETH2_ENDPOINT")},
	}
	c, err := connection.NewConnection(endpoints, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	var start = uint64(0)
	latestBlock, err := c.Eth1LatestBlock()
	if err != nil {
		t.Fatal(err)
	}
	end := latestBlock

	nodeDeposits := []string{
		"0x179386303fC2B51c306Ae9D961C73Ea9a9EA0C8d",
		"0x8A57bC7fB1237f9fBF075A261Ed28F04105Cd89d",
	}

	for _, nodeDepositAddr := range nodeDeposits {
		fmt.Println("nodeDepositAddr:", nodeDepositAddr)

		nodeDepositContract, err := node_deposit.NewNodeDeposit(common.HexToAddress(nodeDepositAddr), c.Eth1Client())
		if err != nil {
			t.Fatal(err)
		}
		iter, err := retry.DoWithData(func() (*node_deposit.NodeDepositStakedIterator, error) {
			return nodeDepositContract.FilterStaked(&bind.FilterOpts{
				Start:   start,
				End:     &end,
				Context: context.Background(),
			})
		}, retry.Delay(time.Second*2), retry.Attempts(150))
		if err != nil {
			t.Fatal(err)
		}

		// for iter.Next() {
		// 	fmt.Println("stake event at:", iter.Event.Raw.BlockNumber)
		// }

		hasEvent := iter.Next()
		iter.Close()
		if hasEvent {
			// found the first node deposit event
			fmt.Println("first stake event", iter.Event.Raw.BlockNumber)
		} else {
			fmt.Println("no node stake event")
		}
	}
	// lsdTokens: 0x37a7BF277f9b1F32296aB595600eA30c55F6eE4B
	// lsdTokens: 0xD2a1e6931e8a41043cE80C4F7EB0F7083E64Bfb8 ( created by robert)
}

func Test_FilterSetMerkleRootEvent(t *testing.T) {
	cc, err := ethclient.Dial(os.Getenv("ETH1_ENDPOINT"))
	if err != nil {
		t.Fatal(err)
	}
	start := uint64(19446894)
	end := uint64(19576494)

	networkWithdrawContract, err := network_withdraw.NewNetworkWithdraw(common.HexToAddress("0x93564AEBDd016c9cF621e7366Ae7F9bFF554C08E"), cc)
	if err != nil {
		t.Fatal(err)
	}
	iter, err := networkWithdrawContract.FilterSetMerkleRoot(&bind.FilterOpts{
		Start:   start,
		End:     &end,
		Context: context.Background(),
	}, nil)

	if err != nil {
		t.Fatal(err)
	}

	for iter.Next() {
		fmt.Println("SetMerkleRoot event at:", iter.Event.Raw.BlockNumber)
		fmt.Println("cid: ", iter.Event.NodeRewardsFileCid)
	}

	// output:
	// SetMerkleRoot event at: 19537813
	// cid:  bafybeigxluiq6pr6vox7n66j2mypgrmuz434njjdobd2y2h3bey5em3zb4
	// SetMerkleRoot event at: 19568460
	// cid:  bafybeic6xhrbflx6jjl4577iitsxs6fl7dufhmfsvcqb2nbpilo4jzlqle
	// SetMerkleRoot event at: 19574490
	// cid:  bafybeifgyaks6tjci7e3p7m5lzxvfb42jmoj2rhmaakc2vku4f7ppgrltm
}

func TestDebug_TraceBlockByNumber(t *testing.T) {
	eth1Client, err := connection.NewEth1Client([]string{os.Getenv("ETH1_ENDPOINT")})
	if err != nil {
		t.Fatal(err)
	}
	blockNumber := big.NewInt(24338297)
	trace, err := eth1Client.Debug_TraceBlockByNumber(context.Background(), blockNumber, connection.Tracer{Tracer: "callTracer"})
	if err != nil {
		t.Fatal(err)
	}

	transferFeeAddresses := []string{
		"0x54da21340773fecaf9a5bad0883a7fc594945d0a",
	}
	to := "0x5ead01d58067a68d0d700374500580ec5c961d0d"
	amount := decimal.Zero
	for _, tx := range trace {
		seekFn := func(tx *connection.TxTrace) bool {
			return utils.In(transferFeeAddresses, strings.ToLower(tx.From)) &&
				strings.EqualFold(tx.To, to)
		}
		amountInTx := service.WalkTrace(seekFn, decimal.Zero, tx.Result).DivRound(decimal.NewFromInt(1e18), 18)
		amount = amount.Add(amountInTx)
		t.Logf("%s: %s", tx.TxHash.Hex(), amountInTx.StringFixed(18))
	}
	t.Logf("total: %s", amount.StringFixed(18))
	// 1608215.139529883832080731
}
