// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"time"

	_ "embed"

	"github.com/MetalBlockchain/metalgo/utils/units"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/reward"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs/fee"
)

var (
	//go:embed genesis_tahoe.json
	tahoeGenesisConfigJSON []byte

	// TahoeParams are the params used for the tahoe testnet
	TahoeParams = Params{
		TxFeeConfig: TxFeeConfig{
			CreateAssetTxFee: 10 * units.MilliAvax,
			StaticFeeConfig: fee.StaticConfig{
				TxFee:                         units.MilliAvax,
				CreateSubnetTxFee:             100 * units.MilliAvax,
				TransformSubnetTxFee:          1 * units.Avax,
				CreateBlockchainTxFee:         100 * units.MilliAvax,
				AddPrimaryNetworkValidatorFee: 0,
				AddPrimaryNetworkDelegatorFee: 0,
				AddSubnetValidatorFee:         units.MilliAvax,
				AddSubnetDelegatorFee:         units.MilliAvax,
			},
		},
		StakingConfig: StakingConfig{
			UptimeRequirement: .8, // 80%
			MinValidatorStake: 1 * units.Avax,
			MaxValidatorStake: 3 * units.MegaAvax,
			MinDelegatorStake: 1 * units.Avax,
			MinDelegationFee:  20000, // 2%
			MinStakeDuration:  24 * time.Hour,
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
