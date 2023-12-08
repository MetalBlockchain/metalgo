// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"go.uber.org/mock/gomock"

	"github.com/MetalBlockchain/metalgo/chains"
	"github.com/MetalBlockchain/metalgo/chains/atomic"
	"github.com/MetalBlockchain/metalgo/codec"
	"github.com/MetalBlockchain/metalgo/codec/linearcodec"
	"github.com/MetalBlockchain/metalgo/database"
	"github.com/MetalBlockchain/metalgo/database/memdb"
	"github.com/MetalBlockchain/metalgo/database/prefixdb"
	"github.com/MetalBlockchain/metalgo/database/versiondb"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/snow/engine/common"
	"github.com/MetalBlockchain/metalgo/snow/uptime"
	"github.com/MetalBlockchain/metalgo/snow/validators"
	"github.com/MetalBlockchain/metalgo/utils"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/crypto/secp256k1"
	"github.com/MetalBlockchain/metalgo/utils/formatting"
	"github.com/MetalBlockchain/metalgo/utils/formatting/address"
	"github.com/MetalBlockchain/metalgo/utils/json"
	"github.com/MetalBlockchain/metalgo/utils/logging"
	"github.com/MetalBlockchain/metalgo/utils/timer/mockable"
	"github.com/MetalBlockchain/metalgo/utils/units"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/api"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/config"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/fx"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/metrics"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/reward"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/state"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/status"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs/executor"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs/mempool"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/utxo"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"

	p_tx_builder "github.com/MetalBlockchain/metalgo/vms/platformvm/txs/builder"
	pvalidators "github.com/MetalBlockchain/metalgo/vms/platformvm/validators"
)

const (
	pending stakerStatus = iota
	current

	defaultWeight = 10000
	trackChecksum = false
)

var (
	_ mempool.BlockTimer = (*environment)(nil)

	defaultMinStakingDuration = 24 * time.Hour
	defaultMaxStakingDuration = 365 * 24 * time.Hour
	defaultGenesisTime        = time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC)
	defaultValidateStartTime  = defaultGenesisTime
	defaultValidateEndTime    = defaultValidateStartTime.Add(10 * defaultMinStakingDuration)
	defaultMinValidatorStake  = 5 * units.MilliAvax
	defaultBalance            = 100 * defaultMinValidatorStake
	preFundedKeys             = secp256k1.TestKeys()
	avaxAssetID               = ids.ID{'y', 'e', 'e', 't'}
	defaultTxFee              = uint64(100)
	xChainID                  = ids.Empty.Prefix(0)
	cChainID                  = ids.Empty.Prefix(1)

	genesisBlkID ids.ID
	testSubnet1  *txs.Tx

	errMissing = errors.New("missing")
)

type stakerStatus uint

type staker struct {
	nodeID             ids.NodeID
	rewardAddress      ids.ShortID
	startTime, endTime time.Time
}

type test struct {
	description           string
	stakers               []staker
	subnetStakers         []staker
	advanceTimeTo         []time.Time
	expectedStakers       map[ids.NodeID]stakerStatus
	expectedSubnetStakers map[ids.NodeID]stakerStatus
}

type environment struct {
	blkManager Manager
	mempool    mempool.Mempool
	sender     *common.SenderTest

	isBootstrapped *utils.Atomic[bool]
	config         *config.Config
	clk            *mockable.Clock
	baseDB         *versiondb.Database
	ctx            *snow.Context
	fx             fx.Fx
	state          state.State
	mockedState    *state.MockState
	atomicUTXOs    avax.AtomicUTXOManager
	uptimes        uptime.Manager
	utxosHandler   utxo.Handler
	txBuilder      p_tx_builder.Builder
	backend        *executor.Backend
}

func (*environment) ResetBlockTimer() {
	// dummy call, do nothing for now
}

