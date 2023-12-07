// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package common

import (
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/snow/engine/common/tracker"
	"github.com/MetalBlockchain/metalgo/snow/validators"
	"github.com/MetalBlockchain/metalgo/utils/constants"
)

// DefaultConfigTest returns a test configuration
func DefaultConfigTest() Config {
	isBootstrapped := false
	bootstrapTracker := &BootstrapTrackerTest{
		IsBootstrappedF: func() bool {
			return isBootstrapped
		},
		BootstrappedF: func(ids.ID) {
			isBootstrapped = true
		},
	}

	beacons := validators.NewManager()

	connectedPeers := tracker.NewPeers()
	startupTracker := tracker.NewStartup(connectedPeers, 0)
	beacons.RegisterCallbackListener(constants.PrimaryNetworkID, startupTracker)

	return Config{
		Ctx:                            snow.DefaultConsensusContextTest(),
		Beacons:                        beacons,
		StartupTracker:                 startupTracker,
		Sender:                         &SenderTest{},
		Bootstrapable:                  &BootstrapableTest{},
		BootstrapTracker:               bootstrapTracker,
		Timer:                          &TimerTest{},
		AncestorsMaxContainersSent:     2000,
		AncestorsMaxContainersReceived: 2000,
		SharedCfg:                      &SharedConfig{},
	}
}
