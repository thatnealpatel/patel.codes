package main

import (
	"bytes"
	"embed"
	"html/template"
	"strings"

	"patel.codes/render"
)

//go:embed templates
var templateFS embed.FS

var (
	pageTmpl    = template.Must(template.ParseFS(templateFS, "templates/page.html"))
	galleryTmpl = template.Must(template.ParseFS(templateFS, "templates/gallery.html"))
)

type pageMeta struct {
	Title string
	URL   string
}

func buildPage(src []byte, meta pageMeta) ([]byte, error) {
	body, err := render.Render(string(src))
	if err != nil {
		return nil, err
	}

	bodyHTML := string(body)

	// One-time exception: the disclosure post explains the shadowing style itself,
	// so the auto-injected notice would be redundant.
	skipNotice := strings.Contains(meta.URL, "/words/patel_codes_llm_disclosure.html")
	bodyHTML = processGenSections(bodyHTML, skipNotice)
	bodyHTML = highlightComments(bodyHTML)

	if meta.URL != "https://patel.codes/" {
		bodyHTML = strings.Replace(bodyHTML, "</p>", ` <span class="home">[<a href="/">home</a>]</span></p>`, 1)
	}

	var buf bytes.Buffer
	err = pageTmpl.ExecuteTemplate(&buf, "page.html", struct {
		Body template.HTML
		OG   pageMeta
	}{
		Body: template.HTML(bodyHTML),
		OG:   meta,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func processGenSections(body string, skipNotice bool) string {
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
	if skipNotice {
		return result
	}
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
