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

	if err := generateWordsIndex(data); err != nil {
		log.Fatal(err)
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

	wordsDir := filepath.Join(data, "words")
	entries, err := os.ReadDir(wordsDir)
	if err != nil {
		log.Fatal(err)
	}

	wordsOutDir := filepath.Join(gen, "words")
	draftsOutDir := filepath.Join(wordsOutDir, "drafts")
	if err := os.MkdirAll(wordsOutDir, 0o755); err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(draftsOutDir, 0o755); err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		src, err := os.ReadFile(filepath.Join(wordsDir, e.Name()))
		if err != nil {
			log.Fatal(err)
		}
		title, _, err := readTitleAndDate(filepath.Join(wordsDir, e.Name()))
		if err != nil {
			log.Fatal(err)
		}

		draft := strings.HasSuffix(e.Name(), ".draft.md")
		var outName, outDir, urlPath string
		if draft {
			outName = strings.TrimSuffix(e.Name(), ".draft.md") + ".html"
			outDir = draftsOutDir
			urlPath = "words/drafts/" + outName
		} else {
			outName = strings.TrimSuffix(e.Name(), ".md") + ".html"
			outDir = wordsOutDir
			urlPath = "words/" + outName
		}

		html, err := buildPage(src, pageMeta{
			Title: title,
			URL:   "https://patel.codes/" + urlPath,
		})
		if err != nil {
			log.Fatalf("building %s: %v", e.Name(), err)
		}
		if err := os.WriteFile(filepath.Join(outDir, outName), html, 0o644); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("wrote gen/%s\n", urlPath)
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

	if err := os.WriteFile(filepath.Join(gen, "robots.txt"), []byte(robotsTxt), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote gen/robots.txt")

	if err := os.WriteFile(filepath.Join(gen, "CNAME"), []byte("patel.codes\n"), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote gen/CNAME")

	sitemap := sitemap(entries, gallerySlugs)
	if err := os.WriteFile(filepath.Join(gen, "sitemap.xml"), []byte(sitemap), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Println("wrote gen/sitemap.xml")

	if err := generateGoImports(gen); err != nil {
		log.Fatal(err)
	}

	const addr = "localhost:9876"
	fmt.Println("serving on", "http://"+addr)
	log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(gen))))
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

func sitemap(blogEntries []os.DirEntry, gallerySlugs []os.DirEntry) string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
<url><loc>https://patel.codes/</loc></url>
`)
	for _, e := range blogEntries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		if strings.HasSuffix(e.Name(), ".draft.md") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".md") + ".html"
		fmt.Fprintf(&b, "<url><loc>https://patel.codes/words/%s</loc></url>\n", name)
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

type wordEntry struct {
	title string
	file  string
	date  string
}

func generateWordsIndex(dataDir string) error {
	wordsDir := filepath.Join(dataDir, "words")
	entries, err := os.ReadDir(wordsDir)
	if err != nil {
		return err
	}

	var words []wordEntry
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") || strings.HasSuffix(e.Name(), ".draft.md") {
			continue
		}
		title, date, err := readTitleAndDate(filepath.Join(wordsDir, e.Name()))
		if err != nil {
			return err
		}
		words = append(words, wordEntry{title: title, file: e.Name(), date: date})
	}

	sort.Slice(words, func(i, j int) bool {
		return words[i].date > words[j].date
	})

	indexPath := filepath.Join(dataDir, "index.md")
	src, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	header := []byte("## words\n")
	start := bytes.Index(src, header)
	if start == -1 {
		return fmt.Errorf("index.md: missing ## words section")
	}
	start += len(header)

	rest := src[start:]
	end := bytes.Index(rest, []byte("\n## "))
	if end == -1 {
		return fmt.Errorf("index.md: missing section after ## words")
	}

	var buf bytes.Buffer
	buf.Write(src[:start])
	buf.WriteByte('\n')
	for _, w := range words {
		htmlName := strings.TrimSuffix(w.file, ".md") + ".html"
		fmt.Fprintf(&buf, "- [%s](./words/%s)\n", w.title, htmlName)
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

// goImportOverrides maps module names to GitHub repos that should always
// get a vanity import page, regardless of the GitHub API heuristic.
var goImportOverrides = map[string]string{
	"proofs": "thatnealpatel/proofs",
	"patel.codes": "thatnealpatel/patel.codes",
}

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
		Name     string      `json:"name"`
		Fork     bool        `json:"fork"`
		Language string      `json:"language"`
		License  interface{} `json:"license"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return fmt.Errorf("decoding github repos: %w", err)
	}

	written := make(map[string]bool)

	for name, ghRepo := range goImportOverrides {
		dir := filepath.Join(gen, name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		html := goImportPageGitHub(name, ghRepo)
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
		html := goImportPageGitHub(repo.Name, "thatnealpatel/"+repo.Name)
		if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(html), 0o644); err != nil {
			return err
		}
		fmt.Printf("wrote gen/%s/index.html (go-import)\n", repo.Name)
	}
	return nil
}

func goImportPageGitHub(module, ghRepo string) string {
	return `<!DOCTYPE html>
<html>
<head>
<meta name="go-import" content="patel.codes/` + module + ` git https://github.com/` + ghRepo + `">
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
