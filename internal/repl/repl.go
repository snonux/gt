package repl

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"codeberg.org/snonux/perc/internal/calculator"
	"codeberg.org/snonux/perc/internal/rpn"
	"github.com/mattn/go-isatty"

	"github.com/c-bata/go-prompt"
)

const historyFile = ".perc_history"

// RPNState holds the state for RPN operations in REPL
type RPNState struct {
	vars   rpn.VariableStore
	rpnCalc *rpn.RPN
}

var rpnState *RPNState

// getRPNState returns or creates the RPN state
func getRPNState() *RPNState {
	if rpnState == nil {
		vars := rpn.NewVariables()
		rpnState = &RPNState{
			vars:    vars,
			rpnCalc: rpn.NewRPN(vars),
		}
	}
	return rpnState
}

// executor runs a calculation command and returns the result
func executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}

	// Check if it's a built-in command
	if cmd, ok := isBuiltinCommand(input); ok {
		output, err := ExecuteCommand(cmd)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		if output != "" {
			fmt.Println(output)
		}
		// Don't add built-in commands to history
		return
	}

	// Check for rpn command prefix
	if strings.HasPrefix(strings.ToLower(input), "rpn ") || strings.HasPrefix(strings.ToLower(input), "calc ") {
		// Extract the expression after rpn/calc
		rest := strings.TrimSpace(strings.TrimPrefix(input, strings.SplitN(input, " ", 2)[0]))
		result, err := runRPN(rest)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println(result)
		return
	}

	// Try RPN parsing first (for bare RPN expressions like "3 4 +")
	rpnResult, rpnErr := runRPN(input)
	if rpnErr == nil {
		// Valid RPN expression - print result
		fmt.Println(rpnResult)
		return
	}

	// Try evaluating as a single operator on the current RPN stack
	// This allows incremental operations like: "1 2 +" then "+"
	state := getRPNState()
	fields := strings.Fields(input)
	if len(fields) == 1 {
		op := strings.ToLower(fields[0])
		// Check if it's a valid operator
		switch op {
		case "+", "-", "*", "/", "^", "%", "dup", "swap", "pop", "show", "clear", "vars":
			result, err := state.rpnCalc.EvalOperator(op)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				fmt.Println(result)
			}
			return
		}
	}

	// Run the percentage calculation
	result, err := calculator.Parse(input)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(result)
}

// runRPN parses and evaluates an RPN expression
func runRPN(input string) (string, error) {
	state := getRPNState()
	return state.rpnCalc.ParseAndEvaluate(input)
}

// isBuiltinCommand checks if input starts with a built-in command
func isBuiltinCommand(input string) (string, bool) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return "", false
	}

	cmd := strings.ToLower(args[0])
	for _, builtin := range builtinCommands {
		if cmd == builtin {
			return input, true
		}
	}
	return "", false
}

// getHistoryPath returns the path to the history file
func getHistoryPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, historyFile)
}

// loadHistory loads history from file
func loadHistory() []string {
	historyPath := getHistoryPath()
	if historyPath == "" {
		return nil
	}

	file, err := os.Open(historyPath)
	if err != nil {
		return nil
	}
	defer file.Close()

	var history []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		history = append(history, scanner.Text())
	}
	return history
}

// saveHistory saves history to file
func saveHistory(history []string) error {
	historyPath := getHistoryPath()
	if historyPath == "" {
		return nil
	}

	// Keep only last 1000 entries to prevent unlimited growth
	if len(history) > 1000 {
		history = history[len(history)-1000:]
	}

	file, err := os.Create(historyPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, entry := range history {
		if _, err := writer.WriteString(entry + "\n"); err != nil {
			return err
		}
	}
	return writer.Flush()
}

// RunREPL starts the interactive REPL
func RunREPL() error {
	// Check if stdin is a TTY
	if !isatty.IsTerminal(os.Stdin.Fd()) {
		fmt.Fprintln(os.Stderr, "REPL mode requires a TTY. Use 'perc <calculation>' for non-interactive mode.")
		return fmt.Errorf("stdin is not a TTY")
	}

	history := loadHistory()

	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("perc - Percentage Calculator"),
		prompt.OptionPrefix("perc> "),
		prompt.OptionLivePrefix(func() (string, bool) {
			return "perc> ", true
		}),
		prompt.OptionHistory(history),
	)

	// Handle SIGINT gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	go func() {
		<-sigChan
		fmt.Println("\nUse 'quit' or 'exit' to exit, or Ctrl+D")
	}()

	// Run the prompt
	p.Run()

	// Note: History is not saved automatically in this version
	// The prompt library stores it in memory but doesn't expose a getter

	return nil
}

// completer provides auto-completion for built-in commands
func completer(d prompt.Document) []prompt.Suggest {
	text := d.GetWordBeforeCursor()
	if text == "" {
		return nil
	}

	var suggestions []prompt.Suggest
	for _, cmd := range builtinCommands {
		if strings.HasPrefix(strings.ToLower(cmd), strings.ToLower(text)) {
			suggestions = append(suggestions, prompt.Suggest{Text: cmd, Description: getCommandDescription(cmd)})
		}
	}
	return suggestions
}

func getCommandDescription(cmd string) string {
	descriptions := map[string]string{
		"help":   "Show help information",
		"clear":  "Clear the screen",
		"quit":   "Exit the REPL",
		"exit":   "Exit the REPL",
		"rpn":    "Evaluate an RPN (postfix notation) expression",
		"calc":   "Same as rpn - evaluate an RPN expression",
	}
	return descriptions[cmd]
}
