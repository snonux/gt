# Logarithm RPN Operators

`gt` provides three logarithm operators: `lg` (base 2), `log` (base 10), and
`ln` (natural log, base *e*). Each is a unary operator that pops a single value
from the stack, computes the logarithm, and pushes the result.

Logarithm results are always unitless (Cool metric), regardless of whether the
input carried a metric.

## How RPN Works

For unary operators, the operand precedes the operator:

```
Input:  8  lg

  Step  Stack
  ----  -----
   8    [8]
   lg   [3]

Result: 3
```

## Operators

### `lg` — Logarithm Base 2

Computes the base-2 logarithm: log₂(*x*).

#### Mathematical explanation

log₂(*x*) = *y* means 2^*y* = *x*. It answers: "To what power must 2 be raised
to produce *x*?"

#### Examples

```
$ echo "8 lg" | gt
3

$ echo "1024 lg" | gt
10

$ echo "1 lg" | gt
0

$ echo "1e-10 lg" | gt
-33.21928095
```

#### Practical use cases

- **Information theory** — Number of bits needed to represent *x* distinct
  values (e.g., log₂(1024) = 10 bits).
- **Algorithms** — Time complexity of divide-and-conquer algorithms (O(log n)).
- **Digital audio** — Bit-depth calculations.
- **Computer science** — Tree height, binary search steps.

---

### `log` — Logarithm Base 10

Computes the base-10 logarithm: log₁₀(*x*).

#### Mathematical explanation

log₁₀(*x*) = *y* means 10^*y* = *x*. It answers: "To what power must 10 be
raised to produce *x*?"

#### Examples

```
$ echo "1000 log" | gt
3

$ echo "100 log" | gt
2

$ echo "10 log" | gt
1

$ echo "1 log" | gt
0

$ echo "1e-10 log" | gt
-10
```

#### Practical use cases

- **Decibels (dB)** — Sound intensity and signal power ratios:
  L = 10 · log₁₀(P/P₀). Example: `1000 10 / log 10 *` → 30 dB.
- **pH scale** — Acidity/alkalinity: pH = -log₁₀([H⁺]). Example:
  `1e-7 log 0 -` → 7 (neutral water).
- **Richter scale** — Earthquake magnitude.
- **Order of magnitude** — Quick estimate of size: log₁₀(500) ≈ 2.7 means
  "between 100 and 1000".
- **Scientific notation** — The integer part of log₁₀(*x*) gives the exponent.

---

### `ln` — Natural Logarithm

Computes the natural logarithm: ln(*x*) = logₑ(*x*), where *e* ≈ 2.718281828.

#### Mathematical explanation

ln(*x*) = *y* means *e*^*y* = *x*. It is the inverse of the exponential function
and the integral of 1/*t* from 1 to *x*.

#### Examples

```
$ echo "e ln" | gt
1

$ echo "2.718281828 ln" | gt
0.9999999998

$ echo "1 ln" | gt
0

$ echo "7.38905609893065 ln" | gt
2
```

#### Practical use cases

- **Continuous growth** — Compound interest, population growth:
  A = A₀ · e^(rt). Solve for t: t = ln(A/A₀) / r.
- **Probability & statistics** — Normal distribution, likelihood functions.
- **Calculus** — Derivative of ln(*x*) is 1/*x*; integral of 1/*x* is ln(|*x*|).
- **Entropy** — Shannon entropy uses natural log in thermodynamics.
- **Chemistry** — Arrhenius equation for reaction rates.

---

## Edge Cases

All three operators are **undefined for non-positive numbers**. Zero and negative
inputs produce an error:

```
$ echo "0 lg" | gt
Error: lg undefined for non-positive numbers

$ echo "-5 log" | gt
Error: log undefined for non-positive numbers

$ echo "0 ln" | gt
Error: ln undefined for non-positive numbers

$ echo "-1 ln" | gt
Error: ln undefined for non-positive numbers
```

### Summary table

| Input    | Result                        |
|----------|-------------------------------|
| *x* > 0  | Computed normally             |
| *x* = 0  | Error (undefined)             |
| *x* < 0  | Error (undefined)             |
| *x* = 1  | 0 (for all three operators)   |
| empty    | Error (stack underflow)       |

---

## Metric Handling

Logarithms always produce unitless results. Even when the input carries a metric,
the result is Cool (no metric):

```
$ echo "100Mbps lg" | gt          # treats as number with metric, result is Cool
```

This is mathematically correct — taking a logarithm of a dimensioned quantity
requires a reference value (e.g., log₁₀(P/P₀) for decibels), so gt treats the
operands as their numeric values.

---

## Implementation

- Source: `internal/rpn/operations_arithmetic.go` — `Log2`, `Log10`, `Ln`
- Uses `math.Log2`, `math.Log10`, `math.Log` from Go's standard library
- Shared helper: `logOp(stack, opName, logFn)` handles pop, validate, compute,
  and push
- Tests: `internal/rpn/operations_test.go` — `TestLog2`, `TestLog10`, `TestLn`
  plus boolean coercion and hyper-operator variants
