# SPDX-License-Identifier: MIT
# Copyright (c) 2026 Paul Buetow

# Fish shell completions for gt — a command-line calculator with RPN and
# percentage arithmetic.
#
# Install:
#   cp completions/gt.fish ~/.config/fish/completions/
#   # or system-wide:
#   sudo cp completions/gt.fish /usr/local/share/fish/vendor_completions.d/

# ── Helpers ──────────────────────────────────────────────────────────────────

# Return true when we are at the top-level (first argument after 'gt').
function __fish_gt_needs_top_level
    __fish_use_subcommand
end

# Return true when the current token is part of an RPN expression (i.e. we are
# NOT inside 'metric …' or 'custom …' subcommand trees that have their own
# completions).  Used as a condition for completing operators / constants /
# units anywhere in the expression.
function __fish_gt_not_in_metric_or_custom
    not __fish_seen_subcommand_from metric custom
end

# ── Flags ────────────────────────────────────────────────────────────────────

# Suppress fish's default file/path completions so only our custom tokens appear.
complete -c gt -k

complete -c gt -n "__fish_gt_needs_top_level" \
    -l log -r -d "Session log file (REPL mode)"

# ── Top-level subcommands ────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_needs_top_level" \
    -a version -d "Show version"

# ── Expression prefixes ──────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "rpn" -d "RPN expression prefix"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "calc" -d "Calculation prefix (alias)"

# ── Arithmetic operators ─────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "+" -d "Addition"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "-" -d "Subtraction"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "*" -d "Multiplication"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "/" -d "Division"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "^" -d "Power"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "**" -d "Fast integer power"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "%" -d "Modulo"

# ── Logarithmic operators ────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "lg" -d "Log base 2"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "log" -d "Log base 10"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "ln" -d "Natural log"

# ── Comparison operators ─────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "gt" -d "Greater than"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "lt" -d "Less than"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "gte" -d "Greater or equal"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "lte" -d "Less or equal"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "eq" -d "Equal"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "==" -d "Equal"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "neq" -d "Not equal"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "!=" -d "Not equal"

# ── Hyper (n-ary) operators ──────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[+]" -d "Hyper add (n-ary)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[-]" -d "Hyper subtract (n-ary)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[*]" -d "Hyper multiply (n-ary)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[/]" -d "Hyper divide (n-ary)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[^]" -d "Hyper power (n-ary)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[%]" -d "Hyper modulo (n-ary)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[lg]" -d "Hyper log2 (n-ary)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[log]" -d "Hyper log10 (n-ary)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[ln]" -d "Hyper ln (n-ary)"

# ── Stack operators ──────────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "dup" -d "Duplicate top of stack"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "swap" -d "Swap top two stack values"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "pop" -d "Discard top of stack"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "d" -d "Delete top two stack values"

# ── Assignment / conversion operators ────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a ":=" -d "Assign right-to-left"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "=:" -d "Assign left-to-right"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "convert" -d "Convert metric unit (with @unit)"

# ── Command operators ────────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "show" -d "Show stack"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "showstack" -d "Show stack (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "print" -d "Show stack (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "vars" -d "List variables"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "constants" -d "List constants"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "clear" -d "Clear all variables"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "clearconstants" -d "Reset constants to defaults"

# ── Boolean literals ─────────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "true" -d "Boolean true (coerces to 1)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "false" -d "Boolean false (coerces to 0)"

# ── Constants ────────────────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "pi" -d "π (3.14159…)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "e" -d "Euler's number (2.71828…)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "euler" -d "Euler's number (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "phi" -d "Golden ratio (1.61803…)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "sqrt2" -d "√2 (1.41421…)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "sqrt3" -d "√3 (1.73205…)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "sqrt5" -d "√5 (2.23606…)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "ln2" -d "ln(2) (0.69314…)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "log2" -d "ln(2) (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "ln10" -d "ln(10) (2.30258…)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "log10" -d "ln(10) (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "log_e" -d "log₁₀(e) (0.43429…)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "log_e10" -d "log₁₀(e) (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "tau" -d "τ = 2π (6.28318…)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "inv_pi" -d "1/π (0.31830…)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "inv_e" -d "1/e (0.36787…)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "inf" -d "Positive infinity"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "infinity" -d "Positive infinity (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "-inf" -d "Negative infinity"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "-infinity" -d "Negative infinity (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "nan" -d "NaN (Not a Number)"

# ── Metric units ─────────────────────────────────────────────────────────────

# DataRate (base: bps)
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "bps" -d "Bits per second"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "Kbps" -d "Kilobits per second"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "Mbps" -d "Megabits per second"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "Gbps" -d "Gigabits per second"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "Tbps" -d "Terabits per second"

