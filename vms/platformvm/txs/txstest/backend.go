// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txstest

import (
	"context"
	"math"

	"github.com/MetalBlockchain/metalgo/chains/atomic"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/set"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/fx"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/state"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs"
	"github.com/MetalBlockchain/metalgo/wallet/chain/p/builder"
	"github.com/MetalBlockchain/metalgo/wallet/chain/p/signer"
)

var (
	_ builder.Backend = (*Backend)(nil)
	_ signer.Backend  = (*Backend)(nil)
)

func newBackend(
	addrs set.Set[ids.ShortID],
	state state.State,
	sharedMemory atomic.SharedMemory,
) *Backend {
	return &Backend{
		addrs:        addrs,
		state:        state,
		sharedMemory: sharedMemory,
	}
}

type Backend struct {
	addrs        set.Set[ids.ShortID]
	state        state.State
	sharedMemory atomic.SharedMemory
}

func (b *Backend) UTXOs(_ context.Context, sourceChainID ids.ID) ([]*avax.UTXO, error) {
	if sourceChainID == constants.PlatformChainID {
		return avax.GetAllUTXOs(b.state, b.addrs)
	}

	utxos, _, _, err := avax.GetAtomicUTXOs(
		b.sharedMemory,
		txs.Codec,
		sourceChainID,
		b.addrs,
		ids.ShortEmpty,
		ids.Empty,
		math.MaxInt,
	)
	return utxos, err
}

func (b *Backend) GetUTXO(_ context.Context, chainID, utxoID ids.ID) (*avax.UTXO, error) {
	if chainID == constants.PlatformChainID {
		return b.state.GetUTXO(utxoID)
	}

	utxoBytes, err := b.sharedMemory.Get(chainID, [][]byte{utxoID[:]})
	if err != nil {
		return nil, err
	}

	utxo := avax.UTXO{}
	if _, err := txs.Codec.Unmarshal(utxoBytes[0], &utxo); err != nil {
		return nil, err
	}
	return &utxo, nil
}

func (b *Backend) GetSubnetOwner(_ context.Context, subnetID ids.ID) (fx.Owner, error) {
	return b.state.GetSubnetOwner(subnetID)
}
