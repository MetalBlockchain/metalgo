// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package timeout

import (
	"sync"
	"testing"
	"time"

	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/message"
	"github.com/MetalBlockchain/metalgo/snow/networking/benchlist"
	"github.com/MetalBlockchain/metalgo/utils/timer"
	"github.com/prometheus/client_golang/prometheus"
)

func TestManagerFire(t *testing.T) {
	benchlist := benchlist.NewNoBenchlist()
	manager, err := NewManager(
		&timer.AdaptiveTimeoutConfig{
			InitialTimeout:     time.Millisecond,
			MinimumTimeout:     time.Millisecond,
			MaximumTimeout:     10 * time.Second,
			TimeoutCoefficient: 1.25,
			TimeoutHalflife:    5 * time.Minute,
		},
		benchlist,
		"",
		prometheus.NewRegistry(),
	)
	if err != nil {
		t.Fatal(err)
	}
	go manager.Dispatch()

	wg := sync.WaitGroup{}
	wg.Add(1)

	manager.RegisterRequest(ids.NodeID{}, ids.ID{}, message.PullQuery, ids.GenerateTestID(), wg.Done)

	wg.Wait()
}

func TestManagerCancel(t *testing.T) {
	benchlist := benchlist.NewNoBenchlist()
	manager, err := NewManager(
		&timer.AdaptiveTimeoutConfig{
			InitialTimeout:     time.Millisecond,
			MinimumTimeout:     time.Millisecond,
			MaximumTimeout:     10 * time.Second,
			TimeoutCoefficient: 1.25,
			TimeoutHalflife:    5 * time.Minute,
		},
		benchlist,
		"",
		prometheus.NewRegistry(),
	)
	if err != nil {
		t.Fatal(err)
	}
	go manager.Dispatch()

	wg := sync.WaitGroup{}
	wg.Add(1)

	fired := new(bool)

	id := ids.GenerateTestID()
	manager.RegisterRequest(ids.NodeID{}, ids.ID{}, message.PullQuery, id, func() { *fired = true })

	manager.RegisterResponse(ids.NodeID{}, ids.ID{}, id, message.Get, 1*time.Second)

	manager.RegisterRequest(ids.NodeID{}, ids.ID{}, message.PullQuery, ids.GenerateTestID(), wg.Done)

	wg.Wait()

	if *fired {
		t.Fatalf("Should have cancelled the function")
	}
}
