// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nftfx

import (
	"github.com/MetalBlockchain/avalanchego/ids"
	"github.com/MetalBlockchain/avalanchego/snow"
	"github.com/MetalBlockchain/avalanchego/vms"
)

var (
	_ vms.Factory = &Factory{}

	// ID that this Fx uses when labeled
	ID = ids.ID{'n', 'f', 't', 'f', 'x'}
)

type Factory struct{}

func (f *Factory) New(*snow.Context) (interface{}, error) { return &Fx{}, nil }
