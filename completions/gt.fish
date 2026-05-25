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

complete -c gt -n "__fish_gt_needs_top_level" \
    -l log -r

# ── Top-level subcommands ────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_needs_top_level" \
    -a version

# ── Expression prefixes ──────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "rpn"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "calc"

# ── Arithmetic operators ─────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "+"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "-"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "*"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "/"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "^"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "**"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "%"

# ── Logarithmic operators ────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "lg"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "log"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "ln"

# ── Comparison operators ─────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "gt"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "lt"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "gte"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "lte"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "eq"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "=="
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "neq"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "!="

# ── Hyper (n-ary) operators ──────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[+]"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[-]"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[*]"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[/]"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[^]"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[%]"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[lg]"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[log]"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "[ln]"

# ── Stack operators ──────────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "dup"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "swap"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "pop"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "d"

# ── Assignment / conversion operators ────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a ":="
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "=:"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "convert"

# ── Command operators ────────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "show"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "showstack"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "print"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "vars"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "constants"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "clear"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "clearconstants"

# ── Boolean literals ─────────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "true"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "false"

# ── Constants ────────────────────────────────────────────────────────────────

complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "pi"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "e"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "euler"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "phi"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "sqrt2"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "sqrt3"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "sqrt5"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "ln2"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "log2"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "ln10"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "log10"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "log_e"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "log_e10"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "tau"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "inv_pi"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "inv_e"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "inf"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "infinity"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "-inf"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "-infinity"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "nan"

# ── Metric units ─────────────────────────────────────────────────────────────

# DataRate (base: bps)
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "bps"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "Kbps"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "Mbps"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "Gbps"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "Tbps"

# DataSize (base: bits)
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "bits"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "bytes"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "KB"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "MB"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "GB"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "TB"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "PB"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "KiB"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "MiB"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "GiB"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "TiB"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "PiB"

# Time (base: seconds)
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "ms"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "s"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "sec"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "secs"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "min"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "hr"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "day"

# Weight (base: kilograms)
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "mg"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "g"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "kg"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "ton"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "lb"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "oz"

# Speed (base: m/s)
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "mps"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "kmh"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "mph"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "knots"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "knot"

# Distance (base: meters)
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "m"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "km"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "mi"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "mile"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "miles"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "ft"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "foot"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "feet"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "in"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "nm"

# DataRate aliases (bit/s notation)
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "bit/s"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "kbit/s"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "mbit/s"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "gbit/s"
complete -c gt -n "__fish_gt_not_in_metric_or_custom" \
    -a "tbit/s"

# ── `metric` subcommand ──────────────────────────────────────────────────────

complete -c gt -n "__fish_seen_subcommand_from metric" \
    -a "show"
complete -c gt -n "__fish_seen_subcommand_from metric" \
    -a "list"
complete -c gt -n "__fish_seen_subcommand_from metric" \
    -a "compatible"
complete -c gt -n "__fish_seen_subcommand_from metric" \
    -a "decimal set"
complete -c gt -n "__fish_seen_subcommand_from metric" \
    -a "binary set"

# metric <category> — list metrics in a category
complete -c gt \
    -n "__fish_seen_subcommand_from metric; and not __fish_seen_subcommand_from show list compatible 'decimal set' 'binary set'" \
    -a "DataRate"
complete -c gt \
    -n "__fish_seen_subcommand_from metric; and not __fish_seen_subcommand_from show list compatible 'decimal set' 'binary set'" \
    -a "DataSize"
complete -c gt \
    -n "__fish_seen_subcommand_from metric; and not __fish_seen_subcommand_from show list compatible 'decimal set' 'binary set'" \
    -a "Distance"
complete -c gt \
    -n "__fish_seen_subcommand_from metric; and not __fish_seen_subcommand_from show list compatible 'decimal set' 'binary set'" \
    -a "Speed"
complete -c gt \
    -n "__fish_seen_subcommand_from metric; and not __fish_seen_subcommand_from show list compatible 'decimal set' 'binary set'" \
    -a "Time"
complete -c gt \
    -n "__fish_seen_subcommand_from metric; and not __fish_seen_subcommand_from show list compatible 'decimal set' 'binary set'" \
    -a "Universal"
complete -c gt \
    -n "__fish_seen_subcommand_from metric; and not __fish_seen_subcommand_from show list compatible 'decimal set' 'binary set'" \
    -a "Weight"

# ── `custom` subcommand ──────────────────────────────────────────────────────

complete -c gt -n "__fish_seen_subcommand_from custom" \
    -a "show"
complete -c gt -n "__fish_seen_subcommand_from custom" \
    -a "list"
complete -c gt -n "__fish_seen_subcommand_from custom" \
    -a "define"
complete -c gt -n "__fish_seen_subcommand_from custom" \
    -a "undefine"

# custom define <name> <factor> <category> — complete categories
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "Custom"
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "DataRate"
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "DataSize"
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "Distance"
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "Speed"
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "Time"
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "Universal"
complete -c gt \
    -n "__fish_seen_subcommand_from define" \
    -a "Weight"
