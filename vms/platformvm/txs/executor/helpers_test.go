// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/MetalBlockchain/metalgo/chains"
	"github.com/MetalBlockchain/metalgo/chains/atomic"
	"github.com/MetalBlockchain/metalgo/codec"
	"github.com/MetalBlockchain/metalgo/codec/linearcodec"
	"github.com/MetalBlockchain/metalgo/database"
	"github.com/MetalBlockchain/metalgo/database/manager"
	"github.com/MetalBlockchain/metalgo/database/prefixdb"
	"github.com/MetalBlockchain/metalgo/database/versiondb"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/snow/uptime"
	"github.com/MetalBlockchain/metalgo/snow/validators"
	"github.com/MetalBlockchain/metalgo/utils"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/crypto"
	"github.com/MetalBlockchain/metalgo/utils/formatting"
	"github.com/MetalBlockchain/metalgo/utils/formatting/address"
	"github.com/MetalBlockchain/metalgo/utils/json"
	"github.com/MetalBlockchain/metalgo/utils/logging"
	"github.com/MetalBlockchain/metalgo/utils/timer/mockable"
	"github.com/MetalBlockchain/metalgo/utils/units"
	"github.com/MetalBlockchain/metalgo/utils/wrappers"
	"github.com/MetalBlockchain/metalgo/version"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/api"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/config"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/fx"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/reward"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/state"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/status"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs/builder"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/utxo"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
)

const (
	testNetworkID = 10 // To be used in tests
	defaultWeight = 10000
)

var (
	defaultMinStakingDuration = 24 * time.Hour
	defaultMaxStakingDuration = 365 * 24 * time.Hour
	defaultGenesisTime        = time.Date(1997, 1, 1, 0, 0, 0, 0, time.UTC)
	defaultValidateStartTime  = defaultGenesisTime
	defaultValidateEndTime    = defaultValidateStartTime.Add(10 * defaultMinStakingDuration)
	defaultMinValidatorStake  = 5 * units.MilliAvax
	defaultBalance            = 100 * defaultMinValidatorStake
	preFundedKeys             = crypto.BuildTestKeys()
	avaxAssetID               = ids.ID{'y', 'e', 'e', 't'}
	defaultTxFee              = uint64(100)
	xChainID                  = ids.Empty.Prefix(0)
	cChainID                  = ids.Empty.Prefix(1)
	lastAcceptedID            = ids.GenerateTestID()

	testSubnet1            *txs.Tx
	testSubnet1ControlKeys = preFundedKeys[0:3]

	// Used to create and use keys.
	testKeyfactory crypto.FactorySECP256K1R
)

type mutableSharedMemory struct {
	atomic.SharedMemory
}

type environment struct {
	isBootstrapped *utils.AtomicBool
	config         *config.Config
	clk            *mockable.Clock
	baseDB         *versiondb.Database
	ctx            *snow.Context
	msm            *mutableSharedMemory
	fx             fx.Fx
	state          state.State
	states         map[ids.ID]state.Chain
	atomicUTXOs    avax.AtomicUTXOManager
	uptimes        uptime.Manager
	utxosHandler   utxo.Handler
	txBuilder      builder.TxBuilder
	backend        Backend
}

func (e *environment) GetState(blkID ids.ID) (state.Chain, bool) {
	if blkID == lastAcceptedID {
		return e.state, true
	}
	chainState, ok := e.states[blkID]
	return chainState, ok
}

func (e *environment) SetState(blkID ids.ID, chainState state.Chain) {
	e.states[blkID] = chainState
}

// TODO: snLookup currently duplicated in vm_test.go. Remove duplication
type snLookup struct {
	chainsToSubnet map[ids.ID]ids.ID
}

func (sn *snLookup) SubnetID(chainID ids.ID) (ids.ID, error) {
	subnetID, ok := sn.chainsToSubnet[chainID]
	if !ok {
		return ids.ID{}, errors.New("")
	}
	return subnetID, nil
}

func newEnvironment() *environment {
	var isBootstrapped utils.AtomicBool
	isBootstrapped.SetValue(true)

	config := defaultConfig()
	clk := defaultClock()

	baseDBManager := manager.NewMemDB(version.CurrentDatabase)
	baseDB := versiondb.New(baseDBManager.Current().Database)
	ctx, msm := defaultCtx(baseDB)

	fx := defaultFx(&clk, ctx.Log, isBootstrapped.GetValue())

	rewards := reward.NewCalculator(config.RewardConfig)
	baseState := defaultState(&config, ctx, baseDB, rewards)

	atomicUTXOs := avax.NewAtomicUTXOManager(ctx.SharedMemory, txs.Codec)
	uptimes := uptime.NewManager(baseState)
	utxoHandler := utxo.NewHandler(ctx, &clk, baseState, fx)

	txBuilder := builder.NewTxBuilder(
		ctx,
		config,
		&clk,
		fx,
		baseState,
		atomicUTXOs,
		utxoHandler,
	)

	backend := Backend{
		Config:       &config,
		Ctx:          ctx,
		Clk:          &clk,
		Bootstrapped: &isBootstrapped,
		Fx:           fx,
		FlowChecker:  utxoHandler,
		Uptimes:      uptimes,
		Rewards:      rewards,
	}

	env := &environment{
		isBootstrapped: &isBootstrapped,
		config:         &config,
		clk:            &clk,
		baseDB:         baseDB,
		ctx:            ctx,
		msm:            msm,
		fx:             fx,
		state:          baseState,
		states:         make(map[ids.ID]state.Chain),
		atomicUTXOs:    atomicUTXOs,
		uptimes:        uptimes,
		utxosHandler:   utxoHandler,
		txBuilder:      txBuilder,
		backend:        backend,
	}

	addSubnet(env, txBuilder, backend)

	return env
}

