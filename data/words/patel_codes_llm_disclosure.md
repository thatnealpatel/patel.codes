# LLM disclosure on patel.codes

2026-08-21

## Motivations

Unfortunately, LLMs appear to be here to stay;
I have found great utility in various side
projects and research efforts. Especially in
the later, I find myself conflicted about how
to honestly communicate results.

While LLMs are horrendus at prose, they are
undeniably useful for programming and maths.
I never thought I would publish LLM-generated
text on my website until
[An accidentally novel combinatorics proof](/words/an_accidentally_novel_combinatorics_proof_1.html).

I had to find a way to disclose the use of
LLM-generated text: I hate reading it, and I
believe it is unfair of me to subject others
to it as well without their knowledge.

I do not want people guessing about if my
words are my own. I want them focused on the
content.

The idealist in me says "I should never publish
anything without having fully written everything
myself." However, the realist in me has to contend
with trading off time spent doing `<cool thing>`
with time spent communicating `<cool thing>`.

The aforementioned post is a good baseline against
which to compare its sharpened counterpart:
[What I actually proved about `A051293`](/words/sharpening_my_a051293_results.html).

## `:::gen`

Any disclosed, generated content does not mean
"untrusted;" instead, it indicates that I did
not author it, but I read and verified it.

My site generator is completely bespoke machinery 
instrumented in Go, so I just added a simple
syntax to fence generated content:

```
:::gen

**This is text I wrote but marked generated.**

:::
```

which renders as

:::gen

**This is text I wrote but marked generated.**

:::

and includes the following banner at the top of the page

**This page contains generated content delineated using
<span class="gen-inline">this shadowing style</span>.**

## Closing thoughts

With several maths results in my pipeline blocked
on my bandwidth to elucidate and review the results,
it is obvious to me that I need to lean (no pun
intended) on the LLMs in lower risk contexts.

Personally, I hope to keep generated content to
<30% of maths posts I publish; for non-maths or
non-technical content, I imagine this will be
closer to <10%.

I hope to continue writing and publishing things
that do not require any LLM-generated text at all.
