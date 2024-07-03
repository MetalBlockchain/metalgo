// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txstest

import (
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/utils/crypto/secp256k1"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/config"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/state"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
	"github.com/MetalBlockchain/metalgo/wallet/chain/p/builder"
	"github.com/MetalBlockchain/metalgo/wallet/chain/p/signer"
)

func NewWalletFactory(
	ctx *snow.Context,
	cfg *config.Config,
	state state.State,
) *WalletFactory {
	return &WalletFactory{
		ctx:   ctx,
		cfg:   cfg,
		state: state,
	}
}

type WalletFactory struct {
	ctx   *snow.Context
	cfg   *config.Config
	state state.State
}

func (w *WalletFactory) NewWallet(keys ...*secp256k1.PrivateKey) (builder.Builder, signer.Signer) {
	var (
		kc      = secp256k1fx.NewKeychain(keys...)
		addrs   = kc.Addresses()
		backend = newBackend(addrs, w.state, w.ctx.SharedMemory)
		context = newContext(w.ctx, w.cfg, w.state.GetTimestamp())
	)

	return builder.New(addrs, context, backend), signer.New(kc, backend)
}
