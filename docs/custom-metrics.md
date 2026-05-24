# Custom Metrics

The `custom` command lets you define your own units of measurement. Custom metrics behave like built-in metrics — they support suffix notation (`10foobar`), metric-aware arithmetic, and conversion to their base unit. They are particularly useful for domain-specific units that gt doesn't ship with out of the box.

Custom metrics are RPN operators — they work on the stack, consume tokens, and produce string output. They can be used in single-command mode (`./gt '<expression>'`) or interactively in REPL mode (`./gt`).

## `custom define`

Define a new custom metric with a name, conversion factor, and category.

```bash
gt 'custom define <name> <factor> <category>'
```

- **name**: Unit name (e.g., `foobar`, `widget`, `mygig`). Must not conflict with existing metrics.
- **factor**: Numeric conversion factor to the category's base unit. Can be positive, negative, or zero.
- **category**: One of the valid metric categories: `Custom`, `DataRate`, `DataSize`, `Distance`, `Speed`, `Time`, `Universal`, or `Weight`.

```bash
gt 'custom define foobar 42 Custom'
# → defined custom metric "foobar" (factor: 42, category: Custom)

gt 'custom define myhour 3600 Time'
# → defined custom metric "myhour" (factor: 3600, category: Time)

gt 'custom define widget 0.5 Universal'
# → defined custom metric "widget" (factor: 0.5, category: Universal)

gt 'custom define negtest -5 Custom'
# → defined custom metric "negtest" (factor: -5, category: Custom)

gt 'custom define zerounit 0 Custom'
# → defined custom metric "zerounit" (factor: 0, category: Custom)
```

The factor determines how many base units one unit of the custom metric equals. For example, `foobar` with factor 42 means `1foobar = 42 Custom_base`.

### REPL workflow

In single-command mode, `custom define` stops token evaluation after registering the metric, so subsequent tokens in the same expression are not processed. Use REPL mode to define a metric and then use it in calculations:

```
> custom define foobar 42 Custom
defined custom metric "foobar" (factor: 42, category: Custom)
> 10foobar
420
> 10foobar 5foobar +
630
> custom undefine foobar
removed custom metric "foobar"
```

In single-command mode, the define command confirms the metric was registered:

```bash
gt 'custom define foobar 42 Custom'
# → defined custom metric "foobar" (factor: 42, category: Custom)
```

## `custom undefine`

Remove a previously defined custom metric. Built-in metrics cannot be removed.

```bash
gt 'custom undefine <name>'
```

```bash
gt 'custom undefine foobar'
# → removed custom metric "foobar"
```

### Errors

Attempting to undefine a built-in metric or a metric that doesn't exist produces an error:

```bash
gt 'custom undefine Cool'
# → Error: rpn: custom undefine: cannot remove built-in metric "Cool"

gt 'custom undefine nonexistent'
# → Error: rpn: custom undefine: metric "nonexistent" not found
```

## `custom list`

List all currently defined custom metrics. Returns an empty message if none are defined.

```bash
gt 'custom list'
# → no custom metrics defined
```

In REPL mode, after defining metrics:

```
> custom define foobar 42 Custom
defined custom metric "foobar" (factor: 42, category: Custom)
> custom define widget 0.5 Universal
defined custom metric "widget" (factor: 0.5, category: Universal)
> custom list
foobar, widget
```

## `custom show`

Show detailed information about a custom metric. Without a name, shows all custom metrics.

```bash
gt 'custom show [name]'
```

### Show all custom metrics

```bash
gt 'custom show'
# → no custom metrics defined
```

In REPL mode after defining metrics:

```
> custom define foobar 42 Custom
defined custom metric "foobar" (factor: 42, category: Custom)
> custom show
  foobar, category: Custom, base: Custom_base, factor: 42
```

### Show a specific metric

```
> custom show foobar
foobar, category: Custom, base: Custom_base, factor: 42
```

### Errors

```
> custom show nonexistent
Error: rpn: custom show: unknown custom metric "nonexistent"

> custom show Mbps
Error: rpn: custom show: metric "Mbps" is not a custom metric
```

## Using Custom Metrics in Calculations

Once defined in REPL mode, custom metrics work with all standard metric operations: suffix notation, arithmetic, and conversion.

### Suffix notation

Attach the custom unit name directly to a number:

```
> custom define foobar 42 Custom
defined custom metric "foobar" (factor: 42, category: Custom)
> 10foobar
420
```

### Addition and subtraction

Custom metrics in the same category are compatible for addition and subtraction. The result uses the unit of the bottom-of-stack operand:

```
> custom define foobar 42 Custom
defined custom metric "foobar" (factor: 42, category: Custom)
> 10foobar 5foobar +
15
> 5foobar 3foobar -
2
```

The result is unitless because the metric system resolves to base units for the calculation. For factor 42: `10foobar = 420`, `5foobar = 210`, so `420 + 210 = 630` which is `15 foobar` in base units. The displayed value is the quotient in terms of the result's unit.

### Multiplication and division

Custom metrics work with multiplication and division:

```
> custom define mul_unit 5 Custom
defined custom metric "mul_unit" (factor: 5, category: Custom)
> 3mul_unit 2 *
30
> custom define div_unit 10 Custom
defined custom metric "div_unit" (factor: 10, category: Custom)
> 20div_unit 4div_unit /
5
```

### Adding Cool (unitless) values

Unitless values (`Cool`/Universal) can be added to custom metrics:

```
> custom define foobar 42 Custom
defined custom metric "foobar" (factor: 42, category: Custom)
> 10foobar 5 +
10
```

### Incompatible categories

