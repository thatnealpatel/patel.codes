# patel.codes is now `<noscript>`

2026-05-10

## Motivations

Recently, I set a goal to more actively write content
for my website. In coming back to the repository, the
friction for writing content was unbearable. While I
initially set out for simplicity, I somehow landed in
the worst of all worlds:

1. JavaScript dependencies (MathJax)
2. Handwriting HTML files
3. Handwriting galleries
4. Forgetting to update sitemap.xml

I wanted to keep everything minimal but functional.

## Sprinkle in a little bit of slop

As a Go programmer, I immediately reached for a minimal
binary that would use `rsc.io/markdown` to convert my
writing into HTML. Similarly, a single Go HTML template
solves the gallery toil.

I also wanted to remove MathJax. This seemed like the
perfect task for an LLM:

1. Self-contained
2. Easily verifiable
3. Replaced typing burden

I booted up Claude Code. After exploring a few Go LaTeX
parsers, I simply told it to generate a recursive descent
parser for LaTeX to preprocess before converting to HTML
with MathML. It happily did so with two bugs in all of 83
seconds. Another couple minutes of active thought and a
usable parser eliminated both my writing friction and gave
me a `<noscript>` website.

The only downside here is that I felt the need to write
the following into my README:

> *There is no promise anything here will be maintained,*
> *will continue compiling, or ever compiled.*

## `./gen`

Now, it is easier than ever for me to write content and
modify my site. I can simply stream my consciousness into
`.drafts/` and programmatically generate and preview my
site before launching: `go run ./cmd/site`.

I really tried to avoid GitHub Actions for various reasons;
however, I settled for the following, notably excluding any
execution of the Go binary for deployment:

1. Checkout
2. Copy `gen` dir
3. Deploy the site

[Source code for patel.codes is on my GitHub.](https://github.com/thatnealpatel/patel.codes)
