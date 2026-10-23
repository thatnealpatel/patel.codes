package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	var (
		root, data, gen string
		err             error
	)
	if root, err = os.Getwd(); err != nil {
		log.Fatal(err)
	}
	data = filepath.Join(root, "data")
	gen = filepath.Join(root, "gen")

	sections := []articleSection{{name: "words", title: "Words"}, {name: "proofs", title: "Proofs"}}
	for i := range sections {
		sections[i].entries, err = os.ReadDir(filepath.Join(data, sections[i].name))
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := validateArticleOutputs(sections); err != nil {
		log.Fatal(err)
	}
	for _, section := range sections {
		if err := generateSectionIndex(data, section); err != nil {
			log.Fatal(err)
		}
	}

	src, err := os.ReadFile(filepath.Join(data, "index.md"))
	if err != nil {
		log.Fatal(err)
	}

	html, err := buildPage(src, pageMeta{
		Title: "patel.codes",
		URL:   "https://patel.codes/",
	})
	if err != nil {
		log.Fatalf("building index: %v", err)
	}

	if err = os.WriteFile(filepath.Join(gen, "index.html"), html, 0o644); err != nil {
		log.Fatal(err)
	}

	fmt.Println("wrote gen/index.html")

	for _, section := range sections {
		if err := generateSectionPages(data, gen, section); err != nil {
			log.Fatal(err)
		}
	}

	galleriesDataDir := filepath.Join(data, "galleries")
	gallerySlugs, err := os.ReadDir(galleriesDataDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, slug := range gallerySlugs {
		if !slug.IsDir() {
			continue
		}
		g, ok := galleryIndex[slug.Name()]
		if !ok {
			log.Printf("warning: no gallery definition for %s, skipping", slug.Name())
			continue
		}

		srcDir := filepath.Join(galleriesDataDir, slug.Name())
		dstDir := filepath.Join(gen, "galleries", slug.Name())
		if err := os.MkdirAll(dstDir, 0o755); err != nil {
			log.Fatal(err)
		}
		if err := cpdir(srcDir, dstDir); err != nil {
			log.Fatal(err)
		}

		wallHTML, err := buildGalleryWall(g)
		if err != nil {
			log.Fatalf("building gallery %s: %v", slug.Name(), err)
		}
		if err := os.WriteFile(filepath.Join(dstDir, "wall.html"), []byte(wallHTML), 0o644); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("wrote gen/galleries/%s/wall.html + %d images\n", slug.Name(), len(g.Images))
	}

	staticDir := filepath.Join(data, "static")
	staticOutDir := filepath.Join(gen, "static")
	if err := os.MkdirAll(staticOutDir, 0o755); err != nil {
		log.Fatal(err)
	}
	if err := cpdir(staticDir, staticOutDir); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote gen/static/")

	cp(filepath.Join(staticDir, "favicon.ico"), filepath.Join(gen, "favicon.ico"))
	fmt.Println("wrote gen/favicon.ico")

	// data/root holds files served verbatim at the site root,
	// e.g. mk.sh.
	if err := cpdir(filepath.Join(data, "root"), gen); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote gen/ (root files)")

	if err := os.WriteFile(filepath.Join(gen, "robots.txt"), []byte(robotsTxt), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote gen/robots.txt")

	if err := os.WriteFile(filepath.Join(gen, "CNAME"), []byte("patel.codes\n"), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote gen/CNAME")

	sitemap := sitemap(sections, gallerySlugs)
	if err := os.WriteFile(filepath.Join(gen, "sitemap.xml"), []byte(sitemap), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote gen/sitemap.xml")

	if err := generateGoImports(gen); err != nil {
		log.Fatal(err)
	}

	const addr = ":9876"
	fmt.Println("serving on", "http://"+addr)
	log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(gen))))
}

type articleSection struct {
	name    string
	title   string
	entries []os.DirEntry
}

func articleOutputPath(name string) string {
	if strings.HasSuffix(name, ".draft.md") {
		return "words/drafts/" + strings.TrimSuffix(name, ".draft.md") + ".html"
	}
	return "words/" + strings.TrimSuffix(name, ".md") + ".html"
}

func validateArticleOutputs(sections []articleSection) error {
	outputs := make(map[string]string)
	for _, section := range sections {
		for _, e := range section.entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			output := articleOutputPath(e.Name())
			source := section.name + "/" + e.Name()
			if previous, ok := outputs[output]; ok {
				return fmt.Errorf("conflicting article output %s: %s and %s", output, previous, source)
			}
			outputs[output] = source
		}
	}
	return nil
}

