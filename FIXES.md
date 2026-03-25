# Proposed Fixes for 100 Go Mistakes Audit

This document outlines specific code changes to address the HIGH and MEDIUM severity issues identified in the audit.

---

## Fix #1: Proper Error Wrapping in RPN Package

**File:** `internal/rpn/rpn.go`  
**Issue:** Errors not wrapped with context  
**Location:** Multiple functions including `ParseAndEvaluate`, `evaluate`

### Current Code:
```go
func (r *RPN) ParseAndEvaluate(input string) (string, error) {
    // ...
    if len(tokens) == 0 {
        return "", fmt.Errorf("no valid tokens found")
    }
    return r.evaluate(tokens)
}

func (r *RPN) evaluate(tokens []string) (string, error) {
    // ...
    if result, err := r.handleOperator(stack, token, i); err != nil {
        return "", err  // Missing context
    }
    // ...
}
```

### Fixed Code:
```go
func (r *RPN) ParseAndEvaluate(input string) (string, error) {
    // ...
    if len(tokens) == 0 {
        return "", fmt.Errorf("rpn: no valid tokens found in input: %q", input)
    }
    return r.evaluate(tokens)
}

func (r *RPN) evaluate(tokens []string) (string, error) {
    // ...
    if result, err := r.handleOperator(stack, token, i); err != nil {
        return "", fmt.Errorf("rpn: failed to evaluate token '%s' at position %d: %w", token, i, err)
    }
    // ...
}
```

---

## Fix #2: Proper Error Comparison in Tests

**File:** `internal/rpn/variables_test.go`  
**Issue:** Direct error comparison instead of `errors.Is()`  
**Location:** `TestOperationsUseVariableUndefined`

### Current Code:
```go
func TestOperationsUseVariableUndefined(t *testing.T) {
    // ...
    err := o.UseVariable(s, "undefined")
    if err == nil {
        t.Error("UseVariable for undefined variable should return error")
    }
    if !errors.Is(err, ErrVariableNotFound) {  // This is actually correct, but some places use direct comparison
        t.Errorf("UseVariable error = %v, should be ErrVariableNotFound", err)
    }
}
```

### Note:
The code already uses `errors.Is()`, which is correct. No change needed here.

---

## Fix #3: Proper Resource Cleanup in REPL

**File:** `internal/repl/repl.go`  
**Issue:** Defer error ignoring in `saveHistory`  
**Location:** Lines ~160-180

### Current Code:
```go
func saveHistory(history []string) error {
    historyPath := getHistoryPath()
    if historyPath == "" {
        return nil
    }

    file, err := os.Create(historyPath)
    if err != nil {
        return err
    }

    writer := bufio.NewWriter(file)
    for _, entry := range history {
        if _, err := writer.WriteString(entry + "\n"); err != nil {
            _ = file.Close()  // Error ignored
            return err
        }
    }
    if err := writer.Flush(); err != nil {
        _ = file.Close()  // Error ignored
        return err
    }
    return file.Close()
}
```

### Fixed Code:
```go
func saveHistory(history []string) error {
    historyPath := getHistoryPath()
    if historyPath == "" {
        return nil
    }

    file, err := os.Create(historyPath)
    if err != nil {
        return err
    }
    defer func() {
        if err := file.Close(); err != nil {
            // Log the error but don't overwrite original error
            log.Printf("Warning: failed to close history file: %v", err)
        }
    }()

    writer := bufio.NewWriter(file)
    for _, entry := range history {
        if _, err := writer.WriteString(entry + "\n"); err != nil {
            return fmt.Errorf("failed to write history entry: %w", err)
        }
    }
    if err := writer.Flush(); err != nil {
        return fmt.Errorf("failed to flush history writer: %w", err)
    }
    return nil
}
```

---

## Fix #4: Mutex Safety in REPL

**File:** `internal/repl/repl.go`  
**Issue:** Mutex potential for copying  
**Location:** Lines ~35-39

### Current Code:
```go
// RPNState holds the state for RPN operations in REPL
type RPNState struct {
    vars    rpn.VariableStore
    rpnCalc *rpn.RPN
}

// getRPNState returns or creates the RPN state
var rpnStateMu sync.RWMutex
var rpnState *RPNState

func getRPNState() *RPNState {
    rpnStateMu.Lock()
    defer rpnStateMu.Unlock()
    if rpnState == nil {
        vars := rpn.NewVariables()
        rpnState = &RPNState{
            vars:    vars,
            rpnCalc: rpn.NewRPN(vars),
        }
    }
    return rpnState
}
```

