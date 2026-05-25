// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"fmt"
	"slices"
	"strings"
)

// GetHelp returns the formatted help text for a topic.
// If topic is empty, it returns the general help overview.
func GetHelp(topic string) string {
	if topic == "" {
		return getGeneralHelp()
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
	if t.Render != nil {
		return t.Render()
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
// Includes operator names and aliases from all registered HelpTopic entries.
func GetCompletionTopics() []string {
	seen := make(map[string]struct{})
	var topics []string

	// Add all operator names (skip "help" itself — no need to complete help help)
	for op, t := range helpByTopic {
		if op == "help" {
			continue
		}
		topics = append(topics, op)
		seen[op] = struct{}{}
		// Add aliases
		for _, a := range t.Aliases {
			if _, ok := seen[a]; !ok {
				topics = append(topics, a)
				seen[a] = struct{}{}
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
