// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package p

import (
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/crypto/bls"
	"github.com/MetalBlockchain/metalgo/utils/crypto/secp256k1"
	"github.com/MetalBlockchain/metalgo/utils/set"
	"github.com/MetalBlockchain/metalgo/utils/units"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/components/gas"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/fx"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/reward"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/signer"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/stakeable"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs/fee"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
	"github.com/MetalBlockchain/metalgo/vms/types"
	"github.com/MetalBlockchain/metalgo/wallet/chain/p/builder"
	"github.com/MetalBlockchain/metalgo/wallet/chain/p/wallet"
	"github.com/MetalBlockchain/metalgo/wallet/subnet/primary/common"
	"github.com/MetalBlockchain/metalgo/wallet/subnet/primary/common/utxotest"
)

var (
	subnetID = ids.GenerateTestID()
	nodeID   = ids.GenerateTestNodeID()

	testKeys       = secp256k1.TestKeys()
	subnetAuthKey  = testKeys[0]
	subnetAuthAddr = subnetAuthKey.Address()
	subnetOwner    = &secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs:     []ids.ShortID{subnetAuthAddr},
	}
	importKey   = testKeys[0]
	importAddr  = importKey.Address()
	importOwner = &secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs:     []ids.ShortID{importAddr},
	}
	rewardKey    = testKeys[0]
	rewardAddr   = rewardKey.Address()
	rewardsOwner = &secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs:     []ids.ShortID{rewardAddr},
	}
	utxoKey   = testKeys[1]
	utxoAddr  = utxoKey.Address()
	utxoOwner = secp256k1fx.OutputOwners{
		Threshold: 1,
		Addrs:     []ids.ShortID{utxoAddr},
	}

	// We hard-code [avaxAssetID] and [subnetAssetID] to make ordering of UTXOs
	// generated by [makeTestUTXOs] reproducible.
	avaxAssetID   = ids.Empty.Prefix(1789)
	subnetAssetID = ids.Empty.Prefix(2024)
	utxos         = makeTestUTXOs(utxoKey)

	avaxOutput = &avax.TransferableOutput{
		Asset: avax.Asset{ID: avaxAssetID},
		Out: &secp256k1fx.TransferOutput{
			Amt:          7 * units.Avax,
			OutputOwners: utxoOwner,
		},
	}

	subnetOwners = map[ids.ID]fx.Owner{
		subnetID: subnetOwner,
	}

	primaryNetworkPermissionlessStaker = &txs.SubnetValidator{
		Validator: txs.Validator{
			NodeID: nodeID,
			End:    uint64(time.Now().Add(time.Hour).Unix()),
			Wght:   2 * units.Avax,
		},
		Subnet: constants.PrimaryNetworkID,
	}

	testContextPreEtna = &builder.Context{
		NetworkID:   constants.UnitTestID,
		AVAXAssetID: avaxAssetID,
		StaticFeeConfig: fee.StaticConfig{
			TxFee:                         units.MicroAvax,
			CreateSubnetTxFee:             19 * units.MicroAvax,
			TransformSubnetTxFee:          789 * units.MicroAvax,
			CreateBlockchainTxFee:         1234 * units.MicroAvax,
			AddPrimaryNetworkValidatorFee: 19 * units.MilliAvax,
			AddPrimaryNetworkDelegatorFee: 765 * units.MilliAvax,
			AddSubnetValidatorFee:         1010 * units.MilliAvax,
			AddSubnetDelegatorFee:         9 * units.Avax,
		},
	}
	staticFeeCalculator = fee.NewStaticCalculator(
		testContextPreEtna.StaticFeeConfig,
	)

	testContextPostEtna = &builder.Context{
		NetworkID:   constants.UnitTestID,
		AVAXAssetID: avaxAssetID,

		ComplexityWeights: gas.Dimensions{
			gas.Bandwidth: 1,
			gas.DBRead:    10,
			gas.DBWrite:   100,
			gas.Compute:   1000,
		},
		GasPrice: 1,
	}
	dynamicFeeCalculator = fee.NewDynamicCalculator(
		testContextPostEtna.ComplexityWeights,
		testContextPostEtna.GasPrice,
	)

	testEnvironmentPreEtna = []environment{
		{
			name:          "Pre-Etna",
			context:       testContextPreEtna,
			feeCalculator: staticFeeCalculator,
		},
		{
			name:          "Pre-Etna with memo",
			context:       testContextPreEtna,
			feeCalculator: staticFeeCalculator,
			memo:          []byte("memo"),
		},
	}
	testEnvironmentPostEtna = []environment{
		{
			name:          "Post-Etna",
			context:       testContextPostEtna,
			feeCalculator: dynamicFeeCalculator,
		},
		{
			name:          "Post-Etna with memo",
			context:       testContextPostEtna,
			feeCalculator: dynamicFeeCalculator,
			memo:          []byte("memo"),
		},
	}
	testEnvironment = slices.Concat(
		testEnvironmentPreEtna,
		testEnvironmentPostEtna,
	)
)

