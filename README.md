# gt

A simple AI-engineered command-line percentage calculator written in Go. No frontier AI models from Claude, OpenAI, Google, etc., were used for this project. The ones used were:

* Qwen 3 Coder Next
* GPT OSS 120b
* Nemotron 3 Super
* Gemma 4 Dense (31B)
* My human brain

This is a toy project created to experience with local LLMs and how good they are at their jobs. 

```
     ┌────────────────────────────────┐
     │       G T Calculator           │
     │                                │
     │       (o.o)                    │
     │      /|   | \   5  6  *        │
     │     / |   |  \  _______        │
     │    /  |   |   \ / G T  \       │
     │   (   |   |    )%  +   |       │
     │    \  |   |  / | RPN   |       │
     │     \ |   | /  | Calc  |       │
     │      \|___|/   |_______|       │
     │        Result: 30              │
     └────────────────────────────────┘
```

## Installation

```bash
go install codeberg.org/snonux/gt/cmd/gt@latest
```

Or using mage:

```bash
mage install
```

## Usage

`gt` supports various percentage calculation formats and RPN (Reverse Polish Notation) stack calculations.

### Percentage Calculations

#### Calculate X% of Y

```bash
gt '20% of 150'
# Output:
# 20.00% of 150.00 = 30.00
#   Steps: (20.00 / 100) * 150.00 = 0.20 * 150.00 = 30.00

gt 'what is 20% of 150'
# Output:
# 20.00% of 150.00 = 30.00
#   Steps: (20.00 / 100) * 150.00 = 0.20 * 150.00 = 30.00
```

#### Find what percentage X is of Y

```bash
gt '30 is what % of 150'
# Output:
# 30.00 is 20.00% of 150.00
#   Steps: (30.00 / 150.00) * 100 = 0.20 * 100 = 20.00%
```

#### Find the whole when X is Y% of it

```bash
gt '30 is 20% of what'
# Output:
# 30.00 is 20.00% of 150.00
#   Steps: (30.00 / 20.00) * 100 = 1.50 * 100 = 150.00
```

### RPN (Reverse Polish Notation) Calculations

RPN (postfix notation) uses a stack-based approach where operators follow their operands. No parentheses needed!

#### Basic Arithmetic

```bash
gt '3 4 +'           # 3 + 4 = 7
# → 7

gt '3 4 -'           # 3 - 4 = -1
# → -1

gt '5 6 *'           # 5 * 6 = 30
# → 30

gt '20 4 /'          # 20 / 4 = 5
# → 5

gt '2 3 ^'           # 2^3 = 8
# → 8

gt '10 3 %'          # 10 % 3 = 1 (modulo)
# → 1
```

#### Expression Chaining

```bash
gt '3 4 + 4 4 - *'   # (3+4) * (4-4) = 0
# → 0

gt '1 2 + 3 *'       # (1+2) * 3 = 9
# → 9
```

#### Variables

```bash
gt 'x 5 ='           # Assign x = 5
# → x = 5

gt 'x 5 = x x +'     # x + x = 10
# → 10

gt 'pi 3.14159 = pi 2 *'  # 2 * π
# → 6.28318

# Note: Variable assignment works in bare mode (e.g., "gt 'x 5 ='").
# gt 'x 5 = x x +'  (works)
```

#### Built-in Constants

The RPN calculator includes several built-in mathematical constants:

```bash
gt 'pi'              # π (pi) - ratio of circumference to diameter
# → 3.141592654

gt 'e'               # Euler's number (base of natural logarithm)
# → 2.718281828

gt 'phi'             # Golden ratio
# → 1.618033989

gt 'sqrt2'           # Square root of 2
# → 1.414213562

gt 'inf'             # Positive infinity
# → +Inf

gt 'nan'             # Not a Number
# → NaN
```

Greek letter variants are also supported:
```bash
gt 'π'               # Same as pi
# → 3.141592654

gt 'φ'               # Same as phi
# → 1.618033989

gt '√2'              # Same as sqrt2
# → 1.414213562
```

Constants can be used in expressions:
```bash
gt 'pi 2 *'          # 2π
# → 6.283185307

gt 'e phi +'         # e + φ
# → 4.336315817

gt 'sqrt2 2 ^'       # (√2)² = 2
# → 2
```

**Note:** If you assign a variable with the same name as a constant (e.g., `pi = 3`), the variable takes precedence and the constant is no longer accessible until you clear the variables.

