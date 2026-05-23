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

// standardAssignHandler handles the standard = operator.
// Format: name value = expression (value on bottom, expression after =)
// Or: name = value (single assignment)
func handleStandardAssign(input string, r *RPN) (string, bool, error) {
	// Check for standard assignment format (name = value or name value = expression)
	hasAssignment := strings.Contains(input, " = ")
	if !hasAssignment {
		// Check for " =" (space before equals) without space after
		hasAssignment = strings.Contains(input, " =")
		// Additional check: the = must not be followed by another = (i.e., not == or !=)
		if hasAssignment && strings.Contains(input, "==") {
			hasAssignment = false
		}
		if hasAssignment && strings.Contains(input, "!=") {
			hasAssignment = false
		}
	}

	if !hasAssignment {
		return "", false, nil
	}

	// Handle single assignment: "name = value"
	if parts := strings.SplitN(input, " = ", 2); len(parts) == 2 {
		name := strings.TrimSpace(parts[0])
		valueStr := strings.TrimSpace(parts[1])

		// Validate name is a single word (variable name)
		nameFields := strings.Fields(name)
		if len(nameFields) == 1 {
			// Validate value is a single number
			valueFields := strings.Fields(valueStr)
			if len(valueFields) == 1 {
				val, err := strconv.ParseFloat(valueFields[0], 64)
				if err != nil {
					return "", false, fmt.Errorf("invalid value '%s' for assignment: %w", valueFields[0], err)
				}
				if err := r.vars.SetVariable(nameFields[0], val); err != nil {
					return "", false, err
				}
				return fmt.Sprintf("%s = %.10g", nameFields[0], val), true, nil
			}
		}
	}

	// Handle assignment with expression: "name value = expression..."
	pos := strings.Index(input, " =")
	if pos >= 0 {
		before := strings.TrimSpace(input[:pos])
		after := strings.TrimSpace(input[pos+2:])

		beforeFields := strings.Fields(before)
		if len(beforeFields) == 2 {
			name := beforeFields[0]
			valueStr := beforeFields[1]

			// Try to parse value as a number
			val, err := strconv.ParseFloat(valueStr, 64)
			if err == nil {
				// Valid assignment pattern: "name value = expr..." or "name value ="
				if err := r.vars.SetVariable(name, val); err != nil {
					return "", false, err
				}

				// If no expression after assignment, just return assignment info
				if after == "" {
					return fmt.Sprintf("%s = %.10g", name, val), true, nil
				}
				result, err := r.evaluate(input, strings.Fields(after))
				return result, true, err
			}
		}
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

	for i, token := range tokens {
		// Check for variable assignment: name value = (but not == or != etc.)
		if token == "=" && (i+1 >= len(tokens) || tokens[i+1] != "=") {
			return "", fmt.Errorf("rpn: invalid assignment syntax at token %d: 'name value =' requires spaces around =", i)
		}

		// Check if it's a boolean literal
		if token == "true" {
			stack.Push(NewFloatFromBool(true))
			continue
		}
		if token == "false" {
			stack.Push(NewFloatFromBool(false))
			continue
		}

		// Check if it's a number
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			if stack.Len() >= r.maxStack {
				return "", fmt.Errorf("stack overflow")
			}
			stack.Push(NewNumber(num, r.mode))
			continue
		}

		// Check if it's a number with a metric suffix (e.g., 100Mbps, 5.5GB, 2hr)
		if num, metric, ok := parseNumberWithMetric(token, r.metricRegistry); ok {
			if stack.Len() >= r.maxStack {
				return "", fmt.Errorf("stack overflow")
			}
			stack.Push(NewNumberWithMetric(num, r.mode, metric))
			continue
		}

		// Check for @ prefix: standalone metric (e.g., @GB, @Mbps)
		// Pushes a Number with value 1 and the looked-up metric
		if len(token) > 1 && token[0] == '@' {
			metricName := token[1:]
			if metric, ok := r.metricRegistry.FindWithAliases(metricName); ok {
				if stack.Len() >= r.maxStack {
					return "", fmt.Errorf("stack overflow")
				}
				stack.Push(NewNumberWithMetric(1, r.mode, metric))
				continue
			}
			return "", fmt.Errorf("unknown metric %q in %q", metricName, token)
		}

		// Handle multi-word metric command: metric <subcommand>
		if token == "metric" && i+1 < len(tokens) {
			subCmd := tokens[i+1]
			switch subCmd {
			case "show":
				result, err := r.ops.MetricShow(stack)
				if err != nil {
					return "", fmt.Errorf("rpn: metric show: %w", err)
				}
				return result, nil
			case "list":
				result, err := r.ops.MetricList(stack)
				if err != nil {
					return "", fmt.Errorf("rpn: metric list: %w", err)
				}
				return result, nil
			case "binary":
				if i+2 < len(tokens) && tokens[i+2] == "set" {
					r.ops.SetPrefixMode(IEC)
					return "prefix mode: IEC", nil
				}
				return "", fmt.Errorf("rpn: metric binary: use 'metric binary set'")
			case "decimal":
				if i+2 < len(tokens) && tokens[i+2] == "set" {
					r.ops.SetPrefixMode(SI)
					return "prefix mode: SI", nil
				}
				return "", fmt.Errorf("rpn: metric decimal: use 'metric decimal set'")
			case "compatible":
				result, err := r.ops.MetricCompatible(stack)
				if err != nil {
					return "", fmt.Errorf("rpn: metric compatible: %w", err)
				}
				return result, nil
			default:
				// Try as a category name
				result, err := r.ops.MetricCategory(stack, subCmd)
				if err != nil {
					return "", fmt.Errorf("rpn: metric %s: %w", subCmd, err)
				}
				return result, nil
			}
		}

		// Handle multi-word custom command: custom <subcommand>
		if token == "custom" && i+1 < len(tokens) {
			subCmd := tokens[i+1]
			switch subCmd {
			case "show":
				result, err := r.ops.MetricShow(stack)
				if err != nil {
					return "", fmt.Errorf("rpn: custom show: %w", err)
				}
				return result, nil
			case "list":
				result, err := r.ops.CustomList(stack)
				if err != nil {
					return "", fmt.Errorf("rpn: custom list: %w", err)
				}
				return result, nil
			case "define":
				if i+4 < len(tokens) {
					name := tokens[i+2]
					factorStr := tokens[i+3]
					category := tokens[i+4]
					factor, err := strconv.ParseFloat(factorStr, 64)
					if err != nil {
						return "", fmt.Errorf("rpn: custom define: invalid factor %q", factorStr)
					}
					err = r.ops.CustomDefine(name, factor, category)
					if err != nil {
						return "", fmt.Errorf("rpn: custom define: %w", err)
					}
					return fmt.Sprintf("defined custom metric %q (factor: %g, category: %s)", name, factor, category), nil
				}
				return "", fmt.Errorf("rpn: custom define: usage: custom define <name> <factor> <category>")
			case "undefine":
				if i+2 < len(tokens) {
					name := tokens[i+2]
					err := r.ops.CustomUndefine(name)
					if err != nil {
						return "", fmt.Errorf("rpn: custom undefine: %w", err)
					}
					return fmt.Sprintf("removed custom metric %q", name), nil
				}
				return "", fmt.Errorf("rpn: custom undefine: usage: custom undefine <name>")
			default:
				return "", fmt.Errorf("rpn: unknown custom subcommand %q. Use: show, list, define, undefine", subCmd)
			}
		}

		// Check if this is a variable name for assignment (:= or =:)
		// For := (right assignment): name value := - first token is always a variable name
		// For =: (left assignment): value name =: - token before =: is a variable name
		shouldPushName := false

		if i+1 < len(tokens) {
			nextToken := tokens[i+1]
			if nextToken == ":=" || nextToken == "=:" {
				// Check if this is a stack assignment (e.g., "x =:" or "x :=")
				// Stack assignment: exactly 2 tokens, first is variable name, second is operator
				if len(tokens) == 2 && i == 0 {
					// This is a stack assignment. Pop the value from stack and assign to variable.
					// Don't push the name as StringNum because the operator expects stack: [value, name]
					// but for stack assignment, the value is already on stack and we just have the name token.
					// Instead, we handle it inline: pop value, assign to name (from token).
					val, err := stack.Pop()
					if err != nil {
						return "", fmt.Errorf("insufficient operands for %s: stack is empty", nextToken)
					}
					valF, err := toFloat64(val, token)
					if err != nil {
						return "", fmt.Errorf("failed to get float64 value for variable %q: %w", token, err)
					}
					if err := r.vars.SetVariable(token, valF); err != nil {
						return "", fmt.Errorf("failed to set variable %q: %w", token, err)
					}
					// Skip the operator token (next one) since we handled it inline
					// We've consumed both tokens, so we're done
					// Return confirmation message showing the assignment
					return fmt.Sprintf("%s = %.10g", token, valF), nil
				} else if _, err := strconv.ParseFloat(token, 64); err != nil && isValidIdentifier(token) {
					// This token is a variable name (not a number)
					shouldPushName = true
				}
			}
		}

		// Special case: first token in := expression (e.g., "x 5 :=")
		// Only push as name if the first token is not a number (it's a variable name)
		if i == 0 && len(tokens) >= 3 && tokens[len(tokens)-1] == ":=" {
			if _, err := strconv.ParseFloat(token, 64); err != nil && isValidIdentifier(token) {
				shouldPushName = true
			}
		}

		if shouldPushName {
			// This token is a variable name, push as StringNum
			stack.Push(NewStringNum(token))
			continue
		}

		// Special case: if token is a defined variable and appears before an assignment operator
		// (within the next few tokens), push the variable NAME (StringNum) instead of VALUE
		// to allow reassignment.
		// For example: "x 5 := x 10 := ..." - the second "x" should be the name, not the value 5.
		// We check if there's an assignment operator within the next 2 tokens (e.g., "x N :=" or "x N =:")
		if isValidIdentifier(token) {
			if _, exists := r.vars.GetVariable(token); exists {
				// Check if there's an assignment operator within the next 2 tokens
				// Format: variable value := or variable value =:
				if i+2 < len(tokens) {
					if tokens[i+2] == ":=" || tokens[i+2] == "=:" {
						// Push the variable name (not value) for assignment
						stack.Push(NewStringNum(token))
						continue
					}
				}
			}
		}

		// Handle special operators and commands
		if result, err := r.handleOperator(stack, token, i); err != nil {
			return "", fmt.Errorf("rpn: failed to handle operator '%s' at position %d: %w", token, i, err)
		} else if result != "" {
			return result, nil
		}
	}

	// Check final stack state
	if stack.Len() == 0 {
		// Empty stack might be valid for assignment operators (:= or =:)
		// Check if the input was an assignment expression
		if strings.Contains(input, ":=") || strings.Contains(input, "=:") {
			// Assignment expression - empty stack is valid (side effect is variable assignment)
			return "", nil
		}
		return "", fmt.Errorf("empty result: expression evaluated to nothing")
	}

	// Save the current stack state for continued operations
	// Create a copy of the stack to preserve it
	r.currentStack = NewStack()
	for _, val := range stack.Values() {
		r.currentStack.Push(val)
	}

	// Get the final result
	if stack.Len() > 1 {
		// Multiple values on stack - show them all
		result, err := r.ops.Show(stack)
		if err != nil {
			return "", fmt.Errorf("final result: %w", err)
		}
		return result, nil
	}

	// Single value - return it
	val, _ := stack.Pop()
	return val.String(), nil
}

