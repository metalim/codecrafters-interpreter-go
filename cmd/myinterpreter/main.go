package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/internal/parser"
	"github.com/codecrafters-io/interpreter-starter-go/internal/scanner"
)

func main() {
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {

	case "tokenize":
		filename := os.Args[2]
		fileContents, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		var exitCode int
		scan := scanner.New(string(fileContents))
		go scan.ScanTokens()

		for tok := range scan.Next {
			if tok.Error != nil {
				fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", tok.Line, tok.Error)
				exitCode = 65
				continue
			}

			if len(tok.Literal) == 0 {
				fmt.Printf("%s %s null\n", tok.Type, tok.Lexeme)
			} else {
				fmt.Printf("%s %s %s\n", tok.Type, tok.Lexeme, tok.Literal)
			}
		}
		if exitCode != 0 {
			os.Exit(exitCode)
		}

	case "parse":
		filename := os.Args[2]
		fileContents, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}

		ast, err := parser.Parse(string(fileContents))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing file: %v\n", err)
			os.Exit(1)
		}
		ast.Write(os.Stdout)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

}
