// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package snowtest

import (
	"context"
	"errors"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"github.com/MetalBlockchain/metalgo/api/metrics"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/snow"
	"github.com/MetalBlockchain/metalgo/snow/validators/validatorstest"
	"github.com/MetalBlockchain/metalgo/upgrade/upgradetest"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/crypto/bls"
	"github.com/MetalBlockchain/metalgo/utils/logging"
)

var (
	XChainID    = ids.GenerateTestID()
	CChainID    = ids.GenerateTestID()
	PChainID    = constants.PlatformChainID
	AVAXAssetID = ids.GenerateTestID()

	errMissing = errors.New("missing")

	_ snow.Acceptor = noOpAcceptor{}
)

type noOpAcceptor struct{}

func (noOpAcceptor) Accept(*snow.ConsensusContext, ids.ID, []byte) error {
	return nil
}

func ConsensusContext(ctx *snow.Context) *snow.ConsensusContext {
	return &snow.ConsensusContext{
		Context:        ctx,
		PrimaryAlias:   ctx.ChainID.String(),
		Registerer:     prometheus.NewRegistry(),
		BlockAcceptor:  noOpAcceptor{},
		TxAcceptor:     noOpAcceptor{},
		VertexAcceptor: noOpAcceptor{},
	}
}

func Context(tb testing.TB, chainID ids.ID) *snow.Context {
	require := require.New(tb)

	secretKey, err := bls.NewSecretKey()
	require.NoError(err)
	publicKey := bls.PublicFromSecretKey(secretKey)

	aliaser := ids.NewAliaser()
	require.NoError(aliaser.Alias(constants.PlatformChainID, "P"))
	require.NoError(aliaser.Alias(constants.PlatformChainID, constants.PlatformChainID.String()))
	require.NoError(aliaser.Alias(XChainID, "X"))
	require.NoError(aliaser.Alias(XChainID, XChainID.String()))
	require.NoError(aliaser.Alias(CChainID, "C"))
	require.NoError(aliaser.Alias(CChainID, CChainID.String()))

	validatorState := &validatorstest.State{
		GetMinimumHeightF: func(context.Context) (uint64, error) {
			return 0, nil
		},
		GetSubnetIDF: func(_ context.Context, chainID ids.ID) (ids.ID, error) {
			subnetID, ok := map[ids.ID]ids.ID{
				constants.PlatformChainID: constants.PrimaryNetworkID,
				XChainID:                  constants.PrimaryNetworkID,
				CChainID:                  constants.PrimaryNetworkID,
			}[chainID]
			if !ok {
				return ids.Empty, errMissing
			}
			return subnetID, nil
		},
	}

	return &snow.Context{
		NetworkID:       constants.UnitTestID,
		SubnetID:        constants.PrimaryNetworkID,
		ChainID:         chainID,
		NodeID:          ids.EmptyNodeID,
		PublicKey:       publicKey,
		NetworkUpgrades: upgradetest.GetConfig(upgradetest.Latest),

		XChainID:    XChainID,
		CChainID:    CChainID,
		AVAXAssetID: AVAXAssetID,

		Log:      logging.NoLog{},
		BCLookup: aliaser,
		Metrics:  metrics.NewPrefixGatherer(),

		ValidatorState: validatorState,
		ChainDataDir:   "",
	}
}
