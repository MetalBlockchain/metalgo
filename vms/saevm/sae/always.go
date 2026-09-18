// Copyright (C) 2019, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package sae

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MetalBlockchain/libevm/core"
	"github.com/MetalBlockchain/libevm/core/rawdb"
	"github.com/MetalBlockchain/libevm/ethdb"
	"github.com/MetalBlockchain/libevm/params"
	"github.com/MetalBlockchain/libevm/triedb"

	"github.com/MetalBlockchain/metalgo/database"
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/vms/saevm/adaptor"
	"github.com/MetalBlockchain/metalgo/vms/saevm/blocks"
	"github.com/MetalBlockchain/metalgo/vms/saevm/hook"
	"github.com/MetalBlockchain/metalgo/vms/saevm/network"
	"github.com/MetalBlockchain/metalgo/vms/saevm/types"

	snowcommon "github.com/MetalBlockchain/metalgo/snow/engine/common"
	ethcommon "github.com/MetalBlockchain/libevm/common"
)

var _ adaptor.ChainVM[*blocks.Block] = (*SinceGenesis[hook.Transaction])(nil)

// SinceGenesis is a harness around a [VM], providing an `Initialize` method
// that treats the chain as being asynchronous since genesis.
type SinceGenesis[T hook.Transaction] struct {
	*VM // created by [SinceGenesis.Initialize]
	*network.Network

	hooks  hook.PointsG[T]
	config Config
}

// NewSinceGenesis constructs a new [SinceGenesis].
func NewSinceGenesis[T hook.Transaction](hooks hook.PointsG[T], c Config) *SinceGenesis[T] {
	return &SinceGenesis[T]{
		hooks:  hooks,
		config: c,
	}
}

// Initialize initializes the VM.
//
//nolint:revive // General-purpose types lose the meaning of args if unused ones are removed
func (vm *SinceGenesis[_]) Initialize(
	ctx context.Context,
	snowCtx *snow.Context,
	avaDB database.Database,
	genesisBytes []byte,
	upgradeBytes []byte,
	configBytes []byte,
	fxs []*snowcommon.Fx,
	appSender snowcommon.AppSender,
) error {
	db := types.NewEthDB(avaDB)
	tdbCfg := vm.config.DBConfig.TrieDBConfig(snowCtx.ChainDataDir, snowCtx.Log)
	config, err := setupGenesis(db, tdbCfg, genesisBytes)
	if err != nil {
		return err
	}

	vm.Network, err = network.New(snowCtx, appSender)
	if err != nil {
		return fmt.Errorf("network.New(...): %v", err)
	}
	vm.VM, err = NewVM(ctx, vm.hooks, vm.config, snowCtx, config, db, vm.Network)
	return err
}

func setupGenesis(db ethdb.Database, tdbConfig *triedb.Config, genesisBytes []byte) (_ *params.ChainConfig, retErr error) {
	tdb := triedb.NewDatabase(db, tdbConfig)
	defer func() {
		retErr = errors.Join(retErr, tdb.Close())
	}()

	genesis := new(core.Genesis)
	if err := json.Unmarshal(genesisBytes, genesis); err != nil {
		return nil, fmt.Errorf("json.Unmarshal(%T): %v", genesis, err)
	}
	config, hash, err := core.SetupGenesisBlock(db, tdb, genesis)
	if err != nil {
		return nil, fmt.Errorf("core.SetupGenesisBlock(...): %v", err)
	}

	// [NewVM] assumes that the genesis block is "finalized", which does not
	// happen in [core.SetupGenesisBlock]. This MUST only happen once.
	if rawdb.ReadFinalizedBlockHash(db) == (ethcommon.Hash{}) {
		rawdb.WriteFinalizedBlockHash(db, hash)
	}

	return config, nil
}

// Shutdown gracefully closes the VM.
func (vm *SinceGenesis[_]) Shutdown(ctx context.Context) error {
	if vm.VM == nil {
		return nil
	}
	return vm.VM.Shutdown(ctx)
}
