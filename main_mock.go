//go:build !wasm
// +build !wasm

package main

import "fmt"

func main() {
	fmt.Println("Hello, world.")
	//jsDoc := js.Global().Get("document")
	//fmt.Printf("123 %a", jsDoc)
	//co := js.Global().Get("console")
	//co.Call("log", "123333")

	fmt.Println("Finished!")
}
