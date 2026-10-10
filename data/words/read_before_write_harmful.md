# Read-before-write is harmful?

2026-10-10

## Motivations

A colleague recently shared [this study with us](https://arxiv.org/abs/2609.40303)
which prompted me to revise this draft and finally
get it out the door. While this particular study
looks at autonomous ML engineering, it does try to
formalize something that I have long suspected and
anecdotally confirmed.

We are presently in an era where quantitative
arguments are difficult to coalesce faithfully:
*It is all about measuring the vibes.* I make
no claim of sound scientific methodology here.
Instead, questions to consider. In an age where
you can *do* more, it is important to take a step
to assess whether or not you should.

Working as a maintainer of the Go project has
significantly re-framed my thinking and behaviors
as a programmer. In many ways, there is still
yet a chasm of expertise that seemingly eludes
me; however, presently, I have increasingly found
it more important to apply the ethos of *prefer
simplicity*.

## My harness

A harness usually implements a `write` and `edit`
tool where the former does an atomic file replace
with content and the latter usually specifies a
set of edits to make using a diff syntax. The
implementation details are not that important. To
ensure that generated content matches what is in
a file *and* that no unsynchronized edits are going
to happen, a `read`-before-`write` check ensures
that the harness can provide feedback to the agent
that something may have changed or that the supplied
edit is incorrect in some way.

I continue to hack on top of a bespoke harness
that is largely a research project that I just
happen to use as my daily driver harness. This
program operates inside of the trust model and
intranet described in [delegated sandboxing strategy](/words/my_thoughts_on_sandboxing_agents.html).

Additionally, I force the usage of `jj` through
a tool I expose called `campaign`. All agents
operate in deterministically allocated `jj`
workspaces. If I want to make changes, I just
operate in a workspace that sits on top of the
main branch in my project. I wrote about the
evolution of this tool [in another post](/words/private_gerrit_for_agents.html).

I was recently doing some routine analysis of
many sessions across many projects to see what
might be interesting to consider changing or
experimenting with. When viewing results about
`edit` tool error rates and usage, I took a step
back and asked myself: *Why am I even enforcing
read-before-write?*

The isolated virtual machine provides me with
the ability to not care about what happens to
any files on disk, and the `jj`-backed tool for
version control ensures that no agent can clobber
another.

## Simplify

Without a sufficiently simple version control model
and a coherent sandboxing strategy, I do not think
these changes make sense to make. However, as I have
elaborated in the past, the ability to obviate the
need for complex implementation for security allows
you to think about making operational simplifications.

To me, RBW is a holdover from when models were less
reliable, prone to hallucinate basic things more, and
operate in a generally unknown environment. Now, with
agents being more prevalent than ever, RBW is mainly
to prevent clobbering in lieu of a version control
system that makes sense for agents.

As you may be aware, certain models (and increasingly
frontier ones too) use a `bash` tool to perform
edits by using the `python` program passed over
`stdin` using a `heredoc`.

This is fundamentally at odds with the RBW since any
observed text through `bash` or other non-`read` tools
cannot be trivially tracked to perform RBW checks.
This implies that agents do more work than required:
Using `bash`, `grep`, or other reading tools, it finds
where it needs to edit then has to launder the same
tokens through a `read` tool to mark observation. If
the stores are not synchronized between all harness
processes, this is pure theatre.

## Results

In my harness, all subagents are given a fixed budget
of 50 turns, and `edit` is always free. This is a fixed
and shared completions and tool calling budget. I
removed RBW complexity, and I derived the following
insights using a two week window centered on the removal
commit date.

I saw `writer` subagents (n=171) go from 16 median read
calls to 1 median read call per dispatch (n=52). Notably,
quality of output remained the same or improved slightly.
Previously, the `edit` tool error rate for this sample
was 21.9% and that fell to 9.8% after removing RBW.

Notably, it appears mangled or incorrect edits that
required repair, always by the authoring agent, was
0.3% (n=692).

The models used were primarily GPT 5.6 / 6.1 Sol with
a smaller portion of them DeepSeek V4.1 Flash. There
was no significant attribution between models for
errors or repairs. I did notice DeepSeek V4.1 Flash
failed `bash` calls meant for editing more often than
Sol; however, recovery was always immediate.

There is no compelling argument to be made about token
usage or cost; however, the win to me is the removal of
complexity that leaves output fidelity unchanged or even
slightly better. The code for `jj`-backed `campaign` is
slightly smaller than than the code for RBW checks, but
it also buys *a lot* more than RBW. For more on why `jj`
is great choice for agents, see my [post about it](/words/private_gerrit_for_agents.html).

## Closing thoughts

I avoid making any statements about token savings
here since they are specific to my setup and the
harness that I am using. There is no way to make a
sound statement. However, I aim to invite others to
experiment.

I suspect that as we see models and their abstractions
improve, we will need two things: Better verification
and simpler operational thinking. The former is a topic
for a future post while the latter is something that
I hope is captured here in spirit.
