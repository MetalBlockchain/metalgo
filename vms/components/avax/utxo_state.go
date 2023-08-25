// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avax

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/MetalBlockchain/metalgo/cache"
	"github.com/MetalBlockchain/metalgo/cache/metercacher"
	"github.com/MetalBlockchain/metalgo/codec"
	"github.com/MetalBlockchain/metalgo/database"
	"github.com/MetalBlockchain/metalgo/database/linkeddb"
	"github.com/MetalBlockchain/metalgo/database/prefixdb"
	"github.com/MetalBlockchain/metalgo/ids"
)

const (
	utxoCacheSize  = 8192
	indexCacheSize = 64
)

var (
	utxoPrefix  = []byte("utxo")
	indexPrefix = []byte("index")
)

// UTXOState is a thin wrapper around a database to provide, caching,
// serialization, and de-serialization for UTXOs.
type UTXOState interface {
	UTXOReader
	UTXOWriter
	Checksum() ids.ID
}

// UTXOReader is a thin wrapper around a database to provide fetching of UTXOs.
type UTXOReader interface {
	UTXOGetter

	// UTXOIDs returns the slice of IDs associated with [addr], starting after
	// [previous].
	// If [previous] is not in the list, starts at beginning.
	// Returns at most [limit] IDs.
	UTXOIDs(addr []byte, previous ids.ID, limit int) ([]ids.ID, error)
}

// UTXOGetter is a thin wrapper around a database to provide fetching of a UTXO.
type UTXOGetter interface {
	// GetUTXO attempts to load a utxo.
	GetUTXO(utxoID ids.ID) (*UTXO, error)
}

type UTXOAdder interface {
	AddUTXO(utxo *UTXO)
}

type UTXODeleter interface {
	DeleteUTXO(utxoID ids.ID)
}

// UTXOWriter is a thin wrapper around a database to provide storage and
// deletion of UTXOs.
type UTXOWriter interface {
	// PutUTXO saves the provided utxo to storage.
	PutUTXO(utxo *UTXO) error

	// DeleteUTXO deletes the provided utxo.
	DeleteUTXO(utxoID ids.ID) error
}

type utxoState struct {
	codec codec.Manager

	// UTXO ID -> *UTXO. If the *UTXO is nil the UTXO doesn't exist
	utxoCache cache.Cacher[ids.ID, *UTXO]
	utxoDB    database.Database

	indexDB    database.Database
	indexCache cache.Cacher[string, linkeddb.LinkedDB]

	checksum ids.ID
}

func NewUTXOState(db database.Database, codec codec.Manager) (UTXOState, error) {
	s := &utxoState{
		codec: codec,

		utxoCache: &cache.LRU[ids.ID, *UTXO]{Size: utxoCacheSize},
		utxoDB:    prefixdb.New(utxoPrefix, db),

		indexDB:    prefixdb.New(indexPrefix, db),
		indexCache: &cache.LRU[string, linkeddb.LinkedDB]{Size: indexCacheSize},
	}
	return s, s.initChecksum()
}

func NewMeteredUTXOState(db database.Database, codec codec.Manager, metrics prometheus.Registerer) (UTXOState, error) {
	utxoCache, err := metercacher.New[ids.ID, *UTXO](
		"utxo_cache",
		metrics,
		&cache.LRU[ids.ID, *UTXO]{Size: utxoCacheSize},
	)
	if err != nil {
		return nil, err
	}

	indexCache, err := metercacher.New[string, linkeddb.LinkedDB](
		"index_cache",
		metrics,
		&cache.LRU[string, linkeddb.LinkedDB]{
			Size: indexCacheSize,
		},
	)
	if err != nil {
		return nil, err
	}

	s := &utxoState{
		codec: codec,

		utxoCache: utxoCache,
		utxoDB:    prefixdb.New(utxoPrefix, db),

		indexDB:    prefixdb.New(indexPrefix, db),
		indexCache: indexCache,
	}
	return s, s.initChecksum()
}

func (s *utxoState) GetUTXO(utxoID ids.ID) (*UTXO, error) {
	if utxo, found := s.utxoCache.Get(utxoID); found {
		if utxo == nil {
			return nil, database.ErrNotFound
		}
		return utxo, nil
	}

	bytes, err := s.utxoDB.Get(utxoID[:])
	if err == database.ErrNotFound {
		s.utxoCache.Put(utxoID, nil)
		return nil, database.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// The key was in the database
	utxo := &UTXO{}
	if _, err := s.codec.Unmarshal(bytes, utxo); err != nil {
		return nil, err
	}

	s.utxoCache.Put(utxoID, utxo)
	return utxo, nil
}

func (s *utxoState) PutUTXO(utxo *UTXO) error {
	utxoBytes, err := s.codec.Marshal(codecVersion, utxo)
	if err != nil {
		return err
	}

	utxoID := utxo.InputID()
	s.utxoCache.Put(utxoID, utxo)
	s.updateChecksum(utxoID)
	if err := s.utxoDB.Put(utxoID[:], utxoBytes); err != nil {
		return err
	}

	addressable, ok := utxo.Out.(Addressable)
	if !ok {
		return nil
	}

	addresses := addressable.Addresses()
	for _, addr := range addresses {
		indexList := s.getIndexDB(addr)
		if err := indexList.Put(utxoID[:], nil); err != nil {
			return err
		}
	}
	return nil
}

func (s *utxoState) DeleteUTXO(utxoID ids.ID) error {
	utxo, err := s.GetUTXO(utxoID)
	if err == database.ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}

	s.utxoCache.Put(utxoID, nil)
	s.updateChecksum(utxoID)
	if err := s.utxoDB.Delete(utxoID[:]); err != nil {
		return err
	}

	addressable, ok := utxo.Out.(Addressable)
	if !ok {
		return nil
	}

	addresses := addressable.Addresses()
	for _, addr := range addresses {
		indexList := s.getIndexDB(addr)
		if err := indexList.Delete(utxoID[:]); err != nil {
			return err
		}
	}
	return nil
}

func (s *utxoState) UTXOIDs(addr []byte, start ids.ID, limit int) ([]ids.ID, error) {
	indexList := s.getIndexDB(addr)
	iter := indexList.NewIteratorWithStart(start[:])
	defer iter.Release()

	utxoIDs := []ids.ID(nil)
	for len(utxoIDs) < limit && iter.Next() {
		utxoID, err := ids.ToID(iter.Key())
		if err != nil {
			return nil, err
		}
		if utxoID == start {
			continue
		}

		start = ids.Empty
		utxoIDs = append(utxoIDs, utxoID)
	}
	return utxoIDs, iter.Error()
}

func (s *utxoState) Checksum() ids.ID {
	return s.checksum
}

func (s *utxoState) getIndexDB(addr []byte) linkeddb.LinkedDB {
	addrStr := string(addr)
	if indexList, exists := s.indexCache.Get(addrStr); exists {
		return indexList
	}

	indexDB := prefixdb.NewNested(addr, s.indexDB)
	indexList := linkeddb.NewDefault(indexDB)
	s.indexCache.Put(addrStr, indexList)
	return indexList
}

func (s *utxoState) initChecksum() error {
	it := s.utxoDB.NewIterator()
	defer it.Release()

	for it.Next() {
		utxoID, err := ids.ToID(it.Key())
		if err != nil {
			return err
		}
		s.updateChecksum(utxoID)
	}
	return it.Error()
}

func (s *utxoState) updateChecksum(modifiedID ids.ID) {
	for i, b := range modifiedID {
		s.checksum[i] ^= b
	}
}
