// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"os"
	"os/signal"
	"syscall"
)

// SignalHandler manages signal handling for the REPL.
// It specifically listens for SIGINT (Ctrl+C) and executes a callback.
type SignalHandler struct {
	sigChan chan os.Signal
}

// NewSignalHandler creates a new signal handler that listens for SIGINT.
// It creates a buffered channel to receive signals.
//
// Returns a new SignalHandler instance
func NewSignalHandler() *SignalHandler {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	return &SignalHandler{
		sigChan: sigChan,
	}
}

// Start launches a goroutine that waits for SIGINT and then runs the callback.
// Start returns immediately; the goroutine blocks until a signal is received or
// Stop is called (which prevents further signals from being delivered).
//
// callback: the function to execute when SIGINT is received
func (s *SignalHandler) Start(callback func()) {
	go func() {
		<-s.sigChan
		callback()
	}()
}

// Stop stops the signal handler by unregistering signals.
// After calling Stop, the signal handler will no longer trigger the callback.
func (s *SignalHandler) Stop() {
	signal.Stop(s.sigChan)
}
