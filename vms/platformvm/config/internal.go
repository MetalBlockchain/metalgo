// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package config

import (
	"time"

	"github.com/MetalBlockchain/metalgo/chains"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow/uptime"
	"github.com/MetalBlockchain/metalgo/snow/validators"
	"github.com/MetalBlockchain/metalgo/upgrade"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/set"
	"github.com/MetalBlockchain/metalgo/vms/components/gas"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/reward"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs"

	txfee "github.com/MetalBlockchain/metalgo/vms/platformvm/txs/fee"
	validatorfee "github.com/MetalBlockchain/metalgo/vms/platformvm/validators/fee"
)

// Internal contains all of the parameters for the PlatformVM that are
// internally set by the node.
type Internal struct {
	// The node's chain manager
	Chains chains.Manager

	// Node's validator set maps subnetID -> validators of the subnet
	//
	// Invariant: The primary network's validator set should have been added to
	//            the manager before calling VM.Initialize.
	// Invariant: The primary network's validator set should be empty before
	//            calling VM.Initialize.
	Validators validators.Manager

	// Static fees are active before Etna
	CreateAssetTxFee uint64 // Override for CreateSubnet and CreateChain before AP3
	StaticFeeConfig  txfee.StaticConfig

	// Dynamic fees are active after Etna
	DynamicFeeConfig gas.Config

	// ACP-77 validator fees are active after Etna
	ValidatorFeeConfig validatorfee.Config

	// Provides access to the uptime manager as a thread safe data structure
	UptimeLockedCalculator uptime.LockedCalculator

	// True if the node is being run with staking enabled
	SybilProtectionEnabled bool

	// If true, only the P-chain will be instantiated on the primary network.
	PartialSyncPrimaryNetwork bool

	// Set of subnets that this node is validating
	TrackedSubnets set.Set[ids.ID]

	// The minimum amount of tokens one must bond to be a validator
	MinValidatorStake uint64

	// The maximum amount of tokens that can be bonded on a validator
	MaxValidatorStake uint64

	// Minimum stake, in nAVAX, that can be delegated on the primary network
	MinDelegatorStake uint64

	// Minimum fee that can be charged for delegation
	MinDelegationFee uint32

	// UptimePercentage is the minimum uptime required to be rewarded for staking
	UptimePercentage float64

	// Minimum amount of time to allow a staker to stake
	MinStakeDuration time.Duration

	// Maximum amount of time to allow a staker to stake
	MaxStakeDuration time.Duration

	// Config for the minting function
	RewardConfig reward.Config

	// All network upgrade timestamps
	UpgradeConfig upgrade.Config

	// UseCurrentHeight forces [GetMinimumHeight] to return the current height
	// of the P-Chain instead of the oldest block in the [recentlyAccepted]
	// window.
	//
	// This config is particularly useful for triggering proposervm activation
	// on recently created subnets (without this, users need to wait for
	// [recentlyAcceptedWindowTTL] to pass for activation to occur).
	UseCurrentHeight bool
}

// Create the blockchain described in [tx], but only if this node is a member of
// the subnet that validates the chain
func (c *Internal) CreateChain(chainID ids.ID, tx *txs.CreateChainTx) {
	if c.SybilProtectionEnabled && // Sybil protection is enabled, so nodes might not validate all chains
		constants.PrimaryNetworkID != tx.SubnetID && // All nodes must validate the primary network
		!c.TrackedSubnets.Contains(tx.SubnetID) { // This node doesn't validate this blockchain
		return
	}

	chainParams := chains.ChainParameters{
		ID:          chainID,
		SubnetID:    tx.SubnetID,
		GenesisData: tx.GenesisData,
		VMID:        tx.VMID,
		FxIDs:       tx.FxIDs,
	}

	c.Chains.QueueChainCreation(chainParams)
}
