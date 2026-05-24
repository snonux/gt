# CLI Usage

gt supports several invocation modes: single-expression evaluation, interactive REPL, piping, and stdin redirection.

## Invocation Modes

### Single-Expression Mode

Pass an expression as command-line arguments:

```bash
gt '3 4 +'                # → 7
gt '20% of 150'           # → 30
gt '100Mbps @Gbps convert' # → 1
```

This is the primary mode for scripts, one-off calculations, and CI pipelines. Each invocation starts fresh — no variables, no history.

### REPL Mode

Run gt with no arguments when stdin is connected to a terminal:

```bash
gt
```

This starts the interactive Read-Eval-Print Loop with command history, tab completion, and persistent variable storage. See [repl-mode.md](repl-mode.md) for full details.

When stdin is not a terminal (e.g., piped input), gt falls back to single-command mode.

### Pipe Mode

Pipe input directly to gt:

```bash
echo '3 4 +' | gt                  # → 7
echo '20% of 150' | gt             # → 30
echo '100kmh @mph convert' | gt    # → 62.13711922
```

gt reads from stdin, trims whitespace, tries RPN parsing first, then falls back to percentage calculation.

### Stdin Redirection

```bash
gt < input.txt
```

gt uses `os.ReadFile("/dev/stdin")` for full reads, with a 4096-byte buffer fallback if `/dev/stdin` is unavailable.

### Empty Input

Running gt with no input and no TTY displays usage information and exits with code 1:

```bash
$ echo '' | gt
Usage: gt <calculation>
       gt version

Percentage calculator examples:
  gt 20% of 150
  ...

Error: no input provided
```

## Version and Help

### Version

```bash
$ gt version
v0.4.2
```

The version string is defined in `internal/version.go` and updated during builds.

### Help

There are no `-h`, `--help`, or `-version` flags. Use `gt version` for the version string. If gt receives unrecognized input (including `-h` or `--help` as arguments), it attempts to evaluate them as expressions, which produces an error:

```bash
$ gt -h
Error: rpn fallback failed for input "-h": perc: unable to parse input "-h": unknown error
```

To see usage information, provide empty input or just run `gt` with no stdin.

### `--log` Flag

The only supported flag is `--log <file>`, which appends a session log of input and output when running in REPL mode:

```bash
gt --log session.log
```

See [repl-mode.md](repl-mode.md) for details on session logging.

## Boolean Coercion

gt treats `true` and `false` as first-class literals on the RPN stack. Comparison operators produce boolean values. Booleans coerce to numbers in arithmetic: `true` is `1`, `false` is `0`.

### Boolean Literals

```bash
gt 'true'                   # → true
gt 'false'                  # → false
```

### Comparison Operators Produce Booleans

```bash
gt '5 5 =='                 # → true
gt '5 6 =='                 # → false
gt '5 3 >'                  # → true
```

### Boolean-Arithmetic Coercion

```bash
gt 'true true +'            # → 2    (1 + 1)
gt 'true false +'           # → 1    (1 + 0)
gt 'false false +'          # → 0    (0 + 0)
gt 'true 2 *'               # → 2    (1 * 2)
```

### Mixed Boolean-Arithmetic Expressions

```bash
gt '5 3 == 1 +'             # → 1    (false + 1 = 0 + 1)
gt '5 5 == 10 *'            # → 10   (true * 10 = 1 * 10)
```

### Comparing Booleans to Numbers

```bash
gt 'true 1 =='              # → true   (1 == 1)
gt 'false 0 =='             # → true   (0 == 0)
gt 'true 0 =='              # → false  (1 != 0)
```

## Exit Codes

| Condition | Exit Code |
|-----------|-----------|
| Successful evaluation | 0 |
| Error (parse failure, invalid expression, etc.) | 1 |
| Empty input (no TTY) | 1 |

Error messages are printed to stdout with the prefix `Error:`.

## Examples

### In Shell Scripts

```bash
THRESHOLD=$(gt '100Mbps 50 /')
if [ "$THRESHOLD" -gt 1000 ]; then
    echo "High throughput"
fi
```

### In CI Pipelines

```yaml
- name: Check bandwidth math
  run: |
    result=$(gt '1Gbps @Mbps convert')
    [ "$result" = "1000" ] || exit 1
```

### In Pipes

```bash
cat sizes.txt | xargs -I{} gt '{} @GB convert'
```

### Quick Conversions at the Prompt

```bash
gt '500mi @km convert'         # → 804.672
gt '15% of 899'                # → discount amount
gt '75 80 >'                   # → false (75 < 80)
```

## Summary of Modes

| Mode | Command | State | Best For |
|------|---------|-------|----------|
| Single-expression | `gt '3 4 +'` | Fresh per invocation | Scripts, CI, one-off math |
| REPL | `gt` (with TTY) | Persistent (saved to disk) | Interactive exploration |
| Pipe | `echo '3 4 +' | gt` | Fresh per invocation | Pipelines, xargs |
| Stdin redirect | `gt < file.txt` | Fresh per invocation | Batch processing |
| Version | `gt version` | N/A | Checking installed version |
