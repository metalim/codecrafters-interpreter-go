package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	line := 1
	var exitCode int
	var state State

	for _, b := range fileContents {
		switch state {
		case EQUAL:
			if b == '=' {
				fmt.Println("EQUAL_EQUAL == null")
				state = DEFAULT
				continue
			}
			fmt.Println("EQUAL = null")
			state = DEFAULT
		}
		switch b {
		case '(':
			fmt.Println("LEFT_PAREN ( null")
		case ')':
			fmt.Println("RIGHT_PAREN ) null")
		case '{':
			fmt.Println("LEFT_BRACE { null")
		case '}':
			fmt.Println("RIGHT_BRACE } null")
		case ',':
			fmt.Println("COMMA , null")
		case '.':
			fmt.Println("DOT . null")
		case '-':
			fmt.Println("MINUS - null")
		case '+':
			fmt.Println("PLUS + null")
		case ';':
			fmt.Println("SEMICOLON ; null")
		case '*':
			fmt.Println("STAR * null")
		case '!':
			fmt.Println("BANG ! null")
		case '=':
			state = EQUAL
		case '<':
			fmt.Println("LESS < null")
		case '>':
			fmt.Println("GREATER > null")
		case '/':
			fmt.Println("SLASH / null")

		default:
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character: %s\n", line, string(b))
			exitCode = 65
		}
	}
	switch state {
	case EQUAL:
		fmt.Println("EQUAL = null")
		state = DEFAULT
	}

	fmt.Println("EOF  null")
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

type State int

const (
	DEFAULT State = iota
	EQUAL
)
