// Copyright (C) 2019, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package statesync

import (
	"fmt"

	"github.com/MetalBlockchain/libevm/core/state/snapshot"
	"github.com/MetalBlockchain/libevm/ethdb"
	"github.com/MetalBlockchain/libevm/triedb"

	"github.com/MetalBlockchain/metalgo/network/p2p"
	"github.com/MetalBlockchain/metalgo/utils/logging"
	"github.com/MetalBlockchain/metalgo/vms/saevm/statesync"

	cchainstate "github.com/MetalBlockchain/metalgo/vms/saevm/cchain/state"
)

// RegisterHandlers registers the SAE state sync handler with the given EVM trie
// database, allowing this node to serve others' state sync requests.
func RegisterHandlers(
	log logging.Logger,
	network *p2p.Network,
	db ethdb.Database,
	tdb *triedb.Database,
	snaps *snapshot.Tree,
	state *cchainstate.State,
) error {
	if err := cchainstate.RegisterSyncHandler(network, state); err != nil {
		return fmt.Errorf("registering C-Chain state handler: %w", err)
	}

	return statesync.RegisterHandlers(log, network, db, tdb, snaps)
}
