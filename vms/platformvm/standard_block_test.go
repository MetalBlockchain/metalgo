// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package platformvm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/MetalBlockchain/metalgo/chains/atomic"
	"github.com/MetalBlockchain/metalgo/database/prefixdb"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils/crypto"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/status"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
)

func TestAtomicTxImports(t *testing.T) {
	assert := assert.New(t)
	vm, baseDB, _, mutableSharedMemory := defaultVM()
	vm.ctx.Lock.Lock()
	defer func() {
		assert.NoError(vm.Shutdown())
		vm.ctx.Lock.Unlock()
	}()

	utxoID := avax.UTXOID{
		TxID:        ids.Empty.Prefix(1),
		OutputIndex: 1,
	}
	amount := uint64(70000)
	recipientKey := keys[1]

	m := atomic.NewMemory(prefixdb.New([]byte{5}, baseDB))

	mutableSharedMemory.SharedMemory = m.NewSharedMemory(vm.ctx.ChainID)
	peerSharedMemory := m.NewSharedMemory(vm.ctx.XChainID)
	utxo := &avax.UTXO{
		UTXOID: utxoID,
		Asset:  avax.Asset{ID: avaxAssetID},
		Out: &secp256k1fx.TransferOutput{
			Amt: amount,
			OutputOwners: secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs:     []ids.ShortID{recipientKey.PublicKey().Address()},
			},
		},
	}
	utxoBytes, err := txs.Codec.Marshal(txs.Version, utxo)
	assert.NoError(err)
	inputID := utxo.InputID()
	err = peerSharedMemory.Apply(map[ids.ID]*atomic.Requests{vm.ctx.ChainID: {PutRequests: []*atomic.Element{{
		Key:   inputID[:],
		Value: utxoBytes,
		Traits: [][]byte{
			recipientKey.PublicKey().Address().Bytes(),
		},
	}}}})
	assert.NoError(err)

	tx, err := vm.txBuilder.NewImportTx(
		vm.ctx.XChainID,
		recipientKey.PublicKey().Address(),
		[]*crypto.PrivateKeySECP256K1R{recipientKey},
		ids.ShortEmpty, // change addr
	)
	assert.NoError(err)

	vm.state.SetTimestamp(vm.ApricotPhase5Time.Add(100 * time.Second))

	vm.mempool.AddDecisionTx(tx)
	b, err := vm.BuildBlock()
	assert.NoError(err)
	// Test multiple verify calls work
	assert.NoError(b.Verify())
	assert.NoError(b.Accept())
	_, txStatus, err := vm.state.GetTx(tx.ID())
	assert.NoError(err)
	// Ensure transaction is in the committed state
	assert.Equal(txStatus, status.Committed)
}
