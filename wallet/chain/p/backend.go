// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package p

import (
	"context"
	"sync"

	"github.com/MetalBlockchain/metalgo/database"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/set"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/fx"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs"
	"github.com/MetalBlockchain/metalgo/wallet/chain/p/builder"
	"github.com/MetalBlockchain/metalgo/wallet/chain/p/signer"
	"github.com/MetalBlockchain/metalgo/wallet/subnet/primary/common"
)

var _ Backend = (*backend)(nil)

// Backend defines the full interface required to support a P-chain wallet.
type Backend interface {
	builder.Backend
	signer.Backend

	AcceptTx(ctx context.Context, tx *txs.Tx) error
}

type backend struct {
	common.ChainUTXOs

	context *builder.Context

	subnetOwnerLock sync.RWMutex
	subnetOwner     map[ids.ID]fx.Owner // subnetID -> owner
}

func NewBackend(context *builder.Context, utxos common.ChainUTXOs, subnetTxs map[ids.ID]*txs.Tx) Backend {
	subnetOwner := make(map[ids.ID]fx.Owner)
	for txID, tx := range subnetTxs { // first get owners from the CreateSubnetTx
		createSubnetTx, ok := tx.Unsigned.(*txs.CreateSubnetTx)
		if !ok {
			continue
		}
		subnetOwner[txID] = createSubnetTx.Owner
	}
	for _, tx := range subnetTxs { // then check for TransferSubnetOwnershipTx
		transferSubnetOwnershipTx, ok := tx.Unsigned.(*txs.TransferSubnetOwnershipTx)
		if !ok {
			continue
		}
		subnetOwner[transferSubnetOwnershipTx.Subnet] = transferSubnetOwnershipTx.Owner
	}
	return &backend{
		ChainUTXOs:  utxos,
		context:     context,
		subnetOwner: subnetOwner,
	}
}

func (b *backend) AcceptTx(ctx context.Context, tx *txs.Tx) error {
	txID := tx.ID()
	err := tx.Unsigned.Visit(&backendVisitor{
		b:    b,
		ctx:  ctx,
		txID: txID,
	})
	if err != nil {
		return err
	}

	producedUTXOSlice := tx.UTXOs()
	return b.addUTXOs(ctx, constants.PlatformChainID, producedUTXOSlice)
}

func (b *backend) addUTXOs(ctx context.Context, destinationChainID ids.ID, utxos []*avax.UTXO) error {
	for _, utxo := range utxos {
		if err := b.AddUTXO(ctx, destinationChainID, utxo); err != nil {
			return err
		}
	}
	return nil
}

func (b *backend) removeUTXOs(ctx context.Context, sourceChain ids.ID, utxoIDs set.Set[ids.ID]) error {
	for utxoID := range utxoIDs {
		if err := b.RemoveUTXO(ctx, sourceChain, utxoID); err != nil {
			return err
		}
	}
	return nil
}

func (b *backend) GetSubnetOwner(_ context.Context, subnetID ids.ID) (fx.Owner, error) {
	b.subnetOwnerLock.RLock()
	defer b.subnetOwnerLock.RUnlock()

	owner, exists := b.subnetOwner[subnetID]
	if !exists {
		return nil, database.ErrNotFound
	}
	return owner, nil
}

func (b *backend) setSubnetOwner(subnetID ids.ID, owner fx.Owner) {
	b.subnetOwnerLock.Lock()
	defer b.subnetOwnerLock.Unlock()

	b.subnetOwner[subnetID] = owner
}
