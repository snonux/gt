// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"fmt"
	"slices"
	"strings"
)

// HelpTopic holds the inline help content for one operator, function, or REPL command.
type HelpTopic struct {
	Category   string // e.g. "Arithmetic", "Stack", "Comparison", "Variables", "Hyper", "REPL"
	Operator   string // the operator or command name (e.g. "+", "help", "rat")
	Aliases    []string
	Description string
	Usage      string // short usage hint
	Examples   []string
}

// helpTopics is the single source of truth for all inline help entries.
// Each entry documents one operator, function, or REPL command with examples.
var helpTopics = []HelpTopic{
	// ── REPL commands ──

	{
		Category:    "REPL",
		Operator:    "help",
		Description: "Show help information",
		Usage:       "help [topic]",
		Examples: []string{
			"help              Show general help overview",
			"help +            Show help for the + operator",
			"help dup          Show help for dup",
		},
	},
	{
		Category:    "REPL",
		Operator:    "clear",
		Description: "Clear the terminal screen",
		Usage:       "clear",
		Examples: []string{
			"clear",
		},
	},
	{
		Category:    "REPL",
		Operator:    "quit",
		Aliases:     []string{"exit"},
		Description: "Exit the REPL",
		Usage:       "quit or exit",
		Examples: []string{
			"quit",
		},
	},
	{
		Category:    "REPL",
		Operator:    "rpn",
		Aliases:     []string{"calc"},
		Description: "Evaluate an RPN (Reverse Polish Notation) expression",
		Usage:       "rpn <expression>",
		Examples: []string{
			"rpn 3 4 +         Evaluate 3 + 4 = 7",
			"rpn 10 2 /        Evaluate 10 / 2 = 5",
			"calc 5 3 *        Same as 'rpn 5 3 *'",
		},
	},
	{
		Category:    "REPL",
		Operator:    "rat",
		Description: "Switch between float64 and rational number modes",
		Usage:       "rat on|off|toggle",
		Examples: []string{
			"rat on      Enable rational mode",
			"rat off     Disable rational mode (float64)",
			"rat toggle  Toggle between modes",
		},
	},
	{
		Category:    "REPL",
		Operator:    "stack",
		Description: "Show current RPN stack state",
		Usage:       "stack",
		Examples: []string{
			"1 2 3 stack     Show stack: 1 2 3",
		},
	},

	// ── Arithmetic operators ──

	{
		Category:    "Arithmetic",
		Operator:    "+",
		Description: "Add two numbers (metric-aware)",
		Usage:       "a b +",
		Examples: []string{
			"3 4 +           3 + 4 = 7",
			"100Mbps 50 +    100Mbps + 50Mbps = 150Mbps (Cool absorbs metric)",
		},
	},
	{
		Category:    "Arithmetic",
		Operator:    "-",
		Description: "Subtract two numbers (metric-aware)",
		Usage:       "a b -",
		Examples: []string{
			"10 3 -          10 - 3 = 7",
			"200Mbps 50 -    200Mbps - 50Mbps = 150Mbps",
		},
	},
	{
		Category:    "Arithmetic",
		Operator:    "*",
		Description: "Multiply two numbers",
		Usage:       "a b *",
		Examples: []string{
			"3 4 *           3 * 4 = 12",
			"100Mbps 2 *     100Mbps * 2 = 200Mbps (Cool preserves metric)",
		},
	},
	{
		Category:    "Arithmetic",
		Operator:    "/",
		Description: "Divide two numbers (a / b)",
		Usage:       "a b /",
		Examples: []string{
			"10 3 /          10 / 3 = 3.333...",
		},
	},
	{
		Category:    "Arithmetic",
		Operator:    "^",
		Description: "Raise to power (a ^ b), result is always unitless",
		Usage:       "a b ^",
		Examples: []string{
			"2 3 ^           2 ^ 3 = 8",
			"10 0.5 ^        10 ^ 0.5 = sqrt(10) = 3.162...",
		},
	},
	{
		Category:    "Arithmetic",
		Operator:    "**",
		Description: "Raise to integer power (a ** b) using fast binary exponentiation",
		Usage:       "a b **",
		Examples: []string{
			"2 10 **         2 ^ 10 = 1024",
			"3 6 **          3 ^ 6 = 729",
		},
	},
	{
		Category:    "Arithmetic",
		Operator:    "%",
		Description: "Modulo (remainder of a / b)",
		Usage:       "a b %",
		Examples: []string{
			"10 3 %          10 mod 3 = 1",
			"7 2 %           7 mod 2 = 1",
		},
	},
	{
		Category:    "Arithmetic",
		Operator:    "lg",
		Description: "Logarithm base 2",
		Usage:       "a lg",
		Examples: []string{
			"8 lg            log2(8) = 3",
			"1024 lg         log2(1024) = 10",
		},
	},
	{
		Category:    "Arithmetic",
		Operator:    "log",
		Description: "Logarithm base 10",
		Usage:       "a log",
		Examples: []string{
			"100 log         log10(100) = 2",
			"1000 log        log10(1000) = 3",
		},
	},
	{
		Category:    "Arithmetic",
		Operator:    "ln",
		Description: "Natural logarithm (base e)",
		Usage:       "a ln",
		Examples: []string{
			"2.71828 ln      ln(2.71828) ~= 1",
			"1 ln            ln(1) = 0",
		},
	},

	// ── Comparison operators ──

	{
		Category:    "Comparison",
		Operator:    ">",
		Aliases:     []string{"gt"},
		Description: "Greater than: pushes true (1) or false (0)",
		Usage:       "a b >",
		Examples: []string{
			"5 3 >           5 > 3 = true (1)",
			"3 5 >           3 > 5 = false (0)",
		},
	},
	{
		Category:    "Comparison",
		Operator:    "<",
		Aliases:     []string{"lt"},
		Description: "Less than: pushes true (1) or false (0)",
		Usage:       "a b <",
		Examples: []string{
			"3 5 <           3 < 5 = true (1)",
			"5 3 <           5 < 3 = false (0)",
		},
	},
	{
		Category:    "Comparison",
		Operator:    ">=",
		Aliases:     []string{"gte"},
		Description: "Greater than or equal: pushes true (1) or false (0)",
		Usage:       "a b >=",
		Examples: []string{
			"5 5 >=          5 >= 5 = true (1)",
			"3 5 >=          3 >= 5 = false (0)",
		},
	},
	{
		Category:    "Comparison",
		Operator:    "<=",
		Aliases:     []string{"lte"},
		Description: "Less than or equal: pushes true (1) or false (0)",
		Usage:       "a b <=",
		Examples: []string{
			"3 3 <=          3 <= 3 = true (1)",
			"5 3 <=          5 <= 3 = false (0)",
		},
	},
	{
		Category:    "Comparison",
		Operator:    "==",
		Aliases:     []string{"eq"},
		Description: "Equal: pushes true (1) or false (0)",
		Usage:       "a b ==",
		Examples: []string{
			"5 5 ==          5 == 5 = true (1)",
			"5 3 ==          5 == 3 = false (0)",
		},
	},
	{
		Category:    "Comparison",
		Operator:    "!=",
		Aliases:     []string{"neq"},
		Description: "Not equal: pushes true (1) or false (0)",
		Usage:       "a b !=",
		Examples: []string{
			"5 3 !=          5 != 3 = true (1)",
			"5 5 !=          5 != 5 = false (0)",
		},
	},

	// ── Stack operators ──

	{
		Category:    "Stack",
		Operator:    "dup",
		Description: "Duplicate the top stack value",
		Usage:       "a dup",
		Examples: []string{
			"5 dup           Stack: 5 5",
			"3 dup dup       Stack: 3 3 3",
		},
	},
	{
		Category:    "Stack",
		Operator:    "swap",
		Description: "Swap the top two stack values",
		Usage:       "a b swap",
		Examples: []string{
			"1 2 swap        Stack: 2 1",
		},
	},
	{
		Category:    "Stack",
		Operator:    "pop",
		Description: "Remove and discard the top stack value",
		Usage:       "pop",
		Examples: []string{
			"1 2 3 pop       Stack: 1 2",
		},
	},
	{
		Category:    "Stack",
		Operator:    "d",
		Description: "Pop a symbol from stack and delete that variable",
		Usage:       ":x d",
		Examples: []string{
			":x d            Delete variable x",
		},
	},
	{
		Category:    "Stack",
		Operator:    "show",
		Aliases:     []string{"showstack", "print"},
		Description: "Display the current stack contents",
		Usage:       "show",
		Examples: []string{
			"1 2 3 show      Prints: 1 2 3",
			"show            Prints: Stack is empty",
		},
	},

	// ── Variable operators ──

	{
		Category:    "Variables",
		Operator:    ":=",
		Description: "Assign value to variable (name on bottom, value on top)",
		Usage:       ":name value :=",
		Examples: []string{
			":x 5 :=         Set variable x to 5",
			":y 3.14 :=      Set variable y to 3.14",
		},
	},
	{
		Category:    "Variables",
		Operator:    "=:",
		Description: "Assign value to variable (value on bottom, name on top)",
		Usage:       "value :name =:",
		Examples: []string{
			"5 :x =:         Set variable x to 5",
		},
	},
	{
		Category:    "Variables",
		Operator:    "vars",
		Description: "List all defined variables and their values",
		Usage:       "vars",
		Examples: []string{
			"vars            List all variables",
		},
	},
	{
		Category:    "Variables",
		Operator:    "clear",
		Description: "Clear all user-defined variables",
		Usage:       "clear",
		Examples: []string{
			"clear           Remove all variables",
		},
	},
	{
		Category:    "Variables",
		Operator:    "convert",
		Description: "Convert a value to a target metric (@X syntax)",
		Usage:       "value @target convert",
		Examples: []string{
			"100Mbps @bps convert    Convert 100Mbps to bps",
		},
	},

	// ── Constants ──

	{
		Category:    "Constants",
		Operator:    "constants",
		Description: "List all available constants",
		Usage:       "constants",
		Examples: []string{
			"constants       List all built-in constants",
		},
	},
	{
		Category:    "Constants",
		Operator:    "clearconstants",
		Description: "Reset all constants to built-in defaults",
		Usage:       "clearconstants",
		Examples: []string{
			"clearconstants  Reset constants",
		},
	},

	// ── Hyper (n-ary) operators ──

	{
		Category:    "Hyper",
		Operator:    "[+]",
		Description: "Add all stack values together (n-ary)",
		Usage:       "a b c [+] ...",
		Examples: []string{
			"1 2 3 [+]       1 + 2 + 3 = 6",
			"10 20 30 40 [+] 100",
		},
	},
	{
		Category:    "Hyper",
		Operator:    "[-]",
		Description: "Subtract all stack values left-associative (n-ary)",
		Usage:       "a b c [-] ...",
		Examples: []string{
			"10 3 2 [-]      10 - 3 - 2 = 5",
		},
	},
	{
		Category:    "Hyper",
		Operator:    "[*]",
		Description: "Multiply all stack values together (n-ary)",
		Usage:       "a b c [*] ...",
		Examples: []string{
			"2 3 4 [*]       2 * 3 * 4 = 24",
		},
	},
	{
		Category:    "Hyper",
		Operator:    "[/]",
		Description: "Divide all stack values left-associative (n-ary)",
		Usage:       "a b c [/] ...",
		Examples: []string{
			"100 2 5 [/]     100 / 2 / 5 = 10",
		},
	},
	{
		Category:    "Hyper",
		Operator:    "[^]",
		Description: "Power all stack values left-associative (n-ary)",
		Usage:       "a b c [^] ...",
		Examples: []string{
			"2 3 2 [^]       2 ^ 3 ^ 2 = 8",
		},
	},
	{
		Category:    "Hyper",
		Operator:    "[%]",
		Description: "Modulo all stack values left-associative (n-ary)",
		Usage:       "a b c [%] ...",
		Examples: []string{
			"100 10 3 [%]    100 % 10 % 3 = 1",
		},
	},
	{
		Category:    "Hyper",
		Operator:    "[lg]",
		Description: "Sum of log2 of all stack values (n-ary)",
		Usage:       "a b [lg] ...",
		Examples: []string{
			"2 4 [lg]        log2(2) + log2(4) = 1 + 2 = 3",
		},
	},
	{
		Category:    "Hyper",
		Operator:    "[log]",
		Description: "Sum of log10 of all stack values (n-ary)",
		Usage:       "a b [log] ...",
		Examples: []string{
			"10 100 [log]    log10(10) + log10(100) = 1 + 2 = 3",
		},
	},
	{
		Category:    "Hyper",
		Operator:    "[ln]",
		Description: "Sum of natural log of all stack values (n-ary)",
		Usage:       "a b [ln] ...",
		Examples: []string{
			"1 2.71828 [ln]  ln(1) + ln(2.71828) = 0 + 1 = 1",
		},
	},
}

