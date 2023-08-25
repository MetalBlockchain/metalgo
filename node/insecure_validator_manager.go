// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package node

import (
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow/networking/router"
	"github.com/MetalBlockchain/metalgo/snow/validators"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/version"
)

type insecureValidatorManager struct {
	router.Router
	vdrs   validators.Set
	weight uint64
}

func (i *insecureValidatorManager) Connected(vdrID ids.NodeID, nodeVersion *version.Application, subnetID ids.ID) {
	if constants.PrimaryNetworkID == subnetID {
		// Sybil protection is disabled so we don't have a txID that added the
		// peer as a validator. Because each validator needs a txID associated
		// with it, we hack one together by padding the nodeID with zeroes.
		dummyTxID := ids.Empty
		copy(dummyTxID[:], vdrID[:])

		// Add will only error here if the total weight of the set would go over
		// [math.MaxUint64]. In this case, we will just not mark this new peer
		// as a validator.
		_ = i.vdrs.Add(vdrID, nil, dummyTxID, i.weight)
	}
	i.Router.Connected(vdrID, nodeVersion, subnetID)
}

func (i *insecureValidatorManager) Disconnected(vdrID ids.NodeID) {
	// RemoveWeight will only error here if there was an error reported during
	// Add.
	_ = i.vdrs.RemoveWeight(vdrID, i.weight)
	i.Router.Disconnected(vdrID)
}
