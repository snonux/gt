# gt

A simple AI-engineered command-line percentage calculator written in Go. No frontier AI models from Claude, OpenAI, Google, etc., were used for this project. The ones used were:

* Qwen 3 Coder Next
* GPT OSS 120b
* Nemotron 3 Super
* Gemma 4 Dense (31B)
* Qwen 3.6 model variations
* My human brain

This is a toy project created to experience with local LLMs and how good they are at their jobs. 

![gt Logo](logo.svg)

## Installation

```bash
go install codeberg.org/snonux/gt/cmd/gt@latest
```

Or using mage:

```bash
mage install
```

[Mage](https://github.com/magefile/mage) is a Makefile replacement written in Go. It uses Go code for task definitions, providing better type safety and dependency management. The project includes a `magefile.go` with common development tasks like `build`, `test`, and `install`.

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

gt '2 3 ^'           # 2^3 = 8 (floating-point power)
# → 8

gt '2 10 **'          # 2^10 = 1024 (Fast binary exponentiation)
# → 1024

gt '2 -3 **'          # 2^(-3) = 0.125 (negative integer exponent)
# → 0.125

gt '5 0 **'           # 5^0 = 1 (zero exponent)
# → 1

gt '10 100 **'        # 10^100 (very large exponent - efficient)
# → 1e+100

gt '10 3 %'          # 10 % 3 = 1 (modulo)
# → 1
```

**Note:** The `**` operator uses binary exponentiation (exponentiation by squaring), which is more efficient than the general `^` operator for large integer exponents. It only works with integer exponents.

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

**Note:** `clearconstants` only clears user-defined constants. To clear all variables (which may include variables overriding built-in constants), use `clear`.

Built-in constants cannot be deleted individually, but you can list variables and clear them to restore access to the original constants:

```bash
gt 'vars'            # List all variables (not constants)
# x = 5

# To access constants again after overriding them with variables:
gt 'clear'           # Clear all variables (restores built-in constants)
gt 'pi'              # Now returns the constant π
# → 3.141592654
```

See the Constant Management section for the `constants` command to list all available constants.

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

## Comparison Operators

Comparison operators push a boolean result (`true`/`false`) onto the stack:

```bash
gt '5 3 gt'          # 5 > 3
# → true

gt '3 5 lt'          # 3 < 5
# → true

gt '5 5 gte'         # 5 >= 5
# → true

gt '5 5 lte'         # 5 <= 5
# → true

gt '5 5 eq'          # 5 == 5
# → true

gt '5 3 neq'         # 5 != 3
# → true
```

Shorthand symbols are also supported:

```bash
gt '5 3 >'           # Same as gt
# → true

gt '3 5 <'           # Same as lt
# → true

gt '5 5 >='          # Same as gte
# → true

gt '5 5 <='          # Same as lte
# → true

gt '5 5 =='          # Same as eq
# → true

gt '5 3 !='          # Same as neq
# → true
```

Comparison operators are metric-aware: they work with numbers that have attached units.

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

**Metric behavior:** HyperAdd (`[+]`), HyperSubtract (`[-]`), and HyperModulo (`[%]`) are metric-aware — they convert all operands to compatible units before computing. HyperMultiply (`[*]`), HyperDivide (`[/]`), and HyperPower (`[^]`) always return `Cool` (unitless).

HyperLog operators compute the sum of logarithms across all stack values:

```bash
gt '2 4 8 [lg]'       # log2(2) + log2(4) + log2(8) = 1 + 2 + 3 = 6
# → 6

gt '10 100 [log]'     # log10(10) + log10(100) = 1 + 2 = 3
# → 3

gt '2.718281828 7.389 [ln]'  # ln(e) + ln(e²) ≈ 1 + 2 = 3
# → 3
```

HyperLog operators always return `Cool` (unitless).

## Log Operators

Single-value logarithm operators:

```bash
gt '8 lg'             # log₂(8) = 3
# → 3

gt '1000 log'         # log₁₀(1000) = 3
# → 3

gt '2.718281828 ln'   # ln(e) ≈ 1
# → 1
```

Log operators return `Cool` (unitless).

## Metrics

Every number on the stack carries a unit of measurement. The default unit is `Cool` (unitless). Numbers with metrics can be converted, combined, and compared with automatic unit handling.

### Suffix Notation

Attach a unit directly to a number (no space):

```bash
gt '100Mbps 1hr *'     # 100 Mbps × 1 hour = 360 Gbits
# → 3.6e+11

gt '5GB 2hr /'         # 5 GB ÷ 2 hours
# → 5.555555556e+09

gt '100km 1hr /'       # Speed from distance/time
# → 27.77777778
```

### Available Metric Categories

| Category | Units |
|----------|-------|
| DataRate | `bps`, `Kbps`, `Mbps`, `Gbps`, `Tbps` (also `bit/s`, `mbit/s`, etc.) |
| DataSize | `bits`, `bytes`, `KB`, `MB`, `GB`, `TB`, `PB`, `KiB`, `MiB`, `GiB`, `TiB`, `PiB` |
| Time | `ms`, `s`, `min`, `hr`, `day` (also `sec`, `secs`) |
| Weight | `mg`, `g`, `kg`, `lb`, `oz`, `ton` |
| Speed | `mps`, `kmh`, `mph`, `knots` (also `knot`) |
| Distance | `m`, `km`, `mi`, `ft`, `in`, `nm` (also `mile`, `miles`, `foot`, `feet`) |
| Universal | `Cool` (default, unitless) |

### Unit Conversion

Convert a value to a different unit using `@` prefix and `convert`:

```bash
gt '1000Mbps @Gbps convert'    # 1000 Mbps → 1 Gbps
# → 1

gt '1hr @min convert'          # 1 hour → 60 min
# → 60

gt '1km @mi convert'           # 1 km → ~0.6214 mi
# → 0.6213711922
```

### Metric-Aware Arithmetic

**Addition and Subtraction** require compatible categories:

```bash
gt '100Mbps 50Mbps +'          # Same category: 150 Mbps
# → 150

gt '1km 500m +'                # Mixed units, same category: 1.5 km
# → 1.5

gt '5 100Mbps +'               # Cool absorbs: 105 Mbps
# → 105

gt '100Mbps 2hr +'             # ERROR: incompatible categories
```

**Multiplication** supports cross-category inference:

```bash
gt '100Mbps 1hr *'             # rate × time = data size
# → 3.6e+11

gt '100kmh 1hr *'              # speed × time = distance
# → 100000
```

**Division** also supports cross-category inference:

```bash
gt '10GB 2hr /'                # data / time = rate
# → 11111111.11

gt '1km 1s /'                  # distance / time = speed
# → 1000
```

**Power** (`^`) always returns `Cool` (unitless) — `x^n` has different units than `x`:

```bash
gt '2hr 3 ^'                   # 2³ = 8 (not 8hr³)
# → 8
```

### Prefix Modes: SI vs IEC

Data size prefixes can use SI (1000-based) or IEC (1024-based) factors. Prefix mode is a REPL setting — in single-command mode the default is SI.

In REPL mode:
```
> metric decimal set     # SI mode (KB = 1000 bytes)
prefix mode: SI
> 1GB @MB convert        # SI: 1000 MB
1000
> metric binary set      # IEC mode (KB = 1024 bytes)
prefix mode: IEC
> 1GB @MB convert        # IEC: 1024 MB
1024
```

### Metric Commands

```bash
gt 'metric show'               # Show metric info for top of stack
# → Mbps, DataRate, base: bps, factor: 1e+06

gt 'metric list'               # List all categories
# → DataRate, DataSize, Distance, Speed, Time, Universal, Weight

gt 'metric DataRate'           # List metrics in a category
# → Gbps, Kbps, Mbps, Tbps, bps

gt 'metric compatible'         # Check if top two stack metrics are compatible
# → Mbps (DataRate) and Gbps (DataRate): true
```

### Custom Metrics

Define your own units:

```bash
gt 'custom define foobar 42 Custom'    # Register custom unit
# → defined custom metric "foobar" (factor: 42, category: Custom)

gt 'custom define foobar 42 Custom 10foobar 5foobar +'  # Use it in arithmetic
# → 15

gt 'custom define foobar 42 Custom custom list'         # List custom metrics
# → foobar

gt 'custom define foobar 42 Custom custom undefine foobar'  # Remove custom unit
# → removed custom metric "foobar"
```

Custom metrics can be defined in any category. They use the specified factor relative to the category's base unit.

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
