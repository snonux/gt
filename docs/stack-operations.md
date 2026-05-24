# Stack Manipulation Operators

`gt` provides a set of operators for directly manipulating the RPN stack. These
let you duplicate values, reorder elements, discard items, inspect the stack,
and clear everything.

## How the Stack Works

The RPN stack grows from bottom to top. The **top** of the stack is the most
recently pushed value. Stack manipulation operators read and modify this
structure without performing arithmetic.

```
Stack (bottom -> top):

  +-------+
  |  value (top)    |
  +-------+
  |  value          |
  +-------+
  |  value          |
  +-------+
  |  value (bottom) |
  +-------+
```

## Operators at a Glance

| Operator     | Operands | Action                              | Example         |
|-------------|----------|-------------------------------------|-----------------|
| `dup`       | 1        | Duplicate the top value             | `5 dup` -> 5 5  |
| `swap`      | 2        | Swap the top two values             | `3 4 swap` -> 4 3 |
| `pop`       | 1        | Discard the top value               | `5 6 pop` -> 5  |
| `show`      | 0        | Print the entire stack              | `1 2 show` -> "1 2" |
| `showstack` | 0        | Alias for `show`                    | Same as `show`  |
| `print`     | 0        | Alias for `show`                    | Same as `show`  |
| `clear`     | 0        | Clear all variables and constants   | See below       |

---

## `dup` — Duplicate Top Value

Pushes a copy of the top stack value onto the stack. Requires at least one
value on the stack.

### Stack visualization

```
Input:  5  dup

  Step  Stack
  ----  -----
   5    [5]
   dup  [5, 5]
```

### Examples

Double a value (multiply by 2):

```
$ gt '5 dup +'
10
```

Square a value:

```
$ gt '7 dup *'
49
```

Triple-duplicate for multi-use:

```
$ gt '7 dup dup +'
7 14
```

Step by step: `7` pushed, `dup` produces `[7, 7]`, `dup` produces `[7, 7, 7]`,
`+` pops the top two producing `[7, 14]`.

### Practical use cases

**Self-comparison** — compare a value against a computed result:

```
$ gt '10 dup 5 gt'
true
```

Checks if `10 > 5` while keeping the original `10` available (the comparison
pops both values, but `dup` ensures the original is preserved before the
operation).

**Reuse without re-entry** — avoid typing the same value multiple times:

```
$ gt '3.14159 dup 2 * dup 2 *'
3.14159 12.56636
```

Computes `2*pi` and `4*pi` from a single entry of pi.

---

## `swap` — Swap Top Two Values

Exchanges the positions of the top two stack values. Requires at least two
values on the stack.

### Stack visualization

```
Input:  3  4  swap

  Step  Stack
  ----  -----
   3    [3]
   4    [3, 4]
   swap [4, 3]
```

### Examples

Reverse subtraction order:

```
$ gt '3 4 swap -'
1
```

Without swap, `3 4 -` = `3 - 4` = -1. With swap, it becomes `4 - 3` = 1.

Reorder for division:

```
$ gt '2 10 swap /'
5
```

Swaps to get `10 / 2` = 5 instead of `2 / 10` = 0.2.

Swap in a multi-element stack (only affects top two):

```
$ gt '3 4 5 swap -'
3 1
```

Step by step: `3` pushed `[3]`, `4` pushed `[3, 4]`, `5` pushed `[3, 4, 5]`,
`swap` swaps top two `[3, 5, 4]`, `-` pops 5 and 4 pushing 1 `[3, 1]`.

### Practical use cases

**Correct operand order** — when the natural reading order differs from RPN
order:

```
$ gt '100 1 swap /'
0.01
```

`swap` turns `[100, 1]` into `[1, 100]` then `/` computes `1 / 100` = 0.01.

**Reordering before operations** — swap to position the right operand on top:

```
$ gt '4 3 swap ^'
81
```

Computes `3^4` = 81 (swap turns `[4, 3]` into `[3, 4]`, then `^` computes
`3^4`).

---

## `pop` — Discard Top Value

Removes and discards the top stack value. Requires at least one value on the
stack.

### Stack visualization

```
Input:  5  6  pop

  Step  Stack
  ----  -----
   5    [5]
   6    [5, 6]
   pop  [5]
```

### Examples

Discard an extra value:

```
$ gt '5 6 pop'
5
```

Use `pop` to keep only the bottom value after a computation:

```
$ gt '10 20 + 30 pop'
30
```

Step by step: `10` pushed, `20` pushed, `+` produces `[30]`, `30` pushed
`[30, 30]`, `pop` discards top `[30]`.

### Practical use cases

**Discard unwanted results** — remove intermediate values:

```
$ gt '1 2 3 + + pop 42'
42
```

**Clean the stack** — when you need a fresh top without clearing everything:

```
$ gt '100 200 pop'
100
```

---

## `show` / `showstack` / `print` — Display the Stack

Prints all values currently on the stack without modifying it. All three are
equivalent. Requires no operands.

### Stack visualization

```
Input:  5  6  show

  Step    Stack    Output
  ----    -----    ------
   5      [5]
   6      [5, 6]
   show   [5, 6]   "5 6"
```

### Examples

Display current stack state:

```
$ gt '5 6 show'
5 6
```

Multiple values:

```
$ gt '10 20 30 show'
10 20 30
```

Empty stack:

```
$ gt 'show'
Stack is empty
```

Using the aliases:

```
$ gt '10 20 30 showstack'
10 20 30

$ gt '10 20 30 print'
10 20 30
```

### Practical use cases

**Debugging** — verify the stack state before the next operation:

```
$ gt '1 2 3 + + show 4 *'
```

This shows `[6]` after the additions, confirming the result before multiplying.

**Intermediate inspection** — chain with show to trace evaluation:

```
$ gt '100 10 / show 5 +'
```

Displays `10` (the result of `100 / 10`) before adding 5.

---

## `clear` — Reset Variables and Constants

Clears all user-defined variables and constants from the session. This affects
the variable table, not the RPN stack directly.

### Examples

Clear all user-defined variables:

```
$ gt 'clear'
All variables cleared
```

### Practical use cases

**Clean slate** — reset a session after defining many temporary variables:

```
$ gt 'clear'
All variables cleared
```

---

## Error Handling

Each stack operator validates that the stack has sufficient operands:

```
$ gt 'dup'
Error: stack is empty

$ gt 'swap'
Error: stack has insufficient operands

$ gt 'pop'
Error: stack is empty
```

`show` / `print` / `showstack` handle an empty stack gracefully:

```
$ gt 'show'
Stack is empty
```

---

## Combined Examples

### Duplicate, compute, and keep original

```
$ gt '5 dup dup + *'
75
```

Step by step: `[5]` -> `dup` -> `[5, 5]` -> `dup` -> `[5, 5, 5]` -> `+` ->
`[5, 10]` -> `*` -> `[50]`. This computes `5 * (5 + 5)` = 50.

### Swap-based reordering in complex expressions

```
$ gt '2 3 4 swap ^ -'
-62
```

Step by step: `[2]` -> `[2, 3]` -> `[2, 3, 4]` -> `swap` -> `[2, 4, 3]` ->
`^` -> `[2, 4^3]` = `[2, 64]` -> `-` -> `[2 - 64]` = `[-62]`.

### Show the intermediate state

```
$ gt '10 20 30 + show 5 *'
```

This pushes 10, 20, 30, adds top two (50), shows the stack `[10, 50]`, then
multiplies to get `500`.
