// Copyright (C) 2019, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package block

import (
	"context"
	"time"

	"github.com/MetalBlockchain/libevm/ethdb"
	"go.uber.org/zap"

	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/network/p2p"
	"github.com/MetalBlockchain/metalgo/snow/engine/common"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/logging"
	"github.com/MetalBlockchain/metalgo/vms/evm/sync/handlers"

	syncpb "github.com/MetalBlockchain/metalgo/proto/pb/sync"
)

// RegisterHandler serves block requests at [p2p.EVMBlockRequestHandlerID] on
// net.
func RegisterHandler(log logging.Logger, net *p2p.Network, db ethdb.Reader) error {
	h := handlers.NewHandler(
		log,
		&responder{
			log: log,
			db:  db,
		},
	)
	return net.AddHandler(p2p.EVMBlockRequestHandlerID, h)
}

var _ handlers.Responder[*syncpb.GetBlockRequest, *syncpb.GetBlockResponse] = (*responder)(nil)

// responder serves the requested block and its accepted ancestors.
type responder struct {
	log logging.Logger
	db  ethdb.Reader
}

const (
	// maxParentsPerRequest caps the parent walk. Blocks vary in size, so
	// [maxResponseBytes] bounds the response itself.
	maxParentsPerRequest = 64

	// maxResponseBytes is a conservative p2p budget. A single block may exceed
	// it, which [GetAncestors] tolerates by exempting the first block.
	maxResponseBytes = constants.MaxContainersLen

	// maxBlocksRetrievalTime matches the node's default for the equivalent
	// bootstrap GetAncestors operation.
	maxBlocksRetrievalTime = 50 * time.Millisecond
)

var errBlockNotFound = &common.AppError{
	Code:    2000,
	Message: "requested block not found",
}

func (r *responder) Respond(ctx context.Context, nodeID ids.NodeID, req *syncpb.GetBlockRequest) (*syncpb.GetBlockResponse, *common.AppError) {
	ctx, cancel := context.WithTimeout(ctx, maxBlocksRetrievalTime)
	defer cancel()

	// The request counts parents, so the response carries one more block.
	height := req.GetHeight()
	numParents := req.GetNumParents()
	blocks, err := GetAncestors(
		ctx,
		r.db,
		height,
		int(min(numParents, maxParentsPerRequest))+1,
		maxResponseBytes,
	)
	if err != nil {
		return nil, handlers.Fault(r.log, nodeID, err)
	}
	if len(blocks) == 0 {
		r.log.Debug("rejecting request, requested block not found",
			zap.Stringer("nodeID", nodeID),
			zap.Uint64("height", height),
			zap.Uint32("parents", numParents),
		)
		return nil, errBlockNotFound
	}
	return &syncpb.GetBlockResponse{Blocks: blocks}, nil
}
