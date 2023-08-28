// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package platformvm

import (
	"context"
	"errors"
	"fmt"

	"github.com/gorilla/rpc/v2"

	"github.com/prometheus/client_golang/prometheus"

	"go.uber.org/zap"

	"github.com/MetalBlockchain/metalgo/cache"
	"github.com/MetalBlockchain/metalgo/codec"
	"github.com/MetalBlockchain/metalgo/codec/linearcodec"
	"github.com/MetalBlockchain/metalgo/database/manager"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/snow/consensus/snowman"
	"github.com/MetalBlockchain/metalgo/snow/engine/common"
	"github.com/MetalBlockchain/metalgo/snow/engine/snowman/block"
	"github.com/MetalBlockchain/metalgo/snow/uptime"
	"github.com/MetalBlockchain/metalgo/snow/validators"
	"github.com/MetalBlockchain/metalgo/utils"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/json"
	"github.com/MetalBlockchain/metalgo/utils/logging"
	"github.com/MetalBlockchain/metalgo/utils/timer/mockable"
	"github.com/MetalBlockchain/metalgo/utils/wrappers"
	"github.com/MetalBlockchain/metalgo/version"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/api"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/blocks"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/config"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/fx"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/metrics"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/reward"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/state"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs/mempool"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/utxo"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"

	blockbuilder "github.com/MetalBlockchain/metalgo/vms/platformvm/blocks/builder"
	blockexecutor "github.com/MetalBlockchain/metalgo/vms/platformvm/blocks/executor"
	txbuilder "github.com/MetalBlockchain/metalgo/vms/platformvm/txs/builder"
	txexecutor "github.com/MetalBlockchain/metalgo/vms/platformvm/txs/executor"
	pvalidators "github.com/MetalBlockchain/metalgo/vms/platformvm/validators"
)

var (
	_ block.ChainVM              = (*VM)(nil)
	_ block.HeightIndexedChainVM = (*VM)(nil)
	_ secp256k1fx.VM             = (*VM)(nil)
	_ validators.State           = (*VM)(nil)
	_ validators.SubnetConnector = (*VM)(nil)

	errMissingValidatorSet = errors.New("missing validator set")
)

type VM struct {
	config.Config
	blockbuilder.Builder
	validators.State

	metrics            metrics.Metrics
	atomicUtxosManager avax.AtomicUTXOManager

	// Used to get time. Useful for faking time during tests.
	clock mockable.Clock

	uptimeManager uptime.Manager

	// The context of this vm
	ctx       *snow.Context
	dbManager manager.Manager

	state state.State

	fx            fx.Fx
	codecRegistry codec.Registry

	// Bootstrapped remembers if this chain has finished bootstrapping or not
	bootstrapped utils.Atomic[bool]

	txBuilder txbuilder.Builder
	manager   blockexecutor.Manager

	// TODO: Remove after v1.11.x is activated
	pruned utils.Atomic[bool]
}

