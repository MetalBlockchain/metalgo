// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package genesis

import (
	"github.com/MetalBlockchain/metalgo/utils/constants"
	"github.com/MetalBlockchain/metalgo/utils/sampler"
)

// getIPs returns the beacon IPs for each network
func getIPs(networkID uint32) []string {
	switch networkID {
	case constants.MainnetID:
		return []string{
			"3.69.62.47:9651",
			"3.127.157.125:9651",
			"18.158.6.202:9651",
			"18.195.232.200:9651",
			"18.198.53.50:9651",
			"44.210.42.105:9651",
			"44.210.42.121:9651",
			"44.210.42.153:9651",
			"44.210.42.120:9651",
			"44.210.42.14:9651",
			"35.163.168.170:9651",
			"34.217.159.211:9651",
			"34.218.36.65:9651",
			"54.190.49.44:9651",
			"44.225.134.136:9651",
			"18.176.244.200:9651",
			"35.72.98.123:9651",
			"3.114.83.143:9651",
			"13.115.26.228:9651",
			"52.194.127.164:9651",
			"43.205.1.217:9651",
			"43.205.34.97:9651",
			"65.2.133.10:9651",
			"13.126.110.239:9651",
			"13.233.75.203:9651",
		}
	case constants.TahoeID:
		return []string{
			"18.158.206.92:9651",
			"18.157.229.117:9651",
			"3.69.225.81:9651",
			"3.73.114.4:9651",
			"3.74.175.192:9651",
			"34.196.246.54:9651",
			"54.174.105.14:9651",
			"34.206.150.45:9651",
			"34.194.108.44:9651",
			"44.195.208.173:9651",
			"34.212.255.6:9651",
			"54.71.4.64:9651",
			"44.226.132.246:9651",
			"52.43.222.212:9651",
			"54.189.216.170:9651",
			"18.178.22.72:9651",
			"35.75.87.63:9651",
			"54.65.199.213:9651",
			"18.177.175.140:9651",
			"3.113.255.59:9651",
		}
	default:
		return nil
	}
}

// getNodeIDs returns the beacon node IDs for each network
func getNodeIDs(networkID uint32) []string {
	switch networkID {
	case constants.MainnetID:
		return []string{
			"NodeID-2cGXDwRBQasgadR7Q1cNLPrTu1CcBbbPg",
			"NodeID-2DaDg8ySpZh4G3pYQDoQL6fBpr6kvAEn7",
			"NodeID-3cR7XPe9cPjXyxe8xSGeHKtZUBW3EY5E9",
			"NodeID-4qGU3jCCskkgD23M5Phs58idoss88gACr",
			"NodeID-6i7bjwpMLjqAh493mVgjSGeWePD5psUmZ",
			"NodeID-6mYSHMtaKhapYefeS33oJtwJRKXqPM1dr",
			"NodeID-21XvgoKWToLv8m2awpP6pjmABYDXWvvu1",
			"NodeID-51tYGGYJQhkXj7korVVmF6dFjapPUkpx2",
			"NodeID-BiG4fTni2A6erA9TD4L867dCzL1ajH9Pz",
			"NodeID-cRLt953CsEA8Hs6mhSzeawsbb335zziH",
			"NodeID-CV7P79ttAXb8vqyc5QoUVxcMrVX4J21Y1",
			"NodeID-FG5jysE61HB8fVg3NucEmX5sXgs5sTpKn",
			"NodeID-GPtuetLJGjtpwoGcvzebj9KXgemhEcAuY",
			"NodeID-HSPLkj13MprfgpEZGYR7Dpm9ptS5m6miV",
			"NodeID-K5yDuQpynevLJWK1iu64ukA9UX566d1Ns",
			"NodeID-PkUiWb8rf9Yh6twJr5RQbMhp1JZpj4W25",
			"NodeID-PMRrRcuXfSjBYcs1EBUayvgcWzfaQUt3p",
			"NodeID-Q8RhVnz4JeRg3s5dQRqnQuc5H4v6Zwrk8",
			"NodeID-QGrgixzuznapYA5LeJB7RvsNWvVMwmofm",
			"NodeID-uhXjVs6ugtw4cStoe4tJa523Bb5kjMQg",
			"NodeID-5cTEiXM4igt3xL9nnX1L3QvEdCAPafTUy",
			"NodeID-3YaV4kst8K3VXwyZpr1XtdLyKZbFQBnRs",
			"NodeID-4mXA7qhwuAsSnco38ZboLc5q6UFM1x8Pf",
			"NodeID-8fWznLr2SqE2BQoBWKmeGRjX5yNqUhy8T",
			"NodeID-KTmn81w4WWwHxV1s19GmL1i2ygYqXjhtR",
		}
	case constants.TahoeID:
		return []string{
			"NodeID-EiRsmffiu3t1MFnhjrrFJEJGqfdtnkEvc",
			"NodeID-6QRpj3qBVPzNwdD2YewfB66Aa9Z9pv243",
			"NodeID-EW9JuwPmgSdt5SXkePz1JnfW3x9zzsqkd",
			"NodeID-AGXv6QpCY4m8Qmug7A9pDYrSSxqdbJAZS",
			"NodeID-3BPMUK91rSbiG1mdBmpissDqTmVFAfc8J",
			"NodeID-B73oENBBKtquFNWkjPfjN9F2pHZp5r16g",
			"NodeID-GJhvteVE4Rn7unxmzZrTwpUmuHvT4ioJU",
			"NodeID-CfsdUtCrAB6grr8kBSLiHSdrKe993Gwvg",
			"NodeID-Ew9uKh7xnyEiZHU1xwi5MMVjHsxkmPScN",
			"NodeID-3uC5oYpq5hqn8Ky46qT13pEYK6avjZY6T",
			"NodeID-GmRnQWwWNjeZ2g7hbaHabHMgpcUABZzXV",
			"NodeID-HSxLZv4zNkRF9vH2xYQW9zkLwtyQdjnne",
			"NodeID-H6xhpvuVgKEXjTYbfAH629SnoZFLtsTMW",
			"NodeID-2eypiRnhHyoZe94PBxbwu1TBDNAP9Te7e",
			"NodeID-AHxkNdDLHMY7x3xFoGXetZgQVfQK9mHpz",
			"NodeID-BJYAmxn5SMHaz3J61fN97r7sAiNGp1o7H",
			"NodeID-971kRywJGDu651QBGDJDPSDyth5NTb2Yy",
			"NodeID-7XmnU74M3APxBV6WLZS51NH9r7KhNoCTr",
			"NodeID-GGCBNGeqoLbHp6jhc4wtn5o92Eszj9HZt",
			"NodeID-M1ivUWVi9Ad2HxJcNcXTwoSfimLBhz8SK",
		}
	default:
		return nil
	}
}

// SampleBeacons returns the some beacons this node should connect to
func SampleBeacons(networkID uint32, count int) ([]string, []string) {
	ips := getIPs(networkID)
	ids := getNodeIDs(networkID)

	if numIPs := len(ips); numIPs < count {
		count = numIPs
	}

	sampledIPs := make([]string, 0, count)
	sampledIDs := make([]string, 0, count)

	s := sampler.NewUniform()
	s.Initialize(uint64(len(ips)))
	indices, _ := s.Sample(count)
	for _, index := range indices {
		sampledIPs = append(sampledIPs, ips[int(index)])
		sampledIDs = append(sampledIDs, ids[int(index)])
	}

	return sampledIPs, sampledIDs
}
