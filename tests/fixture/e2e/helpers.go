// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package e2e

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/MetalBlockchain/coreth/core/types"
	"github.com/MetalBlockchain/coreth/ethclient"
	"github.com/MetalBlockchain/coreth/interfaces"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/MetalBlockchain/metalgo/config"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/tests"
	"github.com/MetalBlockchain/metalgo/tests/fixture/tmpnet"
	"github.com/MetalBlockchain/metalgo/utils/crypto/secp256k1"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/txs/fee"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
	"github.com/MetalBlockchain/metalgo/wallet/chain/p/builder"
	"github.com/MetalBlockchain/metalgo/wallet/subnet/primary"
	"github.com/MetalBlockchain/metalgo/wallet/subnet/primary/common"
)

const (
	DefaultTimeout = tests.DefaultTimeout

	DefaultPollingInterval = tmpnet.DefaultPollingInterval

	// Setting this env will disable post-test bootstrap
	// checks. Useful for speeding up iteration during test
	// development.
	SkipBootstrapChecksEnvName = "E2E_SKIP_BOOTSTRAP_CHECKS"

	DefaultValidatorStartTimeDiff = tmpnet.DefaultValidatorStartTimeDiff

	DefaultGasLimit = uint64(21000) // Standard gas limit

	// An empty string prompts the use of the default path which ensures a
	// predictable target for github's upload-artifact action.
	DefaultNetworkDir = ""

	// Directory used to store private networks (specific to a single test)
	// under the shared network dir.
	PrivateNetworksDirName = "private_networks"
)

// NewPrivateKey returns a new private key.
func NewPrivateKey(tc tests.TestContext) *secp256k1.PrivateKey {
	key, err := secp256k1.NewPrivateKey()
	require.NoError(tc, err)
	return key
}

// Create a new wallet for the provided keychain against the specified node URI.
func NewWallet(tc tests.TestContext, keychain *secp256k1fx.Keychain, nodeURI tmpnet.NodeURI) *primary.Wallet {
	tc.Log().Info("initializing a new wallet",
		zap.Stringer("nodeID", nodeURI.NodeID),
		zap.String("URI", nodeURI.URI),
	)
	baseWallet, err := primary.MakeWallet(
		tc.DefaultContext(),
		nodeURI.URI,
		keychain,
		keychain,
		primary.WalletConfig{},
	)
	require.NoError(tc, err)
	wallet := primary.NewWalletWithOptions(
		baseWallet,
		common.WithPostIssuanceFunc(
			func(id ids.ID) {
				tc.Log().Info("issued transaction",
					zap.Stringer("txID", id),
				)
			},
		),
	)
	OutputWalletBalances(tc, wallet)
	return wallet
}

// OutputWalletBalances outputs the X-Chain and P-Chain balances of the provided wallet.
func OutputWalletBalances(tc tests.TestContext, wallet *primary.Wallet) {
	_, _ = GetWalletBalances(tc, wallet)
}

// GetWalletBalances retrieves the X-Chain and P-Chain balances of the provided wallet.
func GetWalletBalances(tc tests.TestContext, wallet *primary.Wallet) (uint64, uint64) {
	require := require.New(tc)
	var (
		xWallet  = wallet.X()
		xBuilder = xWallet.Builder()
		pWallet  = wallet.P()
		pBuilder = pWallet.Builder()
	)
	xBalances, err := xBuilder.GetFTBalance()
	require.NoError(err, "failed to fetch X-chain balances")
	pBalances, err := pBuilder.GetBalance()
	require.NoError(err, "failed to fetch P-chain balances")
	var (
		xContext    = xBuilder.Context()
		avaxAssetID = xContext.AVAXAssetID
		xAVAX       = xBalances[avaxAssetID]
		pAVAX       = pBalances[avaxAssetID]
	)
	tc.Log().Info("wallet balances in nAVAX",
		zap.Uint64("xChain", xAVAX),
		zap.Uint64("pChain", pAVAX),
	)
	return xAVAX, pAVAX
}

