// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"codeberg.org/snonux/gt/internal/rpn"
)

func TestNewRPNState(t *testing.T) {
	vars := rpn.NewVariables()
	rpnCalc := rpn.NewRPN(vars, nil)
	state := NewRPNState(vars, rpnCalc)
	if state == nil {
		t.Fatal("NewRPNState() returned nil")
	}
}

func TestRPNStateLoadVariablesEmpty(t *testing.T) {
	// When varStoreFile is empty, LoadVariables should return nil
	state := &RPNState{
		vars:      rpn.NewVariables(),
		rpnCalc:   rpn.NewRPN(rpn.NewVariables(), nil),
		varStoreFile: "",
	}
	err := state.LoadVariables()
	if err != nil {
		t.Errorf("LoadVariables() with empty path = %v, want nil", err)
	}
}

func TestRPNStateSaveVariablesEmpty(t *testing.T) {
	state := &RPNState{
		vars:        rpn.NewVariables(),
		rpnCalc:     rpn.NewRPN(rpn.NewVariables(), nil),
		varStoreFile: "",
	}
	err := state.SaveVariables()
	if err != nil {
		t.Errorf("SaveVariables() with empty path = %v, want nil", err)
	}
}

func TestGetVarStoreFilePath(t *testing.T) {
	path := getVarStoreFilePath()
	if path == "" {
		t.Fatal("getVarStoreFilePath() returned empty string")
	}
	if !strings.Contains(path, ".local") {
		t.Errorf("path should contain '.local', got %q", path)
	}
	if !strings.Contains(path, "gt") {
		t.Errorf("path should contain 'gt', got %q", path)
	}
	if !strings.Contains(path, "vars") {
		t.Errorf("path should contain 'vars', got %q", path)
	}
}

func TestRPNStateSaveLoadRoundTrip(t *testing.T) {
	tmpDir := t.TempDir()
	varStorePath := filepath.Join(tmpDir, "vars.json")

	vars := rpn.NewVariables()
	vars.SetVariable("x", 42.0)
	vars.SetVariable("y", 3.14)

	rpnCalc := rpn.NewRPN(vars, nil)
	state := &RPNState{
		vars:         vars,
		rpnCalc:      rpnCalc,
		varStoreFile: varStorePath,
	}

	err := state.SaveVariables()
	if err != nil {
		t.Fatalf("SaveVariables() error: %v", err)
	}

	// Verify the file exists
	if _, err := os.Stat(varStorePath); os.IsNotExist(err) {
		t.Fatal("SaveVariables() did not create file")
	}

	// Load into a new variable store
	newVars := rpn.NewVariables()
	newState := &RPNState{
		vars:         newVars,
		rpnCalc:      rpn.NewRPN(newVars, nil),
		varStoreFile: varStorePath,
	}

	err = newState.LoadVariables()
	if err != nil {
		t.Fatalf("LoadVariables() error: %v", err)
	}

	val, ok := newVars.GetVariable("x")
	if !ok {
		t.Fatal("variable x not found after load")
	}
	if val != 42.0 {
		t.Errorf("x = %v, want 42.0", val)
	}

	val, ok = newVars.GetVariable("y")
	if !ok {
		t.Fatal("variable y not found after load")
	}
	if val != 3.14 {
		t.Errorf("y = %v, want 3.14", val)
	}
}
