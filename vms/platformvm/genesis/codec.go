// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"math"

	"github.com/MetalBlockchain/metalgo/codec"
	"github.com/MetalBlockchain/metalgo/codec/linearcodec"
	"github.com/MetalBlockchain/metalgo/utils/wrappers"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs"
)

var Codec codec.Manager

func init() {
	gc := linearcodec.NewCustomMaxLength(math.MaxInt32)
	Codec = codec.NewManager(math.MaxInt32)

	// To maintain codec type ordering, skip positions
	// for Proposal/Abort/Commit/Standard/Atomic blocks
	gc.SkipRegistrations(5)

	errs := wrappers.Errs{}
	errs.Add(
		txs.RegisterUnsignedTxsTypes(gc),
		Codec.RegisterCodec(txs.Version, gc),
	)
	if errs.Errored() {
		panic(errs.Err)
	}
}