// Initialize this blockchain.
// [vm.ChainManager] and [vm.vdrMgr] must be set before this function is called.
func (vm *VM) Initialize(
	ctx context.Context,
	chainCtx *snow.Context,
	dbManager manager.Manager,
	genesisBytes []byte,
	_ []byte,
	configBytes []byte,
	toEngine chan<- common.Message,
	_ []*common.Fx,
	appSender common.AppSender,
) error {
	chainCtx.Log.Verbo("initializing platform chain")

	execConfig, err := config.GetExecutionConfig(configBytes)
	if err != nil {
		return err
	}
	chainCtx.Log.Info("using VM execution config", zap.Reflect("config", execConfig))

	registerer := prometheus.NewRegistry()
	if err := chainCtx.Metrics.Register(registerer); err != nil {
		return err
	}

	// Initialize metrics as soon as possible
	vm.metrics, err = metrics.New("", registerer)
	if err != nil {
		return fmt.Errorf("failed to initialize metrics: %w", err)
	}

	vm.ctx = chainCtx
	vm.dbManager = dbManager

	vm.codecRegistry = linearcodec.NewDefault()
	vm.fx = &secp256k1fx.Fx{}
	if err := vm.fx.Initialize(vm); err != nil {
		return err
	}

	rewards := reward.NewCalculator(vm.RewardConfig)

	vm.state, err = state.New(
		vm.dbManager.Current().Database,
		genesisBytes,
		registerer,
		&vm.Config,
		execConfig,
		vm.ctx,
		vm.metrics,
		rewards,
		&vm.bootstrapped,
	)
	if err != nil {
		return err
	}

	validatorManager := pvalidators.NewManager(chainCtx.Log, vm.Config, vm.state, vm.metrics, &vm.clock)
	vm.State = validatorManager
	vm.atomicUtxosManager = avax.NewAtomicUTXOManager(chainCtx.SharedMemory, txs.Codec)
	utxoHandler := utxo.NewHandler(vm.ctx, &vm.clock, vm.fx)
	vm.uptimeManager = uptime.NewManager(vm.state)
	vm.UptimeLockedCalculator.SetCalculator(&vm.bootstrapped, &chainCtx.Lock, vm.uptimeManager)

	vm.txBuilder = txbuilder.New(
		vm.ctx,
		&vm.Config,
		&vm.clock,
		vm.fx,
		vm.state,
		vm.atomicUtxosManager,
		utxoHandler,
	)

	txExecutorBackend := &txexecutor.Backend{
		Config:       &vm.Config,
		Ctx:          vm.ctx,
		Clk:          &vm.clock,
		Fx:           vm.fx,
		FlowChecker:  utxoHandler,
		Uptimes:      vm.uptimeManager,
		Rewards:      rewards,
		Bootstrapped: &vm.bootstrapped,
	}

	// Note: There is a circular dependency between the mempool and block
	//       builder which is broken by passing in the vm.
	mempool, err := mempool.NewMempool("mempool", registerer, vm)
	if err != nil {
		return fmt.Errorf("failed to create mempool: %w", err)
	}

	vm.manager = blockexecutor.NewManager(
		mempool,
		vm.metrics,
		vm.state,
		txExecutorBackend,
		validatorManager,
	)
	vm.Builder = blockbuilder.New(
		mempool,
		vm.txBuilder,
		txExecutorBackend,
		vm.manager,
		toEngine,
		appSender,
	)

	// Create all of the chains that the database says exist
	if err := vm.initBlockchains(); err != nil {
		return fmt.Errorf(
			"failed to initialize blockchains: %w",
			err,
		)
	}

	lastAcceptedID := vm.state.GetLastAccepted()
	chainCtx.Log.Info("initializing last accepted",
		zap.Stringer("blkID", lastAcceptedID),
	)
	if err := vm.SetPreference(ctx, lastAcceptedID); err != nil {
		return err
	}

	shouldPrune, err := vm.state.ShouldPrune()
	if err != nil {
		return fmt.Errorf(
			"failed to check if the database should be pruned: %w",
			err,
		)
	}
	if !shouldPrune {
		chainCtx.Log.Info("state already pruned and indexed")
		vm.pruned.Set(true)
		return nil
	}

	go func() {
		err := vm.state.PruneAndIndex(&vm.ctx.Lock, vm.ctx.Log)
		if err != nil {
			vm.ctx.Log.Error("state pruning and height indexing failed",
				zap.Error(err),
			)
		}

		vm.pruned.Set(true)
	}()

	return nil
}

// Create all chains that exist that this node validates.
func (vm *VM) initBlockchains() error {
	if vm.Config.PartialSyncPrimaryNetwork {
		vm.ctx.Log.Info("skipping primary network chain creation")
	} else if err := vm.createSubnet(constants.PrimaryNetworkID); err != nil {
		return err
	}

	if vm.SybilProtectionEnabled {
		for subnetID := range vm.TrackedSubnets {
			if err := vm.createSubnet(subnetID); err != nil {
				return err
			}
		}
	} else {
		subnets, err := vm.state.GetSubnets()
		if err != nil {
			return err
		}
		for _, subnet := range subnets {
			if err := vm.createSubnet(subnet.ID()); err != nil {
				return err
			}
		}
	}
	return nil
}

// Create the subnet with ID [subnetID]
func (vm *VM) createSubnet(subnetID ids.ID) error {
	chains, err := vm.state.GetChains(subnetID)
	if err != nil {
		return err
	}
	for _, chain := range chains {
		tx, ok := chain.Unsigned.(*txs.CreateChainTx)
		if !ok {
			return fmt.Errorf("expected tx type *txs.CreateChainTx but got %T", chain.Unsigned)
		}
		vm.Config.CreateChain(chain.ID(), tx)
	}
	return nil
}

// onBootstrapStarted marks this VM as bootstrapping
func (vm *VM) onBootstrapStarted() error {
	vm.bootstrapped.Set(false)
	return vm.fx.Bootstrapping()
}

// onNormalOperationsStarted marks this VM as bootstrapped
func (vm *VM) onNormalOperationsStarted() error {
	if vm.bootstrapped.Get() {
		return nil
	}
	vm.bootstrapped.Set(true)

	if err := vm.fx.Bootstrapped(); err != nil {
		return err
	}

	primaryVdrIDs, err := validators.NodeIDs(vm.Validators, constants.PrimaryNetworkID)
	if err != nil {
		return err
	}
	if err := vm.uptimeManager.StartTracking(primaryVdrIDs, constants.PrimaryNetworkID); err != nil {
		return err
	}

	for subnetID := range vm.TrackedSubnets {
		vdrIDs, err := validators.NodeIDs(vm.Validators, subnetID)
		if err != nil {
			return err
		}
		if err := vm.uptimeManager.StartTracking(vdrIDs, subnetID); err != nil {
			return err
		}
	}

	if err := vm.state.Commit(); err != nil {
		return err
	}

	// Start the block builder
	vm.Builder.ResetBlockTimer()
	return nil
}

