# Why I stood up a private Gerrit instance

2026-10-10

## Motivations

Large language models, agents, and the harness
that binds them all together is presently a
topic of great importance. As we march through
an age where personalization is more accessible
than ever, I want to share some thoughts about
how it has re-framed my thinking and workflows.

## Vibe coding determinism

For many months, I drove a single agent in a
harness to do things, like many others do. In
particular, the harness I used was Claude Code.
As time went on, I became annoyed at how it
did not quite *fit* me. I wanted to hold a tool
that *just worked* with how I think.

With time, subagents and the set of abstractions
that made LLMs useful became more powerful and
interesting to use. I found myself in a weird
place: I think in terms of graphs, and I wanted
my agents to as well.

I, like many others, built a small program that
would allow my agents to manage their lifecycle:
Planning, coding, and reviewing. The CLI program
was called `tf` and managed the **t**ask **f**low in
the literal sense. A `planner` agent would use
`tf new card` to create managed markdown files
that specified how a self-contained task should
be completed. It would then use `tf wire` with
a wiring syntax to create a DAG that represented
the state. The primary benefit here is that you
had a soft ordering gaurantee when executing. In
some cases, out-of-order or state confusion would
occur. However, agents mostly held the tool in
a sensible manner.

An `orchestrator` was allowed to use `tf` to
see what work is in-flight, pending, blocked,
and resolved. There was a small protocol to
repair the graph if new blockers came up and
this system largely worked. A good mix of skill
and `tf` determinism worked so long as you used
`Fable` and `Opus` class agents.

However, it felt unstable and prone to breaking
at the whim of Anthropic and whatever A/B test
you might be put under. It was also slow in a
sense: Complex parallel work sometimes took
hours to orchestrate and complete and there was
no ability to review the intermediate diffs.

I started wanting something simpler, more in
my control, and purpose-built for me. So,
I used `tf` and Claude Code to bootstrap the
first 500 commits of own bespoke harness written
in Go.

Since July, I have merged 1500 commits to my
harness project using my harness. It is my
daily driver harness and continues to evolve
with me.

## Evolution of `tf`

One of the primary reasons for building my
own harness was codifying this control flow in
a more deterministic way. I spent commits 500
to 750 just observing and collecting data
for how it functioned such that I could reason
about what a coherent abstraction might look
like.

For me, the key problem I wanted to solve
was deterministic, parallel agentic development
such that mechanical blockers, like merge
conflict resolution, were trivially resolvable
that I could free my thinking towards harder
and more complex problems.

Eventually, I took a look at `jj`. I discovered
that `jj` workspaces were effectively the
abstraction I was trying to create. With a very
simple protocol, I created a `campaign` tool
for my harness, and it has been brilliant. It is
by far the most useful tool I have created for
myself.

In short, `campaign` is what allows me to have
many harnesses open in the same project and
jump between them with random ideas or long,
concentrated implementations without any worry
for how it materalizes later. It would just
work when it needed to.

I realized as I was writing and revising this 
post why it is the most important tool I wrote
for myself: I think about 100 things at a time
and jump between them often and carelessly. In
a pre-LLM world, this is what was natural to
me. I am sure many others relate to this sentiment.

When LLMs and agents came along, I felt unsatisifed
because these tools, especially as they became
viable, **forced** me to *serialize my thinking*
in the context of a single agent loop or its
ability to operate in the project I am also trying
to hold in my head. For me, this is *deeply unnatural.*
`campaign` is what allowed me to free myself of
the forced serialization of my thinking process.

## `campaign` API

`campaign` is a `jj`-backed, deterministic tool
that is built into my harness for the main agent.
It is not a CLI program. the harness owns all `jj`
state. Agents merely execute read-only `jj` command.

Its API is remarkably simple: `new`, `status`,
`integrate`, and `adopt` are the only operations.

`new` creates a workspace for changes made
by the orchestrating agent. It returns the
path for which the agent should edit and test.

By default, using the `subagent` tool will
also re-use this same behavior: Subagents
need not call `new` or even become aware of
`campaign` (it is withheld). Subagents are
anchored to the workspace directory. When
they finish, they finalize an intermediate
commit.

`status` reports the canonical campaign state.
With `verbose: true`, it also includes complete
attempt and output inventories. The distinction
between an **output** and an **attempt** is what
makes this work: An output is a self-contained,
finalized unit of work. Using `integrate`, the
orchestrator can `jj describe` it against a base
to produce a change candidate. It **never** can
move the accepted base/trunk.

By default, successful `integrate` implies that
`adopt` is ready to advance the trunk. A machine
wide lockfile is used to synchronize this. The
simplicity of this lockfile accepts that two
unrelated projects may block each other for a
brief amount of time where no synchronization
would otherwise be required.

## So, you stand up Gerrit

It is increasingly becoming more common to write
tools for code review: Agents have made writing
code cheap. For some, code is cheap and also does
not need to be reviewed.

For me, I am still of the opinion that I would
like to review, if only cursory, changes to my
own projects. For collaborative projects, I always
want to review the agent's outputs. For the Go
project, we use [Gerrit](https://www.gerritcodereview.com).
I found it a natural choice: I already know the
UI, its easy to set up (using their container),
and I can slot it directly into my Tailnet under
[my delegated sandboxing strategy](/words/my_thoughts_on_sandboxing_agents.html).

The `campaign` tool described above functions
slightly differently in projects where I define
a `gerrit` remote that points to my private
Gerrit. After `integrate`, the change is ready
to be mailed for human review.

A small program I wrote allows for agents to
interact with Gerrit: [`patel.codes/grfa`](https://patel.codes/grfa).
The delegated authentication by trusted and
untrusted port ensures that the local project
cannot advance the trunk unless the remote has
the merged CL. The Gerrit instance does not allow
submit operations on the untrusted port. Therefore,
only I can ensure it submits.

While there are a few plausible "escapes" an agent
might discover or attempt, they are largely moot
since nowadays a few lines of system prompt and
a simple, deterministic tool will be what an agent
reaches for each time. In practice, I have never
encountered bad behavior. Even if it were to happen,
it is entirely inconsequential and would cause me
all of 30 seconds of repair and annoyance.

## Closing thoughts

I wonder what else has changed in how I do or think
about things that I haven't yet realized are similarly
unnatural. I am curious as to what this might look
like in others as well.
