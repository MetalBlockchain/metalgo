// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package cache

import "github.com/MetalBlockchain/metalgo/utils"

var _ Cacher[struct{}, struct{}] = (*Empty[struct{}, struct{}])(nil)

type Empty[K any, V any] struct{}

func (*Empty[K, V]) Put(K, V) {}

func (*Empty[K, V]) Get(K) (V, bool) {
	return utils.Zero[V](), false
}

func (*Empty[K, _]) Evict(K) {}

func (*Empty[_, _]) Flush() {}

func (*Empty[_, _]) Len() int {
	return 0
}

func (*Empty[_, _]) PortionFilled() float64 {
	return 0
}
