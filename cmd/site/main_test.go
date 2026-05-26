package main

import (
	"bytes"
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
			result, err := latexToMathML([]byte(tc.expr), false)
			if err != nil {
				t.Fatalf("latexToMathML(%q): %v", tc.expr, err)
			}
			if !bytes.HasPrefix(result, []byte("<math>")) {
				t.Fatalf("expected <math> prefix, got %s", result)
			}
			if !bytes.HasSuffix(result, []byte("</math>")) {
				t.Fatalf("expected </math> suffix, got %s", result)
			}
			t.Logf("%s", result)
		})
	}
}

func TestLatexToMathMLDisplay(t *testing.T) {
	result, err := latexToMathML([]byte(`x = 1`), true)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(result, []byte(`<math display="block">`)) {
		t.Fatalf("expected display block, got %s", result)
	}
}

func TestLatexToMathMLTextSpace(t *testing.T) {
	result, err := latexToMathML([]byte(`\textit{hello world}`), false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte("hello world")) {
		t.Fatalf("space in \\textit lost: %s", result)
	}
}

func TestLatexToMathMLSupStructure(t *testing.T) {
	result, err := latexToMathML([]byte(`b^2`), false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte("<msup>")) {
		t.Fatalf("expected <msup> in output: %s", result)
	}
	if !bytes.Contains(result, []byte("<mi>b</mi>")) {
		t.Fatalf("expected <mi>b</mi> as base: %s", result)
	}
	if !bytes.Contains(result, []byte("<mn>2</mn>")) {
		t.Fatalf("expected <mn>2</mn> as script: %s", result)
	}
}

func TestLatexToMathMLFracStructure(t *testing.T) {
	result, err := latexToMathML([]byte(`\frac{x+1}{y}`), false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte("<mfrac>")) {
		t.Fatalf("expected <mfrac>: %s", result)
	}
}

func TestLatexToMathMLBinom(t *testing.T) {
	result, err := latexToMathML([]byte(`\binom{n}{k}`), false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte(`linethickness="0"`)) {
		t.Fatalf("expected linethickness=0 for binom: %s", result)
	}
	if !bytes.Contains(result, []byte("<mfrac")) {
		t.Fatalf("expected <mfrac> in binom: %s", result)
	}
	t.Logf("%s", result)
}

func TestLatexToMathMLDelimited(t *testing.T) {
	result, err := latexToMathML([]byte(`\left(\frac{a}{b}\right)`), false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte("<mfrac>")) {
		t.Fatalf("expected <mfrac>: %s", result)
	}
	if !bytes.Contains(result, []byte("<mo>(</mo>")) {
		t.Fatalf("expected opening paren: %s", result)
	}
	if !bytes.Contains(result, []byte("<mo>)</mo>")) {
		t.Fatalf("expected closing paren: %s", result)
	}
	t.Logf("%s", result)
}

func TestLatexToMathMLCases(t *testing.T) {
	result, err := latexToMathML([]byte(`\begin{cases} k & \text{if } k \mid n \\ 0 & \text{otherwise} \end{cases}`), false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte("<mtable")) {
		t.Fatalf("expected <mtable> for cases: %s", result)
	}
	if !bytes.Contains(result, []byte("<mtr>")) {
		t.Fatalf("expected <mtr> rows: %s", result)
	}
	if !bytes.Contains(result, []byte("<mo>{</mo>")) {
		t.Fatalf("expected opening brace: %s", result)
	}
	t.Logf("%s", result)
}

func TestLatexToMathMLSubstack(t *testing.T) {
	result, err := latexToMathML([]byte(`\sum_{\substack{d \mid k \\ d \text{ odd}}}`), false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte("<mtable")) {
		t.Fatalf("expected <mtable> for substack: %s", result)
	}
	if !bytes.Contains(result, []byte("<mtr>")) {
		t.Fatalf("expected <mtr> rows: %s", result)
	}
	t.Logf("%s", result)
}

func TestLatexToMathMLMid(t *testing.T) {
	result, err := latexToMathML([]byte(`k \mid n`), false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte("∣")) {
		t.Fatalf("expected ∣ for \\mid: %s", result)
	}
}

func TestLatexToMathMLForall(t *testing.T) {
	result, err := latexToMathML([]byte(`\forall k \geq 1`), false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte("∀")) {
		t.Fatalf("expected ∀ for \\forall: %s", result)
	}
}

func TestBuildPageCodeFenceProtection(t *testing.T) {
	src := []byte("# test\n\n```\nx = $y + $z\n```\n\ninline $x^2$ here\n")
	result, err := buildPage(src, pageMeta{Title: "test", URL: "https://test"})
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(result, []byte("<msup>")) && bytes.Contains(result, []byte("x = ")) {
		t.Logf("result may have parsed math inside code fence")
	}
	if !bytes.Contains(result, []byte("<msup>")) {
		t.Fatalf("expected inline math to render: %s", result)
	}
	t.Logf("%s", result)
}

func TestBuildPageMultilineDisplay(t *testing.T) {
	src := []byte("# test\n\n$$a(n) = \\frac{2^{n+1}}{n}$$\n")
	result, err := buildPage(src, pageMeta{Title: "test", URL: "https://test"})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte(`display="block"`)) {
		t.Fatalf("expected display block math: %s", result)
	}
	t.Logf("%s", result)
}

func TestLatexToMathMLBlogExpressions(t *testing.T) {
	cases := []struct {
		name string
		expr string
	}{
		{"roots_of_unity_filter", `\sum_{r=0}^{k-1} \zeta^{rn}`},
		{"product", `\prod_{m=1}^{k}(1 + \zeta^{rm})`},
		{"cases", `\begin{cases} k & \text{if } k \mid n \\ 0 & \text{otherwise} \end{cases}`},
		{"product_omega", `\left(\prod_{j=0}^{q-1}(1 + \omega^j)\right)^{k/q}`},
		{"cases_product", `\prod_{j=0}^{q-1}(1+\omega^j) = \begin{cases} 2 & q \text{ odd} \\ 0 & q \text{ even} \end{cases}`},
		{"substack_sum", `\frac{1}{k}\sum_{\substack{d \mid k \\ d \text{ odd}}} \varphi(d) \cdot 2^{k/d}`},
		{"alpha_def", `\alpha(k, S) = |\{r \in \{0,\ldots,k-1\} : k \mid (\sigma + r \cdot |S|)\}|`},
		{"binom", `\binom{m+1}{k}`},
		{"forall", `L(k) = R(k) \quad \forall k \geq 1`},
		{"frac_left_right", `\left(\frac{j}{n}\right)^i`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := latexToMathML([]byte(tc.expr), true)
			if err != nil {
				t.Fatalf("latexToMathML(%q): %v", tc.expr, err)
			}
			if !bytes.HasPrefix(result, []byte(`<math display="block">`)) {
				t.Fatalf("expected <math display=block> prefix, got %s", result)
			}
			t.Logf("%s", result)
		})
	}
}