// Create a new eth client targeting the specified node URI.
func NewEthClient(tc tests.TestContext, nodeURI tmpnet.NodeURI) ethclient.Client {
	tc.Log().Info("initializing a new eth client",
		zap.Stringer("nodeID", nodeURI.NodeID),
		zap.String("URI", nodeURI.URI),
	)
	nodeAddress := strings.Split(nodeURI.URI, "//")[1]
	uri := fmt.Sprintf("ws://%s/ext/bc/C/ws", nodeAddress)
	client, err := ethclient.Dial(uri)
	require.NoError(tc, err)
	return client
}

// Adds an ephemeral node intended to be used by a single test.
func AddEphemeralNode(tc tests.TestContext, network *tmpnet.Network, flags tmpnet.FlagsMap) *tmpnet.Node {
	require := require.New(tc)

	node := tmpnet.NewEphemeralNode(flags)
	require.NoError(network.StartNode(tc.DefaultContext(), tc.Log(), node))

	tc.DeferCleanup(func() {
		tc.Log().Info("shutting down ephemeral node",
			zap.Stringer("nodeID", node.NodeID),
		)
		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()
		require.NoError(node.Stop(ctx))
	})
	return node
}

// Wait for the given node to report healthy.
func WaitForHealthy(t require.TestingT, node *tmpnet.Node) {
	// Need to use explicit context (vs DefaultContext()) to support use with DeferCleanup
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	require.NoError(t, tmpnet.WaitForHealthy(ctx, node))
}

// Sends an eth transaction, waits for the transaction receipt to be issued
// and checks that the receipt indicates success.
func SendEthTransaction(tc tests.TestContext, ethClient ethclient.Client, signedTx *types.Transaction) *types.Receipt {
	require := require.New(tc)

	txID := signedTx.Hash()
	tc.Log().Info("sending eth transaction",
		zap.Stringer("txID", txID),
	)

	require.NoError(ethClient.SendTransaction(tc.DefaultContext(), signedTx))

	// Wait for the receipt
	var receipt *types.Receipt
	tc.Eventually(func() bool {
		var err error
		receipt, err = ethClient.TransactionReceipt(tc.DefaultContext(), txID)
		if errors.Is(err, interfaces.NotFound) {
			return false // Transaction is still pending
		}
		require.NoError(err)
		return true
	}, DefaultTimeout, DefaultPollingInterval, "failed to see transaction acceptance before timeout")

	require.Equal(types.ReceiptStatusSuccessful, receipt.Status)
	return receipt
}

// Determines the suggested gas price for the configured client that will
// maximize the chances of transaction acceptance.
func SuggestGasPrice(tc tests.TestContext, ethClient ethclient.Client) *big.Int {
	gasPrice, err := ethClient.SuggestGasPrice(tc.DefaultContext())
	require.NoError(tc, err)
	// Double the suggested gas price to maximize the chances of
	// acceptance. Maybe this can be revisited pending resolution of
	// https://github.com/MetalBlockchain/coreth/issues/314.
	gasPrice.Add(gasPrice, gasPrice)
	return gasPrice
}

// Helper simplifying use via an option of a gas price appropriate for testing.
func WithSuggestedGasPrice(tc tests.TestContext, ethClient ethclient.Client) common.Option {
	baseFee := SuggestGasPrice(tc, ethClient)
	return common.WithBaseFee(baseFee)
}

