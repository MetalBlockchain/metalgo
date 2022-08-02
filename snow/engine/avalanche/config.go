// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avalanche

import (
	"github.com/MetalBlockchain/avalanchego/snow"
	"github.com/MetalBlockchain/avalanchego/snow/consensus/avalanche"
	"github.com/MetalBlockchain/avalanchego/snow/engine/avalanche/vertex"
	"github.com/MetalBlockchain/avalanchego/snow/engine/common"
	"github.com/MetalBlockchain/avalanchego/snow/validators"
)

// Config wraps all the parameters needed for an avalanche engine
type Config struct {
	Ctx *snow.ConsensusContext
	common.AllGetsServer
	VM         vertex.DAGVM
	Manager    vertex.Manager
	Sender     common.Sender
	Validators validators.Set

	Params    avalanche.Parameters
	Consensus avalanche.Consensus
}
