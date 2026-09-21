# My thoughts on sandboxing agents

2026-09-21

## Motivations

It is beyond clear that LLMs are here to stay
as a foundational technology. In my mental
model, **LLMs** are made useful by **agents** which in
turn are made useful by a **harness**. However, the
security of these non-deterministic programs
gives me pause.

I aim to present a coherent trust model when
thinking about how to make agents most useful.
Naturally, this trust model is one that both
the human operator and the harness agree on.
I *do not aim* to provide a comprehensive
setup guide here: Use cases vary, and what is
ergonomic for me may not necessarily be what
is ergonomic for you.

## Trust model

It may seem that how you hold the tool and
what you use it for matter in the context of
security, but I assert all use cases are the
same. [From authoring and elucidating proofs](/words/proving_a092482_closed_form.html)
to [using an agent to customize my dotfiles](/words/what_painting_and_dotfiles_have_in_common.html),
there is no room for ambiguity in how I define
and operate inside of my agentic trust model.

LLMs are fundamentally non-deterministic; any
attempt to constrain, monitor, or otherwise
police the program is, in my opinion, theatre.

I firmly take the stance that agents should
never run on **trusted machines**, no matter
what. Here, I define a **trusted machine** as
a credential-bearing machine or a machine with
privileged access to other systems, files, or
what have you.

I trust agents to do good work, but I do not
trust them to do it securely.

## Why?

The convenience of using these tools need not
be a justification for risk-taking behavior.
Furthermore, it prevents you from exploring the
world of open weight models.

While using `bubblewrap` or other sandboxing
solutions on your trusted host is certainly an
option, I ask: Why?

While admittedly rare, agents may find a way
to exploit or break out of these boundaries.
I just do not see the justification for going
through any efforts to sandbox just to let
them run on the same silicon. Why not just
go all the way?

Aside from the security argument, a coherent
sandboxing strategy obviates the need for all
the complexity that both bloats a harness'
codebase and its system prompt: Rolling any
safety measure (read: security theatre) becomes
unnecessary. This allows you to *harness* the
raw utility of the underlying LLM more effectively.

I concede that if your workflow is primarily
GUI based, it is certainly non-trivial to
achieve the same level of convenience as a
TUI user; however, this is a minor detail that
is largely a one-time setup cost.

## Sandboxing

The advent of [Tailscale](https://tailscale.com)
and [exe.dev](https://exe.dev) make sandboxing
trivial.

\$80/month at exe.dev gives you up to 50
VMs with up to `8 vCPU` / `32 GiB` of peak
collective usage at any given time. They
also have \$20 and \$40 options as well as
a \$160 variant (of which I am a subscriber).

For individuals, Tailscale is **free.** You
can even invite up to 6 guest users to your
Tailnet. The `tag` feature provides a nice
interface to instrument security policies.

For example, imagine a network in which you
have the following tags:

- `tag:proxy`
- `tag:sandbox`
- `tag:privileged`

You can tag your personal laptop as `privileged`,
stand up a single `2 vCPU` / `8 GiB` machine
with `proxy`, and finally create a single `8 vCPU`
/ `16 GiB` machine with `sandbox`.

You can then apply the following security
policies:

```
{
	"tagOwners": { ... },
	"ssh": [ ... ],
	"grants": [
		{
			"src": ["tag:privileged"],
			"dst": ["*"],
			"ip":  ["*"],
		},
		{
			"src": ["tag:proxy"],
			"dst": ["tag:sandbox"],
			"ip":  ["tcp:22"],
		},
		{
			"src": ["tag:sandbox"],
			"dst": ["tag:proxy"],
			"ip":  ["tcp:9001"],
		},
		{
			"src": ["tag:proxy"],
			"dst": ["tag:sandbox"],
			"ip":  ["tcp:10000-10099"],
		},
	],
}
```

And just like that you have a coherent model
for operating untrusted agents inside of a
trusted network.

The `proxy` holds all secrets: Inferencing
authentication, GitHub PATs/FGTs, API keys,
etc. The `9001` port is an **untrusted**
ingress point. The proxy can be made aware
of how to serve on it. Additionally, `proxy`
can also serve on `9999` which we may say
is the **trusted** ingress point.

There are different types of **trusted** machines.
For example, the `proxy` is trusted, and it also
contains trusted data. You may wish to roll a
**trusted** machine whose sole purpose is to store
**untrusted** data. This is still a machine that
does not run agents; however, it cleanly isolates
the two roles:

```
{
	"tagOwners": { ... },
	"ssh": [ ... ],
	"grants": [
		...
		{
			"src": ["tag:proxy"],
			"dst": ["tag:library"],
			"ip":  ["tcp:9999", "tcp:9001"],
		},
		{
			"src": ["tag:sandbox"],
			"dst": ["tag:library"],
			"ip":  ["tcp:9001"],
		},
	],
}
```

We have no reason to give `library` access
to the `proxy`. However, we may want all
logging to ingress on `9001` and then expose
a query endpoint on `9999`.

This is largely one-time setup and provides
a reasonable boundary. Now, you can simply
run any model in any harness to do any task
without any cause for concern.

You could, for example, instrument the `proxy`
such that it only allows access selectively to
GitHub: This is what I do. My agent identity on
GitHub is [`notnealpatel`](https://github.com/notnealpatel).

You cannot trust agents to be good citizens on
the internet. Point in case: A recent thread where
[agents opened an issue in the Go project](https://go.dev/issue/81606)
and just used it as a medium to communicate with
each other pointlessly.

I never want my agents assaulting a human's
attention. My `proxy` only allows the agent
to write to an allowlisted set of paths, all
of which are forks it creates or a narrow set
of personal GitHub projects I own.

We are the stewards of what these programs do;
we owe it to each other to operate them in a
respectful manner.

## Sandboxing is... convenient?

While it may seem that you trade convenience
for security, there are material benefits to
setting things up like this. It makes programs
you write simple to reason about.

You do not need to worry about what trusted
machines send to each other: They're trusted.

Inside of this model, you needn't even worry
much about a `sandbox` DoS'ing or otherwise
attacking `proxy` or `library`. It is in the
interest of `sandbox` to keep those machines
up. Spinning up new `sandbox` machines is
trivial. Once Tailscale immutable tags apply,
the machines have instant access to what
they should have access to.

If you had some fun side projects you wanted
to expose to **trusted users** you can do so
by creating another `sandbox` machine and then
splitting off traffic: A new tag `guest` might
be used to route inferencing quota or traffic
at the trusted `proxy`.

The same applies to a website or service
you want to expose to **untrusted users.**
exe.dev makes this boundary dead simple
to instrument using their canonical share
feature.

## Closing thoughts

Some might call this being *paranoid* and
that is perfectly fine: You have a risk
appetite that I simply do not share.

To me, the peace I feel each night when I
go to sleep is worth it. My agents will
continue goal closing and operating with
full autonomy. I choose to spend my energy
reviewing outputs and thinking about what
to do next, not at the approval prompt.
