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
		[]byte(`x := &quot;$notMath$&quot; // comment`),
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
		[]byte(`.markdown .gen-notice a { color: #7b9da5; }`),
		[]byte(`<p class="gen-notice">This page contains <a href="/words/patel_codes_llm_disclosure.html">generated content</a> delineated`),
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

func TestTemplatesUseSharedStylesheet(t *testing.T) {
	page, err := buildPage([]byte("# test\n\n```go\nx := 1 // comment\n```\n"), pageMeta{Title: "test"})
	if err != nil {
		t.Fatal(err)
	}
	wall, err := buildGalleryWall(gallery{Title: "test", Images: []string{"test.png"}})
	if err != nil {
		t.Fatal(err)
	}
	for name, output := range map[string]string{"page": string(page), "gallery": wall} {
		t.Run(name, func(t *testing.T) {
			if !strings.Contains(output, "<style>\n"+string(render.Stylesheet())+"\n") {
				t.Error("shared stylesheet is not embedded verbatim in style element")
			}
			if !strings.Contains(output, `<main class="markdown">`) || !strings.Contains(output, "</main>") {
				t.Error("missing markdown main container")
			}
			for _, legacy := range []string{"PT Mono", "fonts.googleapis.com", "@import", "--paper", "--ink", "--accent", "--code", "--gen-bg", "--font-family", `class="cm"`, "pre code .cm"} {
				if strings.Contains(output, legacy) {
					t.Errorf("output contains legacy style or highlighting %q", legacy)
				}
			}
		})
	}
}
