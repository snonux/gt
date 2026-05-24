# Unit Conversion

gt supports converting values between different units of measurement using the `@` prefix and the `convert` operator. This works across all metric categories — data rates, data sizes, time, weight, speed, and distance.

## The `@` Prefix Syntax

The `@` prefix creates a metric reference — a value of `1` tagged with a specific unit. On its own, it pushes that metric onto the stack:

```bash
gt '@Gbps'             # → 1 (with Gbps metric)
gt '@km'               # → 1 (with km metric)
gt '@min'              # → 1 (with min metric)
```

The `@` prefix resolves the metric name through the same lookup as suffix notation: exact match, then aliases, then case-insensitive. All built-in metrics and aliases are supported.

## The `convert` Operator

The `convert` operator pops two values from the stack:

1. **Top of stack**: target metric (from `@` prefix)
2. **Second on stack**: value to convert

It converts the value through the base unit of its category and pushes the result with the target metric:

```
value → base units → target units
```

### Syntax

```bash
gt '<value> @<target> convert'
```

### Basic Examples

```bash
gt '1000Mbps @Gbps convert'    # → 1
gt '1hr @min convert'          # → 60
gt '1km @mi convert'           # → 0.6213711922
```

## All Category Conversions

### DataRate (base unit: bps)

Converts between bits-per-second units. All use SI (powers of 1000).

| Conversion | Result |
|------------|--------|
| `1Gbps @Mbps convert` | 1000 |
| `1000Mbps @Gbps convert` | 1 |
| `1Tbps @Gbps convert` | 1000 |
| `1Gbps @bps convert` | 1000000000 |
| `100Kbps @bps convert` | 100000 |

### DataSize (base unit: bits)

Converts between data size units. SI units (KB, MB, GB, TB, PB) use powers of 1000; IEC units (KiB, MiB, GiB, TiB, PiB) use powers of 1024. Both coexist and interoperate:

| Conversion | Result |
|------------|--------|
| `1TB @GB convert` | 1000 |
| `1GB @MB convert` | 1000 |
| `1MB @KB convert` | 1000 |
| `1KB @bytes convert` | 1000 |
| `1GB @bytes convert` | 1000000000 |
| `1GiB @MiB convert` | 1024 |
| `1KiB @bytes convert` | 1024 |
| `1PB @TB convert` | 1000 |
| `1GB @GiB convert` | 0.9313225746 |
| `1KiB @KB convert` | 1.024 |
| `1024bytes @KB convert` | 1.024 |
| `1GB @bits convert` | 8000000000 |
| `100bytes @KB convert` | 0.1 |

### Time (base unit: seconds)

| Conversion | Result |
|------------|--------|
| `1hr @min convert` | 60 |
| `1day @hr convert` | 24 |
| `1min @s convert` | 60 |
| `1s @ms convert` | 1000 |
| `1hr @s convert` | 3600 |
| `1day @s convert` | 86400 |
| `500ms @s convert` | 0.5 |
| `90min @hr convert` | 1.5 |
| `72hr @day convert` | 3 |
| `3.5hr @min convert` | 210 |

### Weight (base unit: kilograms)

| Conversion | Result |
|------------|--------|
| `1kg @lb convert` | 2.204622622 |
| `1lb @kg convert` | 0.45359237 |
| `1ton @kg convert` | 1000 |
| `70kg @lb convert` | 154.3235835 |
| `1lb @oz convert` | 16 |
| `1oz @g convert` | 28.34952313 |
| `1g @mg convert` | 1000 |
| `1kg @ton convert` | 0.001 |
| `1ton @lb convert` | 2204.622622 |

### Speed (base unit: m/s)

| Conversion | Result |
|------------|--------|
| `100kmh @mph convert` | 62.13711922 |
| `60mph @kmh convert` | 96.56064 |
| `1kmh @mps convert` | 0.2777777778 |
| `1knots @kmh convert` | 1.852 |
| `100mph @knots convert` | 86.89762419 |
| `340mps @kmh convert` | 1224 |

### Distance (base unit: meters)

| Conversion | Result |
|------------|--------|
| `1km @mi convert` | 0.6213711922 |
| `1mi @km convert` | 1.609344 |
| `1m @ft convert` | 3.280839895 |
| `1ft @m convert` | 0.3048 |
| `1mi @ft convert` | 5280 |
| `1in @m convert` | 0.0254 |
| `1nm @km convert` | 1.852 |
| `5280ft @mi convert` | 1 |
| `3ft @in convert` | 36 |

Note: `nm` stands for **nautical mile** (1852 meters), not nanometer. For nanometer-scale distances, use scientific notation directly (e.g., `1e-9m`).

### Universal (base unit: Cool)

The `Cool` metric is the default unitless type. Converting a Cool value to another metric treats the number as units of the target metric (Cool absorbing):

