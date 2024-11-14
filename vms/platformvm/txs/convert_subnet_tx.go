// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"errors"

	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/vms/components/verify"
	"github.com/MetalBlockchain/metalgo/vms/types"
)

const MaxSubnetAddressLength = 4096

var (
	_ UnsignedTx = (*TransferSubnetOwnershipTx)(nil)

	ErrConvertPermissionlessSubnet = errors.New("cannot convert a permissionless subnet")
	ErrAddressTooLong              = errors.New("address is too long")
)

type ConvertSubnetTx struct {
	// Metadata, inputs and outputs
	BaseTx `serialize:"true"`
	// ID of the Subnet to transform
	Subnet ids.ID `serialize:"true" json:"subnetID"`
	// Chain where the Subnet manager lives
	ChainID ids.ID `serialize:"true" json:"chainID"`
	// Address of the Subnet manager
	Address types.JSONByteSlice `serialize:"true" json:"address"`
	// Authorizes this conversion
	SubnetAuth verify.Verifiable `serialize:"true" json:"subnetAuthorization"`
}

func (tx *ConvertSubnetTx) SyntacticVerify(ctx *snow.Context) error {
	switch {
	case tx == nil:
		return ErrNilTx
	case tx.SyntacticallyVerified:
		// already passed syntactic verification
		return nil
	case tx.Subnet == constants.PrimaryNetworkID:
		return ErrConvertPermissionlessSubnet
	case len(tx.Address) > MaxSubnetAddressLength:
		return ErrAddressTooLong
	}

	if err := tx.BaseTx.SyntacticVerify(ctx); err != nil {
		return err
	}
	if err := tx.SubnetAuth.Verify(); err != nil {
		return err
	}

	tx.SyntacticallyVerified = true
	return nil
}

func (tx *ConvertSubnetTx) Visit(visitor Visitor) error {
	return visitor.ConvertSubnetTx(tx)
}