type environment struct {
	name          string
	context       *builder.Context
	feeCalculator fee.Calculator
	memo          []byte
}

// These tests create a tx, then verify that utxos included in the tx are
// exactly necessary to pay fees for it.

func TestBaseTx(t *testing.T) {
	for _, e := range testEnvironment {
		t.Run(e.name, func(t *testing.T) {
			var (
				require    = require.New(t)
				chainUTXOs = utxotest.NewDeterministicChainUTXOs(t, map[ids.ID][]*avax.UTXO{
					constants.PlatformChainID: utxos,
				})
				backend = wallet.NewBackend(e.context, chainUTXOs, nil)
				builder = builder.New(set.Of(utxoAddr), e.context, backend)
			)

			utx, err := builder.NewBaseTx(
				[]*avax.TransferableOutput{avaxOutput},
				common.WithMemo(e.memo),
			)
			require.NoError(err)
			require.Contains(utx.Outs, avaxOutput)
			require.Equal(types.JSONByteSlice(e.memo), utx.Memo)
			requireFeeIsCorrect(
				require,
				e.feeCalculator,
				utx,
				&utx.BaseTx,
				nil,
				nil,
				nil,
			)
		})
	}
}

func TestAddSubnetValidatorTx(t *testing.T) {
	subnetValidator := &txs.SubnetValidator{
		Validator: txs.Validator{
			NodeID: nodeID,
			End:    uint64(time.Now().Add(time.Hour).Unix()),
		},
		Subnet: subnetID,
	}

	for _, e := range testEnvironment {
		t.Run(e.name, func(t *testing.T) {
			var (
				require    = require.New(t)
				chainUTXOs = utxotest.NewDeterministicChainUTXOs(t, map[ids.ID][]*avax.UTXO{
					constants.PlatformChainID: utxos,
				})
				backend = wallet.NewBackend(e.context, chainUTXOs, subnetOwners)
				builder = builder.New(set.Of(utxoAddr, subnetAuthAddr), e.context, backend)
			)

			utx, err := builder.NewAddSubnetValidatorTx(
				subnetValidator,
				common.WithMemo(e.memo),
			)
			require.NoError(err)
			require.Equal(*subnetValidator, utx.SubnetValidator)
			require.Equal(types.JSONByteSlice(e.memo), utx.Memo)
			requireFeeIsCorrect(
				require,
				e.feeCalculator,
				utx,
				&utx.BaseTx.BaseTx,
				nil,
				nil,
				nil,
			)
		})
	}
}

func TestRemoveSubnetValidatorTx(t *testing.T) {
	for _, e := range testEnvironment {
		t.Run(e.name, func(t *testing.T) {
			var (
				require    = require.New(t)
				chainUTXOs = utxotest.NewDeterministicChainUTXOs(t, map[ids.ID][]*avax.UTXO{
					constants.PlatformChainID: utxos,
				})
				backend = wallet.NewBackend(e.context, chainUTXOs, subnetOwners)
				builder = builder.New(set.Of(utxoAddr, subnetAuthAddr), e.context, backend)
			)

			utx, err := builder.NewRemoveSubnetValidatorTx(
				nodeID,
				subnetID,
				common.WithMemo(e.memo),
			)
			require.NoError(err)
			require.Equal(nodeID, utx.NodeID)
			require.Equal(subnetID, utx.Subnet)
			require.Equal(types.JSONByteSlice(e.memo), utx.Memo)
			requireFeeIsCorrect(
				require,
				e.feeCalculator,
				utx,
				&utx.BaseTx.BaseTx,
				nil,
				nil,
				nil,
			)
		})
	}
}

