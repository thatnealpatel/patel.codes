package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func readTestFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func testArticleSections(t *testing.T, data string) []articleSection {
	t.Helper()
	sections := []articleSection{{name: "words", title: "Words"}, {name: "proofs", title: "Proofs"}}
	for i := range sections {
		dir := filepath.Join(data, sections[i].name)
		writeTestFile(t, filepath.Join(dir, sections[i].name+"-a-old.md"), "# Old\n\n2025-01-01\n\nOld body.\n")
		writeTestFile(t, filepath.Join(dir, sections[i].name+"-z-new.md"), "# New\n\n2026-01-01\n\nNew body.\n")
		writeTestFile(t, filepath.Join(dir, sections[i].name+"-hidden.draft.md"), "# Hidden\n\n2027-01-01\n\nDraft body.\n")
		writeTestFile(t, filepath.Join(dir, "ignored.txt"), "not markdown")
		writeTestFile(t, filepath.Join(dir, "nested.md", "inside.md"), "not a top-level article")
		var err error
		sections[i].entries, err = os.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
	}
	return sections
}

func TestGenerateSectionIndexes(t *testing.T) {
	data := t.TempDir()
	sections := testArticleSections(t, data)
	const prefix = "# Test\n\n## About\n\nKeep me.\n\n"
	const suffix = "## Galleries\n\n- Keep gallery\n\n## Connect\n\nKeep contact.\n"
	writeTestFile(t, filepath.Join(data, "index.md"), prefix+"## Words\n\nStale words.\n\n## Proofs\n\nStale proofs.\n\n"+suffix)
	for _, section := range sections {
		if err := generateSectionIndex(data, section); err != nil {
			t.Fatal(err)
		}
	}
	want := prefix
	for _, section := range sections {
		want += "## " + section.title + "\n\n- [New](./words/" + section.name + "-z-new.html)\n- [Old](./words/" + section.name + "-a-old.html)\n\n"
	}
	want += suffix
	if got := readTestFile(t, filepath.Join(data, "index.md")); got != want {
		t.Fatalf("index = %q, want %q", got, want)
	}
	// Rebuilding either section must leave
	// the other section and surrounding
	// text intact.
	for _, section := range sections {
		if err := generateSectionIndex(data, section); err != nil {
			t.Fatal(err)
		}
		if got := readTestFile(t, filepath.Join(data, "index.md")); got != want {
			t.Fatalf("index changed on rebuild: %q", got)
		}
	}
}

func TestGenerateSectionPages(t *testing.T) {
	data, gen := t.TempDir(), t.TempDir()
	sections := testArticleSections(t, data)
	vanity := goImportPageGitHub("proofs", "thatnealpatel/proofs", "")
	writeTestFile(t, filepath.Join(gen, "proofs", "index.html"), vanity)
	if err := validateArticleOutputs(sections); err != nil {
		t.Fatal(err)
	}
	for _, section := range sections {
		if err := generateSectionPages(data, gen, section); err != nil {
			t.Fatal(err)
		}
		for _, page := range []struct{ path, title, body string }{
			{section.name + "-a-old.html", "Old", "Old body."},
			{section.name + "-z-new.html", "New", "New body."},
			{"drafts/" + section.name + "-hidden.html", "Hidden", "Draft body."},
		} {
			html := readTestFile(t, filepath.Join(gen, "words", page.path))
			for _, want := range []string{
				"<title>" + page.title + "</title>",
				`<meta property="og:url" content="https://patel.codes/words/` + page.path + `">`,
				page.body,
			} {
				if !strings.Contains(html, want) {
					t.Errorf("%s/%s missing %q", section.name, page.path, want)
				}
			}
		}
		for _, excluded := range []string{section.name + "-hidden.html", section.name + "-hidden.draft.html", "ignored.html", "ignored.txt", "nested.html", "nested.md"} {
			if _, err := os.Stat(filepath.Join(gen, "words", excluded)); !os.IsNotExist(err) {
				t.Errorf("unexpected output %s/%s: %v", section.name, excluded, err)
			}
		}
	}
	entries, err := os.ReadDir(filepath.Join(gen, "proofs"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || entries[0].Name() != "index.html" {
		t.Fatalf("unexpected output in proofs vanity directory: %v", entries)
	}
	if got := readTestFile(t, filepath.Join(gen, "proofs", "index.html")); got != vanity {
		t.Fatal("proofs vanity index changed")
	}
}

func TestSitemapArticleSections(t *testing.T) {
	data := t.TempDir()
	sections := testArticleSections(t, data)
	galleries := filepath.Join(data, "galleries")
	for _, slug := range []string{"untitled-no1", "unknown"} {
		if err := os.MkdirAll(filepath.Join(galleries, slug), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeTestFile(t, filepath.Join(galleries, "ignored.txt"), "ignore")
	entries, err := os.ReadDir(galleries)
	if err != nil {
		t.Fatal(err)
	}
	got := sitemap(sections, entries)
	want := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">\n<url><loc>https://patel.codes/</loc></url>\n"
	for _, section := range sections {
		for _, name := range []string{"a-old", "z-new"} {
			want += "<url><loc>https://patel.codes/words/" + section.name + "-" + name + ".html</loc></url>\n"
		}
	}
	want += "<url><loc>https://patel.codes/galleries/untitled-no1/wall.html</loc></url>\n</urlset>\n"
	if got != want {
		t.Fatalf("sitemap = %q, want %q", got, want)
	}
}

func TestValidateArticleOutputs(t *testing.T) {
	for _, tc := range []struct {
		name, wordsFile, proofsFile, conflict string
	}{
		{"published conflict", "same.md", "same.md", "words/same.html"},
		{"draft conflict", "same.draft.md", "same.draft.md", "words/drafts/same.html"},
		{"published and draft", "same.md", "same.draft.md", ""},
		{"different filenames", "first.md", "second.md", ""},
		{"non-markdown ignored", "same.txt", "same.txt", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := t.TempDir()
			sections := []articleSection{{name: "words"}, {name: "proofs"}}
			for i, file := range []string{tc.wordsFile, tc.proofsFile} {
				dir := filepath.Join(data, sections[i].name)
				writeTestFile(t, filepath.Join(dir, file), "# Article\n\n2026-01-01\n")
				var err error
				sections[i].entries, err = os.ReadDir(dir)
				if err != nil {
					t.Fatal(err)
				}
			}
			err := validateArticleOutputs(sections)
			if tc.conflict == "" {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected conflicting output error")
			}
			for _, want := range []string{"conflicting article output", tc.conflict, "words/" + tc.wordsFile, "proofs/" + tc.proofsFile} {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("error %q missing %q", err, want)
				}
			}
		})
	}
}
