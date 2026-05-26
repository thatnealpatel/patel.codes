package parser

import (
	"testing"
)

func TestParseAtoms(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  int // expected node count
	}{
		{"letter", "x", 1},
		{"digit", "7", 1},
		{"multi_digit", "42", 1},
		{"operator", "+", 1},
		{"space", " ", 1},
		{"empty", "", 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nodes, err := Parse(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if len(nodes) != tc.want {
				t.Fatalf("got %d nodes, want %d: %+v", len(nodes), tc.want, nodes)
			}
		})
	}
}

func TestParseAtomTypes(t *testing.T) {
	nodes, err := Parse("x")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := nodes[0].(Letter); !ok {
		t.Fatalf("expected Letter, got %T", nodes[0])
	}

	nodes, err = Parse("42")
	if err != nil {
		t.Fatal(err)
	}
	if n, ok := nodes[0].(Number); !ok || string(n) != "42" {
		t.Fatalf("expected Number(42), got %T %+v", nodes[0], nodes[0])
	}

	nodes, err = Parse("+")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := nodes[0].(Operator); !ok {
		t.Fatalf("expected Operator, got %T", nodes[0])
	}
}

func TestParseSupSub(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"sup_digit", `b^2`},
		{"sup_letter", `x^n`},
		{"sup_group", `e^{i\pi}`},
		{"sub_letter", `x_i`},
		{"sub_group", `a_{n+1}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nodes, err := Parse(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			if len(nodes) != 1 {
				t.Fatalf("expected 1 node (sup/sub pair), got %d: %+v", len(nodes), nodes)
			}
			switch tc.name[:3] {
			case "sup":
				if _, ok := nodes[0].(Sup); !ok {
					t.Fatalf("expected Sup, got %T", nodes[0])
				}
			case "sub":
				if _, ok := nodes[0].(Sub); !ok {
					t.Fatalf("expected Sub, got %T", nodes[0])
				}
			}
		})
	}
}

func TestParseSupBase(t *testing.T) {
	nodes, err := Parse(`ax^2`)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes [a, Sup{x,2}], got %d: %+v", len(nodes), nodes)
	}
	if l, ok := nodes[0].(Letter); !ok || string(l) != "a" {
		t.Fatalf("expected Letter(a), got %T %+v", nodes[0], nodes[0])
	}
	sup, ok := nodes[1].(Sup)
	if !ok {
		t.Fatalf("expected Sup, got %T", nodes[1])
	}
	if l, ok := sup.Base.(Letter); !ok || string(l) != "x" {
		t.Fatalf("expected base Letter(x), got %T %+v", sup.Base, sup.Base)
	}
	if n, ok := sup.Script.(Number); !ok || string(n) != "2" {
		t.Fatalf("expected script Number(2), got %T %+v", sup.Script, sup.Script)
	}
}

func TestParseCommands(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		cmdName string
		numArgs int
	}{
		{"frac", `\frac{a}{b}`, `\frac`, 2},
		{"sqrt", `\sqrt{x}`, `\sqrt`, 1},
		{"textit", `\textit{hi}`, `\textit`, 1},
		{"greek", `\alpha`, `\alpha`, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			nodes, err := Parse(tc.input)
			if err != nil {
				t.Fatal(err)
			}
			cmd, ok := nodes[0].(Command)
			if !ok {
				t.Fatalf("expected Command, got %T", nodes[0])
			}
			if cmd.Name != tc.cmdName {
				t.Fatalf("expected name %q, got %q", tc.cmdName, cmd.Name)
			}
			if len(cmd.Args) != tc.numArgs {
				t.Fatalf("expected %d args, got %d", tc.numArgs, len(cmd.Args))
			}
		})
	}
}

func TestParseSqrtOptArg(t *testing.T) {
	nodes, err := Parse(`\sqrt[3]{x}`)
	if err != nil {
		t.Fatal(err)
	}
	cmd, ok := nodes[0].(Command)
	if !ok {
		t.Fatalf("expected Command, got %T", nodes[0])
	}
	if len(cmd.OptArgs) != 1 {
		t.Fatalf("expected 1 opt arg, got %d", len(cmd.OptArgs))
	}
	if len(cmd.Args) != 1 {
		t.Fatalf("expected 1 required arg, got %d", len(cmd.Args))
	}
}

func TestParseSpacePreserved(t *testing.T) {
	nodes, err := Parse(`\textit{hello world}`)
	if err != nil {
		t.Fatal(err)
	}
	cmd, ok := nodes[0].(Command)
	if !ok {
		t.Fatalf("expected Command, got %T", nodes[0])
	}
	hasSpace := false
	for _, n := range cmd.Args[0] {
		if _, ok := n.(Space); ok {
			hasSpace = true
			break
		}
	}
	if !hasSpace {
		t.Fatalf("space in \\textit arg was lost: %+v", cmd.Args[0])
	}
}

func TestParseGroups(t *testing.T) {
	nodes, err := Parse(`{a+b}^2`)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node (Sup), got %d", len(nodes))
	}
	sup, ok := nodes[0].(Sup)
	if !ok {
		t.Fatalf("expected Sup, got %T", nodes[0])
	}
	if _, ok := sup.Base.(List); !ok {
		t.Fatalf("expected List base, got %T", sup.Base)
	}
}

func TestParseNestedFrac(t *testing.T) {
	nodes, err := Parse(`\frac{\frac{a}{b}}{c}`)
	if err != nil {
		t.Fatal(err)
	}
	outer, ok := nodes[0].(Command)
	if !ok {
		t.Fatalf("expected Command, got %T", nodes[0])
	}
	inner, ok := outer.Args[0][0].(Command)
	if !ok {
		t.Fatalf("expected nested Command, got %T", outer.Args[0][0])
	}
	if inner.Name != `\frac` {
		t.Fatalf("expected \\frac, got %q", inner.Name)
	}
}

func TestParseBackslashSymbols(t *testing.T) {
	nodes, err := Parse(`\{`)
	if err != nil {
		t.Fatal(err)
	}
	op, ok := nodes[0].(Operator)
	if !ok {
		t.Fatalf("expected Operator, got %T", nodes[0])
	}
	if string(op) != `\{` {
		t.Fatalf("expected \\{, got %q", string(op))
	}
}

func TestParseBrackets(t *testing.T) {
	nodes, err := Parse(`[1, 10]`)
	if err != nil {
		t.Fatal(err)
	}
	hasOpen := false
	hasClose := false
	for _, n := range nodes {
		if op, ok := n.(Operator); ok {
			switch string(op) {
			case "[":
				hasOpen = true
			case "]":
				hasClose = true
			}
		}
	}
	if !hasOpen || !hasClose {
		t.Fatalf("expected both [ and ] as operators, got nodes: %+v", nodes)
	}
}

func TestParseSupWithSpace(t *testing.T) {
	nodes, err := Parse(`x ^2`)
	if err != nil {
		t.Fatal(err)
	}
	hasSup := false
	for _, n := range nodes {
		if _, ok := n.(Sup); ok {
			hasSup = true
		}
	}
	if !hasSup {
		t.Fatalf("space before ^ should not prevent sup pairing: %+v", nodes)
	}
}

func TestParseDelimited(t *testing.T) {
	nodes, err := Parse(`\left(\frac{a}{b}\right)`)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d: %+v", len(nodes), nodes)
	}
	d, ok := nodes[0].(Delimited)
	if !ok {
		t.Fatalf("expected Delimited, got %T", nodes[0])
	}
	if d.Open != "(" || d.Close != ")" {
		t.Fatalf("expected ( ), got %q %q", d.Open, d.Close)
	}
}

func TestParseDelimitedInvisible(t *testing.T) {
	nodes, err := Parse(`\left.\frac{a}{b}\right|`)
	if err != nil {
		t.Fatal(err)
	}
	d, ok := nodes[0].(Delimited)
	if !ok {
		t.Fatalf("expected Delimited, got %T", nodes[0])
	}
	if d.Open != "." || d.Close != "|" {
		t.Fatalf("expected . |, got %q %q", d.Open, d.Close)
	}
}

func TestParseDelimitedBrace(t *testing.T) {
	nodes, err := Parse(`\left\{x\right\}`)
	if err != nil {
		t.Fatal(err)
	}
	d, ok := nodes[0].(Delimited)
	if !ok {
		t.Fatalf("expected Delimited, got %T", nodes[0])
	}
	if d.Open != `\{` || d.Close != `\}` {
		t.Fatalf(`expected \{ \}, got %q %q`, d.Open, d.Close)
	}
}

func TestParseEnvCases(t *testing.T) {
	nodes, err := Parse(`\begin{cases} a & b \\ c & d \end{cases}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d: %+v", len(nodes), nodes)
	}
	env, ok := nodes[0].(Env)
	if !ok {
		t.Fatalf("expected Env, got %T", nodes[0])
	}
	if env.Name != "cases" {
		t.Fatalf("expected cases, got %q", env.Name)
	}
}

func TestParseSubstack(t *testing.T) {
	nodes, err := Parse(`\sum_{\substack{d \mid k \\ d > 0}}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d: %+v", len(nodes), nodes)
	}
	sub, ok := nodes[0].(Sub)
	if !ok {
		t.Fatalf("expected Sub, got %T", nodes[0])
	}
	var hasSubstack bool
	switch s := sub.Script.(type) {
	case Command:
		hasSubstack = s.Name == `\substack`
	case List:
		for _, n := range s {
			if cmd, ok := n.(Command); ok && cmd.Name == `\substack` {
				hasSubstack = true
			}
		}
	}
	if !hasSubstack {
		t.Fatalf("expected \\substack in script, got %+v", sub.Script)
	}
}

func TestParseDelimitedSup(t *testing.T) {
	nodes, err := Parse(`\left(\frac{a}{b}\right)^2`)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node (Sup), got %d: %+v", len(nodes), nodes)
	}
	sup, ok := nodes[0].(Sup)
	if !ok {
		t.Fatalf("expected Sup, got %T", nodes[0])
	}
	if _, ok := sup.Base.(Delimited); !ok {
		t.Fatalf("expected Delimited base, got %T", sup.Base)
	}
}

func TestParseSubSup(t *testing.T) {
	nodes, err := Parse(`\sum_{r=0}^{k-1}`)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node (SubSup), got %d: %+v", len(nodes), nodes)
	}
	ss, ok := nodes[0].(SubSup)
	if !ok {
		t.Fatalf("expected SubSup, got %T", nodes[0])
	}
	cmd, ok := ss.Base.(Command)
	if !ok {
		t.Fatalf("expected Command base, got %T", ss.Base)
	}
	if cmd.Name != `\sum` {
		t.Fatalf("expected \\sum, got %q", cmd.Name)
	}
}

func TestParseSupThenSub(t *testing.T) {
	nodes, err := Parse(`x^2_i`)
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d: %+v", len(nodes), nodes)
	}
	if _, ok := nodes[0].(SubSup); !ok {
		t.Fatalf("expected SubSup, got %T", nodes[0])
	}
}

func TestParseRightWithoutLeft(t *testing.T) {
	_, err := Parse(`\right)`)
	if err == nil {
		t.Fatal("expected error for \\right without \\left")
	}
}

func TestParseEndWithoutBegin(t *testing.T) {
	_, err := Parse(`\end{cases}`)
	if err == nil {
		t.Fatal("expected error for \\end without \\begin")
	}
}
