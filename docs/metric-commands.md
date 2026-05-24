# Metric Commands

The `metric` command provides runtime introspection and configuration for gt's metrics system. It lets you inspect metric metadata, explore available categories and units, check compatibility between values, and switch between SI and IEC prefix modes.

All metric commands are RPN operators — they work on the stack, consume tokens, and produce string output. They work in single-command mode (`gt '<expression>'`) or interactively in REPL mode (`gt`).

## `metric show`

Shows metric info for the value at the top of the stack. Displays the metric name, category, base unit, and conversion factor.

```bash
gt '100Mbps metric show'
# → Mbps, DataRate, base: bps, factor: 1e+06

gt '5hr metric show'
# → hr, Time, base: seconds, factor: 4e+03

gt '70kg metric show'
# → kg, Weight, base: kilograms, factor: 1
```

Unitless values (the `Cool` metric) are identified as Universal:

```bash
gt '5 metric show'
# → Cool (Universal)
```

If the stack is empty, `metric show` returns an error:

```bash
gt 'metric show'
# → Error: rpn: metric show: metric show: stack is empty
```

### REPL mode

In REPL mode, `metric show` inspects whatever value is currently on top of the stack:

```
> 42
42
> metric show
Cool (Universal)
> 100kmh
100
> metric show
kmh, Speed, base: m/s, factor: 0.2777777778
```

`metric show` does not consume the value from the stack — it only peeks at the top.

## `metric list`

Lists all available metric categories, alphabetically sorted and comma-separated. Requires no arguments and does not interact with the stack.

```bash
gt 'metric list'
# → DataRate, DataSize, Distance, Speed, Time, Universal, Weight
```

## `metric <category>`

Lists all metric names within a given category. The category name is case-sensitive and must match exactly (e.g., `DataRate`, not `datarate` or `data-rate`).

```bash
gt 'metric DataRate'
# → Gbps, Kbps, Mbps, Tbps, bps

gt 'metric DataSize'
# → GB, GiB, KB, KiB, MB, MiB, PB, PiB, TB, TiB, bits, bytes

gt 'metric Time'
# → day, hr, min, ms, s

gt 'metric Weight'
# → g, kg, lb, mg, oz, ton

gt 'metric Speed'
# → kmh, knots, mph, mps

gt 'metric Distance'
# → ft, in, km, m, mi, nm
```

Unknown category names produce an error:

```bash
gt 'metric Temperature'
# → Error: rpn: metric Temperature: metric: unknown category "Temperature"
```

## `metric compatible`

Checks whether the top two values on the stack have compatible metric categories. Pops nothing from the stack — it only inspects. Pushes a descriptive string result.

Compatible metrics share the same category (or one or both are `Cool`/Universal, which absorbs into any category):

```bash
gt '100Mbps 50Mbps metric compatible'
# → Mbps (DataRate) and Mbps (DataRate): true

gt '100Mbps 2hr metric compatible'
# → Mbps (DataRate) and hr (Time): false

gt '5 10 metric compatible'
# → Cool (Universal) and Cool (Universal): true
```

At least two values are required on the stack:

```bash
gt '100 metric compatible'
# → Error: rpn: metric compatible: metric compatible: need at least 2 values on stack
```

### Practical use

Use `metric compatible` to verify that two values can be added, subtracted, or compared before performing the operation:

```bash
# Same category — addition will work
gt '1km 500m metric compatible'
# → km (Distance) and m (Distance): true

gt '1km 500m +'
# → 1.5

# Different categories — addition will fail
gt '100Mbps 2hr metric compatible'
# → Mbps (DataRate) and hr (Time): false

gt '100Mbps 2hr +'
# → Error: metric arithmetic requires compatible categories
```

## `metric decimal set`

Switches to SI prefix mode (powers of 1000) for data size units. In SI mode, `KB = 1000 bytes`, `MB = 1,000,000 bytes`, etc.

```bash
gt 'metric decimal set'
# → prefix mode: SI
```

This is the default mode. Running `metric decimal set` when already in SI mode is a no-op.

## `metric binary set`

Switches to IEC prefix mode (powers of 1024) for data size units. In IEC mode, `KB = 1024 bytes`, `MB = 1,048,576 bytes`, etc.

```bash
gt 'metric binary set'
# → prefix mode: IEC
```

In IEC mode, the SI-prefixed units (KB, MB, GB, TB, PB) use powers of 1024. The IEC-prefixed units (KiB, MiB, GiB, TiB, PiB) always use powers of 1024 regardless of mode.

### How prefix mode affects conversions

