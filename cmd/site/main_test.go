package main

import (
	"strings"
	"testing"
)

func TestLatexToMathML(t *testing.T) {
	cases := []struct {
		name string
		expr string
	}{
		// basic MathML element mapping
		{"letter", `x`},
		{"number", `42`},
		{"operator", `+`},
		{"greek", `\alpha + \beta = \gamma`},

		// sup/sub produce <msup>/<msub>
		{"sup", `b^2`},
		{"sub", `x_i`},
		{"sup_group", `e^{i\pi}`},
		{"sub_group", `a_{n+1}`},
		{"sup_splits_word", `ax^2`},

		// commands produce correct MathML elements
		{"frac", `\frac{a}{b}`},
		{"frac_nested", `\frac{\frac{a}{b}}{c}`},
		{"sqrt", `\sqrt{x}`},
		{"sqrt_nth", `\sqrt[3]{x}`},
		{"overline", `\overline{x}`},
		{"textit", `\textit{hello world}`},
		{"textbf", `\textbf{bold}`},
		{"text", `\text{conjecture}`},
		{"mod", `(t + D)\mod{7} = d`},
		{"pmod", `a \equiv b \pmod{n}`},

		// compound expressions from the blog
		{"quadratic", `x = \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}`},
		{"poly", `ax^2 + bx + c = 0`},
		{"text_sup", `\text{reasoning}^{\textbf{m}}`},
		{"mixed", `C(1-\epsilon)\textbf{m}+\epsilon\text{h}`},
		{"arrows", `\text{conjecture}^{\text{h}} \rightarrow \text{reasoning}^{\text{h}} \rightarrow \text{outcome}^{\text{h}}`},

		// special characters
		{"backslash_brace", `\{`},
		{"thin_space", `a\,b`},

		// blog-specific expressions
		{"set_notation", `T = \{0, 1, 2, 3, 4\}`},
		{"set_membership", `t \in T`},
		{"set_union", `d \in \{T\cup{W}\}`},
		{"prime_sub", `D_{-1}'`},
		{"iff", `d \in W \iff t \in \{3 , 4 \}`},
		{"epsilon_eq", `\epsilon=0`},
		{"reasoning_sup", `\text{reasoning}^{C(1-\epsilon)\textbf{m}+\epsilon\text{h}}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := latexToMathML(tc.expr, false)
			if err != nil {
				t.Fatalf("latexToMathML(%q): %v", tc.expr, err)
			}
			if !strings.HasPrefix(result, "<math>") {
				t.Fatalf("expected <math> prefix, got %s", result)
			}
			if !strings.HasSuffix(result, "</math>") {
				t.Fatalf("expected </math> suffix, got %s", result)
			}
			t.Logf("%s", result)
		})
	}
}

func TestLatexToMathMLDisplay(t *testing.T) {
	result, err := latexToMathML(`x = 1`, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(result, `<math display="block">`) {
		t.Fatalf("expected display block, got %s", result)
	}
}

func TestLatexToMathMLTextSpace(t *testing.T) {
	result, err := latexToMathML(`\textit{hello world}`, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "hello world") {
		t.Fatalf("space in \\textit lost: %s", result)
	}
}

func TestLatexToMathMLSupStructure(t *testing.T) {
	result, err := latexToMathML(`b^2`, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "<msup>") {
		t.Fatalf("expected <msup> in output: %s", result)
	}
	if !strings.Contains(result, "<mi>b</mi>") {
		t.Fatalf("expected <mi>b</mi> as base: %s", result)
	}
	if !strings.Contains(result, "<mn>2</mn>") {
		t.Fatalf("expected <mn>2</mn> as script: %s", result)
	}
}

func TestLatexToMathMLFracStructure(t *testing.T) {
	result, err := latexToMathML(`\frac{x+1}{y}`, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result, "<mfrac>") {
		t.Fatalf("expected <mfrac>: %s", result)
	}
}
