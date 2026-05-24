// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

package rpn

import (
	"fmt"
	"strconv"
	"strings"
)

// assignmentHandler manages all assignment strategies.
// It delegates assignment parsing to specialized handlers for :=, =:, and = operators.
type assignmentHandler struct {
	registry *assignmentRegistry
}

// assignmentRegistry maintains a registry of assignment strategies.
type assignmentRegistry struct {
	strategies []AssignmentStrategy
}

// AssignmentStrategy represents a function that attempts to parse and handle an assignment.
// Returns (result string, handled bool, error error).
type AssignmentStrategy func(input string, r *RPN) (string, bool, error)

// newAssignmentRegistry creates a new assignment strategy registry.
func newAssignmentRegistry() *assignmentRegistry {
	return &assignmentRegistry{
		strategies: make([]AssignmentStrategy, 0),
	}
}

// register adds an assignment strategy to the registry.
func (r *assignmentRegistry) register(strategy AssignmentStrategy) {
	r.strategies = append(r.strategies, strategy)
}

// parse attempts to parse input using registered strategies in order.
func (r *assignmentRegistry) parse(input string, rn *RPN) (string, bool, error) {
	for _, strategy := range r.strategies {
		if result, handled, err := strategy(input, rn); handled {
			return result, true, err
		}
	}
	return "", false, nil
}

// newAssignmentHandler creates a new assignment handler with all strategies registered.
func newAssignmentHandler() *assignmentHandler {
	h := &assignmentHandler{
		registry: newAssignmentRegistry(),
	}
	h.registry.register(handleAssignRight)
	h.registry.register(handleAssignLeft)
	h.registry.register(handleStandardAssign)
	return h
}

// handle attempts to parse input using registered assignment strategies.
func (h *assignmentHandler) handle(input string, r *RPN) (string, bool, error) {
	return h.registry.parse(input, r)
}

// handleAssignmentOp is the shared implementation for := and =: operators.
// Both share identical logic: check for operator, extract fields, try two
// orderings (value name, then name value), parse, set variable, evaluate remainder.
func handleAssignmentOp(input string, r *RPN, op string) (string, bool, error) {
	if !strings.Contains(input, op) {
		return "", false, nil
	}

	pos := strings.Index(input, op)
	before := strings.TrimSpace(input[:pos])
	after := strings.TrimSpace(input[pos+len(op):])

	beforeFields := strings.Fields(before)
	if len(beforeFields) != 2 {
		return "", false, nil
	}

	// Try value name op format first (stack variant)
	if result, ok, err := tryAssignment(beforeFields[1], beforeFields[0], r, input, after); ok || err != nil {
		return result, ok, err
	}

	// Try name value op format (for backward compatibility)
	return tryAssignment(beforeFields[0], beforeFields[1], r, input, after)
}

// tryAssignment attempts to parse valueStr as a float, set the variable,
// and optionally evaluate remaining tokens.
func tryAssignment(name, valueStr string, r *RPN, input, after string) (string, bool, error) {
	val, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return "", false, nil
	}

	varName := extractVariableName(name)
	if err := r.vars.SetVariable(varName, val); err != nil {
		return "", false, err
	}
	if after == "" {
		return fmt.Sprintf("%s = %.10g", varName, val), true, nil
	}
	result, err := r.evaluate(input, strings.Fields(after))
	return result, true, err
}

// handleAssignRight handles the := operator (right assignment).
// Format: value name := (stack variant) or name value := (direct variant)
func handleAssignRight(input string, r *RPN) (string, bool, error) {
	return handleAssignmentOp(input, r, ":=")
}

// handleAssignLeft handles the =: operator (left assignment).
// Format: value name =: (stack variant) or name value =: (direct variant)
func handleAssignLeft(input string, r *RPN) (string, bool, error) {
	return handleAssignmentOp(input, r, "=:")
}

