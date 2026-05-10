# patel.codes is now `<noscript>`

2026-05-10

## motivations

recently, i set a goal to more actively write content
for my website. in coming back to the repository, the
friction for writing content was unbearable. while i
initially set out for simplicity, i somehow landed in
the worst of all worlds:

1. javascript dependencies (mathjax)
2. handwriting html files
3. handwriting galleries
4. forgetting to update sitemap.xml

i wanted to keep everything minimal but functional.

## sprinkle in a little bit of slop

as a go programmer, i immediately reached for a minimal
binary that would use `rsc.io/markdown` to convert my
writing into html. similarly, a single go html template
solves the gallery toil.

i also wanted to remove mathjax. this seemed like the
perfect task for an llm:

1. self-contained
2. easily verifiable
3. replaced typing burden

i booted up claude code. after exploring a few go LaTeX
parsers, i simply told it to generate a recursive descent
parser for LaTeX to preprocess before converting to html
with mathml. it happily did so with two bugs in all of 83
seconds. another couple minutes of active thought and a
usable parser eliminated both my writing friction and gave
me a `<noscript>` website.

the only downside here is that i felt the need to write
the following into my README:

> *There is no promise anything here will be maintained,*
> *will continue compiling, or ever compiled.*

## `./gen`

now, it is easier than ever for me to write content and
modify my site. i can simply stream my consciousness into
`.drafts/` and programmatically generate and preview my
site before launching: `go run ./cmd/site`.

i really tried to avoid github actions for various reasons;
however, i settled for the following, notably excluding any
execution of the go binary for deployment:

1. checkout
2. copy `gen` dir
3. deploy the site

[source code for patel.codes is on my github.](https://github.com/thatnealpatel/patel.codes)
