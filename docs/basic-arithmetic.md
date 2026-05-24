# Basic Arithmetic RPN Operators

gt evaluates expressions using Reverse Polish Notation (RPN), a stack-based calculation method where operators follow their operands. No parentheses needed — the token sequence determines the order of operations.

## How RPN Works

Each token is processed left to right:

1. **Numbers** are pushed onto the stack.
2. **Operators** pop the required operands, compute the result, and push it back.

For example, `3 4 +` pushes 3, then 4, then `+` pops both, adds them, and pushes 7.

### Stack visualization

```
Input:  3  4  +

  Step  Stack
  ----  -----
   3    [3]
   4    [3, 4]
   +    [7]

Result: 7
```

## Operators

All binary arithmetic operators pop two values and push one result. The **first** value pushed is the left operand (`a`), the **second** is the right operand (`b`). The operation is always `a <op> b`.

### `+` Addition

```
$ gt '3 4 +'
7
```

### `-` Subtraction

Computes `a - b` (first value minus second). Order matters:

```
$ gt '3 4 -'
-1
```

### `*` Multiplication

```
$ gt '5 6 *'
30
```

### `/` Division

Computes `a / b` (first value divided by second).

```
$ gt '20 4 /'
5
```

### `^` Power

Computes `a ^ b` (first value raised to the power of the second). Result is always unitless.

```
$ gt '2 3 ^'
8
$ gt '2 10 ^'
1024
$ gt '5 0 ^'
1
$ gt '2 -3 ^'
0.125
```

### `%` Modulo

Computes `a % b` (remainder of `a / b`).

```
$ gt '10 3 %'
1
$ gt '7 3 %'
1
$ gt '13 5 %'
3
$ gt '100 7 %'
2
```

## Multi-operand Expressions

RPN handles complex expressions without parentheses.

### Chained operations

```
$ gt '1 2 3 + +'
6
```

Step by step: `1` pushed, `2` pushed, `3` pushed, `+` pops 2 and 3 producing 5 (stack: `[1, 5]`), `+` pops 1 and 5 producing 6.

```
$ gt '10 2 3 - *'
-10
```

Step by step: `10` pushed, `2` pushed, `3` pushed, `-` pops 2 and 3 producing -1 (stack: `[10, -1]`), `*` pops 10 and -1 producing -10.

### Nested operations

```
$ gt '3 4 + 5 6 + *'
77
```

This is the RPN equivalent of `(3 + 4) * (5 + 6)` = `7 * 11` = `77`.

```
$ gt '2 3 ^ 4 5 ^ +'
1032
```

This is `2^3 + 4^5` = `8 + 1024` = `1032`.

### Order of operations

In RPN, the order is explicit in the token sequence:

```
$ gt '10 3 2 * /'
1.666666667
```

This computes `10 / (3 * 2)` = `10 / 6`. The inner operation `3 2 *` is evaluated first because `*` comes before `/`.

```
$ gt '100 10 / 5 +'
15
```

This computes `(100 / 10) + 5` = `10 + 5` = `15`.

### Large expressions

```
$ gt '1 2 + 3 4 + * 5 6 + +'
32
```

Breakdown: `(1+2) * (3+4) + (5+6)` = `3 * 7 + 11` = `21 + 11` = `32`.

## Examples

### Compound calculations

Compute `(a + b) * c - d / e` in RPN:

```
$ gt '10 5 + 3 * 20 4 / -'
35
```

### Geometry

Circle area (pi * r^2):

```
$ gt '3.14159 5 2 ^ *'
78.53975
```

Rectangle perimeter (2 * (w + h)):

```
$ gt '10 5 + 2 *'
30
```

Pythagorean theorem (sqrt(a^2 + b^2)):

```
$ gt '3 2 ^ 4 2 ^ + sqrt'
5
```

Average of three numbers:

```
$ gt '10 20 30 + + 3 /'
20
```

## Edge Cases

### Division by zero

Returns an error:

```
$ gt '5 0 /'
Error: division by zero
```

### Modulo by zero

Returns an error:

```
$ gt '5 0 %'
Error: modulo by zero
```

### Negative results

```
$ gt '3 7 -'
-4
$ gt '-5 3 *'
-15
```

### Insufficient operands

Any binary operator on an empty or single-element stack returns an error:

```
$ gt '+'
Error: stack is empty

$ gt '5 +'
Error: stack has insufficient operands
```

## Summary Table

| Operator | RPN    | Infix   | Description          | Example result |
|----------|--------|---------|----------------------|---------------|
| `+`      | `a b +`| `a + b` | Addition             | `3 4 +` = 7   |
| `-`      | `a b -`| `a - b` | Subtraction          | `3 4 -` = -1  |
| `*`      | `a b *`| `a * b` | Multiplication       | `5 6 *` = 30  |
| `/`      | `a b /`| `a / b` | Division             | `20 4 /` = 5  |
| `^`      | `a b ^`| `a ^ b` | Power (exponent)     | `2 3 ^` = 8   |
| `%`      | `a b %`| `a % b` | Modulo (remainder)   | `10 3 %` = 1  |
