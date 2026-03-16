package repl

import (
	"strings"
	"testing"
)

func TestExecuteCommandHelp(t *testing.T) {
	output, err := ExecuteCommand("help")
	if err != nil {
		t.Fatalf("ExecuteCommand('help') returned error: %v", err)
	}
	if output == "" {
		t.Error("ExecuteCommand('help') returned empty output")
	}
	if !strings.Contains(output, "PERC") {
		t.Errorf("ExecuteCommand('help') output should contain 'PERC', got: %s", output[:50])
	}
}

func TestExecuteCommandHelpWithSubcommand(t *testing.T) {
	output, err := ExecuteCommand("help clear")
	if err != nil {
		t.Fatalf("ExecuteCommand('help clear') returned error: %v", err)
	}
	if output == "" {
		t.Error("ExecuteCommand('help clear') returned empty output")
	}
	if !strings.Contains(output, "Clear") {
		t.Errorf("ExecuteCommand('help clear') output should contain 'Clear', got: %s", output[:50])
	}
}

func TestExecuteCommandUnknownCommand(t *testing.T) {
	_, err := ExecuteCommand("unknown")
	if err == nil {
		t.Error("ExecuteCommand('unknown') should return error, got nil")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("Error should mention 'unknown command', got: %v", err)
	}
}

func TestExecuteCommandClear(t *testing.T) {
	_, err := ExecuteCommand("clear")
	if err != nil {
		t.Fatalf("ExecuteCommand('clear') returned error: %v", err)
	}
}

func TestExecuteCommandQuit(t *testing.T) {
	_, err := ExecuteCommand("quit")
	if err != nil {
		t.Fatalf("ExecuteCommand('quit') returned error: %v", err)
	}
}

func TestExecuteCommandExit(t *testing.T) {
	_, err := ExecuteCommand("exit")
	if err != nil {
		t.Fatalf("ExecuteCommand('exit') returned error: %v", err)
	}
}

func TestExecuteCommandEmpty(t *testing.T) {
	_, err := ExecuteCommand("")
	if err != nil {
		t.Fatalf("ExecuteCommand('') returned error: %v", err)
	}
}
