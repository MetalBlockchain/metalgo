// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package propertyfx

import (
	"github.com/MetalBlockchain/metalgo/vms/components/verify"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
)

var _ verify.State = (*OwnedOutput)(nil)

type OwnedOutput struct {
	verify.IsState `json:"-"`

	secp256k1fx.OutputOwners `serialize:"true"`
}
