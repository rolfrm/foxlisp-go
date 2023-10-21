//go:build !wasm
// +build !wasm

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		os.Exit(1)
	}

	// Iterate through command-line arguments
	for _, filePath := range os.Args[1:] {
		// Verify file existence
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			fmt.Printf("File not found: %s\n", filePath)
			continue
		}

		// Read the content of the file
		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			continue
		}

		// Convert file content to string
		lispCode := string(fileContent)

		// Parse Lisp code
		for len(lispCode) > 0 {
			skipWhitespace(&lispCode)
			if len(lispCode) == 0 {
				break
			}

			ast := ParseLisp(&lispCode)

			if cond, ok := ast.(Condition); ok {
				if ast == NothingParsed {

					continue
				}
				fmt.Printf("Error parsing Lisp code in %s: %s\n", filePath, cond.Error())
				return
			}

			result := eval(nil, ast)
			if cond, ok := result.(Condition); ok {
				fmt.Printf("Error executing lisp code: %v", cond.Error())
			}

		}

	}
}
