// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package xsvm

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/rpc/v2"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/MetalBlockchain/metalgo/database"
	"github.com/MetalBlockchain/metalgo/database/versiondb"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/network/p2p"
	"github.com/MetalBlockchain/metalgo/network/p2p/acp118"
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/snow/consensus/snowman"
	"github.com/MetalBlockchain/metalgo/snow/engine/common"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/json"
	"github.com/MetalBlockchain/metalgo/vms/example/xsvm/api"
	"github.com/MetalBlockchain/metalgo/vms/example/xsvm/builder"
	"github.com/MetalBlockchain/metalgo/vms/example/xsvm/chain"
	"github.com/MetalBlockchain/metalgo/vms/example/xsvm/execute"
	"github.com/MetalBlockchain/metalgo/vms/example/xsvm/genesis"
	"github.com/MetalBlockchain/metalgo/vms/example/xsvm/state"

	smblock "github.com/MetalBlockchain/metalgo/snow/engine/snowman/block"
	xsblock "github.com/MetalBlockchain/metalgo/vms/example/xsvm/block"
)

var (
	_ smblock.ChainVM                      = (*VM)(nil)
	_ smblock.BuildBlockWithContextChainVM = (*VM)(nil)
)

type VM struct {
	*p2p.Network

	chainContext *snow.Context
	db           database.Database
	genesis      *genesis.Genesis
	engineChan   chan<- common.Message

	chain   chain.Chain
	builder builder.Builder
}

func (vm *VM) Initialize(
	_ context.Context,
	chainContext *snow.Context,
	db database.Database,
	genesisBytes []byte,
	_ []byte,
	_ []byte,
	engineChan chan<- common.Message,
	_ []*common.Fx,
	appSender common.AppSender,
) error {
	chainContext.Log.Info("initializing xsvm",
		zap.Stringer("version", Version),
	)

	metrics := prometheus.NewRegistry()
	err := chainContext.Metrics.Register("p2p", metrics)
	if err != nil {
		return err
	}

	vm.Network, err = p2p.NewNetwork(
		chainContext.Log,
		appSender,
		metrics,
		"",
	)
	if err != nil {
		return err
	}

	// Allow signing of all warp messages. This is not typically safe, but is
	// allowed for this example.
	acp118Handler := acp118.NewHandler(
		acp118Verifier{},
		chainContext.WarpSigner,
	)
	if err := vm.Network.AddHandler(p2p.SignatureRequestHandlerID, acp118Handler); err != nil {
		return err
	}

	vm.chainContext = chainContext
	vm.db = db
	g, err := genesis.Parse(genesisBytes)
	if err != nil {
		return fmt.Errorf("failed to parse genesis bytes: %w", err)
	}

	vdb := versiondb.New(vm.db)
	if err := execute.Genesis(vdb, chainContext.ChainID, g); err != nil {
		return fmt.Errorf("failed to initialize genesis state: %w", err)
	}
	if err := vdb.Commit(); err != nil {
		return err
	}

	vm.genesis = g
	vm.engineChan = engineChan

	vm.chain, err = chain.New(chainContext, vm.db)
	if err != nil {
		return fmt.Errorf("failed to initialize chain manager: %w", err)
	}

	vm.builder = builder.New(chainContext, engineChan, vm.chain)

	chainContext.Log.Info("initialized xsvm",
		zap.Stringer("lastAcceptedID", vm.chain.LastAccepted()),
	)
	return nil
}

func (vm *VM) SetState(_ context.Context, state snow.State) error {
	vm.chain.SetChainState(state)
	return nil
}

func (vm *VM) Shutdown(context.Context) error {
	if vm.chainContext == nil {
		return nil
	}
	return vm.db.Close()
}

func (*VM) Version(context.Context) (string, error) {
	return Version.String(), nil
}

func (vm *VM) CreateHandlers(context.Context) (map[string]http.Handler, error) {
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	api := api.NewServer(
		vm.chainContext,
		vm.genesis,
		vm.db,
		vm.chain,
		vm.builder,
	)
	return map[string]http.Handler{
		"": server,
	}, server.RegisterService(api, constants.XSVMName)
}

func (*VM) HealthCheck(context.Context) (interface{}, error) {
	return http.StatusOK, nil
}

func (vm *VM) GetBlock(_ context.Context, blkID ids.ID) (snowman.Block, error) {
	return vm.chain.GetBlock(blkID)
}

func (vm *VM) ParseBlock(_ context.Context, blkBytes []byte) (snowman.Block, error) {
	blk, err := xsblock.Parse(blkBytes)
	if err != nil {
		return nil, err
	}
	return vm.chain.NewBlock(blk)
}

func (vm *VM) BuildBlock(ctx context.Context) (snowman.Block, error) {
	return vm.builder.BuildBlock(ctx, nil)
}

func (vm *VM) SetPreference(_ context.Context, preferred ids.ID) error {
	vm.builder.SetPreference(preferred)
	return nil
}

func (vm *VM) LastAccepted(context.Context) (ids.ID, error) {
	return vm.chain.LastAccepted(), nil
}

func (vm *VM) BuildBlockWithContext(ctx context.Context, blockContext *smblock.Context) (snowman.Block, error) {
	return vm.builder.BuildBlock(ctx, blockContext)
}

func (vm *VM) GetBlockIDAtHeight(_ context.Context, height uint64) (ids.ID, error) {
	return state.GetBlockIDByHeight(vm.db, height)
}
