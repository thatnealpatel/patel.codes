# catching llm-induced laziness

2026-06-02

## intro

as i alluded to in
[an accidentally novel combinatorics proof](/words/an_accidentally_novel_combinatorics_proof_1.html),
i have found myself increasingly less capable
of defending the "it is all slop" stance when
it comes to generative ai technologies.

when it concerns the usage of generative ai tools,
the aphorism **"replace typing, not thinking"**
is one that i try to observe daily;
previously, i did not think there was much cause
for concern since i was not blatantly shipping
slop upstream or destroying the attention of others.

i have since reconsidered the cause for concern
after finding myself in a position where, in
hindsight, i was operating in a "replace thinking"
mode.

i wrote a small quality of (my) life program a
few months ago that simply manages a one-to-many
fan out of a set of files in a standard filesystem.
it largely functions as follows:

```shell
% sync     # dry run
% sync -w  # write
```

where `sync` simply shows a colored list of which
destination files will be `added`/`updated`/`deleted`.
in recent weeks, i have added a few more types of
destination files. this naturally gave rise to the
desire to see file diffs _sometimes,_ so the api
evolved:

```shell
% sync     # view file diff
% sync -v  # view full diff
% sync -w
```

## unintentionally brain off

at the time, i found myself multi-tasking; i was
switching between the repository which hosts this
program and another where i was attempting to prove
some small lemmas that i would later need.

i thought about adding support for full content
diffs while my maths harness was spinning on these
lemmas. i was more interested in what the maths
harness was doing (or rather, not doing), so i simply
opened claude code in the `sync` program's directory
in a feeble attempt to multi-task.

### first attempt (fail: >20 min)

my maths harness was getting close to formalizing
something interesting, so i quickly wrote and sent:

```
Add a universal option to 'sync' -v that
shows the actual content diffs in a git-like
colored view for added/removed.
```

after returning, i noticed that the initial
approach was trying to use an older diff algorithm.
it was handling newlines extremely incorrectly and
subsequent prompts to course correct failed.

in total, this took 20 minutes of human and agent
time costing approximately 85,000 tokens. to be
fair, some of this was wall time, but more than
13 minutes was active human or active agent.

### second attempt (fail: 4 min)

the first attempt left me subtly frustrated.
on one hand, _i did not lose anything in the
failed transaction._ however, on the other,
this was a trivial task... why was the agent
failing at such trivial work?

i put around 20 seconds of active thought
into the second attempt, ensuring that i could
quickly divert my attention back to where i
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

my thought process was that the agent could
figure out how to get from naive, unhelpful,
but correct diffs to a correct, minimal diff
program.

after 4 minutes and 55,000 tokens, the result
was semi-correct and unappealing since it did
not make it further than brute force diffs.
this time i did not bother to course correct.
i was clearly frustrated and uttered the
all-too-expected, "i could just do this myself."

### third attempt (success: <2 min)

while frustrated, it was in this moment that
i realized where i erred. i *thought* about
how i would do this myself. this was a program
that hardly anyone but me uses, and it was
not meant to be hardened against adversarial
inputs in any way.

after spending less than one minute turning
my brain on and going to the `diff` man page,
i realized that this is exactly what i wanted:

```
% diff -u \
<(cat ~/.sync/custom/file) \
<(cat ~/path/to/dst/file)
```

all that remained to materialize the program was
some **annoying typing** to apply ANSI coloring
based on the first rune of each line in the `diff`
output. thus, the final attempt began:

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

shamefully, after 103 seconds and 46,000 tokens,
i got the functionality that i wanted with tests
that made sense.

from a strictly self-centered perspective,
prevailing sentiment seems to indicate that these
tools are powerful when used correctly but still
mostly good, even when used suboptimally. if we
define correct usage to be both useful results
and preservation of the operator's cognition, then
suboptimal usage seems to produce "mostly good"
results at the cost of the operator.

## closing thoughts

while the aphorism **"replace typing, not
thinking"** is certainly a nice one, it is
much harder to abide by in practice than
i previously thought.

on one hand, this ordeal captures about 30 minutes
of my life that i am likely to completely forget
about in the next 48 hours; however, it represents
something that i previously wrote off as "affects
other people, surely."

it was slightly worrying to notice myself
replace thinking about programming, especially
in a context where i was having fun.

without changing anything, i fear that there
might come a day where i am no longer able
to tell when this sort of behavior has occurred.

currently, i do not have a convincing argument
for any productive behavioral change beyond simply
using generative ai tooling less frequently. saying
"but i'll be more careful next time" seems like a
stone's throw away from lying to myself.

it appears that this slope is indeed more slippery
than i previously thought.
