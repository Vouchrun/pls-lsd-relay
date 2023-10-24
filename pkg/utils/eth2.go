package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/connection/beacon"
)

// 1 deposited { 2 withdrawl match 3 staked 4 withdrawl unmatch } { 5 offboard 6 OffBoard can withdraw 7 OffBoard withdrawed } 8 waiting 9 active 10 exited 11 withdrawable 12 withdrawdone { 13 distributed }
// 51 active+slash 52 exit+slash 53 withdrawable+slash 54 withdrawdone+slash 55 distributed+slash
const (
	ValidatorStatusUnInitial = uint8(0)
	ValidatorStatusDeposited = uint8(1)

	// lightnode + trust node related
	ValidatorStatusWithdrawMatch   = uint8(2)
	ValidatorStatusStaked          = uint8(3)
	ValidatorStatusWithdrawUnmatch = uint8(4)

	// status on beacon chain
	ValidatorStatusWaiting      = uint8(5)
	ValidatorStatusActive       = uint8(6)
	ValidatorStatusExited       = uint8(7)
	ValidatorStatusWithdrawable = uint8(8)
	ValidatorStatusWithdrawDone = uint8(9)

	// after distribute reward
	ValidatorStatusDistributed = uint8(10) // distribute full withdrawal

	// after slash
	ValidatorStatusActiveSlash       = uint8(51)
	ValidatorStatusExitedSlash       = uint8(52)
	ValidatorStatusWithdrawableSlash = uint8(53)
	ValidatorStatusWithdrawDoneSlash = uint8(54)

	ValidatorStatusDistributedSlash = uint8(55) // distribute full withdrawal
)

// 1 Solo node 2 trust node
const (
	NodeTypeSolo  = uint8(1)
	NodeTypeTrust = uint8(2)
)

const (
	NodeClaimTypeNone         = uint8(0)
	NodeClaimTypeClaimReward  = uint8(1)
	NodeClaimTypeClaimDeposit = uint8(2)
	NodeClaimTypeClaimTotal   = uint8(3)

	DistributeTypeNone        = uint8(0)
	DistributeTypeWithdrawals = uint8(1)
	DistributeTypePriorityFee = uint8(2)
)

var (
	GweiDeci = decimal.NewFromInt(1e9)

	Percent5Deci  = decimal.NewFromFloat(0.05)
	Percent90Deci = decimal.NewFromFloat(0.9)

	StandardEffectiveBalance     = uint64(32e9)
	StandardEffectiveBalanceDeci = decimal.NewFromBigInt(big.NewInt(32), 18)

	StandardTrustNodeFakeDepositBalance = decimal.NewFromInt(1e18)

	MaxPartialWithdrawalAmount     = uint64(8e9)
	MaxPartialWithdrawalAmountDeci = decimal.NewFromInt(8e18)
)

const (
	StakerWithdrawalClaimableTimestamp = uint64(1)
	MinValidatorWithdrawabilityDelay   = uint64(256 + 5)
	MaxDistributeWaitSeconds           = uint64(8 * 60 * 60)
	MaxDistributeWaitEpoch             = uint64(75)

	EjectorUptimeInterval = uint64(10 * 60)
)

// Get an eth2 epoch number by time
func EpochAtTimestamp(config beacon.Eth2Config, time uint64) uint64 {
	return config.GenesisEpoch + (time-config.GenesisTime)/config.SecondsPerEpoch
}

func SlotAtTimestamp(config beacon.Eth2Config, time uint64) uint64 {
	return (time - config.GenesisTime) / config.SecondsPerSlot
}

func StartTimestampOfEpoch(config beacon.Eth2Config, epoch uint64) uint64 {
	return (epoch-config.GenesisEpoch)*config.SecondsPerEpoch + config.GenesisTime
}

func TimestampOfSlot(config beacon.Eth2Config, slot uint64) uint64 {
	return slot*config.SecondsPerSlot + config.GenesisTime
}

// Get an eth2 first slot number by epoch
func StartSlotOfEpoch(config beacon.Eth2Config, epoch uint64) uint64 {
	return config.SlotsPerEpoch * epoch
}
func EndSlotOfEpoch(config beacon.Eth2Config, epoch uint64) uint64 {
	return config.SlotsPerEpoch*(epoch+1) - 1
}

func GetGaspriceFromEthgasstation() (base, priority uint64, err error) {
	rsp, err := http.Get("https://api.ethgasstation.info/api/fee-estimate")
	if err != nil {
		return 0, 0, err
	}
	defer rsp.Body.Close()
	if rsp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("status err %d", rsp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return 0, 0, err
	}
	if len(bodyBytes) == 0 {
		return 0, 0, fmt.Errorf("bodyBytes zero err")
	}
	resGasPrice := ResGasPrice{}
	err = json.Unmarshal(bodyBytes, &resGasPrice)
	if err != nil {
		return 0, 0, err
	}
	return uint64(resGasPrice.BaseFee), uint64(resGasPrice.PriorityFee.Fast), nil
}

