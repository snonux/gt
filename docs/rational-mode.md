# Rational Number Mode

gt supports a rational number calculation mode that uses Go's `math/big.Rat` type for arbitrary-precision rational arithmetic. This mode represents numbers as exact fractions (numerator/denominator) rather than floating-point approximations.

Rational mode is only available in interactive REPL mode. Start gt without arguments to enter the REPL, then use the `rat` command to switch modes.

## Enabling and Disabling

The `rat` command controls rational number mode. It is a REPL-only command — it cannot be used in single-command mode (`gt '...'`).

### `rat on`

Enable rational number mode:

```
> rat on
Rational mode enabled
```

### `rat off`

Disable rational mode and return to float64 calculations (the default):

```
> rat off
Rational mode disabled (using float64)
```

### `rat toggle`

Switch between rational mode and float64 mode:

```
> rat toggle
Rational mode enabled
> rat toggle
Rational mode disabled (using float64)
```

Without an argument, `rat` prints a usage message:

```
> rat
rat command requires an argument: on, off, or toggle
```

Invalid arguments are also reported:

```
> rat enable
Unknown rat mode: enable. Valid modes: on, off, toggle
```

## How Rational Mode Works

gt's rational mode uses Go's `math/big.Rat` type, which represents numbers as exact fractions with arbitrary-precision integer numerator and denominator.

When a number is parsed in rational mode, gt calls `big.Rat.SetString()`, which creates the exact rational representation. For example:

- `0.1` becomes the exact fraction `1/10`
- `0.2` becomes the exact fraction `2/10` (stored as `1/5`)
- `1/3` becomes `1/3` exactly
- `3` becomes `3/1` (integer)

In contrast, float64 mode stores `0.1` as the nearest representable binary fraction `3602879701896397/36028797018963968`, which is not exactly `1/10`.

### Result formatting

Rational numbers are displayed using `big.Rat.FloatString(10)`, which formats the value as a decimal with up to 10 significant digits. This means rational results show trailing zeros to maintain consistent precision:

```
> rat on
Rational mode enabled
> 10 3 -
7.0000000000
```

Compare with float mode which uses `%.10g` formatting (suppressing trailing zeros):

```
> rat off
Rational mode disabled (using float64)
> 10 3 -
7
```

## Precision Comparison

### Division

Division of integers produces the same displayed precision in both modes (10 significant digits), but rational mode stores the exact fraction internally:

| Expression    | Float mode       | Rat mode         | Notes              |
|---------------|------------------|------------------|--------------------|
| `1 3 /`       | `0.3333333333`   | `0.3333333333`   | Both truncated     |
| `1 7 /`       | `0.1428571429`   | `0.1428571429`   | Both truncated     |
| `1 6 /`       | `0.1666666667`   | `0.1666666667`   | Both truncated     |
| `10 3 /`      | `3.3333333333`   | `3.3333333333`   | Both truncated     |

At 10 digits of precision, the displayed output is identical. The difference appears when chaining operations:

### Chained division and multiplication

| Expression     | Float mode       | Rat mode          | Notes                     |
|----------------|------------------|-------------------|---------------------------|
| `3 10 / 3 *`   | `0.9`            | `0.9000000000`    | Float loses precision     |
| `1 3 / 3 *`    | `1`              | `1.0000000000`    | Both exact, different fmt |
| `2 3 / 3 *`    | `2`              | `2.0000000000`    | Both exact, different fmt |

### Integer arithmetic

Integer operations produce correct results in both modes:

| Expression  | Float mode | Rat mode        |
|-------------|------------|-----------------|
| `1 2 +`     | `3`        | `3.0000000000`  |
| `10 3 -`    | `7`        | `7.0000000000`  |
| `3 4 *`     | `12`       | `12.0000000000` |

### Decimal arithmetic limitations

Rational mode stores exact rational representations when parsing decimal literals. However, the arithmetic operations (addition, subtraction, modulo) internally convert values to `float64` for metric-aware computation. This creates a known limitation:

**Some decimal fractions fail in rational mode.** Decimals that cannot be exactly represented as float64 (such as `0.1`, `0.2`, `0.3`) cause errors when used with `+`, `-`, or `%` operators in rational mode:

```
> rat on
Rational mode enabled
> 0.1 0.2 +
Error: rpn operator failed for '+': +: convertToBase: cannot convert rational number to float64
```

This happens because:

1. `0.1` is parsed as the exact fraction `1/10`
2. `big.Rat.Float64()` returns `ok=false` for `1/10` because `1/10` cannot be exactly represented as float64
3. The `convertToBase` metric helper treats `ok=false` as an error

**Decimals that are exact powers of 2 work correctly:** `0.5` (1/2), `0.25` (1/4), `0.75` (3/4), `0.125` (1/8), etc.

```
> rat on
Rational mode enabled
> 0.5 0.5 +
1.0000000000
> 0.25 0.75 +
1.0000000000
> 0.125 0.125 +
0.2500000000
```

