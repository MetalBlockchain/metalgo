// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"time"

	_ "embed"

	"github.com/MetalBlockchain/metalgo/utils/units"
	"github.com/MetalBlockchain/metalgo/vms/components/gas"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/reward"

	txfee "github.com/MetalBlockchain/metalgo/vms/platformvm/txs/fee"
	validatorfee "github.com/MetalBlockchain/metalgo/vms/platformvm/validators/fee"
)

var (
	//go:embed genesis_mainnet.json
	mainnetGenesisConfigJSON []byte

	// MainnetParams are the params used for mainnet
	MainnetParams = Params{
		TxFeeConfig: TxFeeConfig{
			CreateAssetTxFee: 10 * units.MilliAvax,
			StaticFeeConfig: txfee.StaticConfig{
				TxFee:                         units.MilliAvax,
				CreateSubnetTxFee:             1 * units.Avax,
				TransformSubnetTxFee:          10 * units.Avax,
				CreateBlockchainTxFee:         1 * units.Avax,
				AddPrimaryNetworkValidatorFee: 0,
				AddPrimaryNetworkDelegatorFee: 0,
				AddSubnetValidatorFee:         units.MilliAvax,
				AddSubnetDelegatorFee:         units.MilliAvax,
			},
			DynamicFeeConfig: gas.Config{
				Weights: gas.Dimensions{
					gas.Bandwidth: 1,     // Max block size ~1MB
					gas.DBRead:    1_000, // Max reads per block 1,000
					gas.DBWrite:   1_000, // Max writes per block 1,000
					gas.Compute:   4,     // Max compute time per block ~250ms
				},
				MaxCapacity:     1_000_000,
				MaxPerSecond:    100_000, // Refill time 10s
				TargetPerSecond: 50_000,  // Target is half of max
				MinPrice:        1,
				// ExcessConversionConstant = (MaxPerSecond - TargetPerSecond) * NumberOfSecondsPerDoubling / ln(2)
				//
				// ln(2) is a float and the result is consensus critical, so we
				// hardcode the result.
				ExcessConversionConstant: 2_164_043, // Double every 30s
			},
			ValidatorFeeConfig: validatorfee.Config{
				Capacity: 20_000,
				Target:   10_000,
				MinPrice: gas.Price(512 * units.NanoAvax),
				// ExcessConversionConstant = (Capacity - Target) * NumberOfSecondsPerDoubling / ln(2)
				//
				// ln(2) is a float and the result is consensus critical, so we
				// hardcode the result.
				ExcessConversionConstant: 1_246_488_515, // Double every day
			},
		},
		StakingConfig: StakingConfig{
			UptimeRequirement: .8, // 80%
			MinValidatorStake: 2 * units.KiloAvax,
			MaxValidatorStake: 3 * units.MegaAvax,
			MinDelegatorStake: 25 * units.Avax,
			MinDelegationFee:  20000, // 2%
			MinStakeDuration:  2 * 7 * 24 * time.Hour,
			MaxStakeDuration:  365 * 24 * time.Hour,
			RewardConfig: reward.Config{
				MaxConsumptionRate: .12 * reward.PercentDenominator,
				MinConsumptionRate: .10 * reward.PercentDenominator,
				MintingPeriod:      365 * 24 * time.Hour,
				SupplyCap:          666666666 * units.Avax,
			},
		},
	}
)