func (vm *VM) SetState(_ context.Context, state snow.State) error {
	switch state {
	case snow.Bootstrapping:
		return vm.onBootstrapStarted()
	case snow.NormalOp:
		return vm.onNormalOperationsStarted()
	default:
		return snow.ErrUnknownState
	}
}

// Shutdown this blockchain
func (vm *VM) Shutdown(context.Context) error {
	if vm.dbManager == nil {
		return nil
	}

	vm.Builder.Shutdown()

	if vm.bootstrapped.Get() {
		primaryVdrIDs, err := validators.NodeIDs(vm.Validators, constants.PrimaryNetworkID)
		if err != nil {
			return err
		}
		if err := vm.uptimeManager.StopTracking(primaryVdrIDs, constants.PrimaryNetworkID); err != nil {
			return err
		}

		for subnetID := range vm.TrackedSubnets {
			vdrIDs, err := validators.NodeIDs(vm.Validators, subnetID)
			if err != nil {
				return err
			}
			if err := vm.uptimeManager.StopTracking(vdrIDs, subnetID); err != nil {
				return err
			}
		}

		if err := vm.state.Commit(); err != nil {
			return err
		}
	}

	errs := wrappers.Errs{}
	errs.Add(
		vm.state.Close(),
		vm.dbManager.Close(),
	)
	return errs.Err
}

func (vm *VM) ParseBlock(_ context.Context, b []byte) (snowman.Block, error) {
	// Note: blocks to be parsed are not verified, so we must used blocks.Codec
	// rather than blocks.GenesisCodec
	statelessBlk, err := blocks.Parse(blocks.Codec, b)
	if err != nil {
		return nil, err
	}
	return vm.manager.NewBlock(statelessBlk), nil
}

func (vm *VM) GetBlock(_ context.Context, blkID ids.ID) (snowman.Block, error) {
	return vm.manager.GetBlock(blkID)
}

// LastAccepted returns the block most recently accepted
func (vm *VM) LastAccepted(context.Context) (ids.ID, error) {
	return vm.manager.LastAccepted(), nil
}

// SetPreference sets the preferred block to be the one with ID [blkID]
func (vm *VM) SetPreference(_ context.Context, blkID ids.ID) error {
	vm.Builder.SetPreference(blkID)
	return nil
}

func (*VM) Version(context.Context) (string, error) {
	return version.Current.String(), nil
}

// CreateHandlers returns a map where:
// * keys are API endpoint extensions
// * values are API handlers
func (vm *VM) CreateHandlers(context.Context) (map[string]*common.HTTPHandler, error) {
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	server.RegisterInterceptFunc(vm.metrics.InterceptRequest)
	server.RegisterAfterFunc(vm.metrics.AfterRequest)
	if err := server.RegisterService(
		&Service{
			vm:          vm,
			addrManager: avax.NewAddressManager(vm.ctx),
			stakerAttributesCache: &cache.LRU[ids.ID, *stakerAttributes]{
				Size: stakerAttributesCacheSize,
			},
		},
		"platform",
	); err != nil {
		return nil, err
	}

	return map[string]*common.HTTPHandler{
		"": {
			Handler: server,
		},
	}, nil
}

// CreateStaticHandlers returns a map where:
// * keys are API endpoint extensions
// * values are API handlers
func (*VM) CreateStaticHandlers(context.Context) (map[string]*common.HTTPHandler, error) {
	server := rpc.NewServer()
	server.RegisterCodec(json.NewCodec(), "application/json")
	server.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	if err := server.RegisterService(&api.StaticService{}, "platform"); err != nil {
		return nil, err
	}

	return map[string]*common.HTTPHandler{
		"": {
			LockOptions: common.NoLock,
			Handler:     server,
		},
	}, nil
}

func (vm *VM) Connected(_ context.Context, nodeID ids.NodeID, _ *version.Application) error {
	return vm.uptimeManager.Connect(nodeID, constants.PrimaryNetworkID)
}

func (vm *VM) ConnectedSubnet(_ context.Context, nodeID ids.NodeID, subnetID ids.ID) error {
	return vm.uptimeManager.Connect(nodeID, subnetID)
}

func (vm *VM) Disconnected(_ context.Context, nodeID ids.NodeID) error {
	if err := vm.uptimeManager.Disconnect(nodeID); err != nil {
		return err
	}
	return vm.state.Commit()
}

func (vm *VM) CodecRegistry() codec.Registry {
	return vm.codecRegistry
}

func (vm *VM) Clock() *mockable.Clock {
	return &vm.clock
}

func (vm *VM) Logger() logging.Logger {
	return vm.ctx.Log
}

func (vm *VM) VerifyHeightIndex(_ context.Context) error {
	if vm.pruned.Get() {
		return nil
	}

	return block.ErrIndexIncomplete
}

func (vm *VM) GetBlockIDAtHeight(_ context.Context, height uint64) (ids.ID, error) {
	return vm.state.GetBlockIDAtHeight(height)
}