// handleStandardAssign handles the standard = operator.
// Format: name value = expression (name on bottom, value on top, expression after =)
// Or: name = value (single assignment)
func handleStandardAssign(input string, r *RPN) (string, bool, error) {
	// Tokenize and look for standalone "=" token.
	// This avoids false positives from == or != which are separate tokens.
	tokens := Tokenize(input)
	eqIndex := -1
	for i, tok := range tokens {
		if tok == "=" {
			eqIndex = i
			break
		}
	}
	if eqIndex < 0 {
		return "", false, nil
	}

	// Handle single assignment: "name = value" (3 tokens: name, =, value)
	if eqIndex == 1 && len(tokens) == 3 {
		name := tokens[0]
		valueStr := tokens[2]
		val, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return "", false, nil
		}
		if err := r.vars.SetVariable(name, val); err != nil {
			return "", false, err
		}
		return fmt.Sprintf("%s = %.10g", name, val), true, nil
	}

	// Handle assignment with expression: "name value = expression..."
	beforeTokens := tokens[:eqIndex]
	afterTokens := tokens[eqIndex+1:]

	if len(beforeTokens) == 2 {
		name := beforeTokens[0]
		valueStr := beforeTokens[1]

		val, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return "", false, nil
		}
		if err := r.vars.SetVariable(name, val); err != nil {
			return "", false, err
		}

		if len(afterTokens) == 0 {
			return fmt.Sprintf("%s = %.10g", name, val), true, nil
		}
		result, err := r.evaluate(input, afterTokens)
		return result, true, err
	}

	return "", false, nil
}

// ParseAndEvaluate parses and evaluates an RPN expression.
// Returns the result as a formatted string or an error.
// This method is thread-safe for concurrent execution.
func (r *RPN) ParseAndEvaluate(input string) (string, error) {
	// Validate input and initialize
	input = strings.TrimSpace(input)
	if input == "" {
		return "", fmt.Errorf("rpn: empty expression")
	}

	// Lock for write operations on currentStack
	r.mu.Lock()
	if r.currentStack == nil {
		r.currentStack = NewStack()
	}
	r.mu.Unlock()

	// Handle assignment formats using the new assignment handler
	if assignmentResult, isAssignment, err := r.assignHandler.handle(input, r); err != nil {
		return "", fmt.Errorf("rpn: failed to handle assignment: %w", err)
	} else if isAssignment {
		return assignmentResult, nil
	}

	// Evaluate standard RPN expression
	tokens := Tokenize(input)
	if len(tokens) == 0 {
		return "", fmt.Errorf("rpn: no valid tokens found in input: %q", input)
	}

	return r.evaluate(input, tokens)
}

// evaluate evaluates a list of tokens and returns the result.
// This method is thread-safe for concurrent execution.
func (r *RPN) evaluate(input string, tokens []string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Use the current stack for evaluation to preserve state
	// This allows incremental operations in REPL mode
	if r.currentStack == nil {
		r.currentStack = NewStack()
	}
	stack := r.currentStack

	result, handled, err := r.evaluateTokens(stack, tokens)
	if err != nil {
		return "", err
	}
	if handled {
		return result, nil
	}
	return r.processResult(stack, input)
}

// evaluateTokens processes each token through the dispatch loop.
// Returns (result, handled, error). If handled is true, a command consumed tokens.
func (r *RPN) evaluateTokens(stack *Stack, tokens []string) (string, bool, error) {
	for i, token := range tokens {
		result, handled, err := r.dispatchToken(stack, tokens, i, token)
		if err != nil {
			return "", false, err
		}
		if handled {
			return result, true, nil
		}
	}
	return "", false, nil
}

// dispatchToken handles a single token during evaluation.
// Returns (result, handled, error). If handled is true, evaluation stops.
func (r *RPN) dispatchToken(stack *Stack, tokens []string, i int, token string) (string, bool, error) {
	// Check for variable assignment: name value = (but not == or != etc.)
	if token == "=" && (i+1 >= len(tokens) || tokens[i+1] != "=") {
		return "", false, fmt.Errorf("rpn: invalid assignment syntax at token %d: 'name value =' requires spaces around =", i)
	}

	// Push literal values (numbers, booleans, metric values, @metric prefix)
	if pushed, err := r.pushLiteral(stack, token); err != nil {
		return "", false, err
	} else if pushed {
		return "", false, nil
	}

	// Handle multi-word metric command: metric <subcommand>
	if result, handled, err := r.handleMetricCommand(stack, tokens, i); err != nil {
		return "", false, err
	} else if handled {
		return result, true, nil
	}

	// Handle multi-word custom command: custom <subcommand>
	if result, handled, err := r.handleCustomCommand(stack, tokens, i); err != nil {
		return "", false, err
	} else if handled {
		return result, true, nil
	}

	// Check for inline assignment: name value := (exactly 2 tokens)
	if i+1 < len(tokens) {
		nextToken := tokens[i+1]
		if (nextToken == ":=" || nextToken == "=:") && len(tokens) == 2 && i == 0 {
			val, err := stack.Pop()
			if err != nil {
				return "", false, fmt.Errorf("insufficient operands for %s: stack is empty", nextToken)
			}
			valF, err := toFloat64(val, token)
			if err != nil {
				return "", false, fmt.Errorf("failed to get float64 value for variable %q: %w", token, err)
			}
			if err := r.vars.SetVariable(token, valF); err != nil {
				return "", false, fmt.Errorf("failed to set variable %q: %w", token, err)
			}
			return fmt.Sprintf("%s = %.10g", token, valF), true, nil
		}
	}

	// Handle variable name for assignment (shouldPushName or reassignment)
	if handled := r.checkVariableName(stack, tokens, i, token); handled {
		return "", false, nil
	}

	// Handle special operators and commands
	result, err := r.handleOperator(stack, token, i)
	if err != nil {
		return "", false, fmt.Errorf("rpn: failed to handle operator '%s' at position %d: %w", token, i, err)
	}
	return result, result != "", nil
}

