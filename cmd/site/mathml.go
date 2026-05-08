package main

import (
	"bytes"

	"github.com/thatnealpatel/patel.codes/internal/parser"
)

var greekLetters = map[string]string{
	`\alpha`: "α", `\beta`: "β", `\gamma`: "γ", `\delta`: "δ",
	`\epsilon`: "ε", `\varepsilon`: "ε", `\zeta`: "ζ", `\eta`: "η",
	`\theta`: "θ", `\iota`: "ι", `\kappa`: "κ", `\lambda`: "λ",
	`\mu`: "μ", `\nu`: "ν", `\xi`: "ξ", `\pi`: "π",
	`\rho`: "ρ", `\sigma`: "σ", `\tau`: "τ", `\upsilon`: "υ",
	`\phi`: "φ", `\varphi`: "φ", `\chi`: "χ", `\psi`: "ψ", `\omega`: "ω",
	`\Gamma`: "Γ", `\Delta`: "Δ", `\Theta`: "Θ", `\Lambda`: "Λ",
	`\Xi`: "Ξ", `\Pi`: "Π", `\Sigma`: "Σ", `\Phi`: "Φ",
	`\Psi`: "Ψ", `\Omega`: "Ω",
}

var operators = map[string]string{
	`\pm`: "±", `\mp`: "∓", `\times`: "×", `\div`: "÷",
	`\cdot`: "⋅", `\leq`: "≤", `\geq`: "≥", `\neq`: "≠",
	`\approx`: "≈", `\equiv`: "≡", `\in`: "∈", `\notin`: "∉",
	`\subset`: "⊂", `\supset`: "⊃", `\cup`: "∪", `\cap`: "∩",
	`\rightarrow`: "→", `\leftarrow`: "←", `\Rightarrow`: "⇒",
	`\Leftarrow`: "⇐", `\iff`: "⟺",
	`\infty`: "∞", `\partial`: "∂", `\nabla`: "∇",
	`\sum`: "∑", `\prod`: "∏", `\int`: "∫",
	`\ldots`: "…", `\cdots`: "⋯", `\dots`: "…",
	`\quad`: " ", `\qquad`: "  ",
}

var specialChars = map[string]string{
	`\{`: "{", `\}`: "}", `\|`: "‖",
	`\,`: " ", `\;`: " ", `\:`: " ", `\!`: "",
	`\\`: "\n",
}

func latexToMathML(expr string, display bool) (string, error) {
	nodes, err := parser.Parse(expr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if display {
		buf.WriteString(`<math display="block">`)
	} else {
		buf.WriteString(`<math>`)
	}
	writeNodes(&buf, nodes)
	buf.WriteString(`</math>`)
	return buf.String(), nil
}

func writeNodes(buf *bytes.Buffer, nodes []parser.Node) {
	if len(nodes) > 1 {
		buf.WriteString("<mrow>")
	}
	for _, n := range nodes {
		writeNode(buf, n)
	}
	if len(nodes) > 1 {
		buf.WriteString("</mrow>")
	}
}

func writeNode(buf *bytes.Buffer, n parser.Node) {
	switch n := n.(type) {
	case parser.List:
		writeNodes(buf, []parser.Node(n))

	case parser.Letter:
		buf.WriteString("<mi>" + string(n) + "</mi>")

	case parser.Number:
		buf.WriteString("<mn>" + string(n) + "</mn>")

	case parser.Operator:
		s := string(n)
		if repl, ok := specialChars[s]; ok {
			if repl != "" {
				buf.WriteString("<mo>" + repl + "</mo>")
			}
		} else {
			buf.WriteString("<mo>" + s + "</mo>")
		}

	case parser.Sup:
		buf.WriteString("<msup>")
		writeNode(buf, n.Base)
		writeNode(buf, n.Script)
		buf.WriteString("</msup>")

	case parser.Sub:
		buf.WriteString("<msub>")
		writeNode(buf, n.Base)
		writeNode(buf, n.Script)
		buf.WriteString("</msub>")

	case parser.Space:

	case parser.Command:
		writeCommand(buf, n)
	}
}

func writeCommand(buf *bytes.Buffer, cmd parser.Command) {
	switch cmd.Name {
	case `\frac`, `\dfrac`, `\tfrac`:
		buf.WriteString("<mfrac>")
		for _, arg := range cmd.Args {
			writeNodes(buf, arg)
		}
		buf.WriteString("</mfrac>")

	case `\sqrt`:
		if len(cmd.OptArgs) > 0 {
			buf.WriteString("<mroot>")
			if len(cmd.Args) > 0 {
				writeNodes(buf, cmd.Args[0])
			}
			writeNodes(buf, cmd.OptArgs[0])
			buf.WriteString("</mroot>")
		} else {
			buf.WriteString("<msqrt>")
			for _, arg := range cmd.Args {
				writeNodes(buf, arg)
			}
			buf.WriteString("</msqrt>")
		}

	case `\overline`:
		buf.WriteString("<mover><mrow>")
		for _, arg := range cmd.Args {
			writeNodes(buf, arg)
		}
		buf.WriteString("</mrow><mo>¯</mo></mover>")

	case `\hat`:
		buf.WriteString("<mover><mrow>")
		for _, arg := range cmd.Args {
			writeNodes(buf, arg)
		}
		buf.WriteString("</mrow><mo>^</mo></mover>")

	case `\vec`:
		buf.WriteString("<mover><mrow>")
		for _, arg := range cmd.Args {
			writeNodes(buf, arg)
		}
		buf.WriteString("</mrow><mo>→</mo></mover>")

	case `\textit`:
		buf.WriteString(`<mtext mathvariant="italic">`)
		for _, arg := range cmd.Args {
			writeArgText(buf, arg)
		}
		buf.WriteString("</mtext>")

	case `\textbf`, `\mathbf`:
		buf.WriteString(`<mtext mathvariant="bold">`)
		for _, arg := range cmd.Args {
			writeArgText(buf, arg)
		}
		buf.WriteString("</mtext>")

	case `\text`, `\textrm`, `\mathrm`:
		buf.WriteString("<mtext>")
		for _, arg := range cmd.Args {
			writeArgText(buf, arg)
		}
		buf.WriteString("</mtext>")

	case `\mod`:
		buf.WriteString("<mo>mod</mo>")
		for _, arg := range cmd.Args {
			writeNodes(buf, arg)
		}

	case `\pmod`:
		buf.WriteString("<mrow><mo>(</mo><mo>mod</mo>")
		for _, arg := range cmd.Args {
			writeNodes(buf, arg)
		}
		buf.WriteString("<mo>)</mo></mrow>")

	case `\eqref`:
		if len(cmd.Args) > 0 {
			buf.WriteString("<mtext>(</mtext>")
			writeArgText(buf, cmd.Args[0])
			buf.WriteString("<mtext>)</mtext>")
		}

	default:
		if s, ok := greekLetters[cmd.Name]; ok {
			buf.WriteString("<mi>" + s + "</mi>")
		} else if s, ok := operators[cmd.Name]; ok {
			buf.WriteString("<mo>" + s + "</mo>")
		}
	}
}

func writeArgText(buf *bytes.Buffer, nodes []parser.Node) {
	for _, n := range nodes {
		switch n := n.(type) {
		case parser.Letter:
			buf.WriteString(string(n))
		case parser.Number:
			buf.WriteString(string(n))
		case parser.Operator:
			buf.WriteString(string(n))
		case parser.Space:
			buf.WriteString(" ")
		case parser.List:
			writeArgText(buf, []parser.Node(n))
		default:
			writeNode(buf, n)
		}
	}
}
