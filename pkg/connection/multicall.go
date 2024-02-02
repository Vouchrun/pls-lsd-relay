// Copyright 2024 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only
package connection

import (
	"math/big"
	"time"

	"github.com/forta-network/go-multicall"
	"github.com/samber/lo"
	"github.com/stafiprotocol/eth-lsd-relay/pkg/gomicrobee"
)

func (c *Connection) initMulticall() error {
	multiCaller, err := multicall.New(c.Eth1Client())
	if err != nil {
		return err
	}

	c.multiCaller = multiCaller
	c.latestMultiCallMicrobeeSystem = gomicrobee.NewSystem[*multicall.Call, *MultiCall](
		c.doLatestCalls,
		20,
		5*time.Second,
	)
	return nil
}

func (c *Connection) SubmitLatestCallJob(call *multicall.Call) (gomicrobee.JobResult[*MultiCall], error) {
	return c.latestMultiCallMicrobeeSystem.Submit(call)
}

type MultiCall struct {
	*multicall.Call
	BlockNumber uint64
	Err         error
}

func (c *Connection) doLatestCalls(calls []*multicall.Call) []*MultiCall {
	eth1LatestBlock, err := c.Eth1LatestBlock()
	if err != nil {
		return lo.Map(calls, func(call *multicall.Call, index int) *MultiCall {
			call.Failed = true
			return &MultiCall{Call: call, BlockNumber: eth1LatestBlock, Err: err}
		})
	}

	opts := c.CallOpts(big.NewInt(int64(eth1LatestBlock)))
	_, err = c.multiCaller.Call(opts, calls...)
	if err != nil {
		return lo.Map(calls, func(call *multicall.Call, index int) *MultiCall {
			call.Failed = true
			return &MultiCall{Call: call, BlockNumber: eth1LatestBlock, Err: err}
		})
	}

	return lo.Map(calls, func(call *multicall.Call, index int) *MultiCall {
		return &MultiCall{Call: call, BlockNumber: eth1LatestBlock}
	})
}
