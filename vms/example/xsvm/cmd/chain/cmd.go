// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package chain

import (
	"github.com/spf13/cobra"

	"github.com/MetalBlockchain/metalgo/vms/example/xsvm/cmd/chain/create"
	"github.com/MetalBlockchain/metalgo/vms/example/xsvm/cmd/chain/genesis"
)

func Command() *cobra.Command {
	c := &cobra.Command{
		Use:   "chain",
		Short: "Manages XS chains",
	}
	c.AddCommand(
		create.Command(),
		genesis.Command(),
	)
	return c
}
