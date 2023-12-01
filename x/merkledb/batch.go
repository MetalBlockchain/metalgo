// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package merkledb

import "github.com/MetalBlockchain/metalgo/database"

var _ database.Batch = (*batch)(nil)

type batch struct {
	database.BatchOps

	db *merkleDB
}

func (b *batch) Write() error {
	return b.db.commitBatch(b.Ops)
}

func (b *batch) Inner() database.Batch {
	return b
}
