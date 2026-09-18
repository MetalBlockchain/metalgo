// Copyright (C) 2019, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package sae

import (
	"github.com/MetalBlockchain/libevm/common"
	"github.com/MetalBlockchain/libevm/core"
	"github.com/MetalBlockchain/libevm/core/types"
	"github.com/MetalBlockchain/libevm/ethdb"
	"github.com/MetalBlockchain/libevm/event"

	"github.com/MetalBlockchain/metalgo/network/p2p"
	"github.com/MetalBlockchain/metalgo/utils/logging"
	"github.com/MetalBlockchain/metalgo/vms/saevm/blocks"
	"github.com/MetalBlockchain/metalgo/vms/saevm/hook"
	"github.com/MetalBlockchain/metalgo/vms/saevm/saexec"
	"github.com/MetalBlockchain/metalgo/vms/saevm/txgossip"

	saerpc "github.com/MetalBlockchain/metalgo/vms/saevm/sae/rpc"
	saetypes "github.com/MetalBlockchain/metalgo/vms/saevm/types"
)

// GethRPCBackends returns the backing infrastructure for geth's implementations
// of the JSON-RPC namespaces supported by the VM.
func (vm *VM) GethRPCBackends() saerpc.GethBackends {
	return vm.rpcProvider.GethBackends()
}

func (vm *VM) chain() saerpc.Chain {
	return chain{vm, vm.exec}
}

type chain struct {
	*VM
	*saexec.Executor
}

func (c chain) Logger() logging.Logger         { return c.snowCtx.Log }
func (c chain) Hooks() hook.Points             { return c.hooks }
func (c chain) DB() ethdb.Database             { return c.db }
func (c chain) XDB() saetypes.ExecutionResults { return c.xdb }
func (c chain) Mempool() *txgossip.Set         { return c.mempool }
func (c chain) Peers() *p2p.Peers              { return c.network.Peers }
func (c chain) LastAccepted() *blocks.Block    { return c.last.accepted.Load() }
func (c chain) LastSettled() *blocks.Block     { return c.last.settled.Load() }

func (c chain) ConsensusCriticalBlock(h common.Hash) (*blocks.Block, bool) {
	return c.consensusCritical.Load(h)
}

func (c chain) ResolvePendingToLastExecuted() bool {
	return c.VM.config.RPCConfig.ResolvePendingToLastExecuted
}

func (c chain) NewBlock(eth *types.Block, parent, lastSettled *blocks.Block) (*blocks.Block, error) {
	return c.blockBuilder.new(eth, parent, lastSettled)
}

func (c chain) SubscribeAcceptedBlocks(ch chan<- *blocks.Block) event.Subscription {
	return c.acceptedBlocks.Subscribe(ch)
}

func (c chain) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return c.Executor.SubscribeChainHeadEvent(ch)
}
