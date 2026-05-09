package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	src, err := os.ReadFile(filepath.Join(data, "index.md"))
	if err != nil {
		log.Fatal(err)
	}

	html, err := buildPage(string(src))
	if err != nil {
		log.Fatalf("building index: %v", err)
	}

	if err = os.WriteFile(filepath.Join(gen, "index.html"), []byte(html), 0o644); err != nil {
		log.Fatal(err)
	}

	fmt.Println("wrote gen/index.html")

	wordsDir := filepath.Join(data, "words")
	entries, err := os.ReadDir(wordsDir)
	if err != nil {
		log.Fatal(err)
	}

	wordsOutDir := filepath.Join(gen, "words")
	if err := os.MkdirAll(wordsOutDir, 0o755); err != nil {
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
		html, err := buildPage(string(src))
		if err != nil {
			log.Fatalf("building %s: %v", e.Name(), err)
		}
		outName := strings.TrimSuffix(e.Name(), ".md") + ".html"
		if err := os.WriteFile(filepath.Join(wordsOutDir, outName), []byte(html), 0o644); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("wrote gen/words/%s\n", outName)
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
