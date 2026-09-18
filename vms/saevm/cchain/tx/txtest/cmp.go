// Copyright (C) 2019, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txtest

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/saevm/cchain/tx"
	"github.com/MetalBlockchain/metalgo/vms/saevm/cmputils"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
)

// CmpOpt returns a configuration for [cmp.Diff] to compare [tx.Tx] instances.
func CmpOpt() cmp.Option {
	return cmputils.IfIn[tx.Tx](cmp.Options{
		cmpopts.IgnoreUnexported(
			avax.UTXOID{},
			secp256k1fx.OutputOwners{},
		),
		cmpopts.EquateEmpty(),
	})
}

// UTXOCmpOpt returns a configuration for [cmp.Diff] to compare [avax.UTXO]
// instances or slices thereof. Slice order is not considered.
func UTXOCmpOpt() cmp.Option {
	return cmp.Options{
		cmpopts.IgnoreUnexported(
			avax.UTXOID{},
			secp256k1fx.OutputOwners{},
		),
		cmpopts.EquateEmpty(),
		cmpopts.SortSlices(func(a, b *avax.UTXO) bool {
			return a.InputID().Compare(b.InputID()) < 0
		}),
	}
}