// checkVariableName handles variable name detection for assignment contexts.
// Pushes token as StringNum if it's a variable name for reassignment.
// Returns true if handled.
func (r *RPN) checkVariableName(stack *Stack, tokens []string, i int, token string) bool {
	// shouldPushName: token before := or =:
	if r.shouldPushName(tokens, i) {
		stack.Push(NewStringNum(token))
		return true
	}

	// Variable reassignment: push name instead of value for := or =:
	if isValidIdentifier(token) {
		if _, exists := r.vars.GetVariable(token); exists {
			if i+2 < len(tokens) {
				if tokens[i+2] == ":=" || tokens[i+2] == "=:" {
					stack.Push(NewStringNum(token))
					return true
				}
			}
		}
	}

	return false
}

// processResult checks the final stack state, saves it, and formats the result.
func (r *RPN) processResult(stack *Stack, input string) (string, error) {
	if stack.Len() == 0 {
		if strings.Contains(input, ":=") || strings.Contains(input, "=:") {
			return "", nil
		}
		return "", fmt.Errorf("empty result: expression evaluated to nothing")
	}

	// Save the current stack state for continued operations
	r.currentStack = NewStack()
	for _, val := range stack.Values() {
		r.currentStack.Push(val)
	}

	if stack.Len() > 1 {
		result, err := r.ops.Show(stack)
		if err != nil {
			return "", fmt.Errorf("final result: %w", err)
		}
		return result, nil
	}

	val, _ := stack.Pop()
	return val.String(), nil
}

// handleOperator handles operators and special commands using the operator registry.
func (r *RPN) handleOperator(stack *Stack, token string, tokenIndex int) (string, error) {
	// Check if it's a number first
	if _, err := strconv.ParseFloat(token, 64); err == nil {
		return "", nil
	}

	// Check symbol syntax (:x)
	if pushed, err := r.checkAndPushSymbol(stack, token); err != nil {
		return "", err
	} else if pushed {
		return "", nil
	}

	// Resolve variable or constant
	if r.resolveVariableOrConstant(stack, token) {
		return "", nil
	}

	// Handle operators and fallback to symbol for unknown identifiers
	return r.dispatchOperator(stack, token)
}

// checkAndPushSymbol checks for :x syntax and pushes a Symbol if valid.
// Returns (true, nil) if pushed, (false, nil) if not a symbol prefix, or (false, error).
func (r *RPN) checkAndPushSymbol(stack *Stack, token string) (bool, error) {
	if len(token) > 0 && token[0] == ':' {
		symbolName := token[1:] // Remove the leading :
		if symbolName == "" {
			return false, fmt.Errorf("symbol name cannot be empty after colon")
		}
		// Only push as symbol if valid identifier (prevents := and =: as : then =)
		if isValidIdentifier(symbolName) {
			stack.Push(NewSymbol(symbolName))
			return true, nil
		}
	}
	return false, nil
}

// resolveVariableOrConstant looks up token as variable or constant and pushes value.
func (r *RPN) resolveVariableOrConstant(stack *Stack, token string) bool {
	if val, exists := r.vars.GetVariable(token); exists {
		stack.Push(NewNumber(val, r.ops.GetMode()))
		return true
	}
	if val, exists := r.consts.GetConstant(token); exists {
		stack.Push(NewNumber(val, r.ops.GetMode()))
		return true
	}
	return false
}

