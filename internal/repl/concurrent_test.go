package repl

import (
	"sync"
	"testing"
)

func TestConcurrentExecutor(t *testing.T) {
	// Test concurrent calls to executor()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			executor("20% of 150")
		}(i)
	}
	wg.Wait()
}

func TestConcurrentRPN(t *testing.T) {
	// Test concurrent calls to runRPN()
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Ignore error return as the expression "3 4 +" should always succeed
			_, _ = runRPN("3 4 +")
		}(i)
	}
	wg.Wait()
}

func TestConcurrentRatModeToggle(t *testing.T) {
	// Test concurrent calls to executor() that change mode
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			executor("rat toggle")
		}(i)
	}
	wg.Wait()
}

func TestConcurrentExecutorAndRPN(t *testing.T) {
	// Test concurrent calls to executor() and runRPN()
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(2)
		go func(id int) {
			defer wg.Done()
			executor("20% of 150")
		}(i)
		go func(id int) {
			defer wg.Done()
			// Ignore error return as the expression "3 4 +" should always succeed
			_, _ = runRPN("3 4 +")
		}(i)
	}
	wg.Wait()
}