func TestCreateChainTx(t *testing.T) {
	var (
		genesisBytes = []byte{'a', 'b', 'c'}
		vmID         = ids.GenerateTestID()
		fxIDs        = []ids.ID{ids.GenerateTestID()}
		chainName    = "dummyChain"
	)

	for _, e := range testEnvironment {
		t.Run(e.name, func(t *testing.T) {
			var (
				require    = require.New(t)
				chainUTXOs = utxotest.NewDeterministicChainUTXOs(t, map[ids.ID][]*avax.UTXO{
					constants.PlatformChainID: utxos,
				})
				backend = wallet.NewBackend(e.context, chainUTXOs, subnetOwners)
				builder = builder.New(set.Of(utxoAddr, subnetAuthAddr), e.context, backend)
			)

			utx, err := builder.NewCreateChainTx(
				subnetID,
				genesisBytes,
				vmID,
				fxIDs,
				chainName,
				common.WithMemo(e.memo),
			)
			require.NoError(err)
			require.Equal(subnetID, utx.SubnetID)
			require.Equal(genesisBytes, utx.GenesisData)
			require.Equal(vmID, utx.VMID)
			require.ElementsMatch(fxIDs, utx.FxIDs)
			require.Equal(chainName, utx.ChainName)
			require.Equal(types.JSONByteSlice(e.memo), utx.Memo)
			requireFeeIsCorrect(
				require,
				e.feeCalculator,
				utx,
				&utx.BaseTx.BaseTx,
				nil,
				nil,
				nil,
			)
		})
	}
}

func TestCreateSubnetTx(t *testing.T) {
	for _, e := range testEnvironment {
		t.Run(e.name, func(t *testing.T) {
			var (
				require    = require.New(t)
				chainUTXOs = utxotest.NewDeterministicChainUTXOs(t, map[ids.ID][]*avax.UTXO{
					constants.PlatformChainID: utxos,
				})
				backend = wallet.NewBackend(e.context, chainUTXOs, subnetOwners)
				builder = builder.New(set.Of(utxoAddr, subnetAuthAddr), e.context, backend)
			)

			utx, err := builder.NewCreateSubnetTx(
				subnetOwner,
				common.WithMemo(e.memo),
			)
			require.NoError(err)
			require.Equal(subnetOwner, utx.Owner)
			require.Equal(types.JSONByteSlice(e.memo), utx.Memo)
			requireFeeIsCorrect(
				require,
				e.feeCalculator,
				utx,
				&utx.BaseTx.BaseTx,
				nil,
				nil,
				nil,
			)
		})
	}
}

func TestTransferSubnetOwnershipTx(t *testing.T) {
	for _, e := range testEnvironment {
		t.Run(e.name, func(t *testing.T) {
			var (
				require    = require.New(t)
				chainUTXOs = utxotest.NewDeterministicChainUTXOs(t, map[ids.ID][]*avax.UTXO{
					constants.PlatformChainID: utxos,
				})
				backend = wallet.NewBackend(e.context, chainUTXOs, subnetOwners)
				builder = builder.New(set.Of(utxoAddr, subnetAuthAddr), e.context, backend)
			)

			utx, err := builder.NewTransferSubnetOwnershipTx(
				subnetID,
				subnetOwner,
				common.WithMemo(e.memo),
			)
			require.NoError(err)
			require.Equal(subnetID, utx.Subnet)
			require.Equal(subnetOwner, utx.Owner)
			require.Equal(types.JSONByteSlice(e.memo), utx.Memo)
			requireFeeIsCorrect(
				require,
				e.feeCalculator,
				utx,
				&utx.BaseTx.BaseTx,
				nil,
				nil,
				nil,
			)
		})
	}
}

