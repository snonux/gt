# perc

A simple command-line percentage calculator written in Go.

## Installation

```bash
go install codeberg.org/snonux/perc/cmd/perc@latest
```

Or using mage:

```bash
mage install
```

## Usage

`perc` supports various percentage calculation formats:

### Calculate X% of Y

```bash
perc 20% of 150
# Output: 20.00% of 150.00 = 30.00

perc what is 20% of 150
# Output: 20.00% of 150.00 = 30.00
```

### Find what percentage X is of Y

```bash
perc 30 is what % of 150
# Output: 30.00 is 20.00% of 150.00
```

### Find the whole when X is Y% of it

```bash
perc 30 is 20% of what
# Output: 30.00 is 20.00% of 150.00
```

## Building

Using mage:

```bash
mage build
```

Or using go directly:

```bash
go build -o perc ./cmd/perc
```

## Testing

```bash
mage test
```

## License

See LICENSE file for details.
