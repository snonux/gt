# Overall Plan for Implementing Variable Symbol Handling in gt

## Goal
Modify the gt interpreter so that:
1. Users can push a variable *symbol* onto the stack using `:x` syntax.
2. Bare identifier `x` pushes the symbol when the variable is unbound, otherwise pushes its value.
3. Existing assignment operators `:=` and `=:` remain unchanged (they will receive either a symbol from `:x` or from an undefined identifier).

## Changes Required

### 1. Lexer
- Recognize a leading colon (`:`) followed by an identifier as a `TOKEN_SYMBOL`.
- When seen, push a `CELL_SYMBOL` cell holding the identifier string onto the stack.

### 2. Evaluator (identifier handling)
- For a bare identifier token:
  - Attempt `env_get(name)`.
  - If a binding exists → push the value cell.
  - If no binding exists → push a `CELL_SYMBOL` cell (treat as symbol).

### 3. Assignment Operators
- Keep `:=` and `=:`) handlers unchanged.
- They now receive either a symbol cell (from `:x` or from undefined identifier) and a value cell in the expected order.

### 4. REPL / Debugging
- When printing the stack, prefix symbol cells with `:` (e.g., `:x`) to distinguish from values.

### 5. Testing
- Verify the behavior with a test suite covering:
  - `:x` pushes symbol.
  - `x` unbound → pushes symbol.
  - `x` bound → pushes value.
  - `:x 10 :=` binds x to 10.
  - `10 :x =:` binds x to 10 (value‑first).
  - Error cases where needed.

## Tasks (managed via `ask`)
Each task below will be added via `ask add` and tagged with `plan:PLAN.md` for traceability.

