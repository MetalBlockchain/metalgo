// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package chains

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MetalBlockchain/metalgo/api/health"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/subnets"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/logging"
)

type testHealthRegisterer struct {
	healthChecks int
}

func (*testHealthRegisterer) RegisterReadinessCheck(string, health.Checker, ...string) error {
	return nil
}

func (r *testHealthRegisterer) RegisterHealthCheck(string, health.Checker, ...string) error {
	r.healthChecks++
	return nil
}

func (*testHealthRegisterer) RegisterLivenessCheck(string, health.Checker, ...string) error {
	return nil
}

func TestShutdownOnNonCriticalChainCreationFailure(t *testing.T) {
	chainParams := ChainParameters{
		ID:       ids.GenerateTestID(),
		SubnetID: constants.PrimaryNetworkID,
		// The manager rejects a PlatformVM assigned to a non-platform chain
		// before touching VM services, giving this test a deterministic init error.
		VMID: constants.PlatformVMID,
	}

	t.Run("opted in shuts down", func(t *testing.T) {
		subnetManager, err := NewSubnets(ids.EmptyNodeID, map[ids.ID]subnets.Config{
			constants.PrimaryNetworkID: {},
		})
		require.NoError(t, err)
		shutdown := make(chan int, 1)
		m := &manager{
			Aliaser: ids.NewAliaser(),
			ManagerConfig: ManagerConfig{
				Log:                            logging.NoLog{},
				Subnets:                        subnetManager,
				ShutdownOnChainCreationFailure: true,
				ShutdownNodeFunc: func(exitCode int) {
					shutdown <- exitCode
				},
			},
		}

		m.createChain(chainParams)

		select {
		case exitCode := <-shutdown:
			require.Equal(t, 1, exitCode)
		case <-time.After(time.Second):
			require.Fail(t, "node shutdown was not requested")
		}
	})

	t.Run("default keeps node up and reports failed chain health", func(t *testing.T) {
		subnetManager, err := NewSubnets(ids.EmptyNodeID, map[ids.ID]subnets.Config{
			constants.PrimaryNetworkID: {},
		})
		require.NoError(t, err)
		health := &testHealthRegisterer{}
		m := &manager{
			Aliaser: ids.NewAliaser(),
			ManagerConfig: ManagerConfig{
				Log:     logging.NoLog{},
				Subnets: subnetManager,
				Health:  health,
			},
		}

		m.createChain(chainParams)

		require.Equal(t, 1, health.healthChecks)
	})
}
