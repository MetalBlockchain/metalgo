// Copyright (C) 2019, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package types

// This single function is in a standalone file to reduce confusion because
// every required import has something to do with a database!

import (
	"github.com/MetalBlockchain/libevm/core/rawdb"
	"github.com/MetalBlockchain/libevm/ethdb"

	"github.com/MetalBlockchain/metalgo/database"

	evmdb "github.com/MetalBlockchain/metalgo/vms/evm/database"
)

func NewEthDB(db database.Database) ethdb.Database {
	return rawdb.NewDatabase(evmdb.New(db))
}
