package parser

import (
	"io"

	"github.com/codecrafters-io/interpreter-starter-go/internal/scanner"
)

func Parse(source string) (*AST, error) {
	ast := &AST{}
	scan := scanner.New(source)
	go scan.ScanTokens()

	for {
		t := scan.NextToken()
		switch t.Type {
		case scanner.EOF:
			return ast, nil

		default:
			scan.PutBack(t)
			node, err := parseExpression(scan)
			if err != nil {
				return nil, err
			}
			ast.nodes = append(ast.nodes, node)
		}
	}
}

func parseGroup(scan *scanner.Scanner) (*Group, error) {
	g := &Group{}
	for {
		t := scan.NextToken()
		switch t.Type {
		case scanner.RIGHT_PAREN:
			if len(g.nodes) == 0 {
				return nil, &Error{message: "empty group", line: t.Line}
			}
			return g, nil

		default:
			scan.PutBack(t)
			node, err := parseExpression(scan)
			if err != nil {
				return nil, err
			}
			g.nodes = append(g.nodes, node)
		}
	}
}

func parseExpression(scan *scanner.Scanner) (Node, error) {
	var nodes []Node
	for {
		t := scan.NextToken()
		switch t.Type {

		case scanner.EOF:
			if len(nodes) != 1 {
				return nil, &Error{message: "unexpected EOF", line: t.Line}
			}
			scan.PutBack(t)
			return nodes[0], nil

		case scanner.NUMBER:
			nodes = append(nodes, &Number{value: t.Literal})

		case scanner.STRING:
			nodes = append(nodes, &String{value: t.Literal})

		case scanner.LEFT_PAREN:
			group, err := parseGroup(scan)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, group)

		case scanner.RIGHT_PAREN:
			if len(nodes) != 1 {
				return nil, &Error{message: "unexpected )", line: t.Line}
			}
			scan.PutBack(t)
			return nodes[0], nil

		case scanner.BANG:
			if len(nodes) != 0 {
				return nil, &Error{message: "unexpected !", line: t.Line}
			}
			node, err := parseExpression(scan)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, &Unary{op: t.Lexeme, node: node})

		case scanner.MINUS:
			if len(nodes) > 1 {
				return nil, &Error{message: "unexpected !", line: t.Line}
			}
			right, err := parseExpression(scan)
			if err != nil {
				return nil, err
			}
			if len(nodes) == 0 {
				un := &Unary{op: t.Lexeme, node: right}
				nodes = append(nodes, un.Reordered())
			} else {
				bin := &Binary{op: t.Lexeme, left: nodes[0], right: right}
				nodes[0] = bin.Reordered()
			}

		case scanner.PLUS, scanner.STAR, scanner.SLASH:
			if len(nodes) != 1 {
				return nil, &Error{message: "unexpected " + t.Lexeme, line: t.Line}
			}
			right, err := parseExpression(scan)
			if err != nil {
				return nil, err
			}
			bin := &Binary{op: t.Lexeme, left: nodes[0], right: right}
			nodes[0] = bin.Reordered()

		default:
			nodes = append(nodes, &Keyword{value: t.Lexeme})
		}
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

var precedence = map[string]int{
	"+": 1,
	"-": 1, // binary
	"*": 2,
	"/": 2,
	// "!": 3, // unary
	// "-": 3, // unary
}

func (b *Binary) Reordered() Node {
	if inner, ok := b.right.(*Binary); ok {
		// If the current binary operator has a higher or same precedence
		// than the right binary operator, then we need to move the right binary operator up.
		// 61 * 98 / 80
		// parsed as 61 * (98 / 80)
		// but should be parsed as (61 * 98) / 80
		if precedence[b.op] >= precedence[inner.op] {
			b.op, inner.op = inner.op, b.op
			b.left, b.right, inner.left, inner.right =
				inner, inner.right, b.left, inner.left
			b.left = inner.Reordered()
		}
	}
	return b
}

func (u *Unary) Reordered() Node {
	if inner, ok := u.node.(*Binary); ok {
		// If the current unary operator has a higher or same precedence
		// than the inner binary operator, then we need to move the inner binary operator up.
		// -1 * 2
		// parsed as -(1 * 2)
		// but should be parsed as (-1) * 2
		// HINT: unary operators have higher precedence than binary operators
		u.node, inner.left = inner.left, u
		return inner
	}
	return u
}
