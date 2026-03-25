// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package repl

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// HistoryManager handles history file operations for the REPL.
// It provides methods to load, save, and manage command history with a maximum entry limit.
type HistoryManager struct {
	historyFile string
	maxEntries  int
}

// NewHistoryManager creates a new history manager with the given file name.
// The history manager will store up to maxEntries (default: 1000) in the history file.
//
// historyFile: the filename to use for history (without path)
// Returns a new HistoryManager instance
func NewHistoryManager(historyFile string) *HistoryManager {
	return &HistoryManager{
		historyFile: historyFile,
		maxEntries:  1000, // Default max history entries
	}
}

// Path returns the absolute path to the history file.
// The history file is stored in the user's home directory.
//
// Returns the full path to the history file, or empty string if the home directory cannot be determined
func (h *HistoryManager) Path() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, h.historyFile)
}

// Load reads history from the history file.
// It returns all entries from the file, or nil if the file doesn't exist.
//
// Returns a slice of history entries (each line is one entry), or nil on error
func (h *HistoryManager) Load() []string {
	path := h.Path()
	if path == "" {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer func() {
		_ = file.Close()
	}()

	var history []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		history = append(history, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil
	}
	return history
}

// Save writes history to the history file, keeping only the most recent entries.
// It ensures the file doesn't grow unlimited by keeping only the last maxEntries.
// The function creates the file if it doesn't exist and truncates it if needed.
//
// history: the slice of history entries to save
// Returns an error if the file cannot be written
func (h *HistoryManager) Save(history []string) error {
	path := h.Path()
	if path == "" {
		return nil
	}

	// Keep only last maxEntries entries to prevent unlimited growth
	if len(history) > h.maxEntries {
		history = history[len(history)-h.maxEntries:]
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	writer := bufio.NewWriter(file)
	for _, entry := range history {
		if _, err := writer.WriteString(entry + "\n"); err != nil {
			return fmt.Errorf("failed to write history entry: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush history writer: %w", err)
	}
	return nil
}
