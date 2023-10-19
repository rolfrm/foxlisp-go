package main

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestLisp(t *testing.T) {
	got := []LispValue{1, 2}
	if car(got) != 1 {
		t.Errorf("!!")
	}
}
func TestLisp2(t *testing.T) {
	got := []LispValue{1, 2}
	if cadr(got) != 2 {
		t.Errorf("!!")
	}
}

func TestSymId(t *testing.T) {
	a := lisp.sym("A")
	b := lisp.sym("B")
	c := lisp.sym("C")
	a2 := lisp.sym("A")
	b2 := lisp.sym("B")
	c2 := lisp.sym("C")
	if a != a2 || b != b2 || c != c2 {
		t.Errorf(" :( ")
	}
}
func TestSetGlobal(t *testing.T) {
	a := lisp.sym("A")
	lisp.set_value(a, 123)
	r := lisp.get_value(a)
	if r != 123 {
		t.Errorf(":(")
	}
	b := lisp.sym("B")
	lisp.set_value(b, 321)
	if lisp.get_value(a) != 123 || lisp.get_value(b) != 321 {
		t.Errorf(":(")
	}
}

func TestAdd(t *testing.T) {

	code := []LispValue{plus_sym, 1, []LispValue{plus_sym, 3, 2}}

	result := eval(nil, code)
	if result != 6 {
		t.Errorf("Invalid evaluation of result")
	}

	result2 := eval(nil, []LispValue{plus_sym, 1, 10})
	if result2 != 11 {
		t.Errorf("Invalid evaluation of result")
	}
}

func TestLoop(t *testing.T) {
	b := lisp.sym("B")
	lisp.set_value(b, 25_000_000)
	code := []LispValue{loop_sym, []LispValue{gt_sym, b, 0}, []LispValue{set_sym, b, []LispValue{plus_sym, b, -1}}}
	eval(nil, code)
	b_value := lisp.get_value(b)
	if b_value != 0 {
		t.Errorf("Unexpected value for b. %v", b_value)
	}
}

func TestPrint(t *testing.T) {
	code := []LispValue{println_sym, 1, " ", 2, " ", 3}
	eval(nil, code)
}

func TestSizes(t *testing.T) {
	var x int = 1
	fmt.Printf("%v", unsafe.Sizeof(x))
}
