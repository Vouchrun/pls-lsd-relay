package service

import (
	"encoding/json"
	"fmt"
	"math/big"
	"path"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/utils"
)

type NodeRewardsList []*NodeReward
type NodeRewardsMap map[string]*NodeReward

type NodeReward struct {
	Address                string          `json:"address"` // hex with 0x
	Index                  uint32          `json:"index"`
	TotalRewardAmount      decimal.Decimal `json:"totalRewardAmount"`
	TotalExitDepositAmount decimal.Decimal `json:"totalExitDepositAmount"`
	Proof                  string          `json:"proof"`
}

type NodeNewReward struct {
	Address                string          `json:"address"` // hex with 0x
	TotalRewardAmount      decimal.Decimal `json:"totalRewardAmount"`
	TotalExitDepositAmount decimal.Decimal `json:"totalExitDepositAmount"`
}

func (s *Service) voteMerkleRoot() error {
	beaconHead, err := s.connection.Eth2BeaconHead()
	if err != nil {
		return err
	}

	targetEpoch := (beaconHead.FinalizedEpoch / s.merkleRootDuEpochs) * s.merkleRootDuEpochs

	dealedEpochOnchain, err := s.networkWithdrawContract.LatestMerkleRootEpoch(nil)
	if err != nil {
		return err
	}
	if targetEpoch <= dealedEpochOnchain.Uint64() {
		return nil
	}

	targetEth1BlockHeight, err := s.getEpochStartBlocknumber(targetEpoch)
	if err != nil {
		return err
	}

	var dealedEth1BlockHeight uint64
	preNodeRewardList := make(NodeRewardsList, 0)
	preNodeRewardMap := make(NodeRewardsMap)
	if dealedEpochOnchain.Uint64() > 0 {
		merkleRootIter, err := s.networkWithdrawContract.FilterSetMerkleRoot(nil, []*big.Int{dealedEpochOnchain})
		if err != nil {
			return err
		}
		if !merkleRootIter.Next() {
			return fmt.Errorf("SetMerkleRoot event not exit on target epoch %d", dealedEpochOnchain.Uint64())
		}

		preCid := merkleRootIter.Event.NodeRewardsFileCid

		fileBytes, err := utils.DownloadWeb3File(preCid, utils.NodeRewardsFileNameAtEpoch(dealedEpochOnchain.Uint64()))
		if err != nil {
			return err
		}

		err = json.Unmarshal(fileBytes, &preNodeRewardList)
		if err != nil {
			return err
		}

		dealedEth1BlockHeight, err = s.getEpochStartBlocknumber(dealedEpochOnchain.Uint64())
		if err != nil {
			return err
		}
	} else {
		dealedEth1BlockHeight = s.networkCreateBlock
	}

	for _, nodeReward := range preNodeRewardList {
		_, exist := preNodeRewardMap[nodeReward.Address]
		if exist {
			return fmt.Errorf("duplicate node address: %s", nodeReward.Address)
		}
		preNodeRewardMap[nodeReward.Address] = nodeReward
	}

	newNodeRewardsMap, err := s.getNodeNewRewardsBetween(dealedEth1BlockHeight, targetEth1BlockHeight)
	if err != nil {
		return err
	}

	finalNodeRewardsMap := make(NodeRewardsMap, 0)
	for _, node := range preNodeRewardMap {
		f, exist := finalNodeRewardsMap[node.Address]
		if exist {
			f.TotalRewardAmount = f.TotalRewardAmount.Add(node.TotalRewardAmount)
			f.TotalExitDepositAmount = f.TotalExitDepositAmount.Add(node.TotalExitDepositAmount)
		} else {
			finalNodeRewardsMap[node.Address] = &NodeReward{
				Address:                node.Address,
				TotalRewardAmount:      node.TotalRewardAmount,
				TotalExitDepositAmount: node.TotalExitDepositAmount,
			}
		}
	}

	for _, node := range newNodeRewardsMap {
		f, exist := finalNodeRewardsMap[node.Address]
		if exist {
			f.TotalRewardAmount = f.TotalRewardAmount.Add(node.TotalRewardAmount)
			f.TotalExitDepositAmount = f.TotalExitDepositAmount.Add(node.TotalExitDepositAmount)
		} else {
			finalNodeRewardsMap[node.Address] = &NodeReward{
				Address:                node.Address,
				TotalRewardAmount:      node.TotalRewardAmount,
				TotalExitDepositAmount: node.TotalExitDepositAmount,
			}
		}
	}

	// got final list
	finalNodeRewardsList := make(NodeRewardsList, 0)
	for _, node := range finalNodeRewardsMap {
		finalNodeRewardsList = append(finalNodeRewardsList, node)
	}
	sort.SliceStable(finalNodeRewardsList, func(i, j int) bool {
		return finalNodeRewardsList[i].Address < finalNodeRewardsList[j].Address
	})
	for i, node := range finalNodeRewardsList {
		node.Index = uint32(i)
	}

	rootHash := utils.NodeHash{}
	if len(finalNodeRewardsList) > 0 {
		// build merkle tree
		tree, err := buildMerkleTree(finalNodeRewardsList)
		if err != nil {
			return err
		}
		rootHash, err = tree.GetRootHash()
		if err != nil {
			return err
		}

		// calc proof
		for _, nodeReward := range finalNodeRewardsList {
			nodeHash := utils.GetNodeHash(big.NewInt(int64(nodeReward.Index)), common.HexToAddress(nodeReward.Address),
				nodeReward.TotalRewardAmount.BigInt(), nodeReward.TotalExitDepositAmount.BigInt())
			proofList, err := tree.GetProof(nodeHash)
			if err != nil {
				return errors.Wrap(err, "tree.GetProof failed")
			}
			if len(proofList) == 0 {
				return errors.Wrap(err, "tree.GetProof result empty")
			}

			proofStrList := make([]string, len(proofList))
			for i, p := range proofList {
				proofStrList[i] = p.String()
			}
			// set proof
			nodeReward.Proof = strings.Join(proofStrList, ":")
		}
	}

	// upload file
	fileBts, err := json.Marshal(finalNodeRewardsList)
	if err != nil {
		return err
	}
	filePath := path.Join(s.nodeRewardsFilePath, utils.NodeRewardsFileNameAtEpoch(targetEpoch))
	cid, err := utils.UploadFileToWeb3Storage(s.web3Client, fileBts, filePath)
	if err != nil {
		return err
	}

	// check voted
	hasVoted, err := s.networkProposalContract.HasVoted(nil, utils.VoteMerkleRootProposalId(big.NewInt(int64(targetEpoch)),
		rootHash, cid), s.keyPair.CommonAddress())
	if err != nil {
		return fmt.Errorf("networkProposalContract.HasVoted err: %s", err)
	}
	if hasVoted {
		logrus.Debug("networkProposalContract voted")
		return nil
	}

	return s.sendSetMerkleRootTx(rootHash, int64(targetEpoch), cid)
}

