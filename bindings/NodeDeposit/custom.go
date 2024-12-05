// Copyright 2024 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only
package node_deposit

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/forta-network/go-multicall"
	"github.com/samber/lo"
)

type CustomNodeDeposit struct {
	*NodeDeposit
	contractAddress common.Address
	multiCaller     *multicall.Caller
	multiContract   *multicall.Contract
}

// NewNodeDeposit creates a new instance of NodeDeposit, bound to a specific deployed contract.
func NewCustomNodeDeposit(address common.Address, backend bind.ContractBackend, multiCaller *multicall.Caller) (*CustomNodeDeposit, error) {
	ins, err := NewNodeDeposit(address, backend)
	if err != nil {
		return nil, err
	}

	multiContract, err := multicall.NewContract(NodeDepositABI, address.String())
	if err != nil {
		return nil, err
	}

	return &CustomNodeDeposit{
		NodeDeposit:     ins,
		contractAddress: address,
		multiCaller:     multiCaller,
		multiContract:   multiContract,
	}, nil
}

type pubkeyList [][]byte

type GetPubkeysOfNodeOutput struct {
	List pubkeyList
}

func (nodeDeposit *CustomNodeDeposit) getPubkeysOfNodes(opts *bind.CallOpts, _nodes []common.Address) (map[common.Address]pubkeyList, error) {
	if len(_nodes) == 0 {
		return nil, nil
	}

	calls := make([]*multicall.Call, len(_nodes))
	for i, node := range _nodes {
		calls[i] = nodeDeposit.multiContract.NewCall(
			new(GetPubkeysOfNodeOutput),
			"getPubkeysOfNode",
			node,
		)
	}

	_, err := nodeDeposit.multiCaller.Call(opts, calls...)
	if err != nil {
		return nil, err
	}

	return lo.Associate(calls, func(c *multicall.Call) (common.Address, pubkeyList) {
		return c.Inputs[0].(common.Address), c.Outputs.(*GetPubkeysOfNodeOutput).List
	}), nil
}

func (nodeDeposit *CustomNodeDeposit) GetPubkeysOfNodes(opts *bind.CallOpts, _nodes []common.Address) (map[common.Address]pubkeyList, error) {
	if len(_nodes) == 0 {
		return nil, nil
	}
	chunks := lo.Chunk(_nodes, 100)
	result := make(map[common.Address]pubkeyList, len(_nodes))
	for _, chunk := range chunks {
		out, err := nodeDeposit.getPubkeysOfNodes(opts, chunk)
		if err != nil {
			return nil, err
		}
		for k, v := range out {
			result[k] = v
		}
	}
	return result, nil
}

type GetPubkeyInfoListOutput struct {
	Status            uint8
	Owner             common.Address
	NodeDepositAmount *big.Int
	DepositBlock      *big.Int
}

func (nodeDeposit *CustomNodeDeposit) getPubkeyInfoList(opts *bind.CallOpts, pubkeys [][]byte) (map[string]*GetPubkeyInfoListOutput, error) {
	if len(pubkeys) == 0 {
		return nil, nil
	}

	calls := lo.Map(pubkeys, func(pubkey []byte, index int) *multicall.Call {
		return nodeDeposit.multiContract.NewCall(
			new(GetPubkeyInfoListOutput),
			"pubkeyInfoOf",
			pubkey,
		)
	})

	_, err := nodeDeposit.multiCaller.Call(opts, calls...)
	if err != nil {
		return nil, err
	}

	return lo.Associate(calls, func(c *multicall.Call) (string, *GetPubkeyInfoListOutput) {
		return hex.EncodeToString(c.Inputs[0].([]byte)), c.Outputs.(*GetPubkeyInfoListOutput)
	}), nil
}

func (nodeDeposit *CustomNodeDeposit) GetPubkeyInfoList(opts *bind.CallOpts, pubkeys [][]byte) (map[string]*GetPubkeyInfoListOutput, error) {
	if len(pubkeys) == 0 {
		return nil, nil
	}
	chunks := lo.Chunk(pubkeys, 100)
	result := make(map[string]*GetPubkeyInfoListOutput, len(pubkeys))
	for _, chunk := range chunks {
		out, err := nodeDeposit.getPubkeyInfoList(opts, chunk)
		if err != nil {
			return nil, err
		}
		for k, v := range out {
			result[k] = v
		}
	}
	return result, nil
}

type GetNodesLengthMultiCallOutput struct {
	Length *big.Int
}

func (nodeDeposit *CustomNodeDeposit) NewGetNodesLengthMultiCall() *multicall.Call {
	return nodeDeposit.multiContract.NewCall(new(GetNodesLengthMultiCallOutput), "getNodesLength")
}
