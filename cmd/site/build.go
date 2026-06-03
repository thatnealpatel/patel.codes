package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"regexp"
	"strings"

	"rsc.io/markdown"
)

//go:embed templates
var templateFS embed.FS

var (
	pageTmpl    = template.Must(template.ParseFS(templateFS, "templates/page.html"))
	galleryTmpl = template.Must(template.ParseFS(templateFS, "templates/gallery.html"))
)

var (
	reDisplay    = regexp.MustCompile(`(?s)\$\$(.+?)\$\$`)
	reInline     = regexp.MustCompile(`(?s)\$(.+?)\$`)
	reCodeFence  = regexp.MustCompile("(?s)```.*?```")
	reCodeInline = regexp.MustCompile("`[^`]+`")
)

type pageMeta struct {
	Title string
	URL   string
}

func buildPage(src []byte, meta pageMeta) ([]byte, error) {
	// Phase 1: protect code from dollar-sign matching.
	var codeBlocks [][]byte
	src = reCodeFence.ReplaceAllFunc(src, func(m []byte) []byte {
		id := len(codeBlocks)
		codeBlocks = append(codeBlocks, m)
		return []byte(fmt.Sprintf("CODEPH%04d", id))
	})
	src = reCodeInline.ReplaceAllFunc(src, func(m []byte) []byte {
		id := len(codeBlocks)
		codeBlocks = append(codeBlocks, m)
		return []byte(fmt.Sprintf("CODEPH%04d", id))
	})

	// Phase 2: convert LaTeX to MathML, store results, leave
	// text placeholders that the markdown parser won't mangle.
	type mathEntry struct {
		html    []byte
		display bool
	}
	var mathEntries []mathEntry
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
		id := len(mathEntries)
		mathEntries = append(mathEntries, mathEntry{mathml, true})
		return []byte(fmt.Sprintf("\n\nMATHPH%04d\n\n", id))
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
		id := len(mathEntries)
		mathEntries = append(mathEntries, mathEntry{mathml, false})
		return []byte(fmt.Sprintf("MATHPH%04d", id))
	})
	if renderErr != nil {
		return nil, renderErr
	}

	// Phase 3: restore code blocks so the markdown parser sees them.
	for i, block := range codeBlocks {
		ph := []byte(fmt.Sprintf("CODEPH%04d", i))
		src = bytes.Replace(src, ph, block, 1)
	}

	// Phase 4: parse markdown.
	p := markdown.Parser{Table: true}
	doc := p.Parse(string(src))
	body := markdown.ToHTML(doc)

	// Phase 5: restore MathML in the HTML output.
	for i, entry := range mathEntries {
		ph := fmt.Sprintf("MATHPH%04d", i)
		if entry.display {
			body = strings.Replace(body, "<p>"+ph+"</p>\n", string(entry.html)+"\n", 1)
		}
		body = strings.Replace(body, ph, string(entry.html), 1)
	}

	body = processGenSections(body)

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

func processGenSections(body string) string {
	const open = "<p>:::gen</p>"
	const close = "<p>:::</p>"

	if !strings.Contains(body, open) {
		return body
	}

	var out strings.Builder
	rest := body
	for {
		i := strings.Index(rest, open)
		if i < 0 {
			out.WriteString(rest)
			break
		}
		out.WriteString(rest[:i])
		rest = rest[i+len(open):]
		if len(rest) > 0 && rest[0] == '\n' {
			rest = rest[1:]
		}

		j := strings.Index(rest, close)
		if j < 0 {
			out.WriteString("<div class=\"gen\">\n")
			out.WriteString(rest)
			out.WriteString("</div>\n")
			rest = ""
			break
		}

		out.WriteString("<div class=\"gen\">\n")
		out.WriteString(rest[:j])
		out.WriteString("</div>\n")
		rest = rest[j+len(close):]
		if len(rest) > 0 && rest[0] == '\n' {
			rest = rest[1:]
		}
	}

	result := out.String()
	genIdx := strings.Index(result, "<div class=\"gen\">")
	pIdx := strings.Index(result, "</p>")
	if pIdx >= 0 && pIdx < genIdx {
		at := pIdx + len("</p>")
		if at < len(result) && result[at] == '\n' {
			at++
		}
		notice := "<p class=\"gen-notice\">this page contains generated content delineated using <span class=\"gen-inline\">this shadowing style</span>.</p>\n"
		result = result[:at] + notice + result[at:]
	}

	return result
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
