//go:build !wasm
// +build !wasm

package main

import "fmt"

func main() {
	code := ParseLisp("(+ 1 2)")
	result := eval(nil, code)
	fmt.Printf("%v\n", result)
}
