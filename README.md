# perc

A simple AI-engineered command-line percentage calculator written in Go.

## Installation

```bash
go install codeberg.org/snonux/perc/cmd/perc@latest
```

Or using mage:

```bash
mage install
```

## Usage

`perc` supports various percentage calculation formats and RPN (Reverse Polish Notation) stack calculations.

### Percentage Calculations

#### Calculate X% of Y

```bash
perc 20% of 150
# Output:
# 20.00% of 150.00 = 30.00
#   Steps: (20.00 / 100) * 150.00 = 0.20 * 150.00 = 30.00

perc what is 20% of 150
# Output:
# 20.00% of 150.00 = 30.00
#   Steps: (20.00 / 100) * 150.00 = 0.20 * 150.00 = 30.00
```

#### Find what percentage X is of Y

```bash
perc 30 is what % of 150
# Output:
# 30.00 is 20.00% of 150.00
#   Steps: (30.00 / 150.00) * 100 = 0.20 * 100 = 20.00%
```

#### Find the whole when X is Y% of it

```bash
perc 30 is 20% of what
# Output:
# 30.00 is 20.00% of 150.00
#   Steps: (30.00 / 20.00) * 100 = 1.50 * 100 = 150.00
```

### RPN (Reverse Polish Notation) Calculations

RPN (postfix notation) uses a stack-based approach where operators follow their operands. No parentheses needed!

#### Basic Arithmetic

```bash
perc calc 3 4 +           # 3 + 4 = 7
# → 7

perc calc 3 4 -           # 3 - 4 = -1
# → -1

perc calc 5 6 *           # 5 * 6 = 30
# → 30

perc calc 20 4 /          # 20 / 4 = 5
# → 5

perc calc 2 3 ^           # 2^3 = 8
# → 8

perc calc 10 3 %          # 10 % 3 = 1 (modulo)
# → 1
```

#### Expression Chaining

```bash
perc calc 3 4 + 4 4 - *   # (3+4) * (4-4) = 0
# → 0

perc calc 1 2 + 3 *       # (1+2) * 3 = 9
# → 9
```

#### Variables

```bash
perc calc x 5 =           # Assign x = 5
# → x = 5

perc calc x 5 = x x +     # x + x = 10
# → 10

perc calc pi 3.14159 = pi 2 *  # 2 * π
# → 6.28318

# Note: Variable assignment only works with calc/rpn subcommand:
# perc calc x 5 = x x +  (works)
# perc x 5 =             (won't work in bare mode - use "perc calc x 5 =")
```

#### Variable Management

```bash
perc calc vars            # List all variables
# x = 5

perc calc name d          # Delete variable
# Variable removed

perc calc clear           # Clear all variables
# All variables cleared
```

#### Stack Operations

```bash
perc calc 1 2 3 dup       # Duplicate top value
# → 1 2 3 3

perc calc 1 2 swap        # Swap top two values
# → 2 1

perc calc 1 2 3 pop       # Remove top value
# → 1 2

perc calc 1 2 3 show      # Show stack without modifying
# → 1 2 3
```

### REPL Mode Notes

In REPL mode, RPN operations maintain persistent state between commands. This allows you to build up values on the stack across multiple commands.

Example REPL session:
```
perc> rpn 2 3 4 +        # Push 2, 3, 4; add last two
2 7
perc> +                  # Add top two: 2 + 7 = 9
9
perc> 5 *                # Multiply by 5: 9 * 5 = 45
45
```

To show the current stack without modifying it:
```
perc> show               # Show current stack state
45
```

## Hyper Operators

Hyper operators work on all values on the stack simultaneously:

```bash
perc calc 1 2 3 4 5 [+]    # Sum all: 1+2+3+4+5 = 15
# → 15

perc calc 2 3 4 [*]        # Multiply all: 2*3*4 = 24
# → 24

perc calc 10 3 2 [-]       # 10 - 3 - 2 = 5
# → 5

perc calc 100 5 2 [/]      # 100 / 5 / 2 = 10
# → 10

perc calc 2 3 2 [^]        # (2^3)^2 = 64
# → 64

perc calc 100 7 3 [%]      # 100 % 7 % 3 = 2
# → 2
```

## Building

Using mage:

```bash
mage build
```

Or using go directly:

```bash
go build -o perc ./cmd/perc
```

## Testing

```bash
mage test
```

Or for RPN-specific tests:

```bash
mage testRPN
```

## License

See LICENSE file for details.
