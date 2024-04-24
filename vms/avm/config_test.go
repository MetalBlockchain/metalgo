// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package avm

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MetalBlockchain/metalgo/vms/avm/network"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name           string
		configBytes    []byte
		expectedConfig Config
	}{
		{
			name:           "unspecified config",
			configBytes:    nil,
			expectedConfig: DefaultConfig,
		},
		{
			name:        "manually specified checksums enabled",
			configBytes: []byte(`{"checksums-enabled":true}`),
			expectedConfig: Config{
				Network:              network.DefaultConfig,
				IndexTransactions:    DefaultConfig.IndexTransactions,
				IndexAllowIncomplete: DefaultConfig.IndexAllowIncomplete,
				ChecksumsEnabled:     true,
			},
		},
		{
			name:        "manually specified checksums enabled",
			configBytes: []byte(`{"network":{"max-validator-set-staleness":1}}`),
			expectedConfig: Config{
				Network: network.Config{
					MaxValidatorSetStaleness:                    time.Nanosecond,
					TargetGossipSize:                            network.DefaultConfig.TargetGossipSize,
					PullGossipPollSize:                          network.DefaultConfig.PullGossipPollSize,
					PullGossipFrequency:                         network.DefaultConfig.PullGossipFrequency,
					PullGossipThrottlingPeriod:                  network.DefaultConfig.PullGossipThrottlingPeriod,
					PullGossipThrottlingLimit:                   network.DefaultConfig.PullGossipThrottlingLimit,
					ExpectedBloomFilterElements:                 network.DefaultConfig.ExpectedBloomFilterElements,
					ExpectedBloomFilterFalsePositiveProbability: network.DefaultConfig.ExpectedBloomFilterFalsePositiveProbability,
					MaxBloomFilterFalsePositiveProbability:      network.DefaultConfig.MaxBloomFilterFalsePositiveProbability,
					LegacyPushGossipCacheSize:                   network.DefaultConfig.LegacyPushGossipCacheSize,
				},
				IndexTransactions:    DefaultConfig.IndexTransactions,
				IndexAllowIncomplete: DefaultConfig.IndexAllowIncomplete,
				ChecksumsEnabled:     DefaultConfig.ChecksumsEnabled,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require := require.New(t)

			config, err := ParseConfig(test.configBytes)
			require.NoError(err)
			require.Equal(test.expectedConfig, config)
		})
	}
}
