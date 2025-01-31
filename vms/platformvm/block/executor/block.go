// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/MetalBlockchain/metalgo/snow/consensus/snowman"
	"github.com/MetalBlockchain/metalgo/utils"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/block"

	smblock "github.com/MetalBlockchain/metalgo/snow/engine/snowman/block"
)

const ActivationLog = `
 __  __      _        _   _     _           
|  \/  | ___| |_ __ _| | | |   / |___       
| |\/| |/ _ \ __/ _' | | | |   | / __|      
| |  | |  __/ || (_| | | | |___| \__ \      
|_| _|_|\___|\__\__,_|_| |_____|_|___/    _ 
   / \   ___| |_(_)_   ____ _| |_ ___  __| |
  / _ \ / __| __| \ \ / / _' | __/ _ \/ _' |
 / ___ \ (__| |_| |\ V / (_| | ||  __/ (_| |
/_/   \_\___|\__|_| \_/ \__,_|\__\___|\__,_|
`

// TODO: Remove after Etna is activated
var EtnaActivationWasLogged = utils.NewAtomic(false)

var (
	_ snowman.Block             = (*Block)(nil)
	_ snowman.OracleBlock       = (*Block)(nil)
	_ smblock.WithVerifyContext = (*Block)(nil)
)

// Exported for testing in platformvm package.
type Block struct {
	block.Block
	manager *manager
}

func (*Block) ShouldVerifyWithContext(context.Context) (bool, error) {
	return true, nil
}

func (b *Block) VerifyWithContext(ctx context.Context, blockContext *smblock.Context) error {
	blkID := b.ID()
	blkState, previouslyExecuted := b.manager.blkIDToState[blkID]
	warpAlreadyVerified := previouslyExecuted && blkState.verifiedHeights.Contains(blockContext.PChainHeight)

	// If the chain is bootstrapped and the warp messages haven't been verified,
	// we must verify them.
	if !warpAlreadyVerified && b.manager.txExecutorBackend.Bootstrapped.Get() {
		err := VerifyWarpMessages(
			ctx,
			b.manager.ctx.NetworkID,
			b.manager.ctx.ValidatorState,
			blockContext.PChainHeight,
			b,
		)
		if err != nil {
			return err
		}
	}

	// If the block was previously executed, we don't need to execute it again,
	// we can just mark that the warp messages are valid at this height.
	if previouslyExecuted {
		blkState.verifiedHeights.Add(blockContext.PChainHeight)
		return nil
	}

	// Since this is the first time we are verifying this block, we must execute
	// the state transitions to generate the state diffs.
	return b.Visit(&verifier{
		backend:           b.manager.backend,
		txExecutorBackend: b.manager.txExecutorBackend,
		pChainHeight:      blockContext.PChainHeight,
	})
}

func (b *Block) Verify(ctx context.Context) error {
	return b.VerifyWithContext(
		ctx,
		&smblock.Context{
			PChainHeight: 0,
		},
	)
}

func (b *Block) Accept(context.Context) error {
	if err := b.Visit(b.manager.acceptor); err != nil {
		return err
	}

	currentTime := b.manager.state.GetTimestamp()
	if !b.manager.txExecutorBackend.Config.UpgradeConfig.IsEtnaActivated(currentTime) {
		return nil
	}

	if !EtnaActivationWasLogged.Get() && b.manager.txExecutorBackend.Bootstrapped.Get() {
		fmt.Print(ActivationLog)
	}
	EtnaActivationWasLogged.Set(true)
	return nil
}

func (b *Block) Reject(context.Context) error {
	return b.Visit(b.manager.rejector)
}

func (b *Block) Timestamp() time.Time {
	return b.manager.getTimestamp(b.ID())
}

func (b *Block) Options(context.Context) ([2]snowman.Block, error) {
	options := options{
		log:                     b.manager.ctx.Log,
		primaryUptimePercentage: b.manager.txExecutorBackend.Config.UptimePercentage,
		uptimes:                 b.manager.txExecutorBackend.Uptimes,
		state:                   b.manager.backend.state,
	}
	if err := b.Block.Visit(&options); err != nil {
		return [2]snowman.Block{}, err
	}

	return [2]snowman.Block{
		b.manager.NewBlock(options.preferredBlock),
		b.manager.NewBlock(options.alternateBlock),
	}, nil
}
