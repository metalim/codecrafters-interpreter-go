package parser

import (
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

		default:
			node, err := parseValue(t, scan)
			if err != nil {
				return nil, err
			}
			ast.nodes = append(ast.nodes, node)
		}
	}
	return ast, nil
}

func parseGroup(scan *scanner.Scanner) (*Group, error) {
	g := &Group{}
	for t := range scan.Next {
		switch t.Type {
		case scanner.RIGHT_PAREN:
			if len(g.nodes) == 0 {
				return nil, &Error{message: "empty group", line: t.Line}
			}
			return g, nil

		default:
			node, err := parseValue(t, scan)
			if err != nil {
				return nil, err
			}
			g.nodes = append(g.nodes, node)
		}
	}
	return nil, &Error{message: "unmatched '('"}
}

func parseValue(t *scanner.Token, scan *scanner.Scanner) (Node, error) {
	switch t.Type {

	case scanner.EOF:
		return nil, &Error{message: "unexpected EOF", line: t.Line}

	case scanner.NUMBER:
		return &Number{value: t.Literal}, nil

	case scanner.STRING:
		return &String{value: t.Literal}, nil

	case scanner.LEFT_PAREN:
		group, err := parseGroup(scan)
		if err != nil {
			return nil, err
		}
		return group, nil

	case scanner.RIGHT_PAREN:
		return nil, &Error{message: "unexpected ')'", line: t.Line}

	case scanner.BANG, scanner.MINUS:
		node, err := parseValue(<-scan.Next, scan)
		if err != nil {
			return nil, err
		}
		return &Unary{op: t.Lexeme, node: node}, nil

	default:
		return &Keyword{value: t.Lexeme}, nil
	}
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

type Error struct {
	message string
	line    int
}

func (e *Error) Error() string {
	return e.message
}

func (e *Error) LineNumber() int {
	return e.line
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

type Group struct {
	nodes []Node
}

func (g *Group) Write(w io.Writer) {
	io.WriteString(w, "(group")
	for _, n := range g.nodes {
		io.WriteString(w, " ")
		n.Write(w)
	}
	io.WriteString(w, ")")
}

type Unary struct {
	op   string
	node Node
}

func (u *Unary) Write(w io.Writer) {
	io.WriteString(w, "(")
	io.WriteString(w, u.op)
	io.WriteString(w, " ")
	u.node.Write(w)
	io.WriteString(w, ")")
}

type Binary struct {
	op    string
	left  Node
	right Node
}

func (b *Binary) Write(w io.Writer) {
	io.WriteString(w, "(")
	io.WriteString(w, b.op)
	io.WriteString(w, " ")
	b.left.Write(w)
	io.WriteString(w, " ")
	b.right.Write(w)
	io.WriteString(w, ")")
}
