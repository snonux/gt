# Fast Power Operator (`**`)

The `**` operator computes integer exponentiation using **binary exponentiation** (also known as square-and-multiply). It raises the base (first operand) to the integer power of the exponent (second operand) on the RPN stack.

## Binary Exponentiation

Binary exponentiation reduces the number of multiplications from O(n) to O(log n) by exploiting the binary representation of the exponent.

The algorithm works as follows:

1. Initialize `result = 1`
2. While the exponent is positive:
   - If the exponent is **odd**, multiply `result` by the current `base`
   - Square the `base`
   - Divide the exponent by 2 (integer division)
3. For negative exponents, compute `1 / (base^|exp|)`

### Example: 3^13

| Step | exp | exp odd? | result | base |
|------|-----|----------|--------|------|
| 0    | 13  | yes      | 1 * 3 = 3 | 3 |
| 1    | 6   | no       | 3     | 9   |
| 2    | 3   | yes      | 3 * 9 = 27 | 81 |
| 3    | 1   | yes      | 27 * 81 = 2187 | 6561 |
| 4    | 0   | —        | **2187** | — |

Only 4 multiplications instead of 12. For 10^100, that's ~17 multiplications instead of 99.

## Syntax

RPN (Reverse Polish Notation):

```
base exponent **
```

Operands are popped in reverse order: `exponent` (top) then `base`, result is `base^exponent`. Pushed back as a unitless (Cool) number.

## Examples

### Integer exponent

```
2 10 **    → 1024
```

### Negative exponent

```
2 -3 **    → 0.125   (i.e., 1/8)
```

### Zero exponent

```
5 0 **     → 1
```

Any number to the power of 0 is 1 (short-circuited, no multiplication performed).

### Large exponent

```
10 100 **  → 1e+100
```

### Fractional exponent (rejected)

```
2 1.5 **   → Error: exponent must be an integer, got 1.5
```

The `**` operator **requires** an integer exponent. Use `^` for fractional exponents.

## `**` vs `^`

| Feature | `**` (Fast Power) | `^` (Power) |
|---------|-------------------|-------------|
| Exponent type | Integer only | Any float |
| Algorithm | Binary exponentiation (square-and-multiply) | `math.Pow()` |
| Complexity | O(log n) multiplications | Hardware/logarithm-based |
| Fractional exponents | Rejected with error | Supported |
| Negative exponents | Supported (as `1/base^|n|`) | Supported |
| Zero exponent | Short-circuited → 1 | `math.Pow(base, 0)` → 1 |

### When to use `**`

- Exponent is a known integer
- You want a fast, exact integer-power computation
- You want to guard against accidental fractional exponents

### When to use `^`

- Exponent may be fractional (e.g., `9 0.5 ^` → 3 for square root)
- Working with arbitrary floating-point exponents
- General-purpose power computation

## Edge Cases

| Input | Result | Notes |
|-------|--------|-------|
| `x 0 **` | 1 | Short-circuited; no multiplications |
| `x -n **` | 1/(x^n) | Computed recursively via positive exponent |
| `0 0 **` | 1 | By convention (0^0 = 1) |
| `2 1.5 **` | Error | Fractional exponents not supported |
| `x n **` (n large) | float64 result | Result limited by float64 range; may overflow to Inf |
| `-2 3 **` | -8 | Negative base with integer exponent works correctly |

## Performance Notes

- **Time complexity**: O(log |n|) where n is the exponent. An exponent of 100 requires only ~7 loop iterations (vs 99 naive multiplications).
- **Space complexity**: O(1) for positive exponents; O(log |n|) stack depth for negative exponents (recursive call).
- The implementation uses `float64` throughout, so results may lose precision for very large exponents (same limitations as any floating-point computation).
- For the specific case of exponent 0, the function returns 1 immediately without entering the loop.
