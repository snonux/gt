// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"os"
	"strings"
	"testing"
)

func TestHistoryManagerNew(t *testing.T) {
	hm := NewHistoryManager(".test_history")
	if hm == nil {
		t.Fatal("NewHistoryManager returned nil")
	}
	if hm.historyFile != ".test_history" {
		t.Errorf("historyFile = %q, want %q", hm.historyFile, ".test_history")
	}
	if hm.maxEntries != 1000 {
		t.Errorf("maxEntries = %d, want 1000", hm.maxEntries)
	}
}

func TestHistoryManagerPath(t *testing.T) {
	hm := NewHistoryManager(".test_history")
	path := hm.Path()
	if path == "" {
		t.Fatal("Path returned empty string")
	}
	if !strings.HasSuffix(path, ".test_history") {
		t.Errorf("Path should end with .test_history, got %q", path)
	}
}

func TestHistoryManagerPathWithBaseDir(t *testing.T) {
	hm := NewHistoryManager(".test_history").WithBaseDir("/tmp")
	path := hm.Path()
	if path != "/tmp/.test_history" {
		t.Errorf("Path = %q, want /tmp/.test_history", path)
	}
}

func TestHistoryManagerSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	hm := NewHistoryManager(".test_history").WithBaseDir(tmpDir)

	entries := []string{"20% of 150", "rpn 3 4 +", "help"}

	if err := hm.Save(entries); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	loaded := hm.Load()
	if loaded == nil {
		t.Fatal("Load returned nil")
	}
	if len(loaded) != len(entries) {
		t.Fatalf("Load returned %d entries, want %d", len(loaded), len(entries))
	}
	for i, want := range entries {
		if loaded[i] != want {
			t.Errorf("entry %d = %q, want %q", i, loaded[i], want)
		}
	}
}

func TestHistoryManagerLoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	hm := NewHistoryManager(".nonexistent").WithBaseDir(tmpDir)
	history := hm.Load()
	if history != nil {
		t.Errorf("Load from non-existent file returned %d entries, want nil", len(history))
	}
}

func TestHistoryManagerLoadEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	hm := NewHistoryManager(".empty_history").WithBaseDir(tmpDir)

	// Create empty file
	if err := os.WriteFile(hm.Path(), []byte{}, 0644); err != nil {
		t.Fatalf("failed to create empty file: %v", err)
	}

	history := hm.Load()
	// Empty file should return empty slice (not nil) since scanner.Scan returns false immediately
	if len(history) != 0 {
		t.Errorf("Load from empty file returned %d entries, want 0", len(history))
	}
}

func TestHistoryManagerSaveRespectsMaxEntries(t *testing.T) {
	tmpDir := t.TempDir()
	hm := NewHistoryManager(".test_max").WithBaseDir(tmpDir)

	// Create 1500 entries (more than default max of 1000)
	largeEntries := make([]string, 1500)
	for i := range largeEntries {
		largeEntries[i] = strings.Repeat("x", 10)
	}

	if err := hm.Save(largeEntries); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	loaded := hm.Load()
	if loaded == nil {
		t.Fatal("Load returned nil after Save")
	}
	if len(loaded) != 1000 {
		t.Errorf("after Save(1500 entries), Load returned %d, want 1000 (maxEntries)", len(loaded))
	}
	// Verify it keeps the LAST 1000 entries, not the first 1000
	if loaded[0] != largeEntries[500] {
		t.Error("truncated entries should keep the last 1000, not the first 1000")
	}
	if loaded[999] != largeEntries[1499] {
		t.Error("last loaded entry should be last saved entry")
	}
}

func TestHistoryManagerSaveCreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	hm := NewHistoryManager(".gt_test_history").WithBaseDir(tmpDir)

	entries := []string{"entry1", "entry2", "entry3"}
	if err := hm.Save(entries); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	// Verify file was created
	data, err := os.ReadFile(hm.Path())
	if err != nil {
		t.Fatalf("failed to read created file: %v", err)
	}

	content := string(data)
	for _, entry := range entries {
		if !strings.Contains(content, entry) {
			t.Errorf("file should contain %q, got %q", entry, content)
		}
	}
}

func TestHistoryManagerSaveWithEmptySlice(t *testing.T) {
	tmpDir := t.TempDir()
	hm := NewHistoryManager(".empty_save").WithBaseDir(tmpDir)

	if err := hm.Save(nil); err != nil {
		t.Fatalf("Save(nil) returned error: %v", err)
	}

	data, err := os.ReadFile(hm.Path())
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(data) != "" {
		t.Errorf("Save(nil) should produce empty file, got %q", string(data))
	}
}

func TestHistoryManagerSaveWithSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	hm := NewHistoryManager(".special_chars").WithBaseDir(tmpDir)

	entries := []string{
		"help",
		`"quoted entry"`,
		"entry with spaces",
		"rpn 3 4 +",
	}

	if err := hm.Save(entries); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	loaded := hm.Load()
	if len(loaded) != len(entries) {
		t.Fatalf("Load returned %d entries, want %d", len(loaded), len(entries))
	}
	for i, want := range entries {
		if loaded[i] != want {
			t.Errorf("entry %d = %q, want %q", i, loaded[i], want)
		}
	}
}

func TestHistoryManagerSaveOverwritesExisting(t *testing.T) {
	tmpDir := t.TempDir()
	hm := NewHistoryManager(".overwrite_test").WithBaseDir(tmpDir)

	// Save initial entries
	if err := hm.Save([]string{"old1", "old2"}); err != nil {
		t.Fatalf("first Save: %v", err)
	}

	// Save new entries
	if err := hm.Save([]string{"new1"}); err != nil {
		t.Fatalf("second Save: %v", err)
	}

	loaded := hm.Load()
	if len(loaded) != 1 {
		t.Fatalf("Load returned %d entries, want 1", len(loaded))
	}
	if loaded[0] != "new1" {
		t.Errorf("Load returned %q, want new1 (file should be overwritten)", loaded[0])
	}
}

func TestHistoryManagerWithEmptyBaseDir(t *testing.T) {
	// When baseDir is empty, Path() falls back to os.UserHomeDir()
	hm := NewHistoryManager(".test").WithBaseDir("")
	path := hm.Path()
	if path == "" {
		t.Fatal("Path should fall back to os.UserHomeDir when baseDir is empty")
	}
	if !strings.HasSuffix(path, ".test") {
		t.Errorf("Path = %q, should end with .test", path)
	}
}
