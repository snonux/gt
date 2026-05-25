// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"strings"
	"testing"
)

func TestGetHelpEmptyReturnsOverview(t *testing.T) {
	output := GetHelp("")
	if !strings.Contains(output, "gt") {
		t.Error("General help should mention 'gt'")
	}
	if !strings.Contains(output, "help categories") {
		t.Error("General help should mention 'help categories'")
	}
	if !strings.Contains(output, "Arithmetic") {
		t.Error("General help should list Arithmetic section")
	}
}

func TestGetHelpCategories(t *testing.T) {
	output := GetHelp("categories")
	if !strings.Contains(output, "REPL") {
		t.Error("Categories should list REPL")
	}
	if !strings.Contains(output, "Arithmetic") {
		t.Error("Categories should list Arithmetic")
	}
	if !strings.Contains(output, "Hyper") {
		t.Error("Categories should list Hyper")
	}
}

func TestGetHelpClearReturnsREPL(t *testing.T) {
	// "clear" exists in both REPL (screen clear) and Variables (clear vars)
	// helpByTopic should prefer the REPL entry
	output := GetHelp("clear")
	if !strings.Contains(output, "screen") && !strings.Contains(output, "terminal") {
		t.Errorf("'help clear' should show REPL screen clear, got: %s", output[:80])
	}
}

func TestGetHelpKnownOperator(t *testing.T) {
	output := GetHelp("+")
	if !strings.Contains(output, "Add") {
		t.Errorf("'help +' should describe Add, got: %s", output[:80])
	}
	if !strings.Contains(output, "Examples:") {
		t.Error("'help +' should have examples")
	}
	if !strings.Contains(output, "3 4 +") {
		t.Error("'help +' should have example '3 4 +'")
	}
}

func TestGetHelpOperatorWithAliases(t *testing.T) {
	// Test by alias
	output := GetHelp("gt")
	if !strings.Contains(output, ">") {
		t.Errorf("'help gt' should show > operator, got: %s", output[:80])
	}
	if !strings.Contains(output, "Aliases:") {
		t.Error("'help gt' should show aliases")
	}
}

func TestGetHelpHyperOperator(t *testing.T) {
	output := GetHelp("[+]")
	if !strings.Contains(output, "Add all stack values") {
		t.Errorf("'help [+]' should describe hyper add, got: %s", output[:80])
	}
	if !strings.Contains(output, "Hyper") {
		t.Error("'help [+]' should be in Hyper category")
	}
}

func TestGetHelpVariableOperator(t *testing.T) {
	output := GetHelp(":=")
	if !strings.Contains(output, "Assign") {
		t.Errorf("'help :=' should describe assignment, got: %s", output[:80])
	}
	if !strings.Contains(output, "Variables") {
		t.Error("'help :=' should be in Variables category")
	}
}

func TestGetHelpUnknownTopic(t *testing.T) {
	output := GetHelp("nonexistent")
	if !strings.Contains(output, "No help for") {
		t.Errorf("'help nonexistent' should say no help, got: %s", output)
	}
	if !strings.Contains(output, "help categories") {
		t.Error("'help nonexistent' should suggest 'help categories'")
	}
}

func TestGetHelpCaseInsensitive(t *testing.T) {
	output := GetHelp("lg")
	outputUpper := GetHelp("LG")
	if output != outputUpper {
		t.Error("Help should handle case consistently for lg/LG")
	}
}

func TestGetAllTopics(t *testing.T) {
	topics := GetAllTopics()
	if len(topics) == 0 {
		t.Error("GetAllTopics() should return topics")
	}

	// Check known topics are present
	expectedTopics := []string{"+", "-", "*", "/", "dup", "swap", "help", "rat", "[+]"}
	for _, expected := range expectedTopics {
		found := false
		for _, t := range topics {
			if t == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetAllTopics() missing topic %q", expected)
		}
	}
}

func TestGetCompletionTopicsIncludesAliases(t *testing.T) {
	topics := GetCompletionTopics()

	// Should include aliases
	expectedAliases := []string{"exit", "calc", "gt", "lt", "categories"}
	for _, expected := range expectedAliases {
		found := false
		for _, t := range topics {
			if t == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("GetCompletionTopics() missing alias/topic %q", expected)
		}
	}
}

func TestGetCompletionTopicsIsSorted(t *testing.T) {
	topics := GetCompletionTopics()
	for i := 1; i < len(topics); i++ {
		if topics[i] < topics[i-1] {
			t.Errorf("Topics not sorted: %q > %q at index %d", topics[i-1], topics[i], i)
		}
	}
}

func TestFormatTopic(t *testing.T) {
	topic := helpByTopic["+"]
	if topic == nil {
		t.Fatal("helpByTopic[\"+\"] is nil")
	}

	output := formatTopic(topic)
	if !strings.Contains(output, "Topic:") {
		t.Error("formatTopic should contain 'Topic:'")
	}
	if !strings.Contains(output, "Usage:") {
		t.Error("formatTopic should contain 'Usage:'")
	}
	if !strings.Contains(output, "Desc:") {
		t.Error("formatTopic should contain 'Desc:'")
	}
	if !strings.Contains(output, "Examples:") {
		t.Error("formatTopic should contain 'Examples:'")
	}
}

func TestHelpTopicsNoDuplicates(t *testing.T) {
	// Operators can appear in multiple categories (e.g. "clear" in REPL and Variables)
	// Check for true duplicates: same operator within the same category
	seen := make(map[string]map[string]bool) // category -> operators
	for _, topic := range helpTopics {
		if seen[topic.Category] == nil {
			seen[topic.Category] = make(map[string]bool)
		}
		if seen[topic.Category][topic.Operator] {
			t.Errorf("Duplicate topic operator %q in category %q", topic.Operator, topic.Category)
		}
		seen[topic.Category][topic.Operator] = true
	}
}

func TestHelpByAliasResolution(t *testing.T) {
	// Test that aliases resolve to the correct topic
	output := GetHelp("exit")
	if !strings.Contains(output, "Exit") {
		t.Errorf("'help exit' should describe exit, got: %s", output[:80])
	}

	output = GetHelp("calc")
	if !strings.Contains(output, "RPN") {
		t.Errorf("'help calc' should describe RPN, got: %s", output[:80])
	}

	output = GetHelp("showstack")
	if !strings.Contains(output, "stack") {
		t.Errorf("'help showstack' should mention stack, got: %s", output[:80])
	}
}

func TestHelpCompleterIntegration(t *testing.T) {
	adapter := NewAutoCompleter()
	if adapter == nil {
		t.Fatal("NewAutoCompleter returned nil")
	}

	// Test help topic completion
	matches, _ := adapter.Do([]rune("help +"), 6)
	if len(matches) != 1 {
		t.Errorf("'help +' should match +, got %d matches: %v", len(matches), matches)
	}

	// Test partial help topic completion
	matches, _ = adapter.Do([]rune("help du"), 7)
	if len(matches) != 1 {
		t.Errorf("'help du' should match dup, got %d matches: %v", len(matches), matches)
	}
}

func TestCmdHelpIntegration(t *testing.T) {
	// cmdHelp should delegate to GetHelp
	output := cmdHelp(nil)
	if !strings.Contains(output, "gt") {
		t.Error("cmdHelp(nil) should return general help")
	}

	output = cmdHelp([]string{"+"})
	if !strings.Contains(output, "Add") {
		t.Errorf("cmdHelp(['+']) should return + help, got: %s", output[:80])
	}

	output = cmdHelp([]string{"categories"})
	if !strings.Contains(output, "REPL") {
		t.Error("cmdHelp(['categories']) should list categories")
	}
}