type ResGasPrice struct {
	BaseFee     int     `json:"baseFee"`
	BlockNumber int     `json:"blockNumber"`
	BlockTime   float64 `json:"blockTime"`
	GasPrice    struct {
		Fast     int `json:"fast"`
		Instant  int `json:"instant"`
		Standard int `json:"standard"`
	} `json:"gasPrice"`
	NextBaseFee int `json:"nextBaseFee"`
	PriorityFee struct {
		Fast     int `json:"fast"`
		Instant  int `json:"instant"`
		Standard int `json:"standard"`
	} `json:"priorityFee"`
}

func GetGaspriceFromBeacon() (base uint64, err error) {
	rsp, err := http.Get("https://beaconcha.in/api/v1/execution/gasnow")
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status err %d", rsp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(rsp.Body)
	if err != nil {
		return 0, err
	}
	if len(bodyBytes) == 0 {
		return 0, fmt.Errorf("bodyBytes zero err")
	}
	resGasPrice := ResGasPriceFromBeacon{}
	err = json.Unmarshal(bodyBytes, &resGasPrice)
	if err != nil {
		return 0, err
	}

	return uint64(resGasPrice.Data.Standard), nil
}

type ResGasPriceFromBeacon struct {
	Code int `json:"code"`
	Data struct {
		Rapid     int64   `json:"rapid"`
		Fast      int64   `json:"fast"`
		Standard  int64   `json:"standard"`
		Slow      int64   `json:"slow"`
		Timestamp int64   `json:"timestamp"`
		PriceUSD  float64 `json:"priceUSD"`
	} `json:"data"`
}

// bytes32 proposalId = keccak256(abi.encodePacked("execProposal", _to, _callData, _proposalFactor));
func ProposalId(to common.Address, callData []byte, proposalFactor *big.Int) [32]byte {
	return crypto.Keccak256Hash([]byte("execProposal"), to.Bytes(), callData, common.LeftPadBytes(proposalFactor.Bytes(), 32))
}

func WaitTxOkCommon(client *ethclient.Client, txHash common.Hash) (blockNumber uint64, err error) {
	defer func() {
		if err != nil {
			logrus.Errorf("find err: %s, will shutdown.", err.Error())
			ShutdownRequestChannel <- struct{}{}
		}
	}()

	retry := 0
	for {
		if retry > RetryLimit {
			return 0, fmt.Errorf("waitTx %s reach retry limit", txHash.String())
		}
		_, pending, err := client.TransactionByHash(context.Background(), txHash)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"hash": txHash.String(),
				"err":  err.Error(),
			}).Warn("TransactionByHash")

			time.Sleep(RetryInterval)
			retry++
			continue
		} else {
			if pending {
				logrus.WithFields(logrus.Fields{
					"hash":    txHash.String(),
					"pending": pending,
				}).Warn("TransactionByHash")

				time.Sleep(RetryInterval)
				retry++
				continue
			} else {
				// check status
				var receipt *types.Receipt
				subRetry := 0
				for {
					if subRetry > RetryLimit {
						return 0, fmt.Errorf("TransactionReceipt %s reach retry limit", txHash.String())
					}

					receipt, err = client.TransactionReceipt(context.Background(), txHash)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"hash": txHash.String(),
							"err":  err.Error(),
						}).Warn("tx TransactionReceipt")

						time.Sleep(RetryInterval)
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

// user = 90%*(1-nodedeposit/32)
// node = 5% + (90% * nodedeposit/32)
// platform = 5%
// nodeDepositAmount decimals 18
// rewardDeci decimals 18
// return (user reward, node reward, paltform fee) decimals 18
func GetUserNodePlatformReward(nodeCommissionRate, platformCommissionRate, nodeDepositAmountDeci, rewardDeci decimal.Decimal) (decimal.Decimal, decimal.Decimal, decimal.Decimal) {
	if !rewardDeci.IsPositive() || nodeDepositAmountDeci.GreaterThan(StandardEffectiveBalanceDeci) {
		return decimal.Zero, decimal.Zero, decimal.Zero
	}

	// platform Fee
	platformFee := rewardDeci.Mul(platformCommissionRate)

	// node fee
	leftRate := decimal.NewFromInt(1).Sub(nodeCommissionRate.Add(platformCommissionRate))
	nodeTotalRate := nodeCommissionRate.Add(leftRate.Mul(nodeDepositAmountDeci.Div(StandardEffectiveBalanceDeci)))
	nodeFee := rewardDeci.Mul(nodeTotalRate)

	// user fee

	userFee := rewardDeci.Sub(platformFee.Add(nodeFee))

	return userFee, nodeFee, platformFee
}