func TestImportTx(t *testing.T) {
	var (
		sourceChainID = ids.GenerateTestID()
		importedUTXOs = utxos[:1]
	)

	for _, e := range testEnvironment {
		t.Run(e.name, func(t *testing.T) {
			var (
				require    = require.New(t)
				chainUTXOs = utxotest.NewDeterministicChainUTXOs(t, map[ids.ID][]*avax.UTXO{
					constants.PlatformChainID: utxos,
					sourceChainID:             importedUTXOs,
				})
				backend = wallet.NewBackend(e.context, chainUTXOs, nil)
				builder = builder.New(set.Of(utxoAddr), e.context, backend)
			)

			utx, err := builder.NewImportTx(
				sourceChainID,
				importOwner,
				common.WithMemo(e.memo),
			)
			require.NoError(err)
			require.Equal(sourceChainID, utx.SourceChain)
			require.Equal(types.JSONByteSlice(e.memo), utx.Memo)
			require.Empty(utx.Ins)                              // The imported input should be sufficient for fees
			require.Len(utx.ImportedInputs, len(importedUTXOs)) // All utxos should be imported
			requireFeeIsCorrect(
				require,
				e.feeCalculator,
				utx,
				&utx.BaseTx.BaseTx,
				utx.ImportedInputs,
				nil,
				nil,
			)
		})
	}
}

func TestExportTx(t *testing.T) {
	exportedOutputs := []*avax.TransferableOutput{avaxOutput}

	for _, e := range testEnvironment {
		t.Run(e.name, func(t *testing.T) {
			var (
				require    = require.New(t)
				chainUTXOs = utxotest.NewDeterministicChainUTXOs(t, map[ids.ID][]*avax.UTXO{
					constants.PlatformChainID: utxos,
				})
				backend = wallet.NewBackend(e.context, chainUTXOs, nil)
				builder = builder.New(set.Of(utxoAddr), e.context, backend)
			)

			utx, err := builder.NewExportTx(
				subnetID,
				exportedOutputs,
				common.WithMemo(e.memo),
			)
			require.NoError(err)
			require.Equal(subnetID, utx.DestinationChain)
			require.ElementsMatch(exportedOutputs, utx.ExportedOutputs)
			require.Equal(types.JSONByteSlice(e.memo), utx.Memo)
			requireFeeIsCorrect(
				require,
				e.feeCalculator,
				utx,
				&utx.BaseTx.BaseTx,
				nil,
				utx.ExportedOutputs,
				nil,
			)
		})
	}
}

func TestTransformSubnetTx(t *testing.T) {
	const (
		initialSupply                   = 40 * units.MegaAvax
		maxSupply                       = 100 * units.MegaAvax
		minConsumptionRate       uint64 = reward.PercentDenominator
		maxConsumptionRate       uint64 = reward.PercentDenominator
		minValidatorStake        uint64 = 1
		maxValidatorStake               = 100 * units.MegaAvax
		minStakeDuration                = time.Second
		maxStakeDuration                = 365 * 24 * time.Hour
		minDelegationFee         uint32 = 0
		minDelegatorStake        uint64 = 1
		maxValidatorWeightFactor byte   = 5
		uptimeRequirement        uint32 = .80 * reward.PercentDenominator
	)

	// TransformSubnetTx is not valid to be issued post-Etna
	for _, e := range testEnvironmentPreEtna {
		t.Run(e.name, func(t *testing.T) {
			var (
				require    = require.New(t)
				chainUTXOs = utxotest.NewDeterministicChainUTXOs(t, map[ids.ID][]*avax.UTXO{
					constants.PlatformChainID: utxos,
				})
				backend = wallet.NewBackend(e.context, chainUTXOs, subnetOwners)
				builder = builder.New(set.Of(utxoAddr, subnetAuthAddr), e.context, backend)
			)

			utx, err := builder.NewTransformSubnetTx(
				subnetID,
				subnetAssetID,
				initialSupply,
				maxSupply,
				minConsumptionRate,
				maxConsumptionRate,
				minValidatorStake,
				maxValidatorStake,
				minStakeDuration,
				maxStakeDuration,
				minDelegationFee,
				minDelegatorStake,
				maxValidatorWeightFactor,
				uptimeRequirement,
				common.WithMemo(e.memo),
			)
			require.NoError(err)
			require.Equal(subnetID, utx.Subnet)
			require.Equal(subnetAssetID, utx.AssetID)
			require.Equal(initialSupply, utx.InitialSupply)
			require.Equal(maxSupply, utx.MaximumSupply)
			require.Equal(minConsumptionRate, utx.MinConsumptionRate)
			require.Equal(minValidatorStake, utx.MinValidatorStake)
			require.Equal(maxValidatorStake, utx.MaxValidatorStake)
			require.Equal(uint32(minStakeDuration/time.Second), utx.MinStakeDuration)
			require.Equal(uint32(maxStakeDuration/time.Second), utx.MaxStakeDuration)
			require.Equal(minDelegationFee, utx.MinDelegationFee)
			require.Equal(minDelegatorStake, utx.MinDelegatorStake)
			require.Equal(maxValidatorWeightFactor, utx.MaxValidatorWeightFactor)
			require.Equal(uptimeRequirement, utx.UptimeRequirement)
			require.Equal(types.JSONByteSlice(e.memo), utx.Memo)
			requireFeeIsCorrect(
				require,
				e.feeCalculator,
				utx,
				&utx.BaseTx.BaseTx,
				nil,
				nil,
				map[ids.ID]uint64{
					subnetAssetID: maxSupply - initialSupply,
				},
			)
		})
	}
}

