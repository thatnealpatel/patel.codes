package parser

import (
	"fmt"
	"unicode"
)

type Node interface{ node() }

type List []Node

func (List) node() {}

type Letter string

func (Letter) node() {}

type Number string

func (Number) node() {}

type Operator string

func (Operator) node() {}

type Space struct{}

func (Space) node() {}

type Command struct {
	Name    string
	Args    [][]Node
	OptArgs [][]Node
}

func (Command) node() {}

type Sup struct {
	Base   Node
	Script Node
}

func (Sup) node() {}

type Sub struct {
	Base   Node
	Script Node
}

func (Sub) node() {}

var argCount = map[string]int{
	`\frac`: 2, `\dfrac`: 2, `\tfrac`: 2, `\binom`: 2,
	`\sqrt`: 1, `\overline`: 1, `\underline`: 1, `\hat`: 1,
	`\bar`: 1, `\vec`: 1, `\dot`: 1, `\ddot`: 1, `\tilde`: 1,
	`\text`: 1, `\textit`: 1, `\textbf`: 1, `\textmd`: 1,
	`\textrm`: 1, `\mathrm`: 1, `\mathbf`: 1, `\mathit`: 1,
	`\mathcal`: 1, `\mathbb`: 1, `\mathfrak`: 1,
	`\mod`: 1, `\pmod`: 1, `\bmod`: 1,
	`\eqref`: 1, `\label`: 1, `\tag`: 1,
}

var hasOptArg = map[string]bool{
	`\sqrt`: true,
}

type parser struct {
	input []rune
	pos   int
}

func Parse(expr string) ([]Node, error) {
	p := &parser{input: []rune(expr)}
	nodes, err := p.parseExpr()
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %w", expr, err)
	}
	return nodes, nil
}

func (p *parser) parseExpr() ([]Node, error) {
	return p.parseUntil('}')
}

func (p *parser) parseUntil(stop rune) ([]Node, error) {
	var nodes []Node
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == stop {
			break
		}

		node, err := p.parseItem()
		if err != nil {
			return nil, err
		}
		if _, ok := node.(Space); ok {
			nodes = append(nodes, node)
			continue
		}

		saved := p.pos
		p.skipSpaces()
		if p.pos < len(p.input) && p.input[p.pos] == '^' {
			p.pos++
			script, err := p.parseAtom()
			if err != nil {
				return nil, err
			}
			node = Sup{Base: node, Script: script}
		} else if p.pos < len(p.input) && p.input[p.pos] == '_' {
			p.pos++
			script, err := p.parseAtom()
			if err != nil {
				return nil, err
			}
			node = Sub{Base: node, Script: script}
		} else {
			p.pos = saved
		}

		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (p *parser) parseItem() (Node, error) {
	ch := p.input[p.pos]

	switch {
	case ch == '{':
		return p.parseGroup()
	case ch == '\\':
		return p.parseCommand()
	case unicode.IsLetter(ch):
		p.pos++
		return Letter(ch), nil
	case unicode.IsDigit(ch):
		return p.parseNumber(), nil
	case ch == ' ' || ch == '\t' || ch == '\n':
		p.pos++
		return Space{}, nil
	default:
		p.pos++
		return Operator(ch), nil
	}
}

func (p *parser) parseAtom() (Node, error) {
	if p.pos >= len(p.input) {
		return nil, fmt.Errorf("unexpected end of input after ^/_")
	}
	ch := p.input[p.pos]
	if ch == '{' {
		return p.parseGroup()
	}
	if ch == '\\' {
		return p.parseCommand()
	}
	p.pos++
	if unicode.IsDigit(ch) {
		return Number(string(ch)), nil
	}
	if unicode.IsLetter(ch) {
		return Letter(ch), nil
	}
	return Operator(ch), nil
}

func (p *parser) parseGroup() (Node, error) {
	p.pos++ // skip '{'
	nodes, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.pos < len(p.input) && p.input[p.pos] == '}' {
		p.pos++
	}
	if len(nodes) == 1 {
		return nodes[0], nil
	}
	return List(nodes), nil
}

func (p *parser) parseCommand() (Node, error) {
	p.pos++ // skip '\'
	if p.pos >= len(p.input) {
		return Operator(`\`), nil
	}

	if !unicode.IsLetter(p.input[p.pos]) {
		ch := p.input[p.pos]
		p.pos++
		return Operator(string([]rune{'\\', ch})), nil
	}

	start := p.pos
	for p.pos < len(p.input) && unicode.IsLetter(p.input[p.pos]) {
		p.pos++
	}
	name := `\` + string(p.input[start:p.pos])

	nargs, known := argCount[name]
	if !known {
		return Command{Name: name}, nil
	}

	cmd := Command{Name: name}

	if hasOptArg[name] && p.pos < len(p.input) && p.input[p.pos] == '[' {
		p.pos++ // skip '['
		optNodes, err := p.parseUntil(']')
		if err != nil {
			return nil, err
		}
		if p.pos < len(p.input) && p.input[p.pos] == ']' {
			p.pos++
		}
		cmd.OptArgs = append(cmd.OptArgs, optNodes)
	}

	for i := 0; i < nargs; i++ {
		p.skipSpaces()
		if p.pos >= len(p.input) || p.input[p.pos] != '{' {
			break
		}
		p.pos++ // skip '{'
		argNodes, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.pos < len(p.input) && p.input[p.pos] == '}' {
			p.pos++
		}
		cmd.Args = append(cmd.Args, argNodes)
	}

	return cmd, nil
}

func (p *parser) parseNumber() Node {
	start := p.pos
	for p.pos < len(p.input) && unicode.IsDigit(p.input[p.pos]) {
		p.pos++
	}
	return Number(string(p.input[start:p.pos]))
}

func (p *parser) skipSpaces() {
	for p.pos < len(p.input) && p.input[p.pos] == ' ' {
		p.pos++
	}
}
