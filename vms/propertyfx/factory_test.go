// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package propertyfx

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MetalBlockchain/metalgo/utils/logging"
)

func TestFactory(t *testing.T) {
	require := require.New(t)

	factory := Factory{}
	fx, err := factory.New(logging.NoLog{})
	require.NoError(err)
	require.NotNil(fx)
}
