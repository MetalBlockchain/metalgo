// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

// Implements tests for the banff network upgrade.
package banff

import (
	ginkgo "github.com/onsi/ginkgo/v2"

	"github.com/onsi/gomega"

	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/tests"
	"github.com/MetalBlockchain/metalgo/tests/e2e"
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/units"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/components/verify"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
)

var _ = ginkgo.Describe("[Banff]", func() {
	ginkgo.It("can send custom assets X->P and P->X",
		// use this for filtering tests by labels
		// ref. https://onsi.github.io/ginkgo/#spec-labels
		ginkgo.Label(
			"xp",
			"banff",
		),
		func() {
			keychain := e2e.Env.NewKeychain(1)
			wallet := e2e.Env.NewWallet(keychain, e2e.Env.GetRandomNodeURI())

			// Get the P-chain and the X-chain wallets
			pWallet := wallet.P()
			xWallet := wallet.X()

			// Pull out useful constants to use when issuing transactions.
			xChainID := xWallet.BlockchainID()
			owner := &secp256k1fx.OutputOwners{
				Threshold: 1,
				Addrs: []ids.ShortID{
					keychain.Keys[0].Address(),
				},
			}

			var assetID ids.ID
			ginkgo.By("create new X-chain asset", func() {
				assetTx, err := xWallet.IssueCreateAssetTx(
					"RnM",
					"RNM",
					9,
					map[uint32][]verify.State{
						0: {
							&secp256k1fx.TransferOutput{
								Amt:          100 * units.Schmeckle,
								OutputOwners: *owner,
							},
						},
					},
				)
				gomega.Expect(err).Should(gomega.BeNil())
				assetID = assetTx.ID()

				tests.Outf("{{green}}created new X-chain asset{{/}}: %s\n", assetID)
			})

			ginkgo.By("export new X-chain asset to P-chain", func() {
				tx, err := xWallet.IssueExportTx(
					constants.PlatformChainID,
					[]*avax.TransferableOutput{
						{
							Asset: avax.Asset{
								ID: assetID,
							},
							Out: &secp256k1fx.TransferOutput{
								Amt:          100 * units.Schmeckle,
								OutputOwners: *owner,
							},
						},
					},
				)
				gomega.Expect(err).Should(gomega.BeNil())

				tests.Outf("{{green}}issued X-chain export{{/}}: %s\n", tx.ID())
			})

			ginkgo.By("import new asset from X-chain on the P-chain", func() {
				tx, err := pWallet.IssueImportTx(xChainID, owner)
				gomega.Expect(err).Should(gomega.BeNil())

				tests.Outf("{{green}}issued P-chain import{{/}}: %s\n", tx.ID())
			})

			ginkgo.By("export asset from P-chain to the X-chain", func() {
				tx, err := pWallet.IssueExportTx(
					xChainID,
					[]*avax.TransferableOutput{
						{
							Asset: avax.Asset{
								ID: assetID,
							},
							Out: &secp256k1fx.TransferOutput{
								Amt:          100 * units.Schmeckle,
								OutputOwners: *owner,
							},
						},
					},
				)
				gomega.Expect(err).Should(gomega.BeNil())

				tests.Outf("{{green}}issued P-chain export{{/}}: %s\n", tx.ID())
			})

			ginkgo.By("import asset from P-chain on the X-chain", func() {
				tx, err := xWallet.IssueImportTx(constants.PlatformChainID, owner)
				gomega.Expect(err).Should(gomega.BeNil())

				tests.Outf("{{green}}issued X-chain import{{/}}: %s\n", tx.ID())
			})
		})
})