func generateSectionPages(data, gen string, section articleSection) error {
	srcDir := filepath.Join(data, section.name)

	sectionOutDir := filepath.Join(gen, "words")
	draftsOutDir := filepath.Join(sectionOutDir, "drafts")
	if err := os.MkdirAll(sectionOutDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(draftsOutDir, 0o755); err != nil {
		return err
	}

	for _, e := range section.entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		src, err := os.ReadFile(filepath.Join(srcDir, e.Name()))
		if err != nil {
			return err
		}
		title, _, err := readTitleAndDate(filepath.Join(srcDir, e.Name()))
		if err != nil {
			return err
		}

		urlPath := articleOutputPath(e.Name())

		html, err := buildPage(src, pageMeta{
			Title: title,
			URL:   "https://patel.codes/" + urlPath,
		})
		if err != nil {
			return fmt.Errorf("building %s/%s: %w", section.name, e.Name(), err)
		}
		if err := os.WriteFile(filepath.Join(gen, filepath.FromSlash(urlPath)), html, 0o644); err != nil {
			return err
		}
		fmt.Printf("wrote gen/%s\n", urlPath)
	}

	return nil
}

func cp(src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		log.Fatal(err)
	}
}

const robotsTxt = `User-Agent: *
Disallow:

Sitemap: https://patel.codes/sitemap.xml
`