// buildHelpIndex builds lookup maps from helpTopics at package init.
var (
	helpByTopic   = make(map[string]*HelpTopic)     // operator -> topic
	helpByAlias   = make(map[string]*HelpTopic)     // alias -> topic
	helpByCat     = make(map[string][]string)       // category -> []operator names
	categoryOrder []string                          // insertion order of categories
)

func init() {
	seenCat := make(map[string]bool)
	for i := range helpTopics {
		t := &helpTopics[i]
		for _, a := range t.Aliases {
			helpByAlias[a] = t
		}
		if !seenCat[t.Category] {
			seenCat[t.Category] = true
			categoryOrder = append(categoryOrder, t.Category)
		}
		helpByCat[t.Category] = append(helpByCat[t.Category], t.Operator)
	}
	// Iterate in reverse so REPL entries (last in slice) take priority
	// in helpByTopic. This means `help clear` shows the REPL screen-clear
	// help rather than the RPN variable-clear help.
	for i := len(helpTopics) - 1; i >= 0; i-- {
		t := &helpTopics[i]
		helpByTopic[t.Operator] = t
	}
}

// GetHelp returns the formatted help text for a topic.
// If topic is empty, it returns the general help overview.
// If topic is "categories", it lists all categories.
func GetHelp(topic string) string {
	if topic == "" {
		return getGeneralHelp()
	}
	if topic == "categories" {
		return getCategoriesHelp()
	}

	// Look up by operator name or alias
	topic = strings.ToLower(topic)
	t, ok := helpByTopic[topic]
	if !ok {
		t, ok = helpByAlias[topic]
	}
	if !ok {
		return fmt.Sprintf("No help for %q.\n\nType 'help' for an overview, or 'help categories' to list all topics.", topic)
	}
	return formatTopic(t)
}

