# Catching LLM-induced laziness

2026-06-02

## Intro

As I alluded to in
[An accidentally novel combinatorics proof](/words/an_accidentally_novel_combinatorics_proof_1.html),
I have found myself increasingly less capable
of defending the "it is all slop" stance when
it comes to generative AI technologies.

When it concerns the usage of generative AI tools,
the aphorism **"replace typing, not thinking"**
is one that I try to observe daily;
previously, I did not think there was much cause
for concern since I was not blatantly shipping
slop upstream or destroying the attention of others.

I have since reconsidered the cause for concern
after finding myself in a position where, in
hindsight, I was operating in a "replace thinking"
mode.

I wrote a small quality of (my) life program a
few months ago that simply manages a one-to-many
fan out of a set of files in a standard filesystem.
It largely functions as follows:

```shell
% sync     # dry run
% sync -w  # write
```

where `sync` simply shows a colored list of which
destination files will be `added`/`updated`/`deleted`.
In recent weeks, I have added a few more types of
destination files. This naturally gave rise to the
desire to see file diffs _sometimes,_ so the API
evolved:

```shell
% sync     # view file diff
% sync -v  # view full diff
% sync -w
```

## Unintentionally brain off

At the time, I found myself multi-tasking; I was
switching between the repository which hosts this
program and another where I was attempting to prove
some small lemmas that I would later need.

I thought about adding support for full content
diffs while my maths harness was spinning on these
lemmas. I was more interested in what the maths
harness was doing (or rather, not doing), so I simply
opened Claude Code in the `sync` program's directory
in a feeble attempt to multi-task.

### First attempt (fail: >20 min)

My maths harness was getting close to formalizing
something interesting, so I quickly wrote and sent:

```
Add a universal option to 'sync' -v that
shows the actual content diffs in a git-like
colored view for added/removed.
```

After returning, I noticed that the initial
approach was trying to use an older diff algorithm.
It was handling newlines extremely incorrectly and
subsequent prompts to course correct failed.

In total, this took 20 minutes of human and agent
time costing approximately 85,000 tokens. To be
fair, some of this was wall time, but more than
13 minutes was active human or active agent.

### Second attempt (fail: 4 min)

The first attempt left me subtly frustrated.
On one hand, _I did not lose anything in the
failed transaction._ However, on the other,
this was a trivial task... why was the agent
failing at such trivial work?

I put around 20 seconds of active thought
into the second attempt, ensuring that I could
quickly divert my attention back to where I
wanted:

```
Add a naive, brute force diff view when '-v'
is used with 'sync' for any subcommand.
The diff view should simply show naively what
lines in the dst are getting changed based
on the src.

Write some robust, idiomatic test cases in
sync_test.go before implementing to
straighten out the idea in your head.
```

My thought process was that the agent could
figure out how to get from naive, unhelpful,
but correct diffs to a correct, minimal diff
program.

After 4 minutes and 55,000 tokens, the result
was semi-correct and unappealing since it did
not make it further than brute force diffs.
This time I did not bother to course correct.
I was clearly frustrated and uttered the
all-too-expected, "I could just do this myself."

### Third attempt (success: <2 min)

While frustrated, it was in this moment that
I realized where I erred. I *thought* about
how I would do this myself. This was a program
that hardly anyone but me uses, and it was
not meant to be hardened against adversarial
inputs in any way.

After spending less than one minute turning
my brain on and going to the `diff` man page,
I realized that this is exactly what I wanted:

```
% diff -u \
<(cat ~/.sync/custom/file) \
<(cat ~/path/to/dst/file)
```

All that remained to materialize the program was
some **annoying typing** to apply ANSI coloring
based on the first rune of each line in the `diff`
output. Thus, the final attempt began:

```
Add a universally accepted '-v' flag to 'sync'
which when invoked prints a diff for each of the files
in the sync.

To get the diff text, simply use an exec.Command
on:

'''
% diff -u \
<(cat ~/.sync/custom/file) \
<(cat ~/path/to/dst/file)
'''

The output has markers that can help you color the
output if the ansi colors do not make it through
the exec.Command buffer.
```

Shamefully, after 103 seconds and 46,000 tokens,
I got the functionality that I wanted with tests
that made sense.

From a strictly self-centered perspective,
prevailing sentiment seems to indicate that these
tools are powerful when used correctly but still
mostly good, even when used suboptimally. If we
define correct usage to be both useful results
and preservation of the operator's cognition, then
suboptimal usage seems to produce "mostly good"
results at the cost of the operator.

## Closing thoughts

While the aphorism **"replace typing, not
thinking"** is certainly a nice one, it is
much harder to abide by in practice than
I previously thought.

On one hand, this ordeal captures about 30 minutes
of my life that I am likely to completely forget
about in the next 48 hours; however, it represents
something that I previously wrote off as "affects
other people, surely."

It was slightly worrying to notice myself
replace thinking about programming, especially
in a context where I was having fun.

Without changing anything, I fear that there
might come a day where I am no longer able
to tell when this sort of behavior has occurred.

Currently, I do not have a convincing argument
for any productive behavioral change beyond simply
using generative AI tooling less frequently. Saying
"but I'll be more careful next time" seems like a
stone's throw away from lying to myself.

It appears that this slope is indeed more slippery
than I previously thought.
