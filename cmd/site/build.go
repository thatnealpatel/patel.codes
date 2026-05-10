package main

import (
	"bytes"
	"embed"
	"html/template"
	"regexp"

	"rsc.io/markdown"
)

//go:embed templates
var templateFS embed.FS

var pageTmpl = template.Must(template.ParseFS(templateFS, "templates/page.html"))
var galleryTmpl = template.Must(template.ParseFS(templateFS, "templates/gallery.html"))

var reDisplay = regexp.MustCompile(`\$\$(.+?)\$\$`)
var reInline = regexp.MustCompile(`\$(.+?)\$`)

type pageMeta struct {
	Title string
	URL   string
}

func buildPage(src []byte, meta pageMeta) ([]byte, error) {
	var renderErr error

	src = reDisplay.ReplaceAllFunc(src, func(m []byte) []byte {
		if renderErr != nil {
			return m
		}
		expr := bytes.TrimSpace(reDisplay.FindSubmatch(m)[1])
		mathml, err := latexToMathML(expr, true)
		if err != nil {
			renderErr = err
			return m
		}
		return mathml
	})
	if renderErr != nil {
		return nil, renderErr
	}

	src = reInline.ReplaceAllFunc(src, func(m []byte) []byte {
		if renderErr != nil {
			return m
		}
		expr := bytes.TrimSpace(reInline.FindSubmatch(m)[1])
		mathml, err := latexToMathML(expr, false)
		if err != nil {
			renderErr = err
			return m
		}
		return mathml
	})
	if renderErr != nil {
		return nil, renderErr
	}

	var p markdown.Parser
	doc := p.Parse(string(src))
	body := markdown.ToHTML(doc)

	var buf bytes.Buffer
	err := pageTmpl.ExecuteTemplate(&buf, "page.html", struct {
		Body template.HTML
		OG   pageMeta
	}{
		Body: template.HTML(body),
		OG:   meta,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func buildGalleryWall(g gallery) (string, error) {
	gridCSS := "repeat(auto-fit, minmax(calc(var(--page-width) / 4), 1fr))"
	if g.GridCSS != "" {
		gridCSS = g.GridCSS
	}

	var buf bytes.Buffer
	err := galleryTmpl.ExecuteTemplate(&buf, "gallery.html", struct {
		Title       string
		Date        string
		Grid        string
		Zoom        bool
		GridCSSAttr template.CSS
		Images      []string
	}{
		Title:       g.Title,
		Date:        g.Date,
		Grid:        g.Grid,
		Zoom:        g.Zoom,
		GridCSSAttr: template.CSS(gridCSS),
		Images:      g.Images,
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
