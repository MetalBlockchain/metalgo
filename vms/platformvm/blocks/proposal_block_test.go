// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package blocks

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/components/verify"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/validator"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
)

func TestNewBanffProposalBlock(t *testing.T) {
	require := require.New(t)

	timestamp := time.Now().Truncate(time.Second)
	parentID := ids.GenerateTestID()
	height := uint64(1337)

	tx := &txs.Tx{
		Unsigned: &txs.AddValidatorTx{
			BaseTx: txs.BaseTx{
				BaseTx: avax.BaseTx{
					Ins:  []*avax.TransferableInput{},
					Outs: []*avax.TransferableOutput{},
				},
			},
			StakeOuts: []*avax.TransferableOutput{},
			Validator: validator.Validator{},
			RewardsOwner: &secp256k1fx.OutputOwners{
				Addrs: []ids.ShortID{},
			},
		},
		Creds: []verify.Verifiable{},
	}
	require.NoError(tx.Initialize(txs.Codec))

	blk, err := NewBanffProposalBlock(
		timestamp,
		parentID,
		height,
		tx,
	)
	require.NoError(err)

	// Make sure the block and tx are initialized
	require.NotEmpty(blk.Bytes())
	require.NotEmpty(blk.Tx.Bytes())
	require.NotEqual(ids.Empty, blk.Tx.ID())
	require.Equal(tx.Bytes(), blk.Tx.Bytes())
	require.Equal(timestamp, blk.Timestamp())
	require.Equal(parentID, blk.Parent())
	require.Equal(height, blk.Height())
}

func TestNewApricotProposalBlock(t *testing.T) {
	require := require.New(t)

	parentID := ids.GenerateTestID()
	height := uint64(1337)

	tx := &txs.Tx{
		Unsigned: &txs.AddValidatorTx{
			BaseTx: txs.BaseTx{
				BaseTx: avax.BaseTx{
					Ins:  []*avax.TransferableInput{},
					Outs: []*avax.TransferableOutput{},
				},
			},
			StakeOuts: []*avax.TransferableOutput{},
			Validator: validator.Validator{},
			RewardsOwner: &secp256k1fx.OutputOwners{
				Addrs: []ids.ShortID{},
			},
		},
		Creds: []verify.Verifiable{},
	}
	require.NoError(tx.Initialize(txs.Codec))

	blk, err := NewApricotProposalBlock(
		parentID,
		height,
		tx,
	)
	require.NoError(err)

	// Make sure the block and tx are initialized
	require.NotEmpty(blk.Bytes())
	require.NotEmpty(blk.Tx.Bytes())
	require.NotEqual(ids.Empty, blk.Tx.ID())
	require.Equal(tx.Bytes(), blk.Tx.Bytes())
	require.Equal(parentID, blk.Parent())
	require.Equal(height, blk.Height())
}