func newEnvironment(t *testing.T, ctrl *gomock.Controller) *environment {
	res := &environment{
		isBootstrapped: &utils.Atomic[bool]{},
		config:         defaultConfig(),
		clk:            defaultClock(),
	}
	res.isBootstrapped.Set(true)

	res.baseDB = versiondb.New(memdb.New())
	res.ctx = defaultCtx(res.baseDB)
	res.fx = defaultFx(res.clk, res.ctx.Log, res.isBootstrapped.Get())

	rewardsCalc := reward.NewCalculator(res.config.RewardConfig)
	res.atomicUTXOs = avax.NewAtomicUTXOManager(res.ctx.SharedMemory, txs.Codec)

	if ctrl == nil {
		res.state = defaultState(res.config, res.ctx, res.baseDB, rewardsCalc)
		res.uptimes = uptime.NewManager(res.state, res.clk)
		res.utxosHandler = utxo.NewHandler(res.ctx, res.clk, res.fx)
		res.txBuilder = p_tx_builder.New(
			res.ctx,
			res.config,
			res.clk,
			res.fx,
			res.state,
			res.atomicUTXOs,
			res.utxosHandler,
		)
	} else {
		genesisBlkID = ids.GenerateTestID()
		res.mockedState = state.NewMockState(ctrl)
		res.uptimes = uptime.NewManager(res.mockedState, res.clk)
		res.utxosHandler = utxo.NewHandler(res.ctx, res.clk, res.fx)
		res.txBuilder = p_tx_builder.New(
			res.ctx,
			res.config,
			res.clk,
			res.fx,
			res.mockedState,
			res.atomicUTXOs,
			res.utxosHandler,
		)

		// setup expectations strictly needed for environment creation
		res.mockedState.EXPECT().GetLastAccepted().Return(genesisBlkID).Times(1)
	}

	res.backend = &executor.Backend{
		Config:       res.config,
		Ctx:          res.ctx,
		Clk:          res.clk,
		Bootstrapped: res.isBootstrapped,
		Fx:           res.fx,
		FlowChecker:  res.utxosHandler,
		Uptimes:      res.uptimes,
		Rewards:      rewardsCalc,
	}

	registerer := prometheus.NewRegistry()
	res.sender = &common.SenderTest{T: t}

	metrics := metrics.Noop

	var err error
	res.mempool, err = mempool.NewMempool("mempool", registerer, res)
	if err != nil {
		panic(fmt.Errorf("failed to create mempool: %w", err))
	}

	if ctrl == nil {
		res.blkManager = NewManager(
			res.mempool,
			metrics,
			res.state,
			res.backend,
			pvalidators.TestManager,
		)
		addSubnet(res)
	} else {
		res.blkManager = NewManager(
			res.mempool,
			metrics,
			res.mockedState,
			res.backend,
			pvalidators.TestManager,
		)
		// we do not add any subnet to state, since we can mock
		// whatever we need
	}

	return res
}

func addSubnet(env *environment) {
	// Create a subnet
	var err error
	testSubnet1, err = env.txBuilder.NewCreateSubnetTx(
		2, // threshold; 2 sigs from keys[0], keys[1], keys[2] needed to add validator to this subnet
		[]ids.ShortID{ // control keys
			preFundedKeys[0].PublicKey().Address(),
			preFundedKeys[1].PublicKey().Address(),
			preFundedKeys[2].PublicKey().Address(),
		},
		[]*secp256k1.PrivateKey{preFundedKeys[0]},
		preFundedKeys[0].PublicKey().Address(),
	)
	if err != nil {
		panic(err)
	}

	// store it
	genesisID := env.state.GetLastAccepted()
	stateDiff, err := state.NewDiff(genesisID, env.blkManager)
	if err != nil {
		panic(err)
	}

	executor := executor.StandardTxExecutor{
		Backend: env.backend,
		State:   stateDiff,
		Tx:      testSubnet1,
	}
	err = testSubnet1.Unsigned.Visit(&executor)
	if err != nil {
		panic(err)
	}

	stateDiff.AddTx(testSubnet1, status.Committed)
	if err := stateDiff.Apply(env.state); err != nil {
		panic(err)
	}
}

