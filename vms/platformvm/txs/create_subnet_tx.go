// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/fx"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
)

var (
	_ UnsignedTx             = &CreateSubnetTx{}
	_ secp256k1fx.UnsignedTx = &CreateSubnetTx{}
)

// CreateSubnetTx is an unsigned proposal to create a new subnet
type CreateSubnetTx struct {
	// Metadata, inputs and outputs
	BaseTx `serialize:"true"`
	// Who is authorized to manage this subnet
	Owner fx.Owner `serialize:"true" json:"owner"`
}

// InitCtx sets the FxID fields in the inputs and outputs of this
// [CreateSubnetTx]. Also sets the [ctx] to the given [vm.ctx] so that
// the addresses can be json marshalled into human readable format
func (tx *CreateSubnetTx) InitCtx(ctx *snow.Context) {
	tx.BaseTx.InitCtx(ctx)
	tx.Owner.InitCtx(ctx)
}

// SyntacticVerify verifies that this transaction is well-formed
func (tx *CreateSubnetTx) SyntacticVerify(ctx *snow.Context) error {
	switch {
	case tx == nil:
		return ErrNilTx
	case tx.SyntacticallyVerified: // already passed syntactic verification
		return nil
	}

	if err := tx.BaseTx.SyntacticVerify(ctx); err != nil {
		return err
	}
	if err := tx.Owner.Verify(); err != nil {
		return err
	}

	tx.SyntacticallyVerified = true
	return nil
}

func (tx *CreateSubnetTx) Visit(visitor Visitor) error {
	return visitor.CreateSubnetTx(tx)
}