func TestAddPermissionlessValidatorTx(t *testing.T) {
	var utxosOffset uint64 = 2024
	makeUTXO := func(amount uint64) *avax.UTXO {
		utxosOffset++
		return &avax.UTXO{
			UTXOID: avax.UTXOID{
				TxID:        ids.Empty.Prefix(utxosOffset),
				OutputIndex: uint32(utxosOffset),
			},
			Asset: avax.Asset{ID: avaxAssetID},
			Out: &secp256k1fx.TransferOutput{
				Amt:          amount,
				OutputOwners: utxoOwner,
			},
		}
	}

	var (
		utxos = []*avax.UTXO{
			makeUTXO(testContextPreEtna.StaticFeeConfig.AddPrimaryNetworkValidatorFee), // UTXO to pay the fee
			makeUTXO(1 * units.NanoAvax), // small UTXO
			makeUTXO(9 * units.Avax),     // large UTXO
		}

		validationRewardsOwner        = rewardsOwner
		delegationRewardsOwner        = rewardsOwner
		delegationShares       uint32 = reward.PercentDenominator
	)

	sk, err := bls.NewSecretKey()
	require.NoError(t, err)

	pop := signer.NewProofOfPossession(sk)

	for _, e := range testEnvironment {
		t.Run(e.name, func(t *testing.T) {
			var (
				require    = require.New(t)
				chainUTXOs = utxotest.NewDeterministicChainUTXOs(t, map[ids.ID][]*avax.UTXO{
					constants.PlatformChainID: utxos,
				})
				backend = wallet.NewBackend(e.context, chainUTXOs, nil)
				builder = builder.New(set.Of(utxoAddr, rewardAddr), e.context, backend)
			)

			utx, err := builder.NewAddPermissionlessValidatorTx(
				primaryNetworkPermissionlessStaker,
				pop,
				avaxAssetID,
				validationRewardsOwner,
				delegationRewardsOwner,
				delegationShares,
				common.WithMemo(e.memo),
			)
			require.NoError(err)
			require.Equal(primaryNetworkPermissionlessStaker.Validator, utx.Validator)
			require.Equal(primaryNetworkPermissionlessStaker.Subnet, utx.Subnet)
			require.Equal(pop, utx.Signer)
			// Outputs should be merged if possible. For example, if there are two
			// unlocked inputs consumed for staking, this should only produce one staked
			// output.
			require.Len(utx.StakeOuts, 1)
			// check stake amount
			require.Equal(
				map[ids.ID]uint64{
					avaxAssetID: primaryNetworkPermissionlessStaker.Wght,
				},
				addOutputAmounts(utx.StakeOuts),
			)
			require.Equal(validationRewardsOwner, utx.ValidatorRewardsOwner)
			require.Equal(delegationRewardsOwner, utx.DelegatorRewardsOwner)
			require.Equal(delegationShares, utx.DelegationShares)
			require.Equal(types.JSONByteSlice(e.memo), utx.Memo)
			requireFeeIsCorrect(
				require,
				e.feeCalculator,
				utx,
				&utx.BaseTx.BaseTx,
				nil,
				utx.StakeOuts,
				nil,
			)
		})
	}
}

