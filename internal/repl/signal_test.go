// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

func TestNewSignalHandler(t *testing.T) {
	h := NewSignalHandler()
	if h == nil {
		t.Fatal("NewSignalHandler returned nil")
	}
}

func TestSignalHandlerStop(t *testing.T) {
	h := NewSignalHandler()
	// Stop should not panic on a handler that hasn't started
	h.Stop()
}

func TestSignalHandlerStartExecutesCallback(t *testing.T) {
	h := NewSignalHandler()
	defer h.Stop()

	var wg sync.WaitGroup
	wg.Add(1)

	h.Start(func() {
		wg.Done()
	})

	// Send SIGINT to ourselves to trigger the handler
	pid := os.Getpid()
	if err := syscall.Kill(pid, syscall.SIGINT); err != nil {
		t.Skipf("cannot send signal: %v", err)
	}

	// Wait for callback with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Callback was executed
	case <-time.After(2 * time.Second):
		t.Error("callback was not executed within timeout")
	}
}

func TestSignalHandlerStartCallbackRunsInGoroutine(t *testing.T) {
	h := NewSignalHandler()
	defer h.Stop()

	started := make(chan struct{})
	finished := make(chan struct{})

	h.Start(func() {
		close(started)
		// Simulate some work
		time.Sleep(100 * time.Millisecond)
		close(finished)
	})

	// Send SIGINT to trigger
	pid := os.Getpid()
	if err := syscall.Kill(pid, syscall.SIGINT); err != nil {
		t.Skipf("cannot send signal: %v", err)
	}

	// Start() should return immediately (goroutine)
	select {
	case <-started:
		// Good, callback started
	case <-time.After(2 * time.Second):
		t.Error("callback goroutine did not start")
	}

	// Wait for it to finish
	select {
	case <-finished:
		// Good
	case <-time.After(2 * time.Second):
		t.Error("callback goroutine did not finish")
	}
}

func TestSignalHandlerSingleShot(t *testing.T) {
	h := NewSignalHandler()

	var callbackCount atomic.Int32
	h.Start(func() {
		callbackCount.Add(1)
	})

	// First signal should trigger callback
	pid := os.Getpid()
	if err := syscall.Kill(pid, syscall.SIGINT); err != nil {
		t.Skipf("cannot send signal: %v", err)
	}
	time.Sleep(200 * time.Millisecond)

	// The handler is single-shot: after consuming one signal the goroutine exits.
	// Stop() unregisters the signal channel.
	h.Stop()

	if callbackCount.Load() != 1 {
		t.Errorf("expected callback count 1, got %d", callbackCount.Load())
	}
}