// GetAllTopics returns a sorted list of all help topic names.
func GetAllTopics() []string {
	topics := make([]string, 0, len(helpByTopic))
	for op := range helpByTopic {
		topics = append(topics, op)
	}
	slices.Sort(topics)
	return topics
}

// GetCompletionTopics returns all help topics suitable for tab completion.
// Includes operator names, aliases, and special topics like "categories".
func GetCompletionTopics() []string {
	seen := make(map[string]bool)
	var topics []string

	// Add "categories" as a special topic
	topics = append(topics, "categories")

	// Add all operator names (skip "help" itself — no need to complete help help)
	for op, t := range helpByTopic {
		if op == "help" {
			continue
		}
		topics = append(topics, op)
		seen[op] = true
		// Add aliases
		for _, a := range t.Aliases {
			if !seen[a] {
				topics = append(topics, a)
				seen[a] = true
			}
		}
	}
	slices.Sort(topics)
	return topics
}

// getGeneralHelp returns the overview help text.
func getGeneralHelp() string {
	var sb strings.Builder

	sb.WriteString("gt - Reverse Polish Notation (RPN) Calculator\n\n")
	sb.WriteString("Help Topics:\n")
	sb.WriteString("  help [topic]       Show help for a specific topic\n")
	sb.WriteString("  help categories    List all available help topics by category\n\n")

	sb.WriteString("REPL Commands:\n")
	for _, op := range helpByCat["REPL"] {
		t := helpByTopic[op]
		sb.WriteString(fmt.Sprintf("  %-12s %s\n", op, t.Description))
	}

	sb.WriteString("\nArithmetic Operators:\n")
	for _, op := range helpByCat["Arithmetic"] {
		t := helpByTopic[op]
		sb.WriteString(fmt.Sprintf("  %-8s %s\n", op, t.Description))
	}

	sb.WriteString("\nComparison Operators:\n")
	for _, op := range helpByCat["Comparison"] {
		t := helpByTopic[op]
		sb.WriteString(fmt.Sprintf("  %-8s %s\n", op, t.Description))
	}

	sb.WriteString("\nStack Operators:\n")
	for _, op := range helpByCat["Stack"] {
		t := helpByTopic[op]
		sb.WriteString(fmt.Sprintf("  %-12s %s\n", op, t.Description))
	}

	sb.WriteString("\nHyper (n-ary) Operators:\n")
	for _, op := range helpByCat["Hyper"] {
		t := helpByTopic[op]
		sb.WriteString(fmt.Sprintf("  %-8s %s\n", op, t.Description))
	}

	sb.WriteString("\nVariables and Constants:\n")
	for _, cat := range []string{"Variables", "Constants"} {
		for _, op := range helpByCat[cat] {
			t := helpByTopic[op]
			sb.WriteString(fmt.Sprintf("  %-16s %s\n", op, t.Description))
		}
	}

	sb.WriteString("\nKeyboard Shortcuts:\n")
	sb.WriteString("  Ctrl+A         Go to beginning of line\n")
	sb.WriteString("  Ctrl+E         Go to end of line\n")
	sb.WriteString("  Ctrl+L         Clear screen\n")
	sb.WriteString("  Ctrl+D         Exit / delete character\n")
	sb.WriteString("  Up/Down        History navigation\n")
	sb.WriteString("  Tab            Auto-complete\n")

	sb.WriteString("\nUse 'help <topic>' for details on any operator or command.\n")
	sb.WriteString("Use 'help categories' to see all available topics.\n")

	return sb.String()
}

