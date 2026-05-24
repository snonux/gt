# Metrics System

gt attaches units of measurement (metrics) to numbers. Every number on the stack carries a metric — the default is `Cool` (unitless). Metrics enable automatic unit conversion, metric-aware arithmetic, and cross-category inference.

## Suffix Notation

Attach a unit directly to a number with no space:

```bash
gt '100Mbps'        # 100 megabits per second
gt '5GB'            # 5 gigabytes
gt '1hr'            # 1 hour
gt '100kmh'         # 100 kilometers per hour
```

## Available Metric Categories

### DataRate

| Unit | Description |
|------|-------------|
| `bps` | Bits per second (base) |
| `Kbps` | Kilobits per second |
| `Mbps` | Megabits per second |
| `Gbps` | Gigabits per second |
| `Tbps` | Terabits per second |

### DataSize

| Unit | Description |
|------|-------------|
| `bits` / `bytes` | Base units |
| `KB` / `MB` / `GB` / `TB` / `PB` | SI (1000-based) |
| `KiB` / `MiB` / `GiB` / `TiB` / `PiB` | IEC (1024-based) |

### Time

| Unit | Description |
|------|-------------|
| `ms` | Milliseconds |
| `s` | Seconds (base) |
| `min` | Minutes |
| `hr` | Hours |
| `day` | Days |

### Weight

| Unit | Description |
|------|-------------|
| `mg` | Milligrams |
| `g` | Grams (base) |
| `kg` | Kilograms |
| `lb` | Pounds |
| `oz` | Ounces |
| `ton` | Tons |

### Speed

| Unit | Description |
|------|-------------|
| `mps` | Meters per second (base) |
| `kmh` | Kilometers per hour |
| `mph` | Miles per hour |
| `knots` | Knots |

### Distance

| Unit | Description |
|------|-------------|
| `m` | Meters (base) |
| `km` | Kilometers |
| `mi` | Miles |
| `ft` | Feet |
| `in` | Inches |
| `nm` | Nautical miles |

### Universal

| Unit | Description |
|------|-------------|
| `Cool` | Default unitless metric |

## Unit Conversion

Convert between units using `@` prefix and `convert`:

```bash
gt '1000Mbps @Gbps convert'    # → 1
gt '1hr @min convert'          # → 60
gt '1km @mi convert'           # → 0.6213711922
gt '1TB @GB convert'           # → 1000
gt '1KiB @bytes convert'       # → 1024
```

## Metric-Aware Arithmetic

### Addition and Subtraction

Require compatible categories. Mixed units are automatically converted:

```bash
gt '100Mbps 50Mbps +'          # → 150  (same category)
gt '1km 500m +'                # → 1.5  (auto-convert to km)
gt '5 100Mbps +'               # → 105  (Cool absorbs into Mbps)
gt '100Mbps 2hr +'             # ERROR: incompatible categories
```

### Multiplication

Supports cross-category inference:

```bash
gt '100Mbps 1hr *'             # → 3.6e+11  (rate × time = data size)
gt '100kmh 1hr *'              # → 100000   (speed × time = distance)
```

### Division

Also supports cross-category inference:

```bash
gt '10GB 2hr /'                # → 11111111.11  (data / time = rate)
gt '1km 1s /'                  # → 1000   (distance / time = speed)
```

### Power

Always returns `Cool` (unitless):

```bash
gt '2hr 3 ^'                   # → 8  (not 8hr³)
```

## SI vs IEC Prefix Modes

Data size prefixes use SI (1000-based) by default. Switch to IEC (1024-based) in REPL mode:

```
> metric decimal set     # SI mode
prefix mode: SI
> 1GB @MB convert
1000
> metric binary set      # IEC mode
prefix mode: IEC
> 1GB @MB convert
1024
```

In single-command mode, SI is always used.

## Examples

### Bandwidth Planning

```bash
gt '1Gbps 1000 /'              # Gbps to Mbps: 1000
gt '100Mbps 1hr *'             # Bits transferred in 1 hour: 3.6e+11
gt '500GB 1hr /'               # Required throughput: large bps value
```

### Travel Calculations

```bash
gt '60mph @kmh convert'        # mph to km/h: 96.56
gt '500mi @km convert'         # miles to km: 804.67
gt '1000km 60mph /'            # Travel time: hours (after conversion)
```

### Weight Conversions

```bash
gt '70kg @lb convert'          # kg to lbs: 154.32
gt '250lb @kg convert'         # lbs to kg: 113.4
```

## Metric Commands

```bash
gt 'metric show'               # Show metric info for top of stack
gt 'metric list'               # List all categories
gt 'metric DataRate'           # List metrics in a category
gt '100Mbps 1Gbps metric compatible'  # Check compatibility
```

## Custom Metrics

Define your own units:

```bash
gt 'custom define foobar 42 Custom'
# → defined custom metric "foobar" (factor: 42, category: Custom)

gt 'custom define foobar 42 Custom 10foobar 5foobar +'
# → 15

gt 'custom define foobar 42 Custom custom undefine foobar'
# → removed custom metric "foobar"
```

## Notes

- Numbers without suffixes default to `Cool` (unitless)
- Metric names are case-insensitive (`100mbps` = `100Mbps`)
- Metric-aware comparison operators work with compatible categories
- Incompatible category operations produce clear error messages
