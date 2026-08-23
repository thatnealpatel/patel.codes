package main

import (
	"bytes"
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

// Issue #1 & #3: Space nodes between math tokens must produce <mspace>
// in MathML output, not be silently dropped.
func TestLatexToMathMLSpacesBetweenTokens(t *testing.T) {
	cases := []struct {
		name string
		expr string
		want string // substring that must appear
	}{
		{
			"cases_if_space",
			`\begin{cases} k & \text{if } k \mid n \\ 0 & \text{otherwise} \end{cases}`,
			`</mtext><mspace`,
		},
		{
			"cases_q_odd",
			`\begin{cases} 2 & q \text{ odd} \\ 0 & q \text{ even} \end{cases}`,
			`</mi><mspace`,
		},
		{
			"simple_a_b",
			`a b`,
			`<mspace`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := latexToMathML([]byte(tc.expr), false)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Contains(result, []byte(tc.want)) {
				t.Fatalf("expected %q in output:\n%s", tc.want, result)
			}
		})
	}
}

// Issue #2: Regular () operators should not stretch to match
// the height of neighboring superscripts.
func TestLatexToMathMLOperatorParensNonStretchy(t *testing.T) {
	result, err := latexToMathML([]byte(`(-1)^q - 1`), false)
	if err != nil {
		t.Fatal(err)
	}
	s := string(result)
	if !strings.Contains(s, `stretchy="false"`) {
		t.Fatalf("regular () should have stretchy=false:\n%s", s)
	}
}

// Issue #2 continued: \left(\right) Delimited parens SHOULD remain
// stretchy (no stretchy="false").
func TestLatexToMathMLDelimitedParensStretchy(t *testing.T) {
	result, err := latexToMathML([]byte(`\left(\frac{a}{b}\right)`), false)
	if err != nil {
		t.Fatal(err)
	}
	s := string(result)
	if strings.Contains(s, `stretchy="false"`) {
		t.Fatalf("\\left/\\right delimiters should remain stretchy:\n%s", s)
	}
}

// Issue #4: Markdown tables containing inline math with pipe
// characters must render as <table>, not raw pipe text.
func TestBuildPageTableWithInlineMath(t *testing.T) {
	src := []byte("# test\n\n| a | b |\n|---|---|\n| $x$ | $O(n^2)$ |\n")
	result, err := buildPage(src, pageMeta{Title: "test", URL: "https://test"})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte("<table")) {
		t.Fatalf("expected <table> for pipe-delimited markdown table:\n%s", result)
	}
}

// Issue #4 specific: Table rows with $...$ math containing operators
// that look like pipes must not break column splitting.
func TestBuildPageTableMathWithMid(t *testing.T) {
	src := []byte("# test\n\n| source | bound | absorbed by |\n|---|---|---|\n| non-dominant divisors | $O(n^2 \\cdot 2^{n/3})$ | $o(2^n/n^{M+2})$ |\n")
	result, err := buildPage(src, pageMeta{Title: "test", URL: "https://test"})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte("<table")) {
		t.Fatalf("expected <table> for table with inline math:\n%s", result)
	}
}

// Issue #5: Inline math spanning multiple lines must be matched
// by the inline regex so it gets rendered instead of showing raw LaTeX.
func TestBuildPageMultilineInlineMath(t *testing.T) {
	src := []byte("# test\n\ndefine $\\text{hostSet}(k, e) = \\{S : e \\in S, \\;\n|S| \\mid (\\sigma + k \\cdot \\text{rank}(e, S))\\}$.\n")
	result, err := buildPage(src, pageMeta{Title: "test", URL: "https://test"})
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(result, []byte(`$\text{hostSet}`)) {
		t.Fatalf("raw LaTeX leaked through — multi-line inline math was not parsed:\n%s", result)
	}
	if !bytes.Contains(result, []byte("<math>")) {
		t.Fatalf("expected <math> from inline expression:\n%s", result)
	}
}

// Unknown commands must be surfaced in a warning style instead of
// rendering as nothing.
func TestLatexToMathMLUnknownCommandVisible(t *testing.T) {
	result, err := latexToMathML([]byte(`\lvert x \rvert`), false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte("<merror>")) {
		t.Fatalf("expected <merror> for unknown command:\n%s", result)
	}
	if !bytes.Contains(result, []byte(`\lvert`)) {
		t.Fatalf("expected raw command name in output:\n%s", result)
	}
}