func defaultState(
	cfg *config.Config,
	ctx *snow.Context,
	db database.Database,
	rewards reward.Calculator,
) state.State {
	genesisBytes := buildGenesisTest(ctx)
	execCfg, _ := config.GetExecutionConfig([]byte(`{}`))
	state, err := state.New(
		db,
		genesisBytes,
		prometheus.NewRegistry(),
		cfg,
		execCfg,
		ctx,
		metrics.Noop,
		rewards,
	)
	if err != nil {
		panic(err)
	}

	// persist and reload to init a bunch of in-memory stuff
	state.SetHeight(0)
	if err := state.Commit(); err != nil {
		panic(err)
	}
	genesisBlkID = state.GetLastAccepted()
	return state
}

func defaultCtx(db database.Database) *snow.Context {
	ctx := snow.DefaultContextTest()
	ctx.NetworkID = 10
	ctx.XChainID = xChainID
	ctx.CChainID = cChainID
	ctx.AVAXAssetID = avaxAssetID

	atomicDB := prefixdb.New([]byte{1}, db)
	m := atomic.NewMemory(atomicDB)

	ctx.SharedMemory = m.NewSharedMemory(ctx.ChainID)

	ctx.ValidatorState = &validators.TestState{
		GetSubnetIDF: func(_ context.Context, chainID ids.ID) (ids.ID, error) {
			subnetID, ok := map[ids.ID]ids.ID{
				constants.PlatformChainID: constants.PrimaryNetworkID,
				xChainID:                  constants.PrimaryNetworkID,
				cChainID:                  constants.PrimaryNetworkID,
			}[chainID]
			if !ok {
				return ids.Empty, errMissing
			}
			return subnetID, nil
		},
	}

	return ctx
}

func defaultConfig() *config.Config {
	return &config.Config{
		Chains:                 chains.TestManager,
		UptimeLockedCalculator: uptime.NewLockedCalculator(),
		Validators:             validators.NewManager(),
		TxFee:                  defaultTxFee,
		CreateSubnetTxFee:      100 * defaultTxFee,
		CreateBlockchainTxFee:  100 * defaultTxFee,
		MinValidatorStake:      5 * units.MilliAvax,
		MaxValidatorStake:      500 * units.MilliAvax,
		MinDelegatorStake:      1 * units.MilliAvax,
		MinStakeDuration:       defaultMinStakingDuration,
		MaxStakeDuration:       defaultMaxStakingDuration,
		RewardConfig: reward.Config{
			MaxConsumptionRate: .12 * reward.PercentDenominator,
			MinConsumptionRate: .10 * reward.PercentDenominator,
			MintingPeriod:      365 * 24 * time.Hour,
			SupplyCap:          720 * units.MegaAvax,
		},
		ApricotPhase3Time: defaultValidateEndTime,
		ApricotPhase5Time: defaultValidateEndTime,
		BanffTime:         mockable.MaxTime,
	}
}

func defaultClock() *mockable.Clock {
	clk := &mockable.Clock{}
	clk.Set(defaultGenesisTime)
	return clk
}

type fxVMInt struct {
	registry codec.Registry
	clk      *mockable.Clock
	log      logging.Logger
}

func (fvi *fxVMInt) CodecRegistry() codec.Registry {
	return fvi.registry
}

func (fvi *fxVMInt) Clock() *mockable.Clock {
	return fvi.clk
}

func (fvi *fxVMInt) Logger() logging.Logger {
	return fvi.log
}

func defaultFx(clk *mockable.Clock, log logging.Logger, isBootstrapped bool) fx.Fx {
	fxVMInt := &fxVMInt{
		registry: linearcodec.NewDefault(),
		clk:      clk,
		log:      log,
	}
	res := &secp256k1fx.Fx{}
	if err := res.Initialize(fxVMInt); err != nil {
		panic(err)
	}
	if isBootstrapped {
		if err := res.Bootstrapped(); err != nil {
			panic(err)
		}
	}
	return res
}

