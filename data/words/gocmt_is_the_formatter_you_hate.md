# `gocmt` is the formatter you hate

2026-10-20

## Motivations

For years, I would manually tidy my own comments
according to some esoteric spec that lived in my
head. While I generally avoid doing this in the
Go standard library, I know there exists a small
subset of folks who enjoy this, let's call it,
*craft*.

While agents generally *suck* at writing comments,
they are getting better. However, in an ever growing
catalog of personal software, comments-as-a-spec is
helpful for autonomous development. And, I am
still of the opinion that I want to look at code
agents write. The only difference is **how** much
time I spend. In either case, I loathe looking at
both vacuous and ugly formatting.

So, [`gocmt`](https://patel.codes/gocmt) and `gofmt`
were seperated at birth.

## The semantics

Similarly to `gofmt`, the interface is simple:

```shell
% gocmt -d main.go  # inspect
% gocmt -w main.go  # apply
```

From the program itself:

```go
// process returns src with every eligible comment group
// rewrapped; bytes outside comments are never modified.
func process(name string, src []byte) ([]byte, error) {
```

Standalone `//` prose gets rewrapped. Package docs,
comments beginning on line one, trailing comments,
and block comments are left alone for various reasons.
We skip generated files; within blocks in scope, we
similarly passthrough directives, headings, indented
code examples, etc.

For example,

```go
package p

// P is tiny but its comment still gets exactly forty columns of width.
type P struct {
}
```

rewrites to

```go
package p

// P is tiny but its comment still gets
// exactly forty columns of width.
type P struct {
}
```

The next non-blank line anchors the target width,
clamped to 40–80 columns. There is an optimizer
that runs a few different configurations and picks
a local max. Words themselves remain intact.

Short statements are an exception. Despite overhanging
its `{`, the following comment stays put:

```go
// String renders m as its mode name for logs.
func (m Mode) String() string {
}
```

The relevant branch, again from the source:

```go
if joined := strings.Join(ch.words, " "); utf8.RuneCountInString(joined)+pad <= oneLineCols {
    lines = append(lines, "// "+joined)
    continue
}
```

`oneLineCols` is 60; `pad` counts indentation and
`// `. Otherwise, the program tries widths within
two columns of the target and minimizes raggedness
and a penalty for producing a very short final line.

As documented, this is a rather unserious program
that I unironically use. I do not care to convince
you that comments should be written this way: They
probably shouldn't. But I like them like this.

## Closing thoughts

You can install `gocmt`:

```
go install patel.codes/gocmt@latest
```

And, `gocmt`'s package doc proudly reads:

```
// Package main implements a program that
// codifies my idiosyncratic comment prose
// style that most people surely dislike.
//
// I like writing comments that minimize
// the line length difference in a given
// comment block while commonly choosing
// the line length based on the location
// of the `{` on the following line. An
// exception is made when a single-line
// comment can fit in 60 columns or fewer.
//
// All non-terminal comment lines must be
// at least 40 columns wide. The algorithm
// tries a few local orientations and
// picks the global block-optimal.
//
// `gocmt` is a rather unserious program
// that I unironically use.
package main
```