func addSubnet(
	env *environment,
	txBuilder builder.TxBuilder,
	backend Backend,
) {
	// Create a subnet
	var err error
	testSubnet1, err = txBuilder.NewCreateSubnetTx(
		2, // threshold; 2 sigs from keys[0], keys[1], keys[2] needed to add validator to this subnet
		[]ids.ShortID{ // control keys
			preFundedKeys[0].PublicKey().Address(),
			preFundedKeys[1].PublicKey().Address(),
			preFundedKeys[2].PublicKey().Address(),
		},
		[]*crypto.PrivateKeySECP256K1R{preFundedKeys[0]},
		preFundedKeys[0].PublicKey().Address(),
	)
	if err != nil {
		panic(err)
	}

	// store it
	stateDiff, err := state.NewDiff(lastAcceptedID, env)
	if err != nil {
		panic(err)
	}

	executor := StandardTxExecutor{
		Backend: &backend,
		State:   stateDiff,
		Tx:      testSubnet1,
	}
	err = testSubnet1.Unsigned.Visit(&executor)
	if err != nil {
		panic(err)
	}

	stateDiff.AddTx(testSubnet1, status.Committed)
	stateDiff.Apply(env.state)
}

func defaultState(
	cfg *config.Config,
	ctx *snow.Context,
	db database.Database,
	rewards reward.Calculator,
) state.State {
	dummyLocalStake := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "uts",
		Name:      "local_staked",
		Help:      "Total amount of AVAX on this node staked",
	})
	dummyTotalStake := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "uts",
		Name:      "total_staked",
		Help:      "Total amount of AVAX staked",
	})

	genesisBytes := buildGenesisTest(ctx)
	state, err := state.New(
		db,
		genesisBytes,
		prometheus.NewRegistry(),
		cfg,
		ctx,
		dummyLocalStake,
		dummyTotalStake,
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
	state.SetHeight( /*height*/ 0)
	if err := state.Commit(); err != nil {
		panic(err)
	}

	return state
}

func defaultCtx(db database.Database) (*snow.Context, *mutableSharedMemory) {
	ctx := snow.DefaultContextTest()
	ctx.NetworkID = 10
	ctx.XChainID = xChainID
	ctx.AVAXAssetID = avaxAssetID

	atomicDB := prefixdb.New([]byte{1}, db)
	m := atomic.NewMemory(atomicDB)

	msm := &mutableSharedMemory{
		SharedMemory: m.NewSharedMemory(ctx.ChainID),
	}
	ctx.SharedMemory = msm

	ctx.SNLookup = &snLookup{
		chainsToSubnet: map[ids.ID]ids.ID{
			constants.PlatformChainID: constants.PrimaryNetworkID,
			xChainID:                  constants.PrimaryNetworkID,
			cChainID:                  constants.PrimaryNetworkID,
		},
	}

	return ctx, msm
}

func defaultConfig() config.Config {
	return config.Config{
		Chains:                 chains.MockManager{},
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
	}
}

func defaultClock() mockable.Clock {
	clk := mockable.Clock{}
	clk.Set(defaultGenesisTime)
	return clk
}

type fxVMInt struct {
	registry codec.Registry
	clk      *mockable.Clock
	log      logging.Logger
}

func (fvi *fxVMInt) CodecRegistry() codec.Registry { return fvi.registry }
func (fvi *fxVMInt) Clock() *mockable.Clock        { return fvi.clk }
func (fvi *fxVMInt) Logger() logging.Logger        { return fvi.log }

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
	hrp := constants.NetworkIDToHRP[testNetworkID]
	for i, key := range preFundedKeys {
		id := key.PublicKey().Address()
		addr, err := address.FormatBech32(hrp, id.Bytes())
		if err != nil {
			panic(err)
		}
		genesisUTXOs[i] = api.UTXO{
			Amount:  json.Uint64(defaultBalance),
			Address: addr,
		}
	}

	genesisValidators := make([]api.PrimaryValidator, len(preFundedKeys))
	for i, key := range preFundedKeys {
		nodeID := ids.NodeID(key.PublicKey().Address())
		addr, err := address.FormatBech32(hrp, nodeID.Bytes())
		if err != nil {
			panic(err)
		}
		genesisValidators[i] = api.PrimaryValidator{
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
		NetworkID:     json.Uint32(testNetworkID),
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
		panic(fmt.Errorf("problem while building platform chain's genesis state: %v", err))
	}

	genesisBytes, err := formatting.Decode(buildGenesisResponse.Encoding, buildGenesisResponse.Bytes)
	if err != nil {
		panic(err)
	}

	return genesisBytes
}

func shutdownEnvironment(env *environment) error {
	if env.isBootstrapped.GetValue() {
		primaryValidatorSet, exist := env.config.Validators.GetValidators(constants.PrimaryNetworkID)
		if !exist {
			return errors.New("no default subnet validators")
		}
		primaryValidators := primaryValidatorSet.List()

		validatorIDs := make([]ids.NodeID, len(primaryValidators))
		for i, vdr := range primaryValidators {
			validatorIDs[i] = vdr.ID()
		}

		if err := env.uptimes.Shutdown(validatorIDs); err != nil {
			return err
		}
		env.state.SetHeight( /*height*/ math.MaxUint64)
		if err := env.state.Commit(); err != nil {
			return err
		}
	}

	errs := wrappers.Errs{}
	errs.Add(
		env.state.Close(),
		env.baseDB.Close(),
	)
	return errs.Err
}