// Verify that a new node can bootstrap into the network. If the check wasn't skipped,
// the node will be returned to the caller.
func CheckBootstrapIsPossible(tc tests.TestContext, network *tmpnet.Network) *tmpnet.Node {
	require := require.New(tc)

	if len(os.Getenv(SkipBootstrapChecksEnvName)) > 0 {
		tc.Log().Info("skipping bootstrap check due to env var being set",
			zap.String("envVar", SkipBootstrapChecksEnvName),
		)
		return nil
	}
	tc.By("checking if bootstrap is possible with the current network state")

	// Ensure all subnets are bootstrapped
	subnetIDs := make([]string, len(network.Subnets))
	for i, subnet := range network.Subnets {
		subnetIDs[i] = subnet.SubnetID.String()
	}
	flags := tmpnet.FlagsMap{
		config.TrackSubnetsKey: strings.Join(subnetIDs, ","),
	}

	node := tmpnet.NewEphemeralNode(flags)
	require.NoError(network.StartNode(tc.DefaultContext(), tc.Log(), node))
	// StartNode will initiate node stop if an error is encountered during start,
	// so no further cleanup effort is required if an error is seen here.

	// Register a cleanup to ensure the node is stopped at the end of the test
	tc.DeferCleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()
		require.NoError(node.Stop(ctx))
	})

	// Check that the node becomes healthy within timeout
	require.NoError(tmpnet.WaitForHealthy(tc.DefaultContext(), node))

	// Ensure that the primary validators are still healthy
	for _, node := range network.Nodes {
		if node.IsEphemeral {
			continue
		}
		healthy, err := node.IsHealthy(tc.DefaultContext())
		require.NoError(err)
		require.True(healthy, "primary validator %s is not healthy", node.NodeID)
	}

	return node
}

// Start a temporary network with the provided avalanchego binary.
func StartNetwork(
	tc tests.TestContext,
	network *tmpnet.Network,
	avalancheGoExecPath string,
	pluginDir string,
	shutdownDelay time.Duration,
	skipShutdown bool,
	reuseNetwork bool,
) {
	require := require.New(tc)

	err := tmpnet.BootstrapNewNetwork(
		tc.DefaultContext(),
		tc.Log(),
		network,
		DefaultNetworkDir,
		avalancheGoExecPath,
		pluginDir,
	)
	if err != nil {
		// Ensure nodes are stopped if bootstrap fails. The network configuration
		// will remain on disk to enable troubleshooting.
		err := network.Stop(tc.DefaultContext())
		require.NoError(err, "failed to stop network after bootstrap failure")
	}

	tc.Log().Info("network started successfully")

	symlinkPath, err := tmpnet.GetReusableNetworkPathForOwner(network.Owner)
	require.NoError(err)

	if reuseNetwork {
		// Symlink the path of the created network to the default owner path (e.g. latest_avalanchego-e2e)
		// to enable easy discovery for reuse.
		require.NoError(os.Symlink(network.Dir, symlinkPath))
		tc.Log().Info("symlinked network dir for reuse",
			zap.String("networkDir", network.Dir),
			zap.String("symlinkPath", symlinkPath),
		)
	}

	tc.DeferCleanup(func() {
		if reuseNetwork {
			tc.Log().Info("skipping shutdown for network intended for reuse",
				zap.String("networkDir", network.Dir),
				zap.String("symlinkPath", symlinkPath),
			)
			return
		}

		if skipShutdown {
			tc.Log().Info("skipping shutdown for network",
				zap.String("networkDir", network.Dir),
			)
			return
		}

		if shutdownDelay > 0 {
			tc.Log().Info("delaying network shutdown to ensure final metrics scrape",
				zap.Duration("delay", shutdownDelay),
			)
			time.Sleep(shutdownDelay)
		}

		tc.Log().Info("shutting down network")
		ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
		defer cancel()
		require.NoError(network.Stop(ctx))
	})
}

// NewPChainFeeCalculatorFromContext returns either a static or dynamic fee
// calculator depending on the provided context.
func NewPChainFeeCalculatorFromContext(context *builder.Context) fee.Calculator {
	if context.GasPrice != 0 {
		return fee.NewDynamicCalculator(context.ComplexityWeights, context.GasPrice)
	}
	return fee.NewStaticCalculator(context.StaticFeeConfig)
}

// GetRepoRootPath strips the provided suffix from the current working
// directory. If the test binary is executed from the root of the repo, the
// result will be the repo root.
func GetRepoRootPath(suffix string) (string, error) {
	// - When executed via a test binary, the working directory will be wherever
	// the binary is executed from, but scripts should require execution from
	// the repo root.
	//
	// - When executed via ginkgo (nicer for development + supports
	// parallel execution) the working directory will always be the
	// target path (e.g. [repo root]./tests/bootstrap/e2e) and getting the repo
	// root will require stripping the target path suffix.
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(cwd, suffix), nil
}
