package main

import (
	"strconv"
	"unicode"
)

// Helper function to skip whitespace
func skipWhitespace(input *string) {
	for len(*input) > 0 && (unicode.IsSpace(rune((*input)[0]))) {
		*input = (*input)[1:]
	}
}

// Helper function to parse a number
func parseNumber(input *string) (LispValue, bool) {
	value := ""
	for len(*input) > 0 && (unicode.IsDigit(rune((*input)[0])) || (*input)[0] == '-') {
		value += string((*input)[0])
		*input = (*input)[1:]
	}

	if value == "" {
		return nil, false
	}
	if len(*input) > 0 {
		nextchr := (*input)[0]
		if nextchr == ' ' || nextchr == ')' || nextchr == '(' {

		} else {
			return nil, false
		}
	}
	num, err := strconv.Atoi(value)
	if err != nil {
		return nil, false
	}
	return num, true
}

func parseString(input *string) (LispValue, bool) {
	if len(*input) == 0 || (*input)[0] != '"' {
		return nil, false
	}

	*input = (*input)[1:] // Consume opening quote
	value := ""
	escaped := false

	for len(*input) > 0 {
		char := (*input)[0]
		*input = (*input)[1:]

		if char == '\\' && !escaped {
			escaped = true
			continue
		}

		if char == '"' && !escaped {
			// Closing quote found
			return value, true
		}

		value += string(char)
		escaped = false
	}

	return nil, false // Unclosed string
}

type ParserCondition struct {
	err string
}

func (p ParserCondition) Error() string {
	return "Error parsing lisp code"
}

type NothingParsed2 struct {
}

var NothingParsed = ParserCondition{}

// ParseLisp parses Lisp code directly without explicit tokenization.
func ParseLisp(input *string) LispValue {
	var stack []LispValue

	// Parsing loop
	for {
		skipWhitespace(input)
		if len(*input) == 0 {
			if stack != nil {
				var cond Condition = ParserCondition{err: "Incomplete lisp code parsed."}
				return cond
			}

			return NothingParsed
		}

		switch (*input)[0] {
		case ';':
			for len(*input) > 0 && (*input)[0] != '\n' {
				*input = (*input)[1:]
			}
			if len(*input) > 0 {
				*input = (*input)[1:]
			}
			continue
		case '(':
			// Opening parenthesis: push a new slice to the stack
			stack = append(stack, []LispValue{})
		case ')':
			// Closing parenthesis: pop from stack and append to the parent slice
			popped := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if len(stack) > 0 {
				// Append the popped slice to the parent slice
				stack[len(stack)-1] = append(stack[len(stack)-1].([]LispValue), popped)
			} else {
				*input = (*input)[1:]
				// The stack is empty; this is the root slice
				return popped
			}
		default:
			backup := *input
			if num, ok := parseNumber(input); ok {
				if stack == nil {
					return num
				}
				stack[len(stack)-1] = append(stack[len(stack)-1].([]LispValue), num)
				continue
			}

			*input = backup

			if str, ok := parseString(input); ok {
				if stack == nil {
					return str
				}
				stack[len(stack)-1] = append(stack[len(stack)-1].([]LispValue), str)
				continue
			}

			*input = backup

			// Parse symbol
			value := ""
			for len(*input) > 0 && (*input)[0] != ' ' && (*input)[0] != '\n' && (*input)[0] != '\t' && (*input)[0] != '(' && (*input)[0] != ')' {
				value += string((*input)[0])
				*input = (*input)[1:]
			}

			if value != "" {
				// Append the value to the current slice on the stack
				if stack == nil {
					return lisp.sym(value)
				}
				stack[len(stack)-1] = append(stack[len(stack)-1].([]LispValue), lisp.sym(value))
			}
			continue
		}

		*input = (*input)[1:]
	}
}

// ParseNumber parses a string as a Lisp number.
func ParseNumber(code string) (LispValue, bool) {
	// Try to parse the number
	num, err := strconv.Atoi(code)
	if err != nil {
		return nil, false
	}
	return num, true
}
