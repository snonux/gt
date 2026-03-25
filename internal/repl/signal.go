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

// Start starts the signal handler goroutine with the given callback.
// When SIGINT is received, the callback function is executed in the goroutine.
// The function blocks until Stop is called.
//
// callback: the function to execute when SIGINT is received
// Returns: no return value; executes callback in a separate goroutine
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
