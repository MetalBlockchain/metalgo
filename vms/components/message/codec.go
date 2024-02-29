// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package message

import (
	"time"

	"github.com/MetalBlockchain/metalgo/codec"
	"github.com/MetalBlockchain/metalgo/codec/linearcodec"
	"github.com/MetalBlockchain/metalgo/utils"
	"github.com/MetalBlockchain/metalgo/utils/units"
)

const (
	CodecVersion = 0

	maxMessageSize = 512 * units.KiB
)

var Codec codec.Manager

func init() {
	Codec = codec.NewManager(maxMessageSize)
	lc := linearcodec.NewDefault(time.Time{})

	err := utils.Err(
		lc.RegisterType(&Tx{}),
		Codec.RegisterCodec(CodecVersion, lc),
	)
	if err != nil {
		panic(err)
	}
}
