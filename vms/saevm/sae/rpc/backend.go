// Copyright (C) 2019, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rpc

import (
	"errors"

	"github.com/MetalBlockchain/libevm/accounts"
	"github.com/MetalBlockchain/libevm/eth/filters"
	"github.com/MetalBlockchain/libevm/eth/tracers"
	"github.com/MetalBlockchain/libevm/libevm/ethapi"

	"github.com/MetalBlockchain/metalgo/vms/saevm/gasprice"
	"github.com/MetalBlockchain/metalgo/vms/saevm/txgossip"
)

// The GethBackends interface is the union of geth interfaces required by
// their implementations of JSON-RPC namespace handlers.
type GethBackends interface {
	ethapi.Backend
	tracers.Backend
	filters.BloomOverrider
	filters.Backend
}

// GethBackends returns the [GethBackends] that back all JSON-RPC namespace
// handlers registered by [Provider.Server].
func (p *Provider) GethBackends() GethBackends {
	return p.backend
}

var _ GethBackends = (*backend)(nil)

type backend struct {
	Chain
	config         Config
	accountManager *accounts.Manager

	*gasprice.Estimator
	*txgossip.Set
	chainIndexer
	bloomOverrider
	*bloomIndexer
}

func (b *backend) close() error {
	return errors.Join(
		b.accountManager.Close(),
		b.Estimator.Close(),
		b.bloomIndexer.Close(),
	)
}