# DataSize (base: bits)
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "bits" -d "Bits"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "bytes" -d "Bytes"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "KB" -d "Kilobyte (SI / IEC)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "MB" -d "Megabyte (SI / IEC)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "GB" -d "Gigabyte (SI / IEC)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "TB" -d "Terabyte (SI / IEC)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "PB" -d "Petabyte (SI / IEC)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "KiB" -d "Kibibyte (1024 bytes)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "MiB" -d "Mebibyte (1024 KiB)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "GiB" -d "Gibibyte (1024 MiB)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "TiB" -d "Tebibyte (1024 GiB)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "PiB" -d "Pebibyte (1024 TiB)"

# Time (base: seconds)
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "ms" -d "Milliseconds"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "s" -d "Seconds"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "sec" -d "Seconds (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "secs" -d "Seconds (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "min" -d "Minutes"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "hr" -d "Hours"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "day" -d "Days"

# Weight (base: kilograms)
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "mg" -d "Milligrams"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "g" -d "Grams"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "kg" -d "Kilograms"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "ton" -d "Metric tons"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "lb" -d "Pounds"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "oz" -d "Ounces"

# Speed (base: m/s)
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "mps" -d "Meters per second"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "kmh" -d "Kilometers per hour"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "mph" -d "Miles per hour"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "knots" -d "Knots"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "knot" -d "Knots (alias)"

# Distance (base: meters)
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "m" -d "Meters"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "km" -d "Kilometers"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "mi" -d "Miles"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "mile" -d "Miles (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "miles" -d "Miles (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "ft" -d "Feet"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "foot" -d "Feet (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "feet" -d "Feet (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "in" -d "Inches"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "nm" -d "Nautical miles"

# DataRate aliases (bit/s notation)
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "bit/s" -d "Bits per second (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "kbit/s" -d "Kilobits per second (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "mbit/s" -d "Megabits per second (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "gbit/s" -d "Gigabits per second (alias)"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "tbit/s" -d "Terabits per second (alias)"

# ── `metric` subcommand ──────────────────────────────────────────────────────

complete -c gt -n "__fish_seen_subcommand_from metric" \
    -a "show" -d "Show metric info for top of stack"
complete -c gt -n "__fish_seen_subcommand_from metric" \
    -a "list" -d "List all metric categories"
complete -c gt -n "__fish_seen_subcommand_from metric" \
    -a "compatible" -d "Check metric compatibility (top 2 values)"
complete -c gt -n "__fish_seen_subcommand_from metric" \
    -a "decimal set" -d "Switch to SI prefix mode (1000-based)"
complete -c gt -n "__fish_seen_subcommand_from metric" \
    -a "binary set" -d "Switch to IEC prefix mode (1024-based)"

# metric <category> — list metrics in a category
complete -c gt \
    -n "__fish_seen_subcommand_from metric; and not __fish_seen_subcommand_from show list compatible 'decimal set' 'binary set'" \
    -a "DataRate" -d "bps, Kbps, Mbps, Gbps, Tbps"
complete -c gt \
    -n "__fish_seen_subcommand_from metric; and not __fish_seen_subcommand_from show list compatible 'decimal set' 'binary set'" \
    -a "DataSize" -d "bits, bytes, KB, MB, GB, …"
complete -c gt \
    -n "__fish_seen_subcommand_from metric; and not __fish_seen_subcommand_from show list compatible 'decimal set' 'binary set'" \
    -a "Distance" -d "m, km, mi, ft, in, nm"
complete -c gt \
    -n "__fish_seen_subcommand_from metric; and not __fish_seen_subcommand_from show list compatible 'decimal set' 'binary set'" \
    -a "Speed" -d "mps, kmh, mph, knots"
complete -c gt \
    -n "__fish_seen_subcommand_from metric; and not __fish_seen_subcommand_from show list compatible 'decimal set' 'binary set'" \
    -a "Time" -d "ms, s, min, hr, day"
complete -c gt \
    -n "__fish_seen_subcommand_from metric; and not __fish_seen_subcommand_from show list compatible 'decimal set' 'binary set'" \
    -a "Universal" -d "Cool (unitless)"
complete -c gt \
    -n "__fish_seen_subcommand_from metric; and not __fish_seen_subcommand_from show list compatible 'decimal set' 'binary set'" \
    -a "Weight" -d "mg, g, kg, ton, lb, oz"

# ── `custom` subcommand ──────────────────────────────────────────────────────

complete -c gt -n "__fish_seen_subcommand_from custom" \
    -a "show" -d "Show custom metric(s)"
complete -c gt -n "__fish_seen_subcommand_from custom" \
    -a "list" -d "List all custom metrics"
complete -c gt -n "__fish_seen_subcommand_from custom" \
    -a "define" -d "Define a new custom metric"
complete -c gt -n "__fish_seen_subcommand_from custom" \
    -a "undefine" -d "Remove a custom metric"

# custom define <name> <factor> <category> — complete categories
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "Custom" -d "Custom (user-defined)"
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "DataRate" -d "Data rate category"
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "DataSize" -d "Data size category"
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "Distance" -d "Distance category"
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "Speed" -d "Speed category"
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "Time" -d "Time category"
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "Universal" -d "Universal (unitless)"
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "Weight" -d "Weight category"