**Multiplication and division with any decimals work:** These operations bypass the problematic code path:

```
> rat on
Rational mode enabled
> 0.1 0.2 *
0.0200000000
> 0.3 0.1 /
3.0000000000
```

> **Note:** This limitation exists in the current implementation. The `Rat.Float64()` method rejects conversions where the rational number cannot be exactly represented as float64, even though the conversion is valid for computation purposes. A future fix would allow lossy conversions to return the approximate float64 value.

### Float mode precision issues

Float mode exhibits classic IEEE 754 floating-point errors:

```
> rat off
Rational mode disabled (using float64)
> 0.3 0.1 - 0.2 -
-2.775557562e-17
```

The result should be `0`, but floating-point rounding errors accumulate. In rational mode, the same expression fails with an error rather than producing a silently wrong result.

## When to Use Rational Mode

### Recommended use cases

- **Integer arithmetic**: Exact results for integer addition, subtraction, multiplication, and division
- **Dyadic decimals**: Values like `0.5`, `0.25`, `0.125` that are exact powers of 2
- **Division results**: When you need to chain division results into further multiplication
- **Avoiding silent float errors**: Rat mode fails loudly rather than returning an incorrect float result

### Not recommended for

- **General decimal arithmetic**: Values like `0.1`, `0.2`, `0.3` fail with `+`, `-`, `%` due to the float64 conversion limitation described above
- **Single-command mode**: Rational mode is REPL-only; `gt 'rat on 1 2 +'` will not work
- **Performance-sensitive calculations**: See performance trade-offs below

## Performance Trade-offs

Rational mode has higher computational overhead than float64 mode:

| Factor              | Float mode | Rat mode          |
|---------------------|------------|-------------------|
| Number creation     | CPU register | Heap allocation (`*big.Rat`) |
| Memory per number   | 8 bytes    | ~48+ bytes (pointer + big.Ints) |
| Arithmetic          | Hardware FPU | Big integer arithmetic |
| String formatting   | `fmt.Sprintf` | `big.Rat.FloatString` |
| Precision           | ~15-17 digits | Arbitrary (limited by memory) |

For most calculator use cases, the performance difference is negligible. However, for very large expressions or tight REPL loops, float mode will be noticeably faster.

## Edge Cases

### Division by zero

Both modes return an error:

```
> rat on
Rational mode enabled
> 1 0 /
Error: rpn operator failed for '/': /: division by zero
```

### Modulo by zero

Both modes return an error:

```
> rat on
Rational mode enabled
> 10 0 %
Error: rpn operator failed for '%': %: modulo by zero
```

### Insufficient operands

Both modes return the same error:

```
> rat on
Rational mode enabled
> +
Error: rpn operator failed for '+': +: insufficient operands for +: need at least 2 values
```

### Mode persistence in REPL

Rational mode persists across REPL commands within the same session:

```
> rat on
Rational mode enabled
> 1 2 +
3.0000000000
> 3 4 *
12.0000000000
> 10 3 /
3.3333333333
> rat off
Rational mode disabled (using float64)
> 1 2 +
3
```

### Fresh sessions start in float mode

Each new gt session (single-command or new REPL) starts in float mode by default:

```
$ gt '1 2 +'
3
```

### Mode and metrics

Rational mode works with metric-aware operations, subject to the decimal conversion limitation described above. Integers and dyadic decimals with metrics work correctly:

```
> rat on
Rational mode enabled
> 100Mbps 50Mbps +
Error: rpn operator failed for '+': +: convertToBase: cannot convert rational number to float64
```

The metric conversion path has the same float64 conversion issue. For metric calculations, float mode is recommended unless using integer or dyadic decimal values.

### Mode and constants

Constants (pi, e, etc.) are resolved as float64 values. In rational mode, they are wrapped in `Rat` types:

```
> rat on
Rational mode enabled
> pi
3.1415926536
```

Note that `pi` is stored as a float64 constant and converted to `Rat` via `SetFloat64()`, so it carries the same float64 precision as float mode. It does not provide a more precise value of pi.

### Mode and variables

Variables store float64 values. When accessed in rational mode, they are converted to `Rat` via `SetFloat64()`:

```
> x 3.5 =
x = 3.5000000000
> rat on
Rational mode enabled
> x
3.5000000000
```

## Summary

| Feature            | Float mode | Rat mode |
|--------------------|------------|----------|
| Default            | Yes        | No       |
| REPL-only          | No         | Yes      |
| Integer arithmetic | Exact      | Exact    |
| Decimal arithmetic | Approx.    | Exact (with limitations) |
| `+`, `-` with 0.1  | Works (approx.) | Error |
| `*`, `/` with 0.1  | Works      | Works    |
| Metric operations  | Full       | Limited  |
| Output formatting  | `%.10g`    | `FloatString(10)` |
| Memory usage       | Low        | Higher   |
| Performance        | Fast       | Slower   |
