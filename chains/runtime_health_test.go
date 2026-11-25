// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package chains

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MetalBlockchain/metalgo/api/health"
	"github.com/stretchr/testify/require"
)

func TestShutdownOnRepeatedVMHealthFailures(t *testing.T) {
	vmErr := errors.New("VM unhealthy")
	checkErr := vmErr
	shutdown := make(chan int, 2)
	checker := &shutdownOnUnhealthyChecker{
		checker: health.CheckerFunc(func(context.Context) (interface{}, error) {
			return nil, checkErr
		}),
		threshold: 3,
		shutdown:  func(code int) { shutdown <- code },
	}

	_, err := checker.HealthCheck(context.Background())
	require.ErrorIs(t, err, vmErr)
	checkErr = nil
	_, err = checker.HealthCheck(context.Background())
	require.NoError(t, err)
	checkErr = vmErr
	for range 2 {
		_, err = checker.HealthCheck(context.Background())
		require.ErrorIs(t, err, vmErr)
	}
	select {
	case <-shutdown:
		require.Fail(t, "node shut down before the failure threshold")
	default:
	}
	_, err = checker.HealthCheck(context.Background())
	require.ErrorIs(t, err, vmErr)
	select {
	case code := <-shutdown:
		require.Equal(t, 1, code)
	case <-time.After(time.Second):
		require.Fail(t, "node shutdown was not requested after consecutive VM health failures")
	}
	for range 2 {
		_, err = checker.HealthCheck(context.Background())
		require.ErrorIs(t, err, vmErr)
	}
	select {
	case <-shutdown:
		require.Fail(t, "node shutdown was requested more than once")
	default:
	}
}

func TestConcurrentVMHealthFailuresRequestShutdownOnce(t *testing.T) {
	var checks atomic.Int32
	shutdown := make(chan int, 2)
	checker := &shutdownOnUnhealthyChecker{
		checker: health.CheckerFunc(func(context.Context) (interface{}, error) {
			checks.Add(1)
			return nil, errors.New("VM unhealthy")
		}),
		threshold: 3,
		shutdown:  func(code int) { shutdown <- code },
	}

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = checker.HealthCheck(context.Background())
		}()
	}
	wg.Wait()
	require.EqualValues(t, 10, checks.Load())
	select {
	case code := <-shutdown:
		require.Equal(t, 1, code)
	case <-time.After(time.Second):
		require.Fail(t, "concurrent health failures did not request node shutdown")
	}
	select {
	case <-shutdown:
		require.Fail(t, "concurrent health failures requested shutdown more than once")
	default:
	}
}
