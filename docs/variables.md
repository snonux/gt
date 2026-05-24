# Variables and Assignment Operators

The gt calculator supports named variables that store numeric values. Variables persist across expressions within a single RPN session (REPL or multi-expression input), enabling reusable values and multi-step calculations.

## Assignment Operators

Three assignment operators are available, differing in token order:

| Operator | Syntax | Description | Example |
|----------|--------|-------------|---------|
| `:=` | `name value :=` | Right assignment | `x 5 :=` → `x = 5` |
| `=:` | `value name =:` | Left assignment | `10 y =:` → `y = 10` |
| `=` | `name value =` or `name = value` | Standard assignment | `x 10 =` or `x = 10` |

### `:=` (Right Assignment)

The value appears on the right (top of stack). This is the most idiomatic RPN form.

```
x 5 :=          → x = 5
price 9.99 :=   → price = 9.99
```

### `=:` (Left Assignment)

The value appears on the left (bottom of stack). Useful when you think in "value first" order.

```
42 answer =:    → answer = 42
3.14159 pi =:   → pi = 3.14159
```

### `=` (Standard Assignment)

Supports both infix-style (`name = value`) and postfix-style (`name value =`) syntax. Also supports expression continuation — you can append an expression after the assignment to evaluate immediately.

```
x = 5           → x = 5
x 10 =          → x = 10
x 10 = x 5 +    → 15  (assigns x=10, then evaluates x+5)
```

## Variable Management

### `vars` — List Variables

Shows all defined variables, sorted alphabetically with their values.

```
vars
→ No variables defined          (when none exist)
→ a = 10
  b = 3                         (when variables are defined)
```

### `clear` — Clear All Variables

Removes all defined variables at once.

```
clear
→ All variables cleared
```

### `d` — Delete a Variable

Removes a single variable. The variable name must be pushed as a symbol (using `:` prefix) onto the stack before `d`.

```
x 5 :=            → x = 5
:x d              → (deletes x)
x                 → :x    (x is now undefined, shown as symbol)
```

Attempting to delete a non-existent variable returns an error:

```
:nonexistent d    → error: variable not found: nonexistent
```

## Variable Lifecycle

1. **Creation**: A variable is created when a value is assigned using `:=`, `=:`, or `=`.
2. **Usage**: Reference the variable by name in an expression; its value is pushed onto the stack.
3. **Reassignment**: Assigning to an existing variable overwrites its previous value.
4. **Deletion**: Remove with `:name d` or clear all with `clear`.
5. **Session scope**: Variables persist only within a single RPN session. Each CLI invocation starts fresh.

## Using Variables in Expressions

Once defined, variables are used by name in expressions. The variable's value is pushed onto the stack for computation.

### Single Variable Reuse

```
x 5 :=
x x +         → 10
x 2 *         → 10
```

### Multi-Variable Expressions

```
a 10 :=
b 3 :=
a b +         → 13
a b *         → 30
```

### Chained Assignments

Multiple variables can be assigned in a single expression:

```
a 10 := b 3 := c 2 :=
a b + c *     → 26
```

### Variable Reassignment

Variables can be reassigned to new values:

```
x 5 :=          → x = 5
x 10 :=         → x = 10
x               → 10
```

### Assignment with Expression Continuation

The `=` operator supports evaluating an expression immediately after assignment:

```
x 10 = x 5 +    → 15   (assigns x=10, then computes x+5)
```

## Practical Use Cases

### Storing Reusable Values

Store constants or frequently-used values to avoid retyping:

```
pi 3.14159265 :=
radius 7 :=
2 pi radius *          → 43.9822971   (circumference)
```

### Multi-Step Calculations

Break complex calculations into named intermediate steps:

```
base 100 :=
tax_rate 0.08 :=
tax base tax_rate *    → 8
total base tax +       → 108
```

### Iterative Refinement

Update a running value across multiple operations:

```
sum 0 :=
sum 10 +               → 10
sum 20 +               → 30
sum 5 -                → 25
```

## Reference

- **Implementation**: `internal/rpn/operations_variables.go` (assignment operators, variable CRUD)
- **Parsing**: `internal/rpn/rpn_parse.go` (`handleAssignmentOp()`, `handleStandardAssign()`)
- **Registry**: `internal/rpn/operator_registry.go` (operator registration for `:=`, `=:`, `=`, `d`)
- **Variable store**: `internal/rpn/variables.go` (thread-safe variable storage with save/load)