// handleOperator handles operators and special commands using the operator registry.
func (r *RPN) handleOperator(stack *Stack, token string, tokenIndex int) (string, error) {
	// Check if it's a number first
	if _, err := strconv.ParseFloat(token, 64); err == nil {
		return "", nil
	}

	// Check if it's a symbol syntax (:x)
	// Only match :x where x is a valid identifier (not an operator like := or =:)
	if len(token) > 0 && token[0] == ':' {
		symbolName := token[1:] // Remove the leading :
		if symbolName == "" {
			return "", fmt.Errorf("symbol name cannot be empty after colon")
		}
		// Only push as symbol if the remaining part is a valid identifier
		// This prevents := and =: from being treated as : followed by = operator
		if isValidIdentifier(symbolName) {
			stack.Push(NewSymbol(symbolName))
			return "", nil
		}
		// Not a valid symbol, fall through to check for operators
	}

	// Check if it's a variable reference first (before operators)
	if val, exists := r.vars.GetVariable(token); exists {
		stack.Push(NewNumber(val, r.mode))
		return "", nil
	}
	// Check if it's a constant reference (before operators)
	if val, exists := r.consts.GetConstant(token); exists {
		stack.Push(NewNumber(val, r.mode))
		return "", nil
	}

	// Handle standard operators (common logic extracted for DRY)
	// This must be done BEFORE pushing Symbol for unknown identifiers,
	// so that operators are properly handled
	result, handled, err := r.executeOperator(stack, token)
	if err != nil {
		// If it's an unknown token error and we're at the evaluate stage,
		// it might be a bare identifier that should be a symbol
		// Check if the caller is the main evaluate loop
		if strings.Contains(err.Error(), "unknown token") {
			// For bare identifiers, push a Symbol instead of returning error
			// But only if it looks like a valid identifier (alphanumeric/underscore, starts with letter/_)
			// Don't push symbols for tokens with special characters like %, ., etc.
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

	// For bare identifiers that don't exist as variables and aren't operators,
	// push a Symbol (this implements the feature where unbound identifiers act as symbols)
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
