// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"os"

	"go.uber.org/zap"

	"github.com/MetalBlockchain/metalgo/tests"
	"github.com/MetalBlockchain/metalgo/tests/antithesis"
	"github.com/MetalBlockchain/metalgo/tests/fixture/tmpnet"
)

const baseImageName = "antithesis-avalanchego"

// Creates docker-compose.yml and its associated volumes in the target path.
func main() {
	network := tmpnet.LocalNetworkOrPanic()
	if err := antithesis.GenerateComposeConfig(network, baseImageName, "" /* runtimePluginDir */); err != nil {
		tests.NewDefaultLogger("").Fatal("failed to generate compose config",
			zap.Error(err),
		)
		os.Exit(1)
	}
}
