package cssparser

import (
	"bytes"
)

// Position holds a line/column cursor.
type Position struct {
	Line, Column int
}

// PositionRange marks start and end of a node.
type PositionRange struct {
	Start, End Position
}

// Node represents `property: value;`.
type Node struct {
	Property []byte
	Value    []byte
	Position PositionRange
}

func (d *Node) Pos() PositionRange { return d.Position }

// Parser holds the input and a cursor.
type Parser struct {
	input       []byte
	length      int
	pos         int
	lineno, col int
}

// NewParser constructs a fresh parser.
func NewParser(input []byte) *Parser {
	return &Parser{
		input:  input,
		length: len(input),
		lineno: 1,
		col:    1,
	}
}

// peek returns the next byte or 0.
func (p *Parser) peek() byte {
	if p.pos < p.length {
		return p.input[p.pos]
	}
	return 0
}

// consume returns the next byte and advances the cursor,
// updating line/column on '\n'.
func (p *Parser) consume() byte {
	if p.pos >= p.length {
		return 0
	}
	ch := p.input[p.pos]
	p.pos++
	if ch == '\n' {
		p.lineno++
		p.col = 1
	} else {
		p.col++
	}
	return ch
}

// skipWhitespace moves past spaces, tabs, CRs and LFs.
func (p *Parser) skipWhitespace() {
	for {
		switch p.peek() {
		case ' ', '\t', '\r', '\n':
			p.consume()
		default:
			return
		}
	}
}

// parseDeclaration attempts to read a property:value[;].
// Returns nil if it doesn't see a valid prop or colon.
func (p *Parser) parseDeclaration() *Node {
	start := Position{p.lineno, p.col}
	p.skipWhitespace()

	// property prefix '*'?
	if p.peek() == '*' {
		p.consume()
	}
	propStart := p.pos
	for {
		ch := p.peek()
		if ch == 0 || isWS(ch) || ch == ':' {
			break
		}
		// allow [attr] once
		if ch == '[' {
			p.consume()
			for p.peek() != ']' && p.peek() != 0 {
				p.consume()
			}
			if p.peek() == ']' {
				p.consume()
			}
			continue
		}
		p.consume()
	}
	if p.peek() != ':' {
		return nil
	}
	prop := bytes.TrimSpace(p.input[propStart:p.pos])
	// eat ':'
	p.consume()
	p.skipWhitespace()

	// value
	valStart := p.pos
	for {
		ch := p.peek()
		if ch == 0 || ch == ';' || ch == '}' {
			break
		}
		// skip over quoted strings
		if ch == '"' || ch == '\'' {
			quote := ch
			p.consume()
			for p.peek() != quote && p.peek() != 0 {
				if p.peek() == '\\' {
					p.consume()
					p.consume()
				} else {
					p.consume()
				}
			}
			if p.peek() == quote {
				p.consume()
			}
			continue
		}
		// skip over parentheses
		if ch == '(' {
			p.consume()
			for p.peek() != ')' && p.peek() != 0 {
				p.consume()
			}
			if p.peek() == ')' {
				p.consume()
			}
			continue
		}
		p.consume()
	}
	rawVal := bytes.TrimSpace(p.input[valStart:p.pos])
	p.skipWhitespace()
	if p.peek() == ';' {
		p.consume()
	}
	end := Position{p.lineno, p.col}

	return &Node{
		Property: prop,
		Value:    rawVal,
		Position: PositionRange{Start: start, End: end},
	}
}

// Parse walks the entire stylesheet and returns a list of declarations.
func Parse(style []byte) ([]*Node, error) {
	p := NewParser(style)
	var out []*Node

	for p.pos < p.length {
		p.skipWhitespace()
		if p.pos >= p.length {
			break
		}
		if d := p.parseDeclaration(); d != nil {
			out = append(out, d)
			continue
		}
		// unrecognized character: consume to avoid infinite loop
		p.consume()
	}
	return out, nil
}

// isWS reports whether b is ASCII whitespace.
func isWS(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}
