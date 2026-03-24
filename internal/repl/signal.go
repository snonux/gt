package repl

import (
	"os"
	"os/signal"
	"syscall"
)

// SignalHandler manages signal handling for the REPL.
type SignalHandler struct {
	sigChan chan os.Signal
}

// NewSignalHandler creates a new signal handler that listens for SIGINT.
func NewSignalHandler() *SignalHandler {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	return &SignalHandler{
		sigChan: sigChan,
	}
}

// Start starts the signal handler goroutine with the given callback.
func (s *SignalHandler) Start(callback func()) {
	go func() {
		<-s.sigChan
		callback()
	}()
}

// Stop stops the signal handler by unregistering signals.
func (s *SignalHandler) Stop() {
	signal.Stop(s.sigChan)
}
