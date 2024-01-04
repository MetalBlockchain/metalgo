// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package txs

import (
	"math"

	"github.com/MetalBlockchain/metalgo/codec"
	"github.com/MetalBlockchain/metalgo/codec/linearcodec"
	"github.com/MetalBlockchain/metalgo/utils"
	"github.com/MetalBlockchain/metalgo/utils/wrappers"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/signer"
	"github.com/MetalBlockchain/metalgo/vms/platformvm/stakeable"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
)

// Version is the current default codec version
const Version = 0

var (
	Codec codec.Manager

	// GenesisCodec allows txs of larger than usual size to be parsed.
	// While this gives flexibility in accommodating large genesis txs
	// it must not be used to parse new, unverified txs which instead
	// must be processed by Codec
	GenesisCodec codec.Manager
)

func init() {
	c := linearcodec.NewDefault()
	Codec = codec.NewDefaultManager()
	gc := linearcodec.NewCustomMaxLength(math.MaxInt32)
	GenesisCodec = codec.NewManager(math.MaxInt32)

	errs := wrappers.Errs{}
	for _, c := range []linearcodec.Codec{c, gc} {
		// Order in which type are registered affect the byte representation
		// generated by marshalling ops. To maintain codec type ordering,
		// we skip positions for the blocks.
		c.SkipRegistrations(5)

		errs.Add(RegisterUnsignedTxsTypes(c))

		c.SkipRegistrations(4)

		errs.Add(RegisterDUnsignedTxsTypes(c))
	}
	errs.Add(
		Codec.RegisterCodec(Version, c),
		GenesisCodec.RegisterCodec(Version, gc),
	)
	if errs.Errored() {
		panic(errs.Err)
	}
}

// RegisterUnsignedTxsTypes allows registering relevant type of unsigned package
// in the right sequence. Following repackaging of platformvm package, a few
// subpackage-level codecs were introduced, each handling serialization of
// specific types.
//
// RegisterUnsignedTxsTypes is made exportable so to guarantee that other codecs
// are coherent with components one.
func RegisterUnsignedTxsTypes(targetCodec linearcodec.Codec) error {
	errs := wrappers.Errs{}

	// The secp256k1fx is registered here because this is the same place it is
	// registered in the AVM. This ensures that the typeIDs match up for utxos
	// in shared memory.
	errs.Add(targetCodec.RegisterType(&secp256k1fx.TransferInput{}))
	targetCodec.SkipRegistrations(1)
	errs.Add(targetCodec.RegisterType(&secp256k1fx.TransferOutput{}))
	targetCodec.SkipRegistrations(1)
	errs.Add(
		targetCodec.RegisterType(&secp256k1fx.Credential{}),
		targetCodec.RegisterType(&secp256k1fx.Input{}),
		targetCodec.RegisterType(&secp256k1fx.OutputOwners{}),

		targetCodec.RegisterType(&AddValidatorTx{}),
		targetCodec.RegisterType(&AddSubnetValidatorTx{}),
		targetCodec.RegisterType(&AddDelegatorTx{}),
		targetCodec.RegisterType(&CreateChainTx{}),
		targetCodec.RegisterType(&CreateSubnetTx{}),
		targetCodec.RegisterType(&ImportTx{}),
		targetCodec.RegisterType(&ExportTx{}),
		targetCodec.RegisterType(&AdvanceTimeTx{}),
		targetCodec.RegisterType(&RewardValidatorTx{}),

		targetCodec.RegisterType(&stakeable.LockIn{}),
		targetCodec.RegisterType(&stakeable.LockOut{}),

		// Banff additions:
		targetCodec.RegisterType(&RemoveSubnetValidatorTx{}),
		targetCodec.RegisterType(&TransformSubnetTx{}),
		targetCodec.RegisterType(&AddPermissionlessValidatorTx{}),
		targetCodec.RegisterType(&AddPermissionlessDelegatorTx{}),

		targetCodec.RegisterType(&signer.Empty{}),
		targetCodec.RegisterType(&signer.ProofOfPossession{}),
	)
	return errs.Err
}

func RegisterDUnsignedTxsTypes(targetCodec linearcodec.Codec) error {
	return utils.Err(
		targetCodec.RegisterType(&TransferSubnetOwnershipTx{}),
		targetCodec.RegisterType(&BaseTx{}),
	)
}
