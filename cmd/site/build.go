package main

import (
	"bytes"
	"embed"
	"html/template"
	"regexp"
	"strings"

	"rsc.io/markdown"
)

//go:embed templates
var templateFS embed.FS

var pageTmpl = template.Must(template.ParseFS(templateFS, "templates/page.html"))
var galleryTmpl = template.Must(template.ParseFS(templateFS, "templates/gallery.html"))

var reDisplay = regexp.MustCompile(`\$\$(.+?)\$\$`)
var reInline = regexp.MustCompile(`\$(.+?)\$`)

func buildPage(src string) (string, error) {
	var p markdown.Parser
	doc := p.Parse(src)
	body := markdown.ToHTML(doc)

	var renderErr error

	body = reDisplay.ReplaceAllStringFunc(body, func(m string) string {
		if renderErr != nil {
			return m
		}
		expr := strings.TrimSpace(reDisplay.FindStringSubmatch(m)[1])
		mathml, err := latexToMathML(expr, true)
		if err != nil {
			renderErr = err
			return m
		}
		return mathml
	})
	if renderErr != nil {
		return "", renderErr
	}

	body = reInline.ReplaceAllStringFunc(body, func(m string) string {
		if renderErr != nil {
			return m
		}
		expr := strings.TrimSpace(reInline.FindStringSubmatch(m)[1])
		mathml, err := latexToMathML(expr, false)
		if err != nil {
			renderErr = err
			return m
		}
		return mathml
	})
	if renderErr != nil {
		return "", renderErr
	}

	var buf bytes.Buffer
	err := pageTmpl.ExecuteTemplate(&buf, "page.html", struct {
		Body template.HTML
	}{
		Body: template.HTML(body),
	})
	if err != nil {
		return "", err
	}
	return buf.String(), nil
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