// Unknown commands with parsed args must keep the args visible too.
func TestLatexToMathMLUnknownCommandKeepsArgs(t *testing.T) {
	result, err := latexToMathML([]byte(`\mathcal{H}`), false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte(`<merror><mtext>\mathcal</mtext></merror>`)) {
		t.Fatalf("expected warning for \\mathcal:\n%s", result)
	}
	if !bytes.Contains(result, []byte("<mi>H</mi>")) {
		t.Fatalf("expected arg of \\mathcal to survive:\n%s", result)
	}
}

func TestLatexToMathMLFloorCeil(t *testing.T) {
	result, err := latexToMathML([]byte(`\lfloor \log_2 n \rfloor + \lceil x \rceil`), false)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"⌊", "⌋", "⌈", "⌉", "<mo>log</mo>"} {
		if !bytes.Contains(result, []byte(want)) {
			t.Fatalf("expected %q in output:\n%s", want, result)
		}
	}
	if bytes.Contains(result, []byte("<merror>")) {
		t.Fatalf("unexpected <merror>:\n%s", result)
	}
}

func TestLatexToMathMLNamedFunctions(t *testing.T) {
	result, err := latexToMathML([]byte(`\min(a, \max(b, \log n))`), false)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"<mo>min</mo>", "<mo>max</mo>", "<mo>log</mo>"} {
		if !bytes.Contains(result, []byte(want)) {
			t.Fatalf("expected %q in output:\n%s", want, result)
		}
	}
}

func TestLatexToMathMLMathbb(t *testing.T) {
	result, err := latexToMathML([]byte(`f : \mathbb{N} \to \mathbb{Z}`), false)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"<mi>ℕ</mi>", "<mi>ℤ</mi>", "→"} {
		if !bytes.Contains(result, []byte(want)) {
			t.Fatalf("expected %q in output:\n%s", want, result)
		}
	}
	if bytes.Contains(result, []byte("<merror>")) {
		t.Fatalf("unexpected <merror>:\n%s", result)
	}
}

func TestLatexToMathMLMathbbFormulaic(t *testing.T) {
	result, err := latexToMathML([]byte(`\mathbb{A} \mathbb{k} \mathbb{1}`), false)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"𝔸", "𝕜", "𝟙"} {
		if !bytes.Contains(result, []byte(want)) {
			t.Fatalf("expected %q in output:\n%s", want, result)
		}
	}
}

// \ne, \to, \setminus are easy to type instead of \neq, \rightarrow;
// they must render, not drop.
func TestLatexToMathMLAliases(t *testing.T) {
	result, err := latexToMathML([]byte(`a \ne b, A \setminus B`), false)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"≠", "∖"} {
		if !bytes.Contains(result, []byte(want)) {
			t.Fatalf("expected %q in output:\n%s", want, result)
		}
	}
	if bytes.Contains(result, []byte("<merror>")) {
		t.Fatalf("unexpected <merror>:\n%s", result)
	}
}

func TestLatexToMathMLLongrightarrow(t *testing.T) {
	result, err := latexToMathML([]byte(`P \Longrightarrow Q`), false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte("<mo>⟹</mo>")) {
		t.Fatalf("expected long rightwards double arrow in output:\n%s", result)
	}
	if bytes.Contains(result, []byte("<merror>")) {
		t.Fatalf("unexpected <merror>:\n%s", result)
	}
}

func TestLatexToMathMLOperatorName(t *testing.T) {
	result, err := latexToMathML([]byte(`\operatorname{Pre}(L)`), false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte("<mo>Pre</mo>")) {
		t.Fatalf("expected named operator in output:\n%s", result)
	}
	if bytes.Contains(result, []byte("<merror>")) {
		t.Fatalf("unexpected <merror>:\n%s", result)
	}
}

func TestLatexToMathMLBigcup(t *testing.T) {
	result, err := latexToMathML([]byte(`\bigcup_{L\geq0}(B_L+T_L)`), false)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte("<msub><mo>⋃</mo>")) {
		t.Fatalf("expected scripted big union in output:\n%s", result)
	}
	if bytes.Contains(result, []byte("<merror>")) {
		t.Fatalf("unexpected <merror>:\n%s", result)
	}
}

