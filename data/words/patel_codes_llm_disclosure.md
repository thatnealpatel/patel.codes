# llm disclosure on patel.codes

2026-08-21

## motivations

unfortunately, LLMs appear to be here to stay;
i have found great utility in various side
projects and research efforts. especially in
the later, i find myself conflicted about how
to honestly communicate results.

while LLMs are horrendus at prose, they are
undeniably useful for programming and maths.
i never thought i would publish LLM-generated
text on my website until
[an accidentally novel combinatorics proof](/words/an_accidentally_novel_combinatorics_proof_1.html).

i had to find a way to disclose the use of
LLM-generated text: i hate reading it, and i
believe it is unfair of me to subject others
to it as well without their knowledge.

i do not want people guessing about if my
words are my own. i want them focused on the
content.

the idealist in me says "i should never publish
anything without having fully written everything
myself." however, the realist in me has to contend
with trading off time spent doing `<cool thing>`
with time spent communicating `<cool thing>`.

the aforementioned post is a good baseline against
which to compare its sharpened counterpart:
[what i actually proved about `A051293`](/words/sharpening_my_a051293_results.html).

## `:::gen`

any disclosed, generated content does not mean
"untrusted;" instead, it indicates that i did
not author it, but i read and verified it.

my site generator is completely bespoke machinery 
instrumented in Go, so i just added a simple
syntax to fence generated content:

```
:::gen

**this is text i wrote but marked generated.**

:::
```

which renders as

:::gen

**this is text i wrote but marked generated.**

:::

and includes the following banner at the top of the page

**this page contains generated content delineated using
<span class="gen-inline">this shadowing style</span>.**

## closing thoughts

with several maths results in my pipeline blocked
on my bandwidth to elucidate and review the results,
it is obvious to me that i need to lean (no pun
intended) on the LLMs in lower risk contexts.

personally, i hope to keep generated content to
<30% of maths posts i publish; for non-maths or
non-technical content, i imagine this will be
closer to <10%.

i hope to continue writing and publishing things
that do not require any LLM-generated text at all.