// dispatchOperator executes registered operators and falls back to symbol for unknown identifiers.
func (r *RPN) dispatchOperator(stack *Stack, token string) (string, error) {
	result, handled, err := r.executeOperator(stack, token)
	if err != nil {
		// Unknown token: fall back to symbol if valid identifier
		if !r.opRegistry.IsStandardOperator(token) && !r.opRegistry.IsHyperOperator(token) {
			if isValidIdentifier(token) {
				stack.Push(NewSymbol(token))
				return "", nil
			}
		}
		return "", err
	}
	if handled {
		return result, nil
	}

	// Bare identifier — push as symbol
	if isValidIdentifier(token) {
		stack.Push(NewSymbol(token))
	}
	return "", nil
}

// isValidIdentifier checks if a token looks like a valid variable identifier.
// Valid identifiers contain only alphanumeric characters and underscores,
// and start with a letter or underscore (not a digit or special character).
//
// IMPORTANT: To prevent natural language words (like "what", "is", "of") from
// being incorrectly treated as RPN symbols in mixed-mode expressions,
// this function currently restricts valid identifiers to a length of 1.
func isValidIdentifier(token string) bool {
	if len(token) == 0 {
		return false
	}

	// Check first character - must be letter or underscore
	first := token[0]
	if !((first >= 'a' && first <= 'z') || (first >= 'A' && first <= 'Z') || first == '_') {
		return false
	}

	// Check remaining characters - must be alphanumeric or underscore
	for i := 1; i < len(token); i++ {
		c := token[i]
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return false
		}
	}

	// Only allow single-character identifiers for symbol support
	// This prevents words like "what", "is", "of" from becoming symbols
	return len(token) == 1
}

// extractVariableName extracts a variable name from a token, stripping the leading colon if present.
// This allows symbol syntax like :x to be used where the actual variable name is x.
func extractVariableName(token string) string {
	if len(token) > 0 && token[0] == ':' {
		return token[1:]
	}
	return token
}

// checkStackOverflow returns an error if the stack has reached the maximum size.
func (r *RPN) checkStackOverflow(stack *Stack) error {
	if stack.Len() >= r.maxStack {
		return fmt.Errorf("stack overflow")
	}
	return nil
}

// pushLiteral attempts to push a literal value (boolean, number, metric number, @metric) onto the stack.
// Returns (true, nil) if a literal was pushed, (false, nil) if not a literal, or (false, err) on error.
func (r *RPN) pushLiteral(stack *Stack, token string) (bool, error) {
	// Check if it's a boolean literal
	if token == "true" {
		stack.Push(NewFloatFromBool(true))
		return true, nil
	}
	if token == "false" {
		stack.Push(NewFloatFromBool(false))
		return true, nil
	}

	// Check if it's a number
	if num, err := strconv.ParseFloat(token, 64); err == nil {
		if err := r.checkStackOverflow(stack); err != nil {
			return false, err
		}
		if r.ops.GetMode() == RationalMode {
			rat, err := NewRatFromString(token)
			if err != nil {
				return false, err
			}
			stack.Push(rat)
		} else {
			stack.Push(NewFloat(num))
		}
		return true, nil
	}

	// Check if it's a number with a metric suffix (e.g., 100Mbps, 5.5GB, 2hr)
	if num, metric, ok := parseNumberWithMetric(token, r.ops.MetricRegistry()); ok {
		if err := r.checkStackOverflow(stack); err != nil {
			return false, err
		}
		stack.Push(NewNumberWithMetric(num, r.ops.GetMode(), metric))
		return true, nil
	}

	// Check for @ prefix: standalone metric (e.g., @GB, @Mbps)
	// Pushes a Number with value 1 and the looked-up metric
	if len(token) > 1 && token[0] == '@' {
		metricName := token[1:]
		if metric, ok := r.ops.MetricRegistry().FindWithAliases(metricName); ok {
			if err := r.checkStackOverflow(stack); err != nil {
				return false, err
			}
			stack.Push(NewNumberWithMetric(1, r.ops.GetMode(), metric))
			return true, nil
		}
		return false, fmt.Errorf("unknown metric %q in %q", metricName, token)
	}

	return false, nil
}