func buildGenesisTest(ctx *snow.Context) []byte {
	genesisUTXOs := make([]api.UTXO, len(preFundedKeys))
	for i, key := range preFundedKeys {
		id := key.PublicKey().Address()
		addr, err := address.FormatBech32(constants.UnitTestHRP, id.Bytes())
		if err != nil {
			panic(err)
		}
		genesisUTXOs[i] = api.UTXO{
			Amount:  json.Uint64(defaultBalance),
			Address: addr,
		}
	}

	genesisValidators := make([]api.PermissionlessValidator, len(preFundedKeys))
	for i, key := range preFundedKeys {
		nodeID := ids.NodeID(key.PublicKey().Address())
		addr, err := address.FormatBech32(constants.UnitTestHRP, nodeID.Bytes())
		if err != nil {
			panic(err)
		}
		genesisValidators[i] = api.PermissionlessValidator{
			Staker: api.Staker{
				StartTime: json.Uint64(defaultValidateStartTime.Unix()),
				EndTime:   json.Uint64(defaultValidateEndTime.Unix()),
				NodeID:    nodeID,
			},
			RewardOwner: &api.Owner{
				Threshold: 1,
				Addresses: []string{addr},
			},
			Staked: []api.UTXO{{
				Amount:  json.Uint64(defaultWeight),
				Address: addr,
			}},
			DelegationFee: reward.PercentDenominator,
		}
	}

	buildGenesisArgs := api.BuildGenesisArgs{
		NetworkID:     json.Uint32(constants.UnitTestID),
		AvaxAssetID:   ctx.AVAXAssetID,
		UTXOs:         genesisUTXOs,
		Validators:    genesisValidators,
		Chains:        nil,
		Time:          json.Uint64(defaultGenesisTime.Unix()),
		InitialSupply: json.Uint64(360 * units.MegaAvax),
		Encoding:      formatting.Hex,
	}

	buildGenesisResponse := api.BuildGenesisReply{}
	platformvmSS := api.StaticService{}
	if err := platformvmSS.BuildGenesis(nil, &buildGenesisArgs, &buildGenesisResponse); err != nil {
		panic(fmt.Errorf("problem while building platform chain's genesis state: %w", err))
	}

	genesisBytes, err := formatting.Decode(buildGenesisResponse.Encoding, buildGenesisResponse.Bytes)
	if err != nil {
		panic(err)
	}

	return genesisBytes
}

func shutdownEnvironment(t *environment) error {
	if t.mockedState != nil {
		// state is mocked, nothing to do here
		return nil
	}

	if t.isBootstrapped.Get() {
		validatorIDs := t.config.Validators.GetValidatorIDs(constants.PrimaryNetworkID)

		if err := t.uptimes.StopTracking(validatorIDs, constants.PrimaryNetworkID); err != nil {
			return err
		}
		if err := t.state.Commit(); err != nil {
			return err
		}
	}

	var err error
	if t.state != nil {
		err = t.state.Close()
	}
	return utils.Err(
		err,
		t.baseDB.Close(),
	)
}

func addPendingValidator(
	env *environment,
	startTime time.Time,
	endTime time.Time,
	nodeID ids.NodeID,
	rewardAddress ids.ShortID,
	keys []*secp256k1.PrivateKey,
) (*txs.Tx, error) {
	addPendingValidatorTx, err := env.txBuilder.NewAddValidatorTx(
		env.config.MinValidatorStake,
		uint64(startTime.Unix()),
		uint64(endTime.Unix()),
		nodeID,
		rewardAddress,
		reward.PercentDenominator,
		keys,
		ids.ShortEmpty,
	)
	if err != nil {
		return nil, err
	}

	staker, err := state.NewPendingStaker(
		addPendingValidatorTx.ID(),
		addPendingValidatorTx.Unsigned.(*txs.AddValidatorTx),
	)
	if err != nil {
		return nil, err
	}

	env.state.PutPendingValidator(staker)
	env.state.AddTx(addPendingValidatorTx, status.Committed)
	dummyHeight := uint64(1)
	env.state.SetHeight(dummyHeight)
	if err := env.state.Commit(); err != nil {
		return nil, err
	}
	return addPendingValidatorTx, nil
}