// getCategoriesHelp lists all categories and their topics.
func getCategoriesHelp() string {
	var sb strings.Builder
	sb.WriteString("Available help topics by category:\n\n")
	for _, cat := range categoryOrder {
		sb.WriteString(fmt.Sprintf("  %s:\n", cat))
		for _, op := range helpByCat[cat] {
			t := helpByTopic[op]
			aliasStr := ""
			if len(t.Aliases) > 0 {
				aliasStr = " (" + strings.Join(t.Aliases, ", ") + ")"
			}
			sb.WriteString(fmt.Sprintf("    %-16s %s%s\n", op, t.Description, aliasStr))
		}
		sb.WriteString("\n")
	}
	sb.WriteString("Use 'help <topic>' for details on any topic.\n")
	return sb.String()
}

// formatTopic returns the detailed help for a single HelpTopic.
func formatTopic(t *HelpTopic) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Topic:    %s (Category: %s)\n", t.Operator, t.Category))
	if len(t.Aliases) > 0 {
		sb.WriteString(fmt.Sprintf("Aliases:  %s\n", strings.Join(t.Aliases, ", ")))
	}
	sb.WriteString(fmt.Sprintf("Usage:    %s\n", t.Usage))
	sb.WriteString(fmt.Sprintf("Desc:     %s\n", t.Description))

	sb.WriteString("\nExamples:\n")
	for _, ex := range t.Examples {
		sb.WriteString(fmt.Sprintf("  %s\n", ex))
	}

	return sb.String()
}