```bash
gt '5 @kg convert'        # → 5 (treated as 5 kg)
gt '42 @Cool convert'     # → 42 (stays Cool)
```

## Metric Aliases in Conversions

All metric aliases work with the `@` prefix:

```bash
gt '1gbit/s @Kbps convert'   # → 1000000  (gbit/s → Gbps)
gt '1sec @ms convert'        # → 1000     (sec → s)
gt '1mile @km convert'       # → 1.609344 (mile → mi)
gt '1feet @m convert'        # → 0.3048   (feet → ft)
gt '1knot @kmh convert'      # → 1.852    (knot → knots)
```

## Converting Arithmetic Results

`convert` works on any value on the stack, including results of arithmetic operations:

```bash
# Addition then convert
gt '5km 5mi + @km convert'          # → 13.04672

# Subtraction then convert
gt '1Gbps 500Mbps - @Mbps convert'   # → 500

# Division then convert
gt '100km 2hr / @kmh convert'        # → 50

# Cross-category multiplication then convert
gt '100kmh 1hr * @mi convert'        # → 62.13711922

# DataRate × Time → DataSize, then convert
gt '1Gbps 1hr * @GB convert'         # → 450
```

## Practical Use Cases

### Network Bandwidth

```bash
# Convert internet speed to familiar units
gt '500Mbps @Gbps convert'           # → 0.5

# Calculate how many GB you download in an hour at 1 Gbps
gt '1Gbps 1hr * @GB convert'         # → 450

# Server throughput in Mbps
gt '100Kbps @Mbps convert'           # → 0.1
```

### Travel Planning

```bash
# Speed limit conversion (US to metric)
gt '65mph @kmh convert'              # → 104.86

# Distance conversion for flight planning
gt '100nm @km convert'               # → 185.2 (nautical miles to km)

# Flight time: 500 miles at 400 mph
gt '500mi 400mph / @min convert'     # → 75 minutes

# Speed of sound in km/h
gt '340mps @kmh convert'             # → 1224
```

### Weight Conversions

```bash
# Body weight
gt '70kg @lb convert'                # → 154.32

# Shipping: 50 oz to grams
gt '50oz @g convert'                 # → 1417.476156

# Industrial weight
gt '2.5ton @kg convert'              # → 2500
```

### Cooking and Recipes

```bash
# Kitchen weight
gt '250g @oz convert'                # → 8.82

# Small measurements
gt '500mg @g convert'                # → 0.5
```

### Data Storage

```bash
# Disk space: how many GB in a TB
gt '1TB @GB convert'                 # → 1000

# RAM: how many MiB in a GiB
gt '1GiB @MiB convert'               # → 1024

# File size: 500 MB in KB
gt '500MB @KB convert'               # → 500000
```

### Time Planning

```bash
# How many seconds in 3 days
gt '3day @s convert'                 # → 259200

# Half-hour delay in milliseconds
gt '0.5hr @ms convert'               # → 1800000

# Milliseconds to seconds
gt '2500ms @s convert'               # → 2.5
```

## Edge Cases

### Incompatible Categories

Converting between different metric categories produces an error because the units are fundamentally incompatible:

```bash
gt '1km @hr convert'        # Error: incompatible metrics (Distance vs Time)
gt '100Mbps @kg convert'    # Error: incompatible metrics (DataRate vs Weight)
```

The `convert` operator validates that both values belong to the same category (or that the source is `Cool`, which absorbs into any target category).

### Unknown Metrics

Using an unknown metric name with `@` produces a parse error:

```bash
gt '100 @xyz convert'       # Error: unknown metric "xyz"
```

### Same-Unit Conversion

Converting to the same unit is a no-op:

```bash
gt '42km @km convert'       # → 42
gt '1Gbps @Gbps convert'    # → 1
```

### Zero and Negative Values

`convert` handles zero and negative values correctly:

```bash
gt '0km @mi convert'        # → 0
gt '-1kg @lb convert'       # → -2.204622622
gt '1km 1mi - @ft convert'  # → -1999.160105
```

### SI vs IEC Interoperability

SI and IEC data size units coexist and can be cross-converted:

```bash
gt '1GB @GiB convert'       # → 0.9313225746  (1 GB is ~931 MiB)
gt '1KiB @KB convert'       # → 1.024          (1 KiB is 1.024 KB)
```

In single-command mode, SI (powers of 1000) is always used for KB/MB/GB/TB/PB units. IEC units (KiB/MiB/GiB/TiB/PiB) always use powers of 1024 regardless of mode.

### Cool Absorbing

Unitless (`Cool`) values absorb into any metric during `convert`. The number is treated as units of the target metric:

```bash
gt '5 @kg convert'          # → 5 (5 kg)
gt '100 @Mbps convert'      # → 100 (100 Mbps)
```

This is consistent with how Cool absorbs during addition: `5 100Mbps +` treats the `5` as `5Mbps`.
