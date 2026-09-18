// Copyright (C) 2019, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package extstate

import (
	"github.com/MetalBlockchain/libevm/core/state"
	"github.com/MetalBlockchain/libevm/ethdb"
	"github.com/MetalBlockchain/libevm/triedb"

	"github.com/MetalBlockchain/metalgo/graft/evm/firewood"
)

func NewDatabaseWithConfig(db ethdb.Database, config *triedb.Config) state.Database {
	coredb := state.NewDatabaseWithConfig(db, config)
	// NewStateAccessor wraps coredb only when its trie database backend is
	// Firewood; for any other backend it returns coredb unchanged.
	return firewood.NewStateAccessor(coredb)
}

func NewDatabaseWithNodeDB(db ethdb.Database, triedb *triedb.Database) state.Database {
	coredb := state.NewDatabaseWithNodeDB(db, triedb)
	// NewStateAccessor wraps coredb only when its trie database backend is
	// Firewood; for any other backend it returns coredb unchanged.
	return firewood.NewStateAccessor(coredb)
}
