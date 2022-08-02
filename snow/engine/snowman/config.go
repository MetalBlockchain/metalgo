// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package snowman

import (
	"github.com/MetalBlockchain/avalanchego/snow"
	"github.com/MetalBlockchain/avalanchego/snow/consensus/snowball"
	"github.com/MetalBlockchain/avalanchego/snow/consensus/snowman"
	"github.com/MetalBlockchain/avalanchego/snow/engine/common"
	"github.com/MetalBlockchain/avalanchego/snow/engine/snowman/block"
	"github.com/MetalBlockchain/avalanchego/snow/validators"
)

// Config wraps all the parameters needed for a snowman engine
type Config struct {
	common.AllGetsServer

	Ctx        *snow.ConsensusContext
	VM         block.ChainVM
	Sender     common.Sender
	Validators validators.Set
	Params     snowball.Parameters
	Consensus  snowman.Consensus
}
