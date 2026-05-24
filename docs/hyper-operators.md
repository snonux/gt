# Hyper RPN Operators

Hyper operators operate on **all** values currently on the stack, rather than just the top one or two. They pop every value, reduce them left-associatively, and push a single result back.

The syntax uses square brackets: `[+]`, `[*]`, `[-]`, `[/]`, `[^]`, `[%]`, `[lg]`, `[log]`, `[ln]`.

## How Hyper Operators Work

```
Input:  1  2  3  4  5  [+]

  Step  Stack
  ----  -----
   1    [1]
   2    [1, 2]
   3    [1, 2, 3]
   4    [1, 2, 3, 4]
   5    [1, 2, 3, 4, 5]
   [+]  [15]

Result: 15
```

All five values are consumed, summed left-to-right (`((((1 + 2) + 3) + 4) + 5)`), and the single result is pushed back.

### Requirements

- At least **two** values must be on the stack, or the operator returns an error.
- After execution, exactly one value remains on the stack.

## Arithmetic Hyper Operators

### `[+]` Sum

Adds all stack values left-associatively.

```
$ gt '1 2 3 4 5 [+]'
15
```

With two operands, behaves identically to binary `+`:

```
$ gt '10 20 [+]'
30
```

### `[*]` Product

Multiplies all stack values left-associatively. Result is always unitless.

```
$ gt '2 3 4 [*]'
24
$ gt '1 2 3 [*]'
6
```

### `[-]` Subtraction

Subtracts all stack values left-associatively from the first.

```
$ gt '10 3 2 [-]'
5
$ gt '100 10 20 30 [-]'
40
```

### `[/]` Division

Divides all stack values left-associatively from the first. Result is always unitless.

```
$ gt '100 5 2 [/]'
10
$ gt '1000 10 10 [/]'
10
```

### `[^]` Power

Raises left-associatively. The first value is the base, subsequent values are successive exponents. Result is always unitless.

```
$ gt '2 3 2 [^]'
64
```

Equivalently: `(2 ^ 3) ^ 2` = `8 ^ 2` = 64.

With two operands, behaves identically to binary `^`:

```
$ gt '2 10 [^]'
1024
```

### `[%]` Modulo

Computes modulo left-associatively.

```
$ gt '100 7 3 [%]'
2
```

Equivalently: `100 % 7` = 2, then `2 % 3` = 2.

```
$ gt '10 3 2 2 [%]'
0
```

## Logarithmic Hyper Operators

Logarithmic hyper operators compute the **sum** of the log function applied to each stack value. All input values must be positive.

### `[lg]` Base-2 Logarithm Sum

```
$ gt '2 4 8 [lg]'
6
```

Equivalently: `log2(2) + log2(4) + log2(8)` = `1 + 2 + 3` = 6.

```
$ gt '1 2 3 4 5 [lg]'
6.907
```

### `[log]` Base-10 Logarithm Sum

```
$ gt '10 100 [log]'
3
```

Equivalently: `log10(10) + log10(100)` = `1 + 2` = 3.

### `[ln]` Natural Logarithm Sum

```
$ gt '2.718281828 7.389 [ln]'
2.999992408
```

Equivalently: `ln(2.718281828) + ln(7.389)` ≈ `1 + 2` ≈ 3.

## Metric-Aware Behavior

Some hyper operators are metric-aware, meaning they understand units and can operate on values with compatible units.

### Metric-aware operators: `[+]`, `[-]`, `[%]`

These operators:

1. **Validate** that all operands share the same metric category (or are unitless, which is always compatible).
2. **Convert** all values to the result metric's base units before computing.
3. **Push** the result with the first non-unitless metric (or unitless if all inputs are unitless).

```
$ gt '1km 2km 3km [+]'
6
```

Result is in kilometers: `1 + 2 + 3 = 6 km`.

```
$ gt '1km 500m [+]'
1.5
```

`500m` is converted to `0.5km`, then `1km + 0.5km = 1.5km`.

```
$ gt '10m 3m 2m [%]'
1
```

`10 % 3 = 1`, then `1 % 2 = 1 m`.

### Non-metric operators: `[*]`, `[/]`, `[^]`, `[lg]`, `[log]`, `[ln]`

These operators use raw numeric values and always produce unitless (Cool) results, regardless of input units:

- `[*]` — multiplying meters by meters would yield square meters, a different category, so the result is unitless.
- `[/]` — dividing two like units yields a ratio, which is unitless.
- `[^]` — exponents are inherently unitless.
- `[lg]`, `[log]`, `[ln]` — logarithms require dimensionless inputs, so the raw numeric value is used and the result is unitless.

### Mixed unitless and unit values

Unitless (Cool) values can be combined with any metric category — they "absorb" and take on the category of the other operands:

```
$ gt '0 5km [+]'
5
```

## Examples

### Batch aggregation

```
$ gt '12 15 18 14 16 [+]'
75
```

### Running total with metric conversion

```
$ gt '5km 2km 1000m [+]'
8
```

Result is in km: `5 + 2 + 1 = 8 km`.

### Compute a product

```
$ gt '2 5 10 [*]'
100
```

### Cascading subtraction

```
$ gt '100 10 20 30 5 [-]'
35
```

`100 - 10 - 20 - 30 - 5 = 35`.

### Repeated division

```
$ gt '1000 2 2 2 2 [/]'
62.5
```

`1000 / 2 / 2 / 2 / 2 = 62.5`.

### Log sum for information theory

Sum of log2 values for entropy calculations:

```
$ gt '0.25 0.5 0.25 [lg]'
-3
```

`log2(0.25) + log2(0.5) + log2(0.25)` = `-2 + -1 + -2` = -5.

## Edge Cases

### Division by zero

```
$ gt '10 5 0 [/]'
Error: division by zero
```

### Modulo by zero

```
$ gt '100 7 0 [%]'
Error: modulo by zero
```

### Non-positive log inputs

```
$ gt '0 [lg]'
Error: log2 undefined for non-positive numbers

$ gt '-1 [ln]'
Error: ln undefined for non-positive numbers
```

### Incompatible metrics

Mixing different metric categories (e.g., length and weight) in a metric-aware operator returns an error.

### Insufficient operands

All hyper operators require at least two values:

```
$ gt '5 [*]'
Error: insufficient operands for [*]: need at least 2 values

$ gt '[+]'
Error: stack is empty
```

## Summary Table

| Operator | Description                    | Metric-aware | Result metric |
|----------|--------------------------------|--------------|---------------|
| `[+]`    | Sum all values                 | Yes          | First non-Cool|
| `[*]`    | Multiply all values            | No           | Cool (unitless)|
| `[-]`    | Subtract all values            | Yes          | First non-Cool|
| `[/]`    | Divide all values              | No           | Cool (unitless)|
| `[^]`    | Power (left-associative)       | No           | Cool (unitless)|
| `[%]`    | Modulo (left-associative)      | Yes          | First non-Cool|
| `[lg]`   | Sum of log2 for all values     | No           | Cool (unitless)|
| `[log]`  | Sum of log10 for all values    | No           | Cool (unitless)|
| `[ln]`   | Sum of ln for all values       | No           | Cool (unitless)|
