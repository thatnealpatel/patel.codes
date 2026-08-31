package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"patel.codes/render"
)

func TestBuildPageUsesRender(t *testing.T) {
	src := []byte("# test\n\nbefore $x^2$ after\n\n```go\nx := \"$notMath$\" // comment\n```\n\n| expression |\n|---|\n| $\\frac{a}{b}$ |\n")
	page, err := buildPage(src, pageMeta{Title: "test", URL: "https://patel.codes/test.html"})
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range [][]byte{
		[]byte(`<h1>test</h1>`),
		[]byte("<msup>"),
		[]byte("<mfrac>"),
		[]byte("<table>"),
		[]byte(`x := &quot;$notMath$&quot; <span class="cm">// comment</span>`),
		[]byte(`[<a href="/">home</a>]`),
	} {
		if !bytes.Contains(page, want) {
			t.Errorf("page does not contain %q:\n%s", want, page)
		}
	}
	if got, want := bytes.Count(page, []byte("<math")), 2; got != want {
		t.Errorf("math element count = %d, want %d:\n%s", got, want, page)
	}
}

func TestBuildPageReportsInvalidMath(t *testing.T) {
	_, err := buildPage([]byte(`$\left(x$`), pageMeta{})
	if !errors.Is(err, render.ErrInvalidMath) {
		t.Fatalf("error = %v, want render.ErrInvalidMath", err)
	}
}

func TestBuildPageGeneratedSections(t *testing.T) {
	src := []byte("# test\n\nauthored\n\n:::gen\n\ngenerated\n\n:::\n")
	page, err := buildPage(src, pageMeta{Title: "test", URL: "https://patel.codes/words/test.html"})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range [][]byte{
		[]byte(`<p class="gen-notice">`),
		[]byte("<div class=\"gen\">\n<p>generated</p>\n</div>"),
	} {
		if !bytes.Contains(page, want) {
			t.Errorf("page does not contain %q:\n%s", want, page)
		}
	}
}

func TestBuildPageDisclosureSkipsGeneratedNotice(t *testing.T) {
	src := []byte("# disclosure\n\n:::gen\n\ngenerated\n\n:::\n")
	page, err := buildPage(src, pageMeta{
		Title: "disclosure",
		URL:   "https://patel.codes/words/patel_codes_llm_disclosure.html",
	})
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(page, []byte(`class="gen-notice"`)) {
		t.Fatalf("disclosure contains redundant generated-content notice:\n%s", page)
	}
	if !bytes.Contains(page, []byte(`<div class="gen">`)) {
		t.Fatalf("generated section was not wrapped:\n%s", page)
	}
}

func TestMarkdownContentRenders(t *testing.T) {
	err := filepath.WalkDir("../../data", func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		source, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if _, err := render.Render(string(source)); err != nil {
			t.Errorf("rendering %s: %v", path, err)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSpanComments(t *testing.T) {
	goSyntax := commentTokens["go"]
	leanSyntax := commentTokens["lean"]
	tests := []struct {
		name   string
		syntax commentSyntax
		input  string
		want   string
	}{
		{"go single", goSyntax, "x := 1 // half\ny := 2\n", "x := 1 <span class=\"cm\">// half</span>\ny := 2\n"},
		{"go multi", goSyntax, "a /* one\ntwo */ b", "a <span class=\"cm\">/* one\ntwo */</span> b"},
		{"go unterminated multi", goSyntax, "a /* open\nrest", "a <span class=\"cm\">/* open\nrest</span>"},
		{"lean single", leanSyntax, "def f := 0 -- note\n", "def f := 0 <span class=\"cm\">-- note</span>\n"},
		{"lean doc before multi", leanSyntax, "/-- doc -/\ndef f := 0", "<span class=\"cm\">/-- doc -/</span>\ndef f := 0"},
		{"lean multi", leanSyntax, "/- block -/ def f := 0", "<span class=\"cm\">/- block -/</span> def f := 0"},
		{"single at EOF", goSyntax, "x // tail", "x <span class=\"cm\">// tail</span>"},
		{"incomplete multi syntax", commentSyntax{singleStart: "#", multiStart: "=begin"}, "a =begin\n# hi\n", "a =begin\n<span class=\"cm\"># hi</span>\n"},
		{"no comments", goSyntax, "x := 1\n", "x := 1\n"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := spanComments(test.input, test.syntax); got != test.want {
				t.Fatalf("spanComments(%q):\n got %q\nwant %q", test.input, got, test.want)
			}
		})
	}
}

func TestHighlightComments(t *testing.T) {
	input := "<p>hi</p>\n<pre><code class=\"language-go\">x := 1 // c\n</code></pre>\n" +
		"<pre><code class=\"language-py\"># nope\n</code></pre>\n" +
		"<pre><code>// bare\n</code></pre>\n"
	got := highlightComments(input)
	if !strings.Contains(got, `<span class="cm">// c</span>`) {
		t.Fatalf("Go comment not highlighted: %s", got)
	}
	if strings.Contains(got, `<span class="cm"># nope`) {
		t.Fatalf("unlisted language highlighted: %s", got)
	}
	if strings.Contains(got, `<span class="cm">// bare`) {
		t.Fatalf("bare block highlighted: %s", got)
	}
}
