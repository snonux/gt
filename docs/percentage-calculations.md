# Percentage Calculations

gt supports three forms of percentage calculation, each solving a different unknown. All inputs are case-insensitive and support decimal numbers.

## Syntax

### Form 1: X% of Y

Find what a given percentage of a base value equals.

```bash
gt '20% of 150'
gt 'what is 20% of 150'
gt '20% 150'           # "of" is optional
```

### Form 2: X is what % of Y

Find what percentage one number represents of another.

```bash
gt '30 is what % of 150'
gt '30 is what% of 150'    # spaces around % are optional
gt '30 is what % 150'      # "of" is optional
```

### Form 3: X is Y% of what

Find the whole number when given a part and the percentage.

```bash
gt '30 is 20% of what'
gt '30 is 20% what'        # "of" is optional
```

## Examples

### Form 1: X% of Y

```bash
gt '20% of 150'
# 20.00% of 150.00 = 30.00
#   Steps: (20.00 / 100) * 150.00 = 0.20 * 150.00 = 30.00

gt '50% of 200'
# 50.00% of 200.00 = 100.00
#   Steps: (50.00 / 100) * 200.00 = 0.50 * 200.00 = 100.00
```

### Form 2: X is what % of Y

```bash
gt '30 is what % of 150'
# 30.00 is 20.00% of 150.00
#   Steps: (30.00 / 150.00) * 100 = 0.20 * 100 = 20.00%

gt '75 is what % of 300'
# 75.00 is 25.00% of 300.00
#   Steps: (75.00 / 300.00) * 100 = 0.25 * 100 = 25.00%
```

### Form 3: X is Y% of what

```bash
gt '30 is 20% of what'
# 30.00 is 20.00% of 150.00
#   Steps: (30.00 / 20.00) * 100 = 1.50 * 100 = 150.00

gt '50 is 25% of what'
# 50.00 is 25.00% of 200.00
#   Steps: (50.00 / 25.00) * 100 = 2.00 * 100 = 200.00
```

## Output Format

Every percentage calculation returns two lines:

1. **Result line** — the answer formatted to two decimal places
2. **Steps line** — the intermediate arithmetic broken down step by step

| Form | Result format |
|------|--------------|
| X% of Y | `X.00% of Y.00 = Result.00` |
| X is what % of Y | `X.00 is P.00% of Y.00` |
| X is Y% of what | `X.00 is Y.00% of Whole.00` |

The steps line shows the formula used, intermediate values, and the final result so you can verify the math.

## Examples

### Discounts

```bash
gt '15% of 89.99'
# 15.00% of 89.99 = 13.50
#   Steps: (15.00 / 100) * 89.99 = 0.15 * 89.99 = 13.50
# Discount: $13.50 off, you pay $76.49
```

### Tax

```bash
gt '8.25% of 47.50'
# 8.25% of 47.50 = 3.92
# Tax on a $47.50 purchase at 8.25% rate: $3.92
```

### Tips

```bash
gt '18% of 63.40'
# 18.00% of 63.40 = 11.41
# 18% tip on a $63.40 bill: $11.41
```

### Growth Rates

```bash
gt '120 is what % of 100'
# 120.00 is 120.00% of 100.00
# Revenue grew from $100 to $120, that's a 20% increase

gt '85 is what % of 100'
# 85.00 is 85.00% of 100.00
# Revenue dropped from $100 to $85, that's a 15% decrease
```

### Finding the Original Amount

```bash
gt '170 is 85% of what'
# 170.00 is 85.00% of 200.00
# After a 15% discount you paid $170. The original price was $200.
```

## Edge Cases

### Zero percent

```bash
gt '0% of 150'
# 0.00% of 150.00 = 0.00
```

### Percentages over 100%

```bash
gt '150% of 200'
# 150.00% of 200.00 = 300.00

gt '200% of 50'
# 200.00% of 50.00 = 100.00
```

### Decimal values

```bash
gt '7.5% of 200'
# 7.50% of 200.00 = 15.00

gt '3.75 is what % of 150'
# 3.75 is 2.50% of 150.00
```

### Large numbers

```bash
gt '15% of 1000000'
# 15.00% of 1000000.00 = 150000.00
```

### Case insensitivity

```bash
gt 'WHAT IS 20% OF 150'
gt '30 IS WHAT % OF 150'
gt '30 IS 20% OF WHAT'
```

### Whitespace flexibility

```bash
gt '20  %   of   150'
gt '30 is what% of 150'
```

### Not supported

- **Negative numbers** — `-20% of 150` is not recognized; only non-negative numbers are accepted
- **Division by zero** — `30 is what % of 0` and `30 is 0% of what` both produce an error since they require dividing by zero

## Known Limitations

1. **No negative numbers** — The regex patterns only match non-negative numeric values (`\d+(?:\.\d+)?`). Expressions like `-20% of 150` will not be recognized as percentage calculations.

2. **Single-line input** — Percentage calculations are single-command operations; they do not interact with the RPN stack.

3. **No chained operations** — You cannot combine percentage calculations with RPN in a single expression. Use RPN for the arithmetic after getting the percentage result.

4. **Zero percent with "of what"** — `X is 0% of what` is mathematically undefined and returns a division-by-zero error.