func sitemap(sections []articleSection, gallerySlugs []os.DirEntry) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
<url><loc>https://patel.codes/</loc></url>
`)
	for _, section := range sections {
		for _, e := range section.entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			if strings.HasSuffix(e.Name(), ".draft.md") {
				continue
			}
			fmt.Fprintf(&b, "<url><loc>https://patel.codes/%s</loc></url>\n", articleOutputPath(e.Name()))
		}
	}
	for _, s := range gallerySlugs {
		if !s.IsDir() {
			continue
		}
		if _, ok := galleryIndex[s.Name()]; !ok {
			continue
		}
		fmt.Fprintf(&b, "<url><loc>https://patel.codes/galleries/%s/wall.html</loc></url>\n", s.Name())
	}
	b.WriteString("</urlset>\n")
	return b.String()
}

type articleEntry struct {
	title string
	file  string
	date  string
}

func generateSectionIndex(dataDir string, section articleSection) error {
	sectionDir := filepath.Join(dataDir, section.name)

	var articles []articleEntry
	for _, e := range section.entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || strings.HasSuffix(e.Name(), ".draft.md") {
			continue
		}
		title, date, err := readTitleAndDate(filepath.Join(sectionDir, e.Name()))
		if err != nil {
			return err
		}
		articles = append(articles, articleEntry{title: title, file: e.Name(), date: date})
	}

	sort.Slice(articles, func(i, j int) bool {
		return articles[i].date > articles[j].date
	})

	indexPath := filepath.Join(dataDir, "index.md")
	src, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	header := []byte("## " + section.title + "\n")
	start := bytes.Index(src, header)
	if start == -1 {
		return fmt.Errorf("index.md: missing ## %s section", section.title)
	}
	start += len(header)

	rest := src[start:]
	end := bytes.Index(rest, []byte("\n## "))
	if end == -1 {
		return fmt.Errorf("index.md: missing section after ## %s", section.title)
	}

	var buf bytes.Buffer
	buf.Write(src[:start])
	buf.WriteByte('\n')
	for _, w := range articles {
		fmt.Fprintf(&buf, "- [%s](./%s)\n", w.title, articleOutputPath(w.file))
	}
	buf.Write(rest[end:])

	return os.WriteFile(indexPath, buf.Bytes(), 0o644)
}

func readTitleAndDate(path string) (title, date string, err error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return "", "", err
	}
	lines := bytes.SplitN(src, []byte("\n"), 4)
	if len(lines) < 3 {
		return "", "", fmt.Errorf("%s: expected at least 3 lines", path)
	}
	title = string(bytes.TrimPrefix(lines[0], []byte("# ")))
	date = string(bytes.TrimSpace(lines[2]))
	return title, date, nil
}

type goImport struct {
	repo   string
	subdir string
}

// goImportOverrides describes modules that
// exist on my GitHub that should always be
// mapped to name import pages.
var goImportOverrides = map[string]goImport{
	"proofs":      {repo: "thatnealpatel/proofs"},
	"patel.codes": {repo: "thatnealpatel/patel.codes"},
	"jsonldb":     {repo: "thatnealpatel/mono", subdir: "jsonldb"},
	"ranking":     {repo: "thatnealpatel/mono", subdir: "ranking"},
	"unsafe":      {repo: "thatnealpatel/mono", subdir: "unsafe"},
	"gh":          {repo: "thatnealpatel/mono", subdir: "gh"}, // TODO: remove.
	"gerrit":      {repo: "thatnealpatel/mono", subdir: "gerrit"},
	"srvfs":       {repo: "thatnealpatel/mono", subdir: "srvfs"},
	"ghdl":        {repo: "thatnealpatel/mono", subdir: "ghdl"},
	"rfc":         {repo: "thatnealpatel/mono", subdir: "rfc"},
	"whatwg":      {repo: "thatnealpatel/mono", subdir: "whatwg"},
	"oeis":        {repo: "thatnealpatel/mono", subdir: "oeis"},
	"wiki":        {repo: "thatnealpatel/mono", subdir: "wiki"},
	"erdos":       {repo: "thatnealpatel/mono", subdir: "erdos"},
	"cve":         {repo: "thatnealpatel/mono", subdir: "cve"},
	"leandoc":     {repo: "thatnealpatel/mono", subdir: "leandoc"},
	"plsno":       {repo: "thatnealpatel/mono", subdir: "plsno"}, // TODO: rename.
	"fetch":       {repo: "thatnealpatel/mono", subdir: "fetch"},
	"ghfa":        {repo: "thatnealpatel/mono", subdir: "ghfa"},
	"retrieval":   {repo: "thatnealpatel/mono", subdir: "retrieval"},
	"render":      {repo: "thatnealpatel/mono", subdir: "render"},
	"gocmt":       {repo: "thatnealpatel/mono", subdir: "gocmt"},
	"indexing":    {repo: "thatnealpatel/mono", subdir: "indexing"},
	"sagedoc":     {repo: "thatnealpatel/mono", subdir: "sagedoc"},
	"grfa":        {repo: "thatnealpatel/mono", subdir: "grfa"},
}

// generateGoImports uses a lazy method to
// derive any new Go repositories that need
// their own generated vanity page.
func generateGoImports(gen string) error {
	resp, err := http.Get("https://api.github.com/users/thatnealpatel/repos?per_page=100")
	if err != nil {
		return fmt.Errorf("fetching github repos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github api: %s", resp.Status)
	}

	var repos []struct {
		Name     string `json:"name"`
		Fork     bool   `json:"fork"`
		Language string `json:"language"`
		License  any    `json:"license"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return fmt.Errorf("decoding github repos: %w", err)
	}

	written := map[string]bool{
		// We need to pre-seed mono to ensure
		// that it does not get auto-generated
		// as a Go module since it technically
		// does meet all the critiera.
		"mono": true,
	}

	for name, imp := range goImportOverrides {
		dir := filepath.Join(gen, name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		html := goImportPageGitHub(name, imp.repo, imp.subdir)
		if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(html), 0o644); err != nil {
			return err
		}
		fmt.Printf("wrote gen/%s/index.html (go-import, override)\n", name)
		written[name] = true
	}

	for _, repo := range repos {
		if written[repo.Name] {
			continue
		}
		if repo.Fork || (repo.Language != "Go" && repo.Language != "Go Template") || repo.License == nil {
			continue
		}
		dir := filepath.Join(gen, repo.Name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		html := goImportPageGitHub(repo.Name, "thatnealpatel/"+repo.Name, "")
		if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(html), 0o644); err != nil {
			return err
		}
		fmt.Printf("wrote gen/%s/index.html (go-import)\n", repo.Name)
	}
	return nil
}

func goImportPageGitHub(module, ghRepo, subdir string) string {
	extra := ""
	if subdir != "" {
		extra = " " + subdir
	}
	return `<!DOCTYPE html>
<html>
<head>
<meta name="go-import" content="patel.codes/` + module + ` git https://github.com/` + ghRepo + extra + `">
<meta http-equiv="refresh" content="0; url=https://pkg.go.dev/patel.codes/` + module + `">
</head>
<body>
Redirecting to <a href="https://pkg.go.dev/patel.codes/` + module + `">pkg.go.dev/patel.codes/` + module + `</a>...
</body>
</html>
`
}

func cpdir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == ".DS_Store" {
			return nil
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, rel)

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		out, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		return err
	})
}