func TestAddPermissionlessDelegatorTx(t *testing.T) {
	for _, e := range testEnvironment {
		t.Run(e.name, func(t *testing.T) {
			var (
				require    = require.New(t)
				chainUTXOs = utxotest.NewDeterministicChainUTXOs(t, map[ids.ID][]*avax.UTXO{
					constants.PlatformChainID: utxos,
				})
				backend = wallet.NewBackend(e.context, chainUTXOs, nil)
				builder = builder.New(set.Of(utxoAddr, rewardAddr), e.context, backend)
			)

			utx, err := builder.NewAddPermissionlessDelegatorTx(
				primaryNetworkPermissionlessStaker,
				avaxAssetID,
				rewardsOwner,
				common.WithMemo(e.memo),
			)
			require.NoError(err)
			require.Equal(primaryNetworkPermissionlessStaker.Validator, utx.Validator)
			require.Equal(primaryNetworkPermissionlessStaker.Subnet, utx.Subnet)
			// check stake amount
			require.Equal(
				map[ids.ID]uint64{
					avaxAssetID: primaryNetworkPermissionlessStaker.Wght,
				},
				addOutputAmounts(utx.StakeOuts),
			)
			require.Equal(rewardsOwner, utx.DelegationRewardsOwner)
			require.Equal(types.JSONByteSlice(e.memo), utx.Memo)
			requireFeeIsCorrect(
				require,
				e.feeCalculator,
				utx,
				&utx.BaseTx.BaseTx,
				nil,
				utx.StakeOuts,
				nil,
			)
		})
	}
}

func TestConvertSubnetTx(t *testing.T) {
	var (
		chainID = ids.GenerateTestID()
		address = utils.RandomBytes(32)
	)
	for _, e := range testEnvironmentPostEtna {
		t.Run(e.name, func(t *testing.T) {
			var (
				require    = require.New(t)
				chainUTXOs = utxotest.NewDeterministicChainUTXOs(t, map[ids.ID][]*avax.UTXO{
					constants.PlatformChainID: utxos,
				})
				backend = wallet.NewBackend(e.context, chainUTXOs, subnetOwners)
				builder = builder.New(set.Of(utxoAddr, subnetAuthAddr), e.context, backend)
			)

			utx, err := builder.NewConvertSubnetTx(
				subnetID,
				chainID,
				address,
				common.WithMemo(e.memo),
			)
			require.NoError(err)
			require.Equal(subnetID, utx.Subnet)
			require.Equal(chainID, utx.ChainID)
			require.Equal(types.JSONByteSlice(address), utx.Address)
			require.Equal(types.JSONByteSlice(e.memo), utx.Memo)
			requireFeeIsCorrect(
				require,
				e.feeCalculator,
				utx,
				&utx.BaseTx.BaseTx,
				nil,
				nil,
				nil,
			)
		})
	}
}

