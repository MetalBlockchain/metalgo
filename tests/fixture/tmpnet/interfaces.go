// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package tmpnet

import (
	"context"
	"io"

	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/node"
)

// Defines network capabilities supportable regardless of how a network is orchestrated.
type Network interface {
	GetConfig() NetworkConfig
	GetNodes() []Node
	AddEphemeralNode(w io.Writer, flags FlagsMap) (Node, error)
}

// Defines node capabilities supportable regardless of how a network is orchestrated.
type Node interface {
	GetID() ids.NodeID
	GetConfig() NodeConfig
	GetProcessContext() node.NodeProcessContext
	IsHealthy(ctx context.Context) (bool, error)
	Stop() error
}
