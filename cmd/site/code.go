package main

import (
	"regexp"
	"strings"
)

// commentSyntax describes the comment delimiters for one language.
// A multi or doc pair is only used when both its start and end are
// nonempty; an incomplete pair is silently ignored.
type commentSyntax struct {
	singleStart          string
	multiStart, multiEnd string
	docStart, docEnd     string
}

// commentTokens is the opt-in allowlist mapping a fenced block's
// ```<lang> info string to its comment delimiters. Languages not
// listed get no comment highlighting. The scanner runs over
// HTML-escaped code text, so any delimiter containing & < > "
// must be stored pre-escaped.
var commentTokens = map[string]commentSyntax{
	"lean": {
		singleStart: "--",
		multiStart:  "/-", multiEnd: "-/",
		docStart: "/--", docEnd: "-/",
	},
	"go": {
		singleStart: "//",
		multiStart:  "/*", multiEnd: "*/",
	},
}

// The code text is HTML-escaped, so a literal </code></pre> cannot
// appear inside the block and the non-greedy match is safe.
var reCodeBlock = regexp.MustCompile(`(?s)<pre><code class="language-([^"]+)">(.*?)</code></pre>`)

// highlightComments wraps comment runs inside fenced code blocks
// in <span class="cm"> for languages listed in commentTokens.
func highlightComments(body string) string {
	return reCodeBlock.ReplaceAllStringFunc(body, func(m string) string {
		sub := reCodeBlock.FindStringSubmatch(m)
		cs, ok := commentTokens[sub[1]]
		if !ok {
			return m
		}
		return `<pre><code class="language-` + sub[1] + `">` +
			spanComments(sub[2], cs) + `</code></pre>`
	})
}

// spanComments scans HTML-escaped code text and wraps each comment
// in <span class="cm">. At any position a doc comment wins over a
// multi comment (its start is a prefix, e.g. Lean /-- vs /-), which
// wins over a single comment. Single comments run to end of line;
// multi and doc comments run to their end token, or to the end of
// the block when unterminated. Nested block comments are not
// tracked.
func spanComments(code string, cs commentSyntax) string {
	multiOK := cs.multiStart != "" && cs.multiEnd != ""
	docOK := cs.docStart != "" && cs.docEnd != ""
	var b strings.Builder
	for i := 0; i < len(code); {
		var start, end string
		switch {
		case docOK && strings.HasPrefix(code[i:], cs.docStart):
			start, end = cs.docStart, cs.docEnd
		case multiOK && strings.HasPrefix(code[i:], cs.multiStart):
			start, end = cs.multiStart, cs.multiEnd
		case cs.singleStart != "" && strings.HasPrefix(code[i:], cs.singleStart):
			start, end = cs.singleStart, "\n"
		default:
			b.WriteByte(code[i])
			i++
			continue
		}
		j := len(code)
		if k := strings.Index(code[i+len(start):], end); k >= 0 {
			j = i + len(start) + k + len(end)
			if end == "\n" {
				j-- // keep the newline outside the span
			}
		}
		b.WriteString(`<span class="cm">`)
		b.WriteString(code[i:j])
		b.WriteString(`</span>`)
		i = j
	}
	return b.String()
}
