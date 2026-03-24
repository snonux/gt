package repl

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// HistoryManager handles history file operations.
type HistoryManager struct {
	historyFile string
	maxEntries  int
}

// NewHistoryManager creates a new history manager with the given file name.
func NewHistoryManager(historyFile string) *HistoryManager {
	return &HistoryManager{
		historyFile: historyFile,
		maxEntries:  1000, // Default max history entries
	}
}

// Path returns the path to the history file.
func (h *HistoryManager) Path() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, h.historyFile)
}

// Load reads history from file.
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

// Save writes history to file, keeping only the most recent entries.
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