func makeTestUTXOs(utxosKey *secp256k1.PrivateKey) []*avax.UTXO {
	// Note: we avoid ids.GenerateTestNodeID here to make sure that UTXO IDs
	// won't change run by run. This simplifies checking what utxos are included
	// in the built txs.
	const utxosOffset uint64 = 2024

	utxosAddr := utxosKey.Address()
	return []*avax.UTXO{
		{ // a small UTXO first, which should not be enough to pay fees
			UTXOID: avax.UTXOID{
				TxID:        ids.Empty.Prefix(utxosOffset),
				OutputIndex: uint32(utxosOffset),
			},
			Asset: avax.Asset{ID: avaxAssetID},
			Out: &secp256k1fx.TransferOutput{
				Amt: 2 * units.MilliAvax,
				OutputOwners: secp256k1fx.OutputOwners{
					Locktime:  0,
					Addrs:     []ids.ShortID{utxosAddr},
					Threshold: 1,
				},
			},
		},
		{ // a locked, small UTXO
			UTXOID: avax.UTXOID{
				TxID:        ids.Empty.Prefix(utxosOffset + 1),
				OutputIndex: uint32(utxosOffset + 1),
			},
			Asset: avax.Asset{ID: avaxAssetID},
			Out: &stakeable.LockOut{
				Locktime: uint64(time.Now().Add(time.Hour).Unix()),
				TransferableOut: &secp256k1fx.TransferOutput{
					Amt: 3 * units.MilliAvax,
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs:     []ids.ShortID{utxosAddr},
					},
				},
			},
		},
		{ // a subnetAssetID denominated UTXO
			UTXOID: avax.UTXOID{
				TxID:        ids.Empty.Prefix(utxosOffset + 2),
				OutputIndex: uint32(utxosOffset + 2),
			},
			Asset: avax.Asset{ID: subnetAssetID},
			Out: &secp256k1fx.TransferOutput{
				Amt: 99 * units.MegaAvax,
				OutputOwners: secp256k1fx.OutputOwners{
					Locktime:  0,
					Addrs:     []ids.ShortID{utxosAddr},
					Threshold: 1,
				},
			},
		},
		{ // a locked, large UTXO
			UTXOID: avax.UTXOID{
				TxID:        ids.Empty.Prefix(utxosOffset + 3),
				OutputIndex: uint32(utxosOffset + 3),
			},
			Asset: avax.Asset{ID: avaxAssetID},
			Out: &stakeable.LockOut{
				Locktime: uint64(time.Now().Add(time.Hour).Unix()),
				TransferableOut: &secp256k1fx.TransferOutput{
					Amt: 88 * units.Avax,
					OutputOwners: secp256k1fx.OutputOwners{
						Threshold: 1,
						Addrs:     []ids.ShortID{utxosAddr},
					},
				},
			},
		},
		{ // a large UTXO last, which should be enough to pay any fee by itself
			UTXOID: avax.UTXOID{
				TxID:        ids.Empty.Prefix(utxosOffset + 4),
				OutputIndex: uint32(utxosOffset + 4),
			},
			Asset: avax.Asset{ID: avaxAssetID},
			Out: &secp256k1fx.TransferOutput{
				Amt: 9 * units.Avax,
				OutputOwners: secp256k1fx.OutputOwners{
					Locktime:  0,
					Addrs:     []ids.ShortID{utxosAddr},
					Threshold: 1,
				},
			},
		},
	}
}

// requireFeeIsCorrect calculates the required fee for the unsigned transaction
// and verifies that the burned amount is exactly the required fee.
func requireFeeIsCorrect(
	require *require.Assertions,
	feeCalculator fee.Calculator,
	utx txs.UnsignedTx,
	baseTx *avax.BaseTx,
	additionalIns []*avax.TransferableInput,
	additionalOuts []*avax.TransferableOutput,
	additionalFee map[ids.ID]uint64,
) {
	amountConsumed := addInputAmounts(baseTx.Ins, additionalIns)
	amountProduced := addOutputAmounts(baseTx.Outs, additionalOuts)

	expectedFee, err := feeCalculator.CalculateFee(utx)
	require.NoError(err)
	expectedAmountBurned := addAmounts(
		map[ids.ID]uint64{
			avaxAssetID: expectedFee,
		},
		additionalFee,
	)
	expectedAmountConsumed := addAmounts(amountProduced, expectedAmountBurned)
	require.Equal(expectedAmountConsumed, amountConsumed)
}

func addAmounts(allAmounts ...map[ids.ID]uint64) map[ids.ID]uint64 {
	amounts := make(map[ids.ID]uint64)
	for _, amountsToAdd := range allAmounts {
		for assetID, amount := range amountsToAdd {
			amounts[assetID] += amount
		}
	}
	return amounts
}

func addInputAmounts(inputSlices ...[]*avax.TransferableInput) map[ids.ID]uint64 {
	consumed := make(map[ids.ID]uint64)
	for _, inputs := range inputSlices {
		for _, in := range inputs {
			consumed[in.AssetID()] += in.In.Amount()
		}
	}
	return consumed
}

func addOutputAmounts(outputSlices ...[]*avax.TransferableOutput) map[ids.ID]uint64 {
	produced := make(map[ids.ID]uint64)
	for _, outputs := range outputSlices {
		for _, out := range outputs {
			produced[out.AssetID()] += out.Out.Amount()
		}
	}
	return produced
}
