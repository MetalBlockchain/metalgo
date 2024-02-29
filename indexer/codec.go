// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package indexer

import (
	"math"
	"time"

	"github.com/MetalBlockchain/metalgo/codec"
	"github.com/MetalBlockchain/metalgo/codec/linearcodec"
)

const CodecVersion = 0

var Codec codec.Manager

func init() {
	lc := linearcodec.NewDefault(time.Time{})
	Codec = codec.NewManager(math.MaxInt)

	if err := Codec.RegisterCodec(CodecVersion, lc); err != nil {
		panic(err)
	}
}
