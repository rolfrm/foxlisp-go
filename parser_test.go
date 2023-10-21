package main

import (
	"fmt"
	"testing"
)

func TestParseLisp(t *testing.T) {
	testCases := []struct {
		input     string
		expected  LispValue
		expectErr bool
	}{
		{"(+ 2 3)", []LispValue{lisp.sym("+"), 2, 3}, false},
		{"(+ 2 (* 3 4) -5)", []LispValue{lisp.sym("+"), 2, []LispValue{lisp.sym("*"), 3, 4}, -5}, false},
		{"(+ 2(* 3 4)-5)", []LispValue{lisp.sym("+"), 2, []LispValue{lisp.sym("*"), 3, 4}, -5}, false},
		{"(invalid)", []LispValue{lisp.sym("invalid")}, false},
		{"(invalid", nil, true},
		{"123abc", lisp.sym("123abc"), false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Parsing: %s", tc.input), func(t *testing.T) {
			result := ParseLisp(tc.input)

			if tc.expectErr {
				if result != nil {
					t.Errorf("Expected an error, but got result: %+v", result)
				}
			} else {
				if !equal(result, tc.expected) {
					t.Errorf("Result does not match expected. Expected: %+v, Got: %+v", tc.expected, result)
				}
			}
		})
	}
}

func TestParseNumber(t *testing.T) {
	testCases := []struct {
		input     string
		expected  LispValue
		expectErr bool
	}{
		{"123", 123, false},
		{"-456", -456, false},
		{"invalid", nil, true}, // Invalid number
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Parsing number: %s", tc.input), func(t *testing.T) {
			result, ok := ParseNumber(tc.input)

			if tc.expectErr {
				if ok {
					t.Errorf("Expected an error, but got result: %+v", result)
				}
			} else {
				if result != tc.expected {
					t.Errorf("Result does not match expected. Expected: %+v, Got: %+v", tc.expected, result)
				}
			}
		})
	}
}

// equal checks if two LispValues are equal.
func equal(a, b LispValue) bool {
	switch aVal := a.(type) {
	case int:
		bVal, ok := b.(int)
		return ok && aVal == bVal
	case string:
		bVal, ok := b.(string)
		return ok && aVal == bVal
	case float32:
		bVal, ok := b.(float32)
		return ok && aVal == bVal
	case float64:
		bVal, ok := b.(float64)
		return ok && aVal == bVal
	case Symbol:
		bVal, ok := b.(Symbol)
		return ok && bVal.id == aVal.id
	case []LispValue:
		bVal, ok := b.([]LispValue)
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for i := range aVal {
			if !equal(aVal[i], bVal[i]) {
				return false
			}
		}
		return true
	default:
		return false
	}
}
