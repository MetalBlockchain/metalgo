// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package metervm

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/MetalBlockchain/metalgo/utils/metric"
	"github.com/MetalBlockchain/metalgo/utils/wrappers"
)

func newAverager(name string, reg prometheus.Registerer, errs *wrappers.Errs) metric.Averager {
	return metric.NewAveragerWithErrs(
		name,
		"time (in ns) of a "+name,
		reg,
		errs,
	)
}
