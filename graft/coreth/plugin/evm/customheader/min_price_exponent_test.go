// Copyright (C) 2019, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package customheader

import (
	"testing"

	"github.com/MetalBlockchain/libevm/core/types"
	"github.com/stretchr/testify/require"

	"github.com/MetalBlockchain/metalgo/graft/coreth/plugin/evm/customtypes"
	"github.com/MetalBlockchain/metalgo/utils"
	"github.com/MetalBlockchain/metalgo/vms/saevm/cchain/dynamic"
)

func TestVerifyMinPriceExponent(t *testing.T) {
	require.NoError(t, VerifyMinPriceExponent(&types.Header{Time: 1001}))

	withExponent := customtypes.WithHeaderExtra(
		&types.Header{Time: 1001},
		&customtypes.HeaderExtra{MinPriceExponent: utils.PointerTo(dynamic.PriceExponent(1000))},
	)
	require.ErrorIs(t, VerifyMinPriceExponent(withExponent), errRemoteMinPriceExponentSet)
}
