# RPN (Postfix Notation) Stack Calculator Implementation Plan

## Context

The `perc` project is a percentage calculator that currently supports three formats:
- `20% of 150`
- `30 is what % of 150`
- `30 is 20% of what`

Users want to extend it to support postfix notation (Reverse Polish Notation) stack-based calculations like `3 4 + 4 4 - *`, with variable assignments and reuse capabilities.

## Requirements

### Core Features
- **Postfix Notation Parser**: Parse space-separated tokens where numbers are pushed to a stack and operators pop operands and push results
- **Arithmetic Operations**: Addition (+), subtraction (-), multiplication (*), division (/), power (^), modulo (%)
- **Variable Support**:
  - Assign: `varname value =` stores value in variable
  - Reuse: `varname` pushes stored value onto stack
  - Delete: `varname d` removes variable
  - List: `vars` shows all variables
  - Clear: `clear` removes all variables
- **Stack Inspection**: `dup` (duplicate top), `swap` (swap top two), `pop` (remove top), `show` (print stack)
- **Error Handling**: Division by zero, invalid operators, insufficient operands, undefined variables
- **Input Methods**: Support both `perc calc 3 4 +` and `perc rpn 3 4 +` syntax

## Task Structure

### Core Implementation (Tasks 401-403)
| Task ID | Description | Dependencies |
|---------|-------------|--------------|
| 401 | Create `internal/rpn/variables.go` - Variable storage and management | None |
| 402 | Create `internal/rpn/operations.go` - Operator implementations and stack manipulation | 401 |
| 403 | Create `internal/rpn/rpn.go` - RPN parser and evaluator | 401, 402 |

### Integration (Tasks 404-408)
| Task ID | Description | Dependencies |
|---------|-------------|--------------|
| 404 | Add `ParseRPN()` to `internal/calculator/calculator.go` | 403 |
| 405 | Update `cmd/perc/main.go` - Add calc/rpn subcommands | 404 |
| 406 | Update `internal/repl/repl.go` - Handle RPN input | 405 |
| 407 | Update `internal/repl/commands.go` - Add rpn command | 406 |
| 408 | Update `Magefile.go` - Add RPN() target | 407 |

## Implementation Approach

### File Structure
```
internal/
├── calculator/
│   ├── calculator.go          # Existing - add RPN support
│   └── calculator_test.go     # Add RPN tests
├── rpn/
│   ├── rpn.go                 # New: RPN parser and evaluator
│   ├── operations.go          # New: Operator implementations
│   └── variables.go           # New: Variable storage and management
└── repl/
    ├── repl.go                # Update to support RPN
    └── commands.go            # Add rpn command

cmd/perc/
└── main.go                    # Add calc/rpn subcommand support
```

### Key Components

**`internal/rpn/rpn.go`**:
- `ParseAndEvaluate(input string) (string, error)` - Main entry point
- Tokenize input into numbers, operators, variables
- Execute RPN evaluation using stack

**`internal/rpn/operations.go`**:
- Operator functions: add, subtract, multiply, divide, power, modulo
- Stack manipulation: dup, swap, pop, show
- Error handling for invalid operations

**`internal/rpn/variables.go`**:
- Variable storage map
- `SetVariable(name string, value float64)` - Assign: `name value =`
- `GetVariable(name string)` - Retrieve variable value
- `DeleteVariable(name string)` - Delete: `name d`
- `ListVariables()` - List all: `vars`
- `ClearVariables()` - Clear all: `clear`

**`internal/calculator/calculator.go`**:
- Add `ParseRPN(input string) (string, error)` function
- Integrate with existing `Parse()` to detect RPN format

**`cmd/perc/main.go`**:
- Add `calc` and `rpn` subcommand support
- Route to appropriate parser based on command

**`internal/repl/repl.go`**:
- Update executor to handle RPN input
- Update commands.go to include `rpn` as built-in command

**`Magefile.go`**:
- Add `RPN()` target for testing

## Usage Examples

```bash
# Basic arithmetic
perc calc 3 4 +           # → 7
perc calc 3 4 + 4 4 - *   # → 0

# Power and modulo
perc calc 2 3 ^           # → 8
perc calc 10 3 %          # → 1

# Variables
perc calc x 5 = x x +     # → 10
perc calc pi 3.14159 = pi 2 *   # → 6.28318

# Stack operations
perc calc 1 2 3 dup       # → 1 2 3 3
perc calc 1 2 swap        # → 2 1
perc calc 1 2 3 pop       # → 1 2

# Variable management
perc calc vars            # List all variables
perc calc x d             # Delete variable x
perc calc clear           # Clear all variables
```

## Testing Strategy

### Unit Tests
- Each operation function (add, sub, mul, div, pow, mod)
- Variable operations (set, get, delete, list, clear)
- RPN tokenization and evaluation

### Integration Tests
- Full RPN evaluation with mixed operations
- Variable assignment and reuse
- Error cases (division by zero, undefined variables, insufficient operands)

### Manual Testing
1. `perc calc 3 4 +` - should output `7`
2. `perc calc 3 4 + 4 4 - *` - should output `0`
3. `perc calc 2 3 ^` - should output `8`
4. `perc calc x 5 = x x +` - should output `10`
5. `perc calc vars` - should list variables
6. `perc calc clear` - should clear all variables
7. `mage rpn` - verify Mage integration

## Task UUIDs for Reference

| Task ID | UUID |
|---------|------|
| 401 | (see `task 401 _uuid`) |
| 402 | (see `task 402 _uuid`) |
| 403 | (see `task 403 _uuid`) |
| 404 | (see `task 404 _uuid`) |
| 405 | (see `task 405 _uuid`) |
| 406 | (see `task 406 _uuid`) |
| 407 | (see `task 407 _uuid`) |
| 408 | (see `task 408 _uuid`) |