#### Constant Management

List all constants:
```bash
gt 'constants'       # List all constants
# e = 2.718281828
# pi = 3.141592654
# phi = 1.618033989
# sqrt2 = 1.414213562
# ...
```

Clear user-defined constants (built-in constants are preserved):
```bash
gt 'clearconstants'  # Clear user-defined constants
# All constants cleared
```

Built-in constants cannot be deleted individually, but you can list them and clear variables that might override them:

```bash
gt 'vars'            # List all variables (not constants)
# x = 5

# To access constants again after overriding them:
gt 'clear'           # Clear all variables
gt 'pi'              # Now returns the constant π
# → 3.141592654
```

```bash
gt 'vars'            # List all variables
# x = 5

gt 'delete x'        # Delete a specific variable
# Variable removed

gt 'clear'           # Clear all variables
# All variables cleared
```

### Working with Variables

Variables persist across commands in REPL mode but are cleared when exiting. In bare mode (single command), variables are only available within that command's execution context.

Example:
```bash
# In bare mode, variables don't persist between commands
gt 'x 5 ='           # Assign x = 5 (in this command only)
# → x = 5

gt x               # x is not defined in this separate command
# Error: variable not found

# In REPL mode, variables persist
> x 5 =            # Assign x = 5
> x                # x is still 5
5
> clear            # Clear all variables
> x                # x is now undefined
# Error: variable not found
```

#### Stack Operations

```bash
gt '1 2 3 dup'       # Duplicate top value
# → 1 2 3 3

gt '1 2 swap'        # Swap top two values
# → 2 1

gt '1 2 3 pop'       # Remove top value
# → 1 2

gt '1 2 3 show'      # Show stack without modifying
# → 1 2 3
```

### REPL Mode Notes

In REPL mode, RPN operations maintain persistent state between commands. This allows you to build up values on the stack across multiple commands.

Example REPL session:
```
> 2 3 4 +        # Push 2, 3, 4; add last two
2 7
> +                  # Add top two: 2 + 7 = 9
9
> 5 *                # Multiply by 5: 9 * 5 = 45
45
```

To show the current stack without modifying it:
```
> show               # Show current stack state
45
```

## Boolean-to-Number Coercion

Boolean values are automatically coerced to numbers when used in arithmetic operations:
- `true` is treated as `1`
- `false` is treated as `0`

This enables mixed boolean-numeric expressions:

```bash
gt 5 3 == 1 +      # 5 == 3 is false (0), 0 + 1 = 1
# → 1

gt 0 false +       # false is 0, 0 + 0 = 0
# → 0

gt true 2 *        # true is 1, 1 * 2 = 2
# → 2

gt 9 3 > 4 5 < +   # 9 > 3 is true (1), 4 < 5 is true (1), 1 + 1 = 2
# → 2
```

Note: The boolean result is shown as `true`/`false` when printed, but when used as an operand it behaves as the corresponding numeric value.

## Hyper Operators

Hyper operators work on all values on the stack simultaneously:

```bash
gt '1 2 3 4 5 [+]'    # Sum all: 1+2+3+4+5 = 15
# → 15

gt '2 3 4 [*]'        # Multiply all: 2*3*4 = 24
# → 24

gt '10 3 2 [-]'       # 10 - 3 - 2 = 5
# → 5

gt '100 5 2 [/]'      # 100 / 5 / 2 = 10
# → 10

gt '2 3 2 [^]'        # (2^3)^2 = 64
# → 64

gt '100 7 3 [%]'      # 100 % 7 % 3 = 2
# → 2
```

## Building

Using mage:

```bash
mage build
```

Or using go directly:

```bash
go build -o gt ./cmd/gt
```

## Testing

```bash
mage test
```

Or for RPN-specific tests:

```bash
mage testRPN
```

## Rational Number Mode (Optional)

The calculator supports precise rational number calculations using Go's `*big.Rat` type. By default, calculations use float64 for performance.

### Enabling Rational Mode

In REPL mode, you can switch between float64 and rational number modes:

```
> rat on           # Enable rational number mode
Rational mode enabled

> rat off          # Disable rational number mode (default)
Rational mode disabled (using float64)

> rat toggle       # Switch to the other mode
Rational mode enabled
```

When rational mode is enabled:
- Results are calculated with arbitrary precision
- Output is displayed as a decimal approximation
- Use `rat off` to return to standard float64 calculations

## License

See LICENSE file for details.
