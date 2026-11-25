// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package chains

import (
	"context"
	"sync"

	"github.com/MetalBlockchain/metalgo/api/health"
)

// shutdownOnUnhealthyChecker shuts down the node after repeated VM health
// failures. A successful check clears the failure streak so transient failures
// do not take down the node.
type shutdownOnUnhealthyChecker struct {
	checker   health.Checker
	threshold int
	shutdown  func(int)

	mu                  sync.Mutex
	consecutiveFailures int
	shutdownRequested   bool
}

func (c *shutdownOnUnhealthyChecker) HealthCheck(ctx context.Context) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	details, err := c.checker.HealthCheck(ctx)
	if err == nil {
		c.consecutiveFailures = 0
		return details, nil
	}

	c.consecutiveFailures++
	if !c.shutdownRequested && c.consecutiveFailures >= c.threshold {
		c.shutdownRequested = true
		go c.shutdown(1)
	}
	return details, err
}
