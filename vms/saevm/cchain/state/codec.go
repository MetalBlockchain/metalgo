// Copyright (C) 2019, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"math"

	"github.com/MetalBlockchain/metalgo/codec"
	"github.com/MetalBlockchain/metalgo/codec/linearcodec"
)

const codecVersion uint16 = 0

var c codec.Manager

func init() {
	c = codec.NewManager(math.MaxInt)
	lc := linearcodec.NewDefault()
	if err := c.RegisterCodec(codecVersion, lc); err != nil {
		panic(err)
	}
}
