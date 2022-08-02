// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"time"

	"github.com/MetalBlockchain/avalanchego/utils/units"
	"github.com/MetalBlockchain/avalanchego/vms/platformvm/reward"
)

var (
	mainnetGenesisConfigJSON = `{
		"networkID": 1,
		"allocations": [
			{
				"ethAddr": "0x15884b0f70b7c0db084e5aa738605f7a681a5d6e",
				"avaxAddr": "X-metal1922357kfy7ygftn9ahahtp6s4jcv5umsrv6dk6",
				"initialAmount": 120000000000000000
			},
			{
				"ethAddr": "0x15884b0f70b7c0db084e5aa738605f7a681a5d6e",
				"avaxAddr": "X-metal1wulualvtpydn9ur64zwm6v69vut8f89mpf5ag6",
				"initialAmount": 144333333000000000
			},
			{
				"ethAddr": "0x15884b0f70b7c0db084e5aa738605f7a681a5d6e",
				"avaxAddr": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"initialAmount": 0,
				"unlockSchedule": [
					{
						"amount": 10000000000000000
					}
				]
			},
			{
				"ethAddr": "0xb3d82b1367d362de99ab59a658165aff520cbd4d",
				"avaxAddr": "X-metal15s7ephn84n4gkpqrncpqpm2cfr72z7490tzmqy",
				"initialAmount": 0,
				"unlockSchedule": [
					{
						"amount": 1750000000000000,
						"locktime": 1658987891
					},
					{
						"amount": 1750000000000000,
						"locktime": 1661666291
					},
					{
						"amount": 1750000000000000,
						"locktime": 1664344691
					},
					{
						"amount": 1750000000000000,
						"locktime": 1667023091
					},
					{
						"amount": 1750000000000000,
						"locktime": 1669701491
					},
					{
						"amount": 1750000000000000,
						"locktime": 1672379891
					},
					{
						"amount": 1750000000000000,
						"locktime": 1675058291
					},
					{
						"amount": 1750000000000000,
						"locktime": 1677736691
					},
					{
						"amount": 1750000000000000,
						"locktime": 1680415091
					},
					{
						"amount": 1750000000000000,
						"locktime": 1683093491
					},
					{
						"amount": 1750000000000000,
						"locktime": 1685771891
					},
					{
						"amount": 1750000000000000,
						"locktime": 1688450291
					}
				]
			},
			{
				"ethAddr": "0xb3d82b1367d362de99ab59a658165aff520cbd4d",
				"avaxAddr": "X-metal17qq8ygxm5yuer3q2fj3z4jz379ktcjlragyxqg",
				"initialAmount": 0,
				"unlockSchedule": [
					{
						"amount": 1500000000000000,
						"locktime": 1658987891
					},
					{
						"amount": 1500000000000000,
						"locktime": 1661666291
					},
					{
						"amount": 1500000000000000,
						"locktime": 1664344691
					},
					{
						"amount": 1500000000000000,
						"locktime": 1667023091
					},
					{
						"amount": 1500000000000000,
						"locktime": 1669701491
					},
					{
						"amount": 1500000000000000,
						"locktime": 1672379891
					},
					{
						"amount": 1500000000000000,
						"locktime": 1675058291
					},
					{
						"amount": 1500000000000000,
						"locktime": 1677736691
					},
					{
						"amount": 1500000000000000,
						"locktime": 1680415091
					},
					{
						"amount": 1500000000000000,
						"locktime": 1683093491
					},
					{
						"amount": 1500000000000000,
						"locktime": 1685771891
					},
					{
						"amount": 1500000000000000,
						"locktime": 1688450291
					}
				]
			},
			{
				"ethAddr": "0xb3d82b1367d362de99ab59a658165aff520cbd4d",
				"avaxAddr": "X-metal1qup7x7fx9tf2jft2wue56t4hvepsqalcp24qf0",
				"initialAmount": 0,
				"unlockSchedule": [
					{
						"amount": 1500000000000000,
						"locktime": 1658987891
					},
					{
						"amount": 1500000000000000,
						"locktime": 1661666291
					},
					{
						"amount": 1500000000000000,
						"locktime": 1664344691
					},
					{
						"amount": 1500000000000000,
						"locktime": 1667023091
					},
					{
						"amount": 1500000000000000,
						"locktime": 1669701491
					},
					{
						"amount": 1500000000000000,
						"locktime": 1672379891
					},
					{
						"amount": 1500000000000000,
						"locktime": 1675058291
					},
					{
						"amount": 1500000000000000,
						"locktime": 1677736691
					},
					{
						"amount": 1500000000000000,
						"locktime": 1680415091
					},
					{
						"amount": 1500000000000000,
						"locktime": 1683093491
					},
					{
						"amount": 1500000000000000,
						"locktime": 1685771891
					},
					{
						"amount": 1500000000000000,
						"locktime": 1688450291
					}
				]
			},
			{
				"ethAddr": "0xb3d82b1367d362de99ab59a658165aff520cbd4d",
				"avaxAddr": "X-metal10tqqmmymlcpyggxjfth36nwdszcynwr7nz8mh8",
				"initialAmount": 0,
				"unlockSchedule": [
					{
						"amount": 1000000000000000,
						"locktime": 1658987891
					},
					{
						"amount": 1000000000000000,
						"locktime": 1661666291
					},
					{
						"amount": 1000000000000000,
						"locktime": 1664344691
					},
					{
						"amount": 1000000000000000,
						"locktime": 1667023091
					},
					{
						"amount": 1000000000000000,
						"locktime": 1669701491
					},
					{
						"amount": 1000000000000000,
						"locktime": 1672379891
					},
					{
						"amount": 1000000000000000,
						"locktime": 1675058291
					},
					{
						"amount": 1000000000000000,
						"locktime": 1677736691
					},
					{
						"amount": 1000000000000000,
						"locktime": 1680415091
					},
					{
						"amount": 1000000000000000,
						"locktime": 1683093491
					},
					{
						"amount": 1000000000000000,
						"locktime": 1685771891
					},
					{
						"amount": 1000000000000000,
						"locktime": 1688450291
					}
				]
			},
			{
				"ethAddr": "0xb3d82b1367d362de99ab59a658165aff520cbd4d",
				"avaxAddr": "X-metal1sfgp2glrc4wz37f8e028chrht86hqtwzajpe0u",
				"initialAmount": 0,
				"unlockSchedule": [
					{
						"amount": 83333333333333,
						"locktime": 1658987891
					},
					{
						"amount": 83333333333333,
						"locktime": 1661666291
					},
					{
						"amount": 83333333333333,
						"locktime": 1664344691
					},
					{
						"amount": 83333333333333,
						"locktime": 1667023091
					},
					{
						"amount": 83333333333333,
						"locktime": 1669701491
					},
					{
						"amount": 83333333333333,
						"locktime": 1672379891
					},
					{
						"amount": 83333333333333,
						"locktime": 1675058291
					},
					{
						"amount": 83333333333333,
						"locktime": 1677736691
					},
					{
						"amount": 83333333333333,
						"locktime": 1680415091
					},
					{
						"amount": 83333333333333,
						"locktime": 1683093491
					},
					{
						"amount": 83333333333333,
						"locktime": 1685771891
					},
					{
						"amount": 83333333333333,
						"locktime": 1688450291
					}
				]
			}
		],
		"startTime": 1658987891,
		"initialStakeDuration": 31536000,
		"initialStakeDurationOffset": 54000,
		"initialStakedFunds": [
			"X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd"
		],
		"initialStakers": [
			{
				"nodeID": "NodeID-21XvgoKWToLv8m2awpP6pjmABYDXWvvu1",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 1000000
			},
			{
				"nodeID": "NodeID-2DaDg8ySpZh4G3pYQDoQL6fBpr6kvAEn7",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 500000
			},
			{
				"nodeID": "NodeID-2cGXDwRBQasgadR7Q1cNLPrTu1CcBbbPg",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 250000
			},
			{
				"nodeID": "NodeID-3cR7XPe9cPjXyxe8xSGeHKtZUBW3EY5E9",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 125000
			},
			{
				"nodeID": "NodeID-4qGU3jCCskkgD23M5Phs58idoss88gACr",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 62500
			},
			{
				"nodeID": "NodeID-51tYGGYJQhkXj7korVVmF6dFjapPUkpx2",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 31250
			},
			{
				"nodeID": "NodeID-6i7bjwpMLjqAh493mVgjSGeWePD5psUmZ",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-6mYSHMtaKhapYefeS33oJtwJRKXqPM1dr",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-BiG4fTni2A6erA9TD4L867dCzL1ajH9Pz",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-CV7P79ttAXb8vqyc5QoUVxcMrVX4J21Y1",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-FG5jysE61HB8fVg3NucEmX5sXgs5sTpKn",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-GPtuetLJGjtpwoGcvzebj9KXgemhEcAuY",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-HSPLkj13MprfgpEZGYR7Dpm9ptS5m6miV",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-K5yDuQpynevLJWK1iu64ukA9UX566d1Ns",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-PMRrRcuXfSjBYcs1EBUayvgcWzfaQUt3p",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-PkUiWb8rf9Yh6twJr5RQbMhp1JZpj4W25",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-Q8RhVnz4JeRg3s5dQRqnQuc5H4v6Zwrk8",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-QGrgixzuznapYA5LeJB7RvsNWvVMwmofm",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-cRLt953CsEA8Hs6mhSzeawsbb335zziH",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-uhXjVs6ugtw4cStoe4tJa523Bb5kjMQg",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-8fWznLr2SqE2BQoBWKmeGRjX5yNqUhy8T",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-KTmn81w4WWwHxV1s19GmL1i2ygYqXjhtR",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-5cTEiXM4igt3xL9nnX1L3QvEdCAPafTUy",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-3YaV4kst8K3VXwyZpr1XtdLyKZbFQBnRs",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			},
			{
				"nodeID": "NodeID-4mXA7qhwuAsSnco38ZboLc5q6UFM1x8Pf",
				"rewardAddress": "X-metal1432ng5f83lf99c8jqykfksqccq6nrwwnqdrdzd",
				"delegationFee": 20000
			}
		],
		"cChainGenesis": "{\"config\":{\"chainId\":381931,\"homesteadBlock\":0,\"daoForkBlock\":0,\"daoForkSupport\":true,\"eip150Block\":0,\"eip150Hash\":\"0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0\",\"eip155Block\":0,\"eip158Block\":0,\"byzantiumBlock\":0,\"constantinopleBlock\":0,\"petersburgBlock\":0,\"istanbulBlock\":0,\"muirGlacierBlock\":0,\"apricotPhase1BlockTimestamp\":0,\"apricotPhase2BlockTimestamp\":0,\"apricotPhase3BlockTimestamp\":0,\"apricotPhase4BlockTimestamp\":0,\"apricotPhase5BlockTimestamp\":0},\"nonce\":\"0x0\",\"timestamp\":\"0x0\",\"extraData\":\"0x00\",\"gasLimit\":\"0x5f5e100\",\"difficulty\":\"0x0\",\"mixHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\",\"coinbase\":\"0x0000000000000000000000000000000000000000\",\"alloc\":{\"0100000000000000000000000000000000000000\":{\"code\":\"0x7300000000000000000000000000000000000000003014608060405260043610603d5760003560e01c80631e010439146042578063b6510bb314606e575b600080fd5b605c60048036036020811015605657600080fd5b503560b1565b60408051918252519081900360200190f35b818015607957600080fd5b5060af60048036036080811015608e57600080fd5b506001600160a01b03813516906020810135906040810135906060013560b6565b005b30cd90565b836001600160a01b031681836108fc8690811502906040516000604051808303818888878c8acf9550505050505015801560f4573d6000803e3d6000fd5b505050505056fea26469706673582212201eebce970fe3f5cb96bf8ac6ba5f5c133fc2908ae3dcd51082cfee8f583429d064736f6c634300060a0033\",\"balance\":\"0x0\"}},\"number\":\"0x0\",\"gasUsed\":\"0x0\",\"parentHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\"}",
		"message": "Always act in good faith with full transparency and accountability"
	}`

	// MainnetParams are the params used for mainnet
	MainnetParams = Params{
		TxFeeConfig: TxFeeConfig{
			TxFee:                 units.MilliAvax,
			CreateAssetTxFee:      10 * units.MilliAvax,
			CreateSubnetTxFee:     1 * units.Avax,
			CreateBlockchainTxFee: 1 * units.Avax,
		},
		StakingConfig: StakingConfig{
			UptimeRequirement: .8, // 80%
			MinValidatorStake: 2 * units.KiloAvax,
			MaxValidatorStake: 3 * units.MegaAvax,
			MinDelegatorStake: 25 * units.Avax,
			MinDelegationFee:  20000, // 2%
			MinStakeDuration:  2 * 7 * 24 * time.Hour,
			MaxStakeDuration:  365 * 24 * time.Hour,
			RewardConfig: reward.Config{
				MaxConsumptionRate: .12 * reward.PercentDenominator,
				MinConsumptionRate: .10 * reward.PercentDenominator,
				MintingPeriod:      365 * 24 * time.Hour,
				SupplyCap:          666666666 * units.Avax,
			},
		},
	}
)
