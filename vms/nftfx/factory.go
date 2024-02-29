// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package nftfx

import (
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils/logging"
	"github.com/MetalBlockchain/metalgo/vms"
)

var (
	_ vms.Factory = (*Factory)(nil)

	// ID that this Fx uses when labeled
	ID = ids.ID{'n', 'f', 't', 'f', 'x'}
)

type Factory struct{}

func (*Factory) New(logging.Logger) (interface{}, error) {
	return &Fx{}, nil
}