func buildMerkleTree(nodelist NodeRewardsList) (*utils.MerkleTree, error) {
	if len(nodelist) == 0 {
		return nil, fmt.Errorf("proof list empty")
	}
	list := make(utils.NodeHashList, len(nodelist))
	for i, data := range nodelist {

		list[i] = utils.GetNodeHash(big.NewInt(int64(data.Index)), common.HexToAddress(data.Address),
			data.TotalRewardAmount.BigInt(), data.TotalExitDepositAmount.BigInt())
	}
	mt := utils.NewMerkleTree(list)
	return mt, nil
}

func (s *Service) sendSetMerkleRootTx(rootHash utils.NodeHash, targetEpoch int64, cid string) error {

	var merkleTreeRootHash [32]byte
	copy(merkleTreeRootHash[:], rootHash)

	err := s.connection.LockAndUpdateTxOpts()
	if err != nil {
		return fmt.Errorf("LockAndUpdateTxOpts err: %s", err)
	}
	defer s.connection.UnlockTxOpts()

	tx, err := s.networkWithdrawContract.SetMerkleRoot(s.connection.TxOpts(), big.NewInt(targetEpoch), merkleTreeRootHash, cid)
	if err != nil {
		return err
	}

	logrus.Info("send SetMerkleRoot tx hash: ", tx.Hash().String())

	return s.waitTxOk(tx.Hash())
}
