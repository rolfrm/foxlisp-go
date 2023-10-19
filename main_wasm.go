//go:build wasm
// +build wasm

package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	jsDoc := js.Global().Get("document")
	fmt.Printf("123 %a", jsDoc)

}
