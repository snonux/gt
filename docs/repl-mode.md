# REPL Mode

gt has an interactive Read-Eval-Print Loop (REPL) mode that provides a persistent calculation session with command history, tab completion, and persistent state. Use it for exploratory calculations, multi-step RPN work, and configuring settings like rational mode and metric prefix mode.

## Starting the REPL

### `gt` (no arguments)

Run gt without arguments when stdin is connected to a terminal:

```bash
gt
```

This starts the interactive prompt:

```
>
```

### `mage repl`

From the gt source directory, start a REPL session without building:

```bash
mage repl
```

This runs the Go source directly via `go run`.

### `--log <file>`

Append a session log of input and output to a file:

```bash
gt --log session.log
```

When stdin is not a terminal (e.g., piped input), gt falls back to single-command mode instead of starting the REPL. In that case it reads from stdin, tries RPN parsing first, and falls back to percentage calculation:

```bash
echo '3 4 +' | gt          # → 7
echo '20% of 150' | gt     # → 30
```

### No-input usage

Running gt with no input and no TTY displays usage information:

```bash
$ gt
Usage: gt <calculation>
       gt version
...
```

## Navigation and History

The REPL uses [chzyer/readline](https://github.com/chzyer/readline) for an Emacs-style interactive line editor.

### History

Commands are saved to `~/.gt_history` and persist across sessions. A maximum of 1000 entries are kept.

| Shortcut | Action |
|----------|--------|
| Up arrow | Previous command |
| Down arrow | Next command |
| Ctrl+R | Reverse history search (type to filter) |

### Line editing

| Shortcut | Action |
|----------|--------|
| Ctrl+A | Move to beginning of line |
| Ctrl+E | Move to end of line |
| Ctrl+F | Forward one character |
| Ctrl+B | Backward one character |
| Ctrl+H | Delete character before cursor (Backspace) |
| Ctrl+D | Delete character under cursor |
| Ctrl+W | Cut word before cursor |
| Ctrl+K | Cut from cursor to end of line |
| Ctrl+U | Cut from beginning to cursor |
| Ctrl+L | Clear the screen |

## Tab Completion

Pressing Tab completes built-in REPL commands based on what you've typed. Completion is case-insensitive.

The eight built-in commands are: `help`, `clear`, `quit`, `exit`, `rpn`, `calc`, `rat`, `stack`.

```
> h<TAB>
help
> q<TAB>
quit
> r<TAB>    # shows all 'r' commands
rat  rpn
```

When multiple commands match the prefix, cycling through Tab shows all options.

## State Persistence

### Variables

Variables defined in one session are saved to `~/.local/state/gt/vars` (XDG State directory) and automatically restored when you start a new session:

```
> x 42 =
x = 42
> quit
Goodbye!

# In a new session:
> x
42
```

### History

Command history is persisted to `~/.gt_history` and loaded on startup.

## Built-in REPL Commands

The REPL has eight built-in commands. They are recognized case-insensitively.

### `help`

Display general help or help for a specific command:

```
> help
PERC - Percentage Calculator REPL

Built-in Commands:
  help             Show this help message
  help <command>   Show help for a specific topic
  clear            Clear the screen
  quit / exit      Exit the REPL
  rpn / calc       Evaluate an RPN (postfix notation) expression
  rat on/off/toggle Switch between float64 and rational number modes
  stack            Show current stack state (same as 'rpn show')
...

> help clear
clear - Clear the screen
Usage: clear
```

### `clear`

Clear the terminal screen using ANSI escape sequences:

```
> clear
```

(Alternatively, press Ctrl+L.)

### `quit` / `exit`

Exit the REPL and print a farewell message:

```
> quit
Goodbye!
```

Both `quit` and `exit` are equivalent.

### `rpn` / `calc`

Prefix an RPN expression explicitly. These are optional — bare RPN expressions work without them:

```
> rpn 3 4 +
7
> calc 3 4 +
7
> 3 4 +
7
```

### `rat`

Switch between float64 and rational number modes. See [rational-mode.md](rational-mode.md) for full details.

```
> rat on
Rational mode enabled
> 1 3 /
0.3333333333
> rat off
Rational mode disabled (using float64)
> 1 3 /
0.3333333333
```

`rat` requires an argument: `on`, `off`, or `toggle`.

### `stack`

Show a hint for viewing the RPN stack. In the REPL, use `show` directly on the input line:

```
> stack
Use 'rpn show' to view the current stack state
```

## Calculation Modes

The REPL dispatches input through a chain of handlers:

1. **Built-in commands** — help, clear, quit, exit, rpn, calc, rat, stack
2. **RPN expressions** — bare RPN (`3 4 +`), single operators (`dup`), single numbers (`42`), and symbol syntax (`:x`)
3. **Percentage calculations** — `20% of 150`, `what is 20% of 150`, etc.
4. **Error** — unknown commands produce an error message

This means you can freely mix RPN and percentage calculations in a session:

```
> 20% of 150
20.00% of 150.00 = 30.00
  Steps: (20.00 / 100) * 150.00 = 0.20 * 150.00 = 30.00
> 3 4 +
7
> 30 20 +
50
```

## REPL-Only Features

Some features require the persistent state of REPL mode and are not available in single-command mode (`gt '<expression>'`).

### Rational mode (`rat`)

The `rat` command is REPL-only. See [rational-mode.md](rational-mode.md).

```
> rat on
Rational mode enabled
> 1 2 +
3.0000000000
```

### Metric prefix mode (`metric decimal set` / `metric binary set`)

Prefix mode changes persist across REPL lines. See [metric-commands.md](metric-commands.md).

```
> metric binary set
prefix mode: IEC
> 1GB @MB convert
1024
> metric decimal set
prefix mode: SI
> 1GB @MB convert
1000
```

In single-command mode, prefix mode is always SI and resets on every invocation.

### Custom metrics (`custom`)

Define your own units and use them in subsequent calculations. See [custom-metrics.md](custom-metrics.md).

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

In single-command mode, `custom define` stops token evaluation after registering the metric, so subsequent tokens are not processed.

### Variables

Variables persist across REPL lines and are saved to disk between sessions:

```
> x 10 =
x = 10
> y 20 =
y = 20
> x y +
30
> quit
Goodbye!

# New session — variables are restored:
> x
10
```

### Metric commands (`metric show`, `metric list`, etc.)

All metric commands work in the REPL, inspecting the persistent stack:

```
> 100Mbps
100
> metric show
Mbps, DataRate, base: bps, factor: 1e+06
> metric list
DataRate, DataSize, Distance, Speed, Time, Universal, Weight
> 1km 5mi metric compatible
km (Distance) and mi (Distance): true
```

See [metric-commands.md](metric-commands.md) for full details.

## Exiting the REPL

Five methods exit the REPL:

| Method | Behavior |
|--------|----------|
| `quit` | Prints "Goodbye!" and exits |
| `exit` | Same as `quit` |
| Ctrl+D | EOF — exits cleanly |
| Ctrl+C | SIGINT — prints "Use 'quit' or 'exit' to exit, or Ctrl+D" and continues |
| | (pressing Ctrl+C again repeats the hint) |

Ctrl+C does not exit by default — it gives you a chance to cancel. Press Ctrl+D or type `quit` to leave.

## Examples

### Multi-step RPN calculation

```
> 100Mbps
100
> 1hr *
360000000000
> metric show
bps, DataRate, base: bps, factor: 1
> @GB convert
360000000000
```

### Exploring metric categories

```
> metric list
DataRate, DataSize, Distance, Speed, Time, Universal, Weight
> metric DataRate
Gbps, Kbps, Mbps, Tbps, bps
> 1Gbps @Mbps convert
1000
> metric Distance
ft, in, km, m, mi, nm
> 100km @mi convert
62.13711922
```

### Defining and using custom units

```
> custom define reel 304.8 Distance
defined custom metric "reel" (factor: 304.8, category: Distance)
> 5reel @m convert
1524
> custom list
reel
> custom undefine reel
removed custom metric "reel"
```

### Switching between float and rational modes

```
> rat on
Rational mode enabled
> 10 3 -
7.0000000000
> 1 3 / 3 *
1.0000000000
> rat off
Rational mode disabled (using float64)
> 10 3 -
7
> 1 3 / 3 *
1
```

### Variable-driven calculations

```
> rate 100Mbps =
rate = 100
> time 2hr =
time = 2
> rate time *
200
> show
Stack (bottom to top):
  200 Mbps
```

## REPL vs Single-Command Mode

| Feature | REPL | Single-command |
|---------|------|----------------|
| Command history | Yes (persisted) | No |
| Tab completion | Yes | No |
| Arrow key navigation | Yes | No |
| Ctrl+R search | Yes | No |
| Variable persistence | Yes (saved to disk) | No (fresh per invocation) |
| `rat` (rational mode) | Yes | No |
| `metric binary set` persists | Yes | No (always SI) |
| Custom metrics across lines | Yes | No (define stops evaluation) |
| Line editing (Ctrl+A/E/W/K/U) | Yes | No |
| Session logging (`--log`) | Yes | No |

Use the REPL for interactive exploration, multi-step calculations, and any workflow involving state changes. Use single-command mode for scripts, one-off calculations, and piping.
