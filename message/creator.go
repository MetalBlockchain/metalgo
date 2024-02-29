// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package message

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/MetalBlockchain/metalgo/utils/compression"
	"github.com/MetalBlockchain/metalgo/utils/logging"
)

var _ Creator = (*creator)(nil)

type Creator interface {
	OutboundMsgBuilder
	InboundMsgBuilder
}

type creator struct {
	OutboundMsgBuilder
	InboundMsgBuilder
}

func NewCreator(
	log logging.Logger,
	metrics prometheus.Registerer,
	parentNamespace string,
	compressionType compression.Type,
	maxMessageTimeout time.Duration,
) (Creator, error) {
	namespace := fmt.Sprintf("%s_codec", parentNamespace)
	builder, err := newMsgBuilder(
		log,
		namespace,
		metrics,
		maxMessageTimeout,
	)
	if err != nil {
		return nil, err
	}

	return &creator{
		OutboundMsgBuilder: newOutboundBuilder(compressionType, builder),
		InboundMsgBuilder:  newInboundBuilder(builder),
	}, nil
}
