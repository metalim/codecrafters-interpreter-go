package parser

import (
	"errors"
	"io"

	"github.com/codecrafters-io/interpreter-starter-go/internal/scanner"
)

func Parse(source string) (*AST, error) {
	ast := &AST{}
	scan := scanner.New(source)
	go scan.ScanTokens()

	for t := range scan.Next {
		switch t.Type {
		case scanner.EOF:
			return ast, nil

		case scanner.NUMBER:
			ast.nodes = append(ast.nodes, &Number{value: t.Literal})

		case scanner.STRING:
			ast.nodes = append(ast.nodes, &String{value: t.Literal})

		default:
			ast.nodes = append(ast.nodes, &Keyword{value: t.Lexeme})
		}
	}
	return nil, errors.New("unreachable")
}

type AST struct {
	nodes []Node
}

func (a *AST) Write(w io.Writer) {
	for _, n := range a.nodes {
		n.Write(w)
		io.WriteString(w, "\n")
	}
}

type Node interface {
	Write(io.Writer)
}

type Keyword struct {
	value string
}

func (k *Keyword) Write(w io.Writer) {
	io.WriteString(w, k.value)
}

type Number struct {
	value string
}

func (n *Number) Write(w io.Writer) {
	io.WriteString(w, n.value)
}

type String struct {
	value string
}

func (s *String) Write(w io.Writer) {
	io.WriteString(w, s.value)
}

type Binary struct {
	left  Node
	right Node
	op    string
}

func (b *Binary) Write(w io.Writer) {
	io.WriteString(w, "(")
	io.WriteString(w, b.op)
	b.left.Write(w)
	io.WriteString(w, " ")
	b.right.Write(w)
	io.WriteString(w, ")")
}