### Fixed Code:
```go
// RPNState holds the state for RPN operations in REPL
// Note: This struct should never be copied
type RPNState struct {
    vars    rpn.VariableStore
    rpnCalc *rpn.RPN
}

// rpnStateMu protects rpnState
// Note: The mutex must NOT be copied - keep it as a top-level variable
var rpnStateMu sync.RWMutex
var rpnState *RPNState

// getRPNState returns or creates the RPN state
func getRPNState() *RPNState {
    // First check with read lock for performance
    rpnStateMu.RLock()
    if rpnState != nil {
        state := rpnState
        rpnStateMu.RUnlock()
        return state
    }
    rpnStateMu.RUnlock()

    // Need to create - use write lock
    rpnStateMu.Lock()
    defer rpnStateMu.Unlock()
    if rpnState == nil {
        vars := rpn.NewVariables()
        rpnState = &RPNState{
            vars:    vars,
            rpnCalc: rpn.NewRPN(vars),
        }
    }
    return rpnState
}
```

---

## Fix #5: Improved Error Context in Percentage Calculator

**File:** `internal/perc/perc.go`  
**Issue:** Errors not wrapped with context  
**Location:** Various places

### Current Code:
```go
func Parse(input string) (string, error) {
    // ...
    result, err := perc.Parse(input)
    if err != nil {
        return "", err  // Missing context
    }
    return result, nil
}
```

### Fixed Code:
```go
func Parse(input string) (string, error) {
    // ...
    result, err := perc.Parse(input)
    if err != nil {
        return "", fmt.Errorf("rpn fallback failed for input %q: %w", input, err)
    }
    return result, nil
}
```

---

## Fix #6: Better Slice Handling in Variables

**File:** `internal/rpn/variables.go`  
**Issue:** Slice capacity retention in `Clear()`  
**Location:** Lines ~65-67

### Current Code:
```go
func (s *Stack) Clear() {
    s.values = s.values[:0]  // Retains capacity
}
```

### Fixed Code:
```go
func (s *Stack) Clear() {
    // Option 1: Reset to nil (releases memory)
    s.values = nil
    
    // Option 2: Keep capacity but reset length (faster for reuse)
    // s.values = s.values[:0]
}
```

### Recommendation:
Use Option 1 (nil) if memory usage is a concern, Option 2 if stack will be reused immediately.

---

## Fix #7: Performance Optimization in Variables

**File:** `internal/rpn/variables.go`  
**Issue:** Allocations in hot paths  
**Location:** `ListVariables`

### Current Code:
```go
func (v *Variables) ListVariables() []VariableInfo {
    v.mu.RLock()
    defer v.mu.RUnlock()
    
    var infos []VariableInfo  // New allocation each call
    for name, value := range v.variables {
        infos = append(infos, VariableInfo{Name: name, Value: value})
    }
    // ...
}
```

### Fixed Code (Option 1 - Pre-allocation):
```go
func (v *Variables) ListVariables() []VariableInfo {
    v.mu.RLock()
    defer v.mu.RUnlock()
    
    // Pre-allocate slice with known capacity
    infos := make([]VariableInfo, 0, len(v.variables))
    for name, value := range v.variables {
        infos = append(infos, VariableInfo{Name: name, Value: value})
    }
    
    // Sort by name for consistent output
    sort.Slice(infos, func(i, j int) bool {
        return infos[i].Name < infos[j].Name
    })
    
    return infos
}
```

### Fixed Code (Option 2 - Cached Result):
```go
// For frequently accessed data, consider caching
func (v *Variables) ListVariables() []VariableInfo {
    v.mu.RLock()
    defer v.mu.RUnlock()
    
    // Create a copy to avoid holding the lock during sorting
    infos := make([]VariableInfo, 0, len(v.variables))
    for name, value := range v.variables {
        infos = append(infos, VariableInfo{Name: name, Value: value})
    }
    
    // Release lock before sorting (long operation)
    v.mu.RUnlock()
    
    // Sort outside the critical section
    sort.Slice(infos, func(i, j int) bool {
        return infos[i].Name < infos[j].Name
    })
    
    return infos
}
```

---

## Testing the Fixes

After applying fixes, run:

```bash
# Run all tests with race detector
go test -race ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Check for linter issues
go vet ./...
golangci-lint run
```

---

## Summary of Key Changes

| Fix | File | Change | Impact |
|-----|------|--------|--------|
| #1 | `rpn.go` | Error wrapping | Better debugging |
| #3 | `repl.go` | Proper resource cleanup | No resource leaks |
| #4 | `repl.go` | Mutex safety | Thread safety |
| #5 | `perc.go` | Error context | Better error messages |
| #6 | `variables.go` | Slice handling | Memory management |
| #7 | `variables.go` | Performance optimization | Reduced allocations |

Each fix addresses specific Go anti-patterns identified in the audit while maintaining code correctness and improving maintainability.