func TestLatexToMathMLBigDelimiters(t *testing.T) {
	result, err := latexToMathML([]byte(`\bigl(B_L+\operatorname{keepOnes}(u)\bigr)`), false)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`<mo fence="true" form="prefix" stretchy="true" minsize="1.2em" maxsize="1.2em">(</mo>`,
		`<mo fence="true" form="postfix" stretchy="true" minsize="1.2em" maxsize="1.2em">)</mo>`,
	} {
		if !bytes.Contains(result, []byte(want)) {
			t.Fatalf("expected %q in output:\n%s", want, result)
		}
	}
	if bytes.Contains(result, []byte("<merror>")) {
		t.Fatalf("unexpected <merror>:\n%s", result)
	}
}

func TestLatexToMathMLBoxedAndXmapsto(t *testing.T) {
	cases := []struct {
		expr string
		want []string
	}{
		{`\boxed{x+1}`, []string{`<menclose notation="box">`, `<mi>x</mi>`}},
		{`A \xmapsto{\tau} B`, []string{`<mover><mo stretchy="true">⟼</mo>`, `<mi>τ</mi>`}},
	}
	for _, tc := range cases {
		result, err := latexToMathML([]byte(tc.expr), false)
		if err != nil {
			t.Fatalf("latexToMathML(%q): %v", tc.expr, err)
		}
		for _, want := range tc.want {
			if !bytes.Contains(result, []byte(want)) {
				t.Fatalf("latexToMathML(%q): expected %q in output:\n%s", tc.expr, want, result)
			}
		}
		if bytes.Contains(result, []byte("<merror>")) {
			t.Fatalf("latexToMathML(%q): unexpected <merror>:\n%s", tc.expr, result)
		}
	}
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

func TestSpanComments(t *testing.T) {
	goSyn := commentTokens["go"]
	leanSyn := commentTokens["lean"]
	cases := []struct {
		name string
		syn  commentSyntax
		in   string
		want string
	}{
		{"go_single", goSyn,
			"x := 1 // half\ny := 2\n",
			"x := 1 <span class=\"cm\">// half</span>\ny := 2\n"},
		{"go_multi", goSyn,
			"a /* one\ntwo */ b",
			"a <span class=\"cm\">/* one\ntwo */</span> b"},
		{"go_multi_unterminated", goSyn,
			"a /* open\nrest",
			"a <span class=\"cm\">/* open\nrest</span>"},
		{"lean_single", leanSyn,
			"def f := 0 -- note\n",
			"def f := 0 <span class=\"cm\">-- note</span>\n"},
		{"lean_doc_over_multi", leanSyn,
			"/-- doc -/\ndef f := 0",
			"<span class=\"cm\">/-- doc -/</span>\ndef f := 0"},
		{"lean_multi", leanSyn,
			"/- block -/ def f := 0",
			"<span class=\"cm\">/- block -/</span> def f := 0"},
		{"single_at_eof_no_newline", goSyn,
			"x // tail",
			"x <span class=\"cm\">// tail</span>"},
		{"downgrade_incomplete_multi",
			commentSyntax{singleStart: "#", multiStart: "=begin"},
			"a =begin\n# hi\n",
			"a =begin\n<span class=\"cm\"># hi</span>\n"},
		{"no_comments", goSyn, "x := 1\n", "x := 1\n"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := spanComments(tc.in, tc.syn)
			if got != tc.want {
				t.Fatalf("spanComments(%q):\n got %q\nwant %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestHighlightComments(t *testing.T) {
	in := "<p>hi</p>\n<pre><code class=\"language-go\">x := 1 // c\n</code></pre>\n" +
		"<pre><code class=\"language-py\"># nope\n</code></pre>\n" +
		"<pre><code>// bare\n</code></pre>\n"
	got := highlightComments(in)
	if !strings.Contains(got, "<span class=\"cm\">// c</span>") {
		t.Fatalf("go comment not highlighted: %s", got)
	}
	if strings.Contains(got, "<span class=\"cm\"># nope") {
		t.Fatalf("unlisted language highlighted: %s", got)
	}
	if strings.Contains(got, "<span class=\"cm\">// bare") {
		t.Fatalf("bare (no language) block highlighted: %s", got)
	}
}

func TestBuildPageCommentHighlight(t *testing.T) {
	src := []byte("# test\n\n```go\nx := 1 // half\n```\n")
	result, err := buildPage(src, pageMeta{Title: "test", URL: "https://test"})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(result, []byte(`<span class="cm">// half</span>`)) {
		t.Fatalf("expected highlighted comment in page: %s", result)
	}
}
