# Symbols (:x Syntax)

Symbols are named placeholders on the RPN stack. Unlike variables, which store resolved numeric values, symbols carry a name and remain unresolved until explicitly bound through an assignment operation.

Every bare identifier that has no matching variable or constant is automatically pushed as a symbol. The `:` prefix makes this explicit.

## Syntax

### `:x` — Push a Symbol

Prefix a single-letter identifier with `:` to push it as a symbol onto the stack:

```
:x          → :x
:A          → :A
:_          → :_
```

The identifier after `:` must be a single letter, uppercase or lowercase, or an underscore. Multi-character names and names starting with digits are not allowed:

```
:abc        → error: unknown token ':abc'
:1          → error: unknown token ':1'
:x!         → error: unknown token ':x!'
```

### `:` (bare colon) — Error

A colon with no identifier following it is rejected:

```
:           → error: symbol name cannot be empty after colon
```

## What Symbols Are

A symbol is a stack value that holds a name rather than a numeric value. When displayed, symbols are always prefixed with `:` to distinguish them from plain numbers and other values.

Internally, symbols implement the `StackValue` interface (so they can be placed on the stack) but not `NumericValue` (so they cannot participate in arithmetic). This means symbols work with stack-manipulation operators but not with arithmetic operators:

```
:x dup              → :x :x         (dup works on symbols)
:x 5 :x swap        → 5 :x :x       (swap works on symbols)
:x :y :z show       → :x :y :z      (show displays all stack values)
:x 1 +              → error: value ":x" is not numeric
```

## Symbols vs Variables

The key difference is **binding time**:

| | Variable | Symbol |
|---|---|---|
| **Binding** | Immediate — resolved to a value at use | Delayed — name pushed as-is |
| **Stack value** | Numeric (`Float` or `Rat`) | `Symbol` |
| **Arithmetic** | Works | Not supported |
| **Assignment** | Assigned with `:=`, `=:`, `=` | Used as the name side of `:=` or `=:` |

### Immediate Binding (Variables)

When a variable is defined and then referenced by bare name, the variable's current value is pushed:

```
x 5 :=          → x = 5
x               → 5           (pushes the value 5)
x x +           → 10          (each x resolves to 5)
```

### Delayed Binding (Symbols)

When the same name is used as a symbol (with `:` prefix), the name itself is pushed regardless of whether the variable exists:

```
x 5 :=          → x = 5
:x              → :x          (pushes the symbol, not the value 5)
:x dup          → :x :x       (symbol can be duplicated)
```

An undefined bare name also pushes a symbol — `x` and `:x` are equivalent when `x` is unbound:

```
x               → :x          (unbound, shown as symbol)
:x              → :x          (explicit symbol)
```

## Using Symbols with Assignment

Symbols are primarily useful as the name side of assignment operators. The `:` prefix makes the intent clear:

### Right Assignment (`:=`)

```
:x 5 :=         → x = 5
:a 10 :=        → a = 10
```

This is equivalent to `x 5 :=` when `x` is unbound (the bare name also pushes a symbol), but the `:` prefix is explicit and works the same way regardless of whether a variable named `x` already exists.

### Left Assignment (`=:`)

```
5 :x =:         → x = 5
10 :a =:        → a = 10
```

## Using Symbols for Variable Deletion

The `d` operator deletes a variable by name. Push the variable name as a symbol, then apply `d`:

```
x 5 :=            → x = 5
:x d              → (deletes x, no output)
x                 → :x        (x is now unbound)
```

Attempting to delete a variable that does not exist produces an error:

```
:x d              → error: variable not found: x
```

## Symbols on the Stack

Multiple symbols can coexist on the stack alongside numbers and other values:

```
:x :y :z          → :x :y :z
5 :x              → 5 :x
true :x show      → true :x
:x 5 3 +          → :x 8         (3+4=7 happens, :x stays)
```

### Stack Operations

Stack operators that work on any `StackValue` work with symbols:

```
:x dup              → :x :x
:x 5 :x swap        → 5 :x :x     (swaps top two items)
5 :x :y swap        → 5 :y :x
```

`pop` removes the top item but leaves the stack empty, which is reported as an error unless combined with an assignment:

```
:x pop              → error: empty result: expression evaluated to nothing
```

## Practical Use Cases

### Explicit Variable Assignment

Using `:` prefix makes assignment intent clear and unambiguous:

```
:x 100 :=
:y 200 :=
:x :y show          → :x :y        (symbols pushed, not values)
x y +               → 300          (variables resolved to values)
```

### Delayed Binding for Symbolic Manipulation

Symbols can be pushed and rearranged on the stack before being used. This is useful when you want to defer resolution or manipulate the stack structure:

```
:x :y :z            → :x :y :z
:x :y :z dup        → :x :y :z :z
```

### Safe Variable Reference

When a variable might or might not be defined, using `:` prefix always produces a symbol, making the behavior predictable:

```
x                   → :x or 5      (depends on whether x was defined)
:x                  → :x           (always a symbol)
```

### Variable Cleanup

Combine symbols with `d` for explicit variable deletion:

```
x 10 :=             → x = 10
:x d                → (x deleted)
vars                → No variables defined
```

## Error Summary

| Input | Error | Reason |
|-------|-------|--------|
| `:` | symbol name cannot be empty after colon | No identifier after `:` |
| `:abc` | unknown token ':abc' | Multi-character identifiers not allowed |
| `:1` | unknown token ':1' | Identifier starts with digit |
| `:x!` | unknown token ':x!' | Special character in name |
| `:x 1 +` | value ":x" is not numeric | Symbols can't be used in arithmetic |
| `:x d` | variable not found: x | Variable x not defined |
| `:x pop` | empty result: expression evaluated to nothing | Stack empty after pop |

## Reference

- **Type**: `internal/rpn/number.go` — `Symbol` struct and `NewSymbol()`
- **Parsing**: `internal/rpn/rpn_parse.go` — `checkAndPushSymbol()`, `dispatchOperator()`
- **Assignment**: `internal/rpn/operations_variables.go` — `AssignLeft()`, `AssignRight()` (both accept `*Symbol` as the name argument)
- **Deletion**: `internal/rpn/operations_variables.go` — `DeleteVariable()` (called via the `d` operator, expects `*Symbol` or `*StringNum`)
