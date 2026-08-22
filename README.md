# gt

A command-line calculator with percentage calculations and RPN (Reverse Polish Notation) stack arithmetic, written in Go.

![gt Logo](logo.svg)

This is a toy project to experiment with local LLMs (Qwen, Gemma, Nemotron, GPT OSS) as pair programmers.

## Installation

```bash
go install github.com/snonux/gt/cmd/gt@latest
```

Or using [mage](https://github.com/magefile/mage):

```bash
mage install
```

## Quick Start

```bash
gt '20% of 150'         # → 30
gt '3 4 +'              # RPN: 3 + 4 → 7
gt 'pi 2 *'             # 2π → 6.283185307
gt '100Mbps 1hr *'      # rate × time → 3.6e+11 bits
gt                      # interactive REPL
```

## Shell Completions

Fish completions are available in [`completions/`](completions/). To install:

```bash
cp completions/gt.fish ~/.config/fish/completions/
```

Or system-wide:

```bash
sudo cp completions/gt.fish /usr/local/share/fish/vendor_completions.d/
```

## Feature Guide

All features are documented in [`docs/`](docs/):

| Topic | Doc |
|-------|-----|
| CLI modes, piping, REPL | [cli-usage](docs/cli-usage.md) |
| Percentage calculations | [percentage-calculations](docs/percentage-calculations.md) |
| RPN basics | [basic-arithmetic](docs/basic-arithmetic.md) |
| Fast integer power (`**`) | [fast-power](docs/fast-power.md) |
| Logarithm operators | [log-operators](docs/log-operators.md) |
| Hyper (n-ary) operators | [hyper-operators](docs/hyper-operators.md) |
| Comparisons and booleans | [comparisons](docs/comparisons.md) |
| Variables and assignment | [variables](docs/variables.md) |
| Symbols (`:x` syntax) | [symbols](docs/symbols.md) |
| Built-in constants | [constants](docs/constants.md) |
| Stack manipulation | [stack-operations](docs/stack-operations.md) |
| Metrics and units | [metrics](docs/metrics.md) |
| Unit conversion | [unit-conversion](docs/unit-conversion.md) |
| Metric commands | [metric-commands](docs/metric-commands.md) |
| Custom metrics | [custom-metrics](docs/custom-metrics.md) |
| REPL mode | [repl-mode](docs/repl-mode.md) |
| Rational number mode | [rational-mode](docs/rational-mode.md) |
| Fish shell completions | [completions](completions/)

## Building and Testing

```bash
mage build      # Build the binary
mage test       # Run all tests
```

## License

See [LICENSE](LICENSE) for details.
