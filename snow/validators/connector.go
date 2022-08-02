// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package validators

import (
	"github.com/MetalBlockchain/avalanchego/ids"
	"github.com/MetalBlockchain/avalanchego/version"
)

// Connector represents a handler that is called when a connection is marked as
// connected or disconnected
type Connector interface {
	Connected(id ids.ShortID, nodeVersion version.Application) error
	Disconnected(id ids.ShortID) error
}
