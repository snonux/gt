# Comparison RPN Operators

gt provides six comparison operators, each with both a named and a symbolic alias. All operators are metric-aware — they convert values to base units before comparing. Results are boolean (`true` / `false`) and coerce into arithmetic (true → 1, false → 0).

## Operators

| Operator | Alias | RPN         | Infix   | Description     |
|----------|-------|-------------|---------|-----------------|
| `gt`     | `>`   | `a b gt`    | `a > b` | Greater than    |
| `lt`     | `<`   | `a b lt`    | `a < b` | Less than       |
| `gte`    | `>=`  | `a b gte`   | `a >= b`| Greater/equal   |
| `lte`    | `<=`  | `a b lte`   | `a <= b`| Less/equal      |
| `eq`     | `==`  | `a b eq`    | `a == b`| Equal           |
| `neq`    | `!=`  | `a b neq`   | `a != b`| Not equal       |

All operators pop two values (`a` first, then `b`) from the stack, compare them, and push the boolean result.

## Truth Table

### `gt` / `>`

| a   | b   | Result |
|-----|-----|--------|
| 5   | 3   | true   |
| 3   | 5   | false  |
| 5   | 5   | false  |

```
$ gt '5 3 gt'
true
$ gt '5 3 >'
true
```

### `lt` / `<`

| a   | b   | Result |
|-----|-----|--------|
| 3   | 5   | true   |
| 5   | 3   | false  |
| 5   | 5   | false  |

### `gte` / `>=`

| a   | b   | Result |
|-----|-----|--------|
| 5   | 3   | true   |
| 5   | 5   | true   |
| 3   | 5   | false  |

### `lte` / `<=`

| a   | b   | Result |
|-----|-----|--------|
| 3   | 5   | true   |
| 5   | 5   | true   |
| 5   | 3   | false  |

### `eq` / `==`

| a   | b   | Result |
|-----|-----|--------|
| 5   | 5   | true   |
| 5   | 3   | false  |

### `neq` / `!=`

| a   | b   | Result |
|-----|-----|--------|
| 5   | 3   | true   |
| 5   | 5   | false  |

## Metric-Aware Comparison

Comparison operators convert operands to base units before comparing, so different units within the same category can be compared directly.

### Same category, different units

```
$ gt '1km 1000m eq'
true
$ gt '1km 500m gt'
true
$ gt '500m 1km lt'
true
```

### Data storage (decimal SI prefixes)

```
$ gt '1GB 1000MB gte'
true
$ gt '1GB 1000MB gt'
false
$ gt '1GB 1024MB eq'
false
```

`1GB` equals exactly `1000MB` (decimal SI), so `gt` is false and `gte` is true. `1024MB` is not equal to `1GB`.

### Network throughput

```
$ gt '100Mbps 50Mbps gt'
true
$ gt '1Gbps 1000Mbps eq'
true
```

### Incompatible categories

Comparing values from different metric categories produces an error:

```
$ gt '5m 3kg gt'
Error: incompatible metric categories
```

Unitless numbers are always compatible with any metric value.

## Boolean Coercion in Arithmetic

Boolean results coerce to numbers: `true` → 1, `false` → 0. This lets comparison results flow directly into arithmetic.

### Conditional addition

```
$ gt '5 3 gt 1 +'
2
```

`5 > 3` is true (1), so `1 + 1 = 2`.

```
$ gt '3 5 gt 1 +'
1
```

`3 > 5` is false (0), so `0 + 1 = 1`.

### Conditional multiplication

```
$ gt '5 3 gt 10 *'
10
$ gt '3 5 gt 10 *'
0
```

### Chained comparisons

```
$ gt '9 3 gt 4 5 lt +'
2
```

`9 > 3` (true/1) + `4 < 5` (true/1) = `1 + 1 = 2`.

### Boolean arithmetic with subtraction

```
$ gt '5 3 gt 0 -'
1
$ gt '5 5 eq 1 -'
0
```

## Examples

### Threshold checks

```
$ gt '85 80 gt'
true
```

CPU at 85%, threshold at 80% — alert.

### Range validation

Check that a value falls within acceptable bounds:

```
$ gt '72 68 gte 100 lte +'
2
```

Temperature 72°F: `72 >= 68` (true/1) + `72 <= 100` (true/1) = 2 (both checks pass).

```
$ gt '105 68 gte 100 lte +'
1
```

Temperature 105°F: `105 >= 68` (true/1) + `105 <= 100` (false/0) = 1 (outside range).

When the sum is 2, the value is within range. Any other result means at least one check failed.

### Metric validation

```
$ gt '100Mbps 50Mbps gt'
true
```

Download speed 100Mbps exceeds minimum 50Mbps — pass.

### Counting conditions

Sum boolean results to count how many conditions are met:

```
$ gt '10 5 gt 20 15 gt 30 25 gt +'
3
```

All three conditions are true: `1 + 1 + 1 = 3`.

### Guard expressions

Multiply a value by a condition to guard it:

```
$ gt '99 80 gt *'
99
```

Value 99 is used because `99 > 80` is true (1): `99 * 1 = 99`.

```
$ gt '50 80 gt *'
0
```

Value 50 is discarded because `50 > 80` is false (0): `50 * 0 = 0`.

## Edge Cases

### Zero comparisons

```
$ gt '0 0 eq'
true
$ gt '0 0 neq'
false
$ gt '0 1 lt'
true
```

### Negative numbers

```
$ gt '-1 0 lt'
true
$ gt '-5 -3 lt'
true
```

### Insufficient operands

Any comparison on an empty or single-element stack returns an error:

```
$ gt 'gt'
Error: stack is empty
$ gt '5 gt'
Error: stack has insufficient operands
```