// handleMetricCommand dispatches 'metric <subcmd>' commands.
// Returns (result, handled, error). If handled is true, result may be non-empty.
func (r *RPN) handleMetricCommand(stack *Stack, tokens []string, i int) (string, bool, error) {
	if token := tokens[i]; token != "metric" || i+1 >= len(tokens) {
		return "", false, nil
	}

	subCmd := tokens[i+1]
	switch subCmd {
	case "show":
		result, err := r.ops.MetricShow(stack)
		if err != nil {
			return "", true, fmt.Errorf("rpn: metric show: %w", err)
		}
		return result, true, nil
	case "list":
		result, err := r.ops.MetricList(stack)
		if err != nil {
			return "", true, fmt.Errorf("rpn: metric list: %w", err)
		}
		return result, true, nil
	case "binary":
		if i+2 < len(tokens) && tokens[i+2] == "set" {
			r.ops.SetPrefixMode(IEC)
			return "prefix mode: IEC", true, nil
		}
		return "", true, fmt.Errorf("rpn: metric binary: use 'metric binary set'")
	case "decimal":
		if i+2 < len(tokens) && tokens[i+2] == "set" {
			r.ops.SetPrefixMode(SI)
			return "prefix mode: SI", true, nil
		}
		return "", true, fmt.Errorf("rpn: metric decimal: use 'metric decimal set'")
	case "compatible":
		result, err := r.ops.MetricCompatible(stack)
		if err != nil {
			return "", true, fmt.Errorf("rpn: metric compatible: %w", err)
		}
		return result, true, nil
	default:
		// Try as a category name
		result, err := r.ops.MetricCategory(stack, subCmd)
		if err != nil {
			return "", true, fmt.Errorf("rpn: metric %s: %w", subCmd, err)
		}
		return result, true, nil
	}
}

// handleCustomCommand dispatches 'custom <subcmd>' commands.
// Returns (result, handled, error). If handled is true, result may be non-empty.
func (r *RPN) handleCustomCommand(stack *Stack, tokens []string, i int) (string, bool, error) {
	if token := tokens[i]; token != "custom" || i+1 >= len(tokens) {
		return "", false, nil
	}

	subCmd := tokens[i+1]
	switch subCmd {
	case "show":
		name := ""
		if i+2 < len(tokens) {
			name = tokens[i+2]
		}
		result, err := r.ops.CustomShow(stack, name)
		if err != nil {
			return "", true, fmt.Errorf("rpn: custom show: %w", err)
		}
		return result, true, nil
	case "list":
		result, err := r.ops.CustomList(stack)
		if err != nil {
			return "", true, fmt.Errorf("rpn: custom list: %w", err)
		}
		return result, true, nil
	case "define":
		if i+4 < len(tokens) {
			name := tokens[i+2]
			factorStr := tokens[i+3]
			category := tokens[i+4]
			factor, err := strconv.ParseFloat(factorStr, 64)
			if err != nil {
				return "", true, fmt.Errorf("rpn: custom define: invalid factor %q", factorStr)
			}
			err = r.ops.CustomDefine(name, factor, category)
			if err != nil {
				return "", true, fmt.Errorf("rpn: custom define: %w", err)
			}
			return fmt.Sprintf("defined custom metric %q (factor: %g, category: %s)", name, factor, category), true, nil
		}
		return "", true, fmt.Errorf("rpn: custom define: usage: custom define <name> <factor> <category>")
	case "undefine":
		if i+2 < len(tokens) {
			name := tokens[i+2]
			err := r.ops.CustomUndefine(name)
			if err != nil {
				return "", true, fmt.Errorf("rpn: custom undefine: %w", err)
			}
			return fmt.Sprintf("removed custom metric %q", name), true, nil
		}
		return "", true, fmt.Errorf("rpn: custom undefine: usage: custom undefine <name>")
	default:
		return "", true, fmt.Errorf("rpn: unknown custom subcommand %q. Use: show, list, define, undefine", subCmd)
	}
}

// shouldPushName determines whether a token should be pushed as a variable name (StringNum)
// rather than evaluated as a value. Returns true if the token is part of an := or =: assignment.
func (r *RPN) shouldPushName(tokens []string, i int) bool {
	token := tokens[i]

	// Check if next token is := or =:
	if i+1 < len(tokens) {
		nextToken := tokens[i+1]
		if nextToken == ":=" || nextToken == "=:" {
			// Stack assignment: exactly 2 tokens, handled inline (return true won't be reached)
			// Non-stack assignment: check if token is a variable name
			if _, err := strconv.ParseFloat(token, 64); err != nil && isValidIdentifier(token) {
				return true
			}
		}
	}

	// Special case: first token in := expression (e.g., "x 5 :=")
	if i == 0 && len(tokens) >= 3 && tokens[len(tokens)-1] == ":=" {
		if _, err := strconv.ParseFloat(token, 64); err != nil && isValidIdentifier(token) {
			return true
		}
	}

	return false
}