Custom metrics in different categories cannot be mixed:

```
> custom define timetest 3600 Time
defined custom metric "timetest" (factor: 3600, category: Time)
> custom define disttest 1609 Distance
defined custom metric "disttest" (factor: 1609, category: Distance)
> 1timetest 1disttest +
Error: metric arithmetic requires compatible categories (Time vs Distance)
```

### Conversion

Custom metrics support `@<unit> convert` within their category. Conversion to a different category fails:

```
> custom define foobar 42 Custom
defined custom metric "foobar" (factor: 42, category: Custom)
> 10foobar @Custom_base convert
420
> 10foobar @Mbps convert
Error: incompatible categories for conversion
```

### Hyper operations

Custom metrics work with hyper operations like `[+]`:

```
> custom define hyper_test 100 Custom
defined custom metric "hyper_test" (factor: 100, category: Custom)
> 10hyper_test 5hyper_test 3hyper_test [+]
18
```

## Practical Use Cases

### Company-specific units

Define internal measurement units that aren't standard:

```
> custom define reel 304.8 meters
defined custom metric "reel" (factor: 304.8, category: Distance)
> 5reel @m convert
1524
```

### Recipe and cooking units

Create convenient cooking measurements:

```
> custom define cup 240 Weight
defined custom metric "cup" (factor: 240, category: Weight)
> 3cup 2cup +
5
> 3cup @g convert
720
```

### Game currencies and points

Define game-specific units for calculating rewards:

```
> custom define gold 1000 Universal
defined custom metric "gold" (factor: 1000, category: Universal)
> 5gold 300 +
5300
```

### Time conversions with custom units

Create convenient time units:

```
> custom define fortnight 1209600 Time
defined custom metric "fortnight" (factor: 1209600, category: Time)
> 2fortnight @day convert
28
```

## Edge Cases

### Duplicate names

Defining a metric with a name that already exists (including built-in metrics) fails:

```bash
gt 'custom define Mbps 1000 DataRate'
# → Error: rpn: custom define: metric "Mbps" already exists

gt 'custom define foobar 42 Custom'
gt 'custom define foobar 99 Custom'
# → Error: rpn: custom define: metric "foobar" already exists
```

In REPL mode, the second define on the same RPN instance fails:

```
> custom define foobar 42 Custom
defined custom metric "foobar" (factor: 42, category: Custom)
> custom define foobar 99 Custom
Error: rpn: custom define: metric "foobar" already exists
```

### Invalid category

The category must be one of the valid categories. Case-sensitive matching is required:

```bash
gt 'custom define foo 1 Nope'
# → Error: rpn: custom define: unknown category "Nope"

gt 'custom define foo 1 custom'
# → Error: rpn: custom define: unknown category "custom"
```

Valid categories: `Custom`, `DataRate`, `DataSize`, `Distance`, `Speed`, `Time`, `Universal`, `Weight`.

### Invalid factor

The factor must be a valid number:

```bash
gt 'custom define foo abc Custom'
# → Error: rpn: custom define: invalid factor "abc"
```

### Factor zero

Defining a metric with factor 0 is allowed, but arithmetic involving it produces `NaN` (division by zero when converting from base units):

```
> custom define zero 0 Custom
defined custom metric "zero" (factor: 0, category: Custom)
> 5zero 3 +
NaN
```

### Negative factor

Negative factors are allowed:

```
> custom define negtest -5 Custom
defined custom metric "negtest" (factor: -5, category: Custom)
> 2negtest 1 +
3
```

### Missing arguments

```bash
gt 'custom define'
# → Error: rpn: custom define: usage: custom define <name> <factor> <category>

gt 'custom undefine'
# → Error: rpn: custom undefine: usage: custom undefine <name>
```

### Unknown subcommand

```bash
gt 'custom rename foobar'
# → Error: rpn: unknown custom subcommand "rename". Use: show, list, define, undefine
```

## REPL vs Single-Command Mode

Custom metric commands behave differently depending on how gt is invoked:

### Single-command mode (`./gt '<expression>'`)

Each invocation creates a fresh process with a clean metric registry. `custom define` confirms the metric was registered but stops evaluation, so subsequent tokens are not processed:

```bash
gt 'custom define foobar 42 Custom'
# → defined custom metric "foobar" (factor: 42, category: Custom)

# This does NOT work — custom define stops evaluation:
gt 'custom define foobar 42 Custom 10foobar 5foobar +'
# → defined custom metric "foobar" (factor: 42, category: Custom)
```

### REPL mode (`./gt` interactively)

The RPN engine and metric registry persist across lines. Define metrics first, then use them in subsequent calculations:

```
> custom define foobar 42 Custom
defined custom metric "foobar" (factor: 42, category: Custom)
> 10foobar 5foobar +
15
> 10foobar @Custom_base convert
420
> metric show
foobar, category: Custom, base: Custom_base, factor: 42
> custom list
foobar
> custom undefine foobar
removed custom metric "foobar"
> custom list
no custom metrics defined
```

REPL mode is the recommended way to work with custom metrics, as it allows defining, using, and managing metrics in a single session.

## Summary

| Command | Description | Arguments |
|---------|-------------|-----------|
| `custom define <name> <factor> <category>` | Define a custom metric | name, factor, category |
| `custom undefine <name>` | Remove a custom metric | name |
| `custom list` | List all custom metrics | None |
| `custom show` | Show all custom metrics | None |
| `custom show <name>` | Show a specific metric | name |

Valid categories: `Custom`, `DataRate`, `DataSize`, `Distance`, `Speed`, `Time`, `Universal`, `Weight`