| Unit | SI mode factor | IEC mode factor |
|------|---------------|-----------------|
| KB | 1,000 bytes | 1,024 bytes |
| MB | 1,000,000 bytes | 1,048,576 bytes |
| GB | 1,000,000,000 bytes | 1,073,741,824 bytes |
| TB | 10^12 bytes | 10^12 bytes (approx) |
| KiB | 1,024 bytes | 1,024 bytes |
| MiB | 1,048,576 bytes | 1,048,576 bytes |

The prefix mode only affects SI-prefixed data size units (KB, MB, GB, TB, PB). All other metrics — DataRate, Time, Weight, Speed, Distance — are unaffected.

## REPL vs Single-Command Mode

### Single-command mode (`gt '<expression>'`)

A fresh RPN engine is created for each invocation. State changes like `metric binary set` only apply to that single expression and are discarded afterward:

```bash
# This works: mode is set and used in the same expression
gt 'metric binary set 1GB @MB convert'
# → prefix mode: SI    (set runs first, consumes tokens, stops)

# This is always SI because each invocation starts fresh
gt '1GB @MB convert'
# → 1000
```

Note: In single-command mode, `metric binary set` and subsequent tokens in the same expression may not all execute — metric commands return immediately after handling their tokens. Use REPL mode for mode changes followed by calculations.

### REPL mode (`gt` interactively)

The RPN engine persists across lines, so prefix mode changes remain in effect:

```
> metric binary set
prefix mode: IEC
> 1GB @MB convert
1024
> 1GB @GiB convert
1
> metric decimal set
prefix mode: SI
> 1GB @MB convert
1000
> metric list
DataRate, DataSize, Distance, Speed, Time, Universal, Weight
```

REPL mode is the recommended way to use `metric binary set` and `metric decimal set`, as well as for exploring categories and inspecting metrics interactively.

## SI vs IEC Prefix Modes

Data size prefixes in gt have two modes:

- **SI (decimal)**: KB = 10^3, MB = 10^6, GB = 10^9, TB = 10^12, PB = 10^15
- **IEC (binary)**: KB = 2^10, MB = 2^20, GB = 2^30, TB = 2^40, PB = 2^50

The mode only affects the SI-prefixed units (KB, MB, GB, TB, PB). The dedicated IEC units (KiB, MiB, GiB, TiB, PiB) always use powers of 1024 regardless of the current prefix mode. This means you can always use KiB/MiB/GiB for unambiguous binary sizes:

```bash
gt '1GiB @MiB convert'    # Always 1024
gt '1KiB @bytes convert'  # Always 1024
```

### When to use each mode

- **SI mode** (default): Matches networking and storage marketing (disks, bandwidth). Use when you need 1 GB = 1,000 MB.
- **IEC mode**: Matches operating system memory reporting. Use when you need 1 GB = 1,024 MB.
- **IEC units** (KiB/MiB/GiB): Always unambiguous — use these in scripts or documentation where precision matters.

## Examples

### Exploring the metrics system

```
> metric list
DataRate, DataSize, Distance, Speed, Time, Universal, Weight
> metric DataRate
Gbps, Kbps, Mbps, Tbps, bps
> 100Mbps metric show
Mbps, DataRate, base: bps, factor: 1e+06
```

### Verifying operations will work

```
> 1km 5mi metric compatible
km (Distance) and mi (Distance): true
> 1km 5mi +
6.04672
> 100Mbps 2hr metric compatible
Mbps (DataRate) and hr (Time): false
> 100Mbps 2hr +
Error: rpn operator failed for '+': metric arithmetic requires compatible categories (DataRate vs Time)
```

### Bandwidth calculations in IEC mode

```
> metric binary set
prefix mode: IEC
> 1GB @MB convert
1024
> 1GB @GB convert
1
> metric decimal set
prefix mode: SI
> 1GB @MB convert
1000
```

### Inspecting a value's metric

```
> 3.5hr
3.5
> metric show
hr, Time, base: seconds, factor: 4e+03
> 1000 @Mbps convert
1000
> metric show
Mbps, DataRate, base: bps, factor: 1e+06
```

## Summary

| Command | Description | Stack needed |
|---------|-------------|-------------|
| `metric show` | Show metric info for top of stack | 1 value |
| `metric list` | List all category names | None |
| `metric <category>` | List metrics in a category | None |
| `metric compatible` | Check compatibility of top two values | 2 values |
| `metric decimal set` | Switch to SI prefix mode | None |
| `metric binary set` | Switch to IEC prefix mode | None |
