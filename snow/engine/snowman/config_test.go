// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package snowman

import (
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/snow/consensus/snowball"
	"github.com/MetalBlockchain/metalgo/snow/consensus/snowman"
	"github.com/MetalBlockchain/metalgo/snow/engine/common"
	"github.com/MetalBlockchain/metalgo/snow/engine/common/tracker"
	"github.com/MetalBlockchain/metalgo/snow/engine/snowman/block"
	"github.com/MetalBlockchain/metalgo/snow/validators"
)

func DefaultConfig() Config {
	return Config{
		Ctx:                 snow.DefaultConsensusContextTest(),
		VM:                  &block.TestVM{},
		Sender:              &common.SenderTest{},
		Validators:          validators.NewManager(),
		ConnectedValidators: tracker.NewPeers(),
		Params: snowball.Parameters{
			K:                     1,
			AlphaPreference:       1,
			AlphaConfidence:       1,
			BetaVirtuous:          1,
			BetaRogue:             2,
			ConcurrentRepolls:     1,
			OptimalProcessing:     100,
			MaxOutstandingItems:   1,
			MaxItemProcessingTime: 1,
		},
		Consensus: &snowman.Topological{},
	}
}
