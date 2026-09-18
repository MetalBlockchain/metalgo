// Copyright (C) 2019, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rpc

import (
	"github.com/MetalBlockchain/libevm/core/rawdb"
	"github.com/MetalBlockchain/libevm/core/types"
	"github.com/MetalBlockchain/libevm/event"
	"github.com/MetalBlockchain/libevm/rpc"

	"github.com/MetalBlockchain/metalgo/vms/saevm/blocks"
	"github.com/MetalBlockchain/metalgo/vms/saevm/gasprice"
)

type estimatorBackend struct {
	chain Chain
}

var _ gasprice.Backend = (*estimatorBackend)(nil)

func (e *estimatorBackend) BlockByNumber(n rpc.BlockNumber) (*types.Block, error) {
	return readByNumber(e.chain, n, rawdb.ReadBlock)
}

func (e *estimatorBackend) LastAcceptedBlock() *blocks.Block {
	return e.chain.LastAccepted()
}

func (e *estimatorBackend) ResolveBlockNumber(bn rpc.BlockNumber) (uint64, error) {
	return blocks.ResolveRPCNumber(e.chain, bn)
}

func (e *estimatorBackend) SubscribeAcceptedBlocks(ch chan<- *blocks.Block) event.Subscription {
	return e.chain.SubscribeAcceptedBlocks(ch)
}
