# Constants

The gt calculator includes a comprehensive set of built-in mathematical constants that can be used directly in RPN expressions. Constants are resolved automatically — just use the name as a token.

## Built-in Constants

### Fundamental Constants

| Constant | Value | Description |
|----------|-------|-------------|
| `pi` / `π` | 3.141592654 | Ratio of circle circumference to diameter |
| `tau` / `τ` | 6.283185307 | Full circle (2π) |
| `e` / `euler` | 2.718281828 | Base of natural logarithm |
| `phi` / `φ` | 1.618033989 | Golden ratio |

### Square Roots

| Constant | Value | Description |
|----------|-------|-------------|
| `sqrt2` / `√2` | 1.414213562 | √2 |
| `sqrt3` / `√3` | 1.732050808 | √3 |
| `sqrt5` / `√5` | 2.236067977 | √5 |

### Logarithms

| Constant | Value | Description |
|----------|-------|-------------|
| `ln2` / `log2` | 0.6931471806 | Natural log of 2 |
| `ln10` / `log10` | 2.302585093 | Natural log of 10 |
| `log_e` | 0.4342944819 | Log base e of 10 |
| `log_e10` | 0.4342944819 | Same as log_e |

### Reciprocals

| Constant | Value | Description |
|----------|-------|-------------|
| `1/e` / `inv_e` | 0.3678794412 | Reciprocal of e |
| `1/π` / `inv_pi` | 0.3183098862 | Reciprocal of π |

### Special Values

| Constant | Value | Description |
|----------|-------|-------------|
| `inf` / `infinity` | +Inf | Positive infinity |
| `-inf` / `-infinity` | -Inf | Negative infinity |
| `nan` | NaN | Not a Number |

## Usage

Constants are used directly as tokens in RPN expressions:

```bash
gt 'pi'              # → 3.141592654
gt 'pi 2 *'          # → 6.283185307  (2π)
gt 'e phi +'         # → 4.336315817  (e + φ)
gt 'sqrt2 2 ^'       # → 2  ((√2)²)
gt 'tau 4 /'         # → 1.570796327  (π/2)
```

Greek letter aliases work identically:

```bash
gt 'π'               # → 3.141592654
gt 'φ 2 *'           # → 3.236067977
gt '√2 √3 *'         # → 2.449489743  (√6)
```

## Commands

### List all constants

```bash
gt 'constants'
```

Lists all 36 built-in constants with their values, sorted alphabetically.

### Clear user-defined constants

```bash
gt 'clearconstants'
```

Removes all user-defined constants and resets built-in constants to their default values. Built-in constants cannot be permanently deleted.

## Practical Use Cases

### Geometry

```bash
gt 'pi 5 5 * *'           # Circle area: π × 5² = 78.54
gt '2 pi 5 *'             # Circle circumference: 2π × 5 = 31.42
gt '4 3 pi 5 5 5 * * * /' # Sphere volume: 4/3 π r³ = 523.6
```

### Engineering

```bash
gt 'e 1 -'                # 1/e ≈ 0.632 (RC circuit time constant)
gt 'sqrt2 1000 *'         # RMS voltage: √2 × 1000 = 1414.2
gt 'phi 10 *'             # Golden rectangle: φ × 10 = 16.18
```

### Information Theory

```bash
gt 'ln2'                  # ≈ 0.693 (bits to nats conversion)
gt '2 ln2 /'              # ≈ 1.443 (nats to bits)
```

## Edge Cases

### Infinity arithmetic

```bash
gt 'inf 5 +'              # → +Inf
gt 'inf inf +'            # → +Inf
gt 'inf -inf +'           # → NaN
```

### NaN propagation

```bash
gt 'nan 5 +'              # → NaN
gt '0 0 /'                # → NaN (if supported)
```

### Variable shadowing

If you assign a variable with the same name as a constant, the variable takes precedence:

```bash
gt 'pi 3 = pi'           # → 3  (variable shadows constant)
gt 'clear'               # Clear variables
gt 'pi'                  # → 3.141592654  (constant restored)
```

## Notes

- Constants are resolved before variables in the lookup order
- All constants use float64 precision (10 significant digits in output)
- Greek letter aliases are fully supported for the most common constants
- Built-in constants cannot be individually deleted, only overwritten with variables
