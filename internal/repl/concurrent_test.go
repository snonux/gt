package repl

import (
	"sync"
	"testing"

	"codeberg.org/snonux/gt/internal/rpn"
)

// TestConcurrentExecutor tests concurrent calls to defaultExecutor with fresh state
func TestConcurrentExecutor(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			vars := rpn.NewVariables()
			rpnCalc := rpn.NewRPN(vars)
			calculator := NewRPNCalculator(rpnCalc)
			rpl := &REPL{
				ttyChecker:    &TTYChecker{},
				historyMgr:    NewHistoryManager(".gt_history"),
				signalHandler: NewSignalHandler(),
				commandChain:  NewCommandChain(),
				rpnState:      &RPNState{vars: vars, calculator: calculator},
			}
			defaultExecutor(rpl, "20% of 150")
		}(i)
	}
	wg.Wait()
}

// TestConcurrentRPN tests concurrent inline RPN evaluation
func TestConcurrentRPN(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			vars := rpn.NewVariables()
			rpnCalc := rpn.NewRPN(vars)
			_, _ = rpnCalc.ParseAndEvaluate("3 4 +")
		}(i)
	}
	wg.Wait()
}

// TestConcurrentRatModeToggle tests concurrent rat mode toggles with fresh state
func TestConcurrentRatModeToggle(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			vars := rpn.NewVariables()
			rpnCalc := rpn.NewRPN(vars)
			calculator := NewRPNCalculator(rpnCalc)
			rpl := &REPL{
				ttyChecker:    &TTYChecker{},
				historyMgr:    NewHistoryManager(".gt_history"),
				signalHandler: NewSignalHandler(),
				commandChain:  NewCommandChain(),
				rpnState:      &RPNState{vars: vars, calculator: calculator},
			}
			defaultExecutor(rpl, "rat toggle")
		}(i)
	}
	wg.Wait()
}

// TestConcurrentExecutorAndRPN tests concurrent executor and RPN calls with fresh state
func TestConcurrentExecutorAndRPN(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(2)
		go func(id int) {
			defer wg.Done()
			vars := rpn.NewVariables()
			rpnCalc := rpn.NewRPN(vars)
			calculator := NewRPNCalculator(rpnCalc)
			rpl := &REPL{
				ttyChecker:    &TTYChecker{},
				historyMgr:    NewHistoryManager(".gt_history"),
				signalHandler: NewSignalHandler(),
				commandChain:  NewCommandChain(),
				rpnState:      &RPNState{vars: vars, calculator: calculator},
			}
			defaultExecutor(rpl, "20% of 150")
		}(i)
		go func(id int) {
			defer wg.Done()
			vars := rpn.NewVariables()
			rpnCalc := rpn.NewRPN(vars)
			_, _ = rpnCalc.ParseAndEvaluate("3 4 +")
		}(i)
	}
	wg.Wait()
}
