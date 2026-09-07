# What painting and dotfiles have in common

2026-09-07

## Motivations

I have found myself growing more comfortable
with generating programs intended only for my
consumption in the [age of personalized software](https://blog.exe.dev/devtools-must-be-open-source).

I do a lot of maths research in my free time,
and I grew frustrated with rendering LaTeX in
my hard when appearing in the terminal. This
is a monumentally stupid thing to tolerate. I
extended my TUI program to simply bind a small
Go server to use another piece of personalized
software for rendering Markdown+LaTex:
[patel.codes/render](/render). This is now the
same library that renders my website. It is a
completely bespoke LaTeX parser interlaced with
a fork of [rsc.io/markdown](https://rsc.io/markdown).
A few years ago, this would be a deeply misguided
thing to do; however, is that still true today?

## Slop for me is not slop for thee

We have all been there: You setup `$SHELL_OF_CHOICE`
and `$EDITOR_OF_CHOICE`, optionally with your
`$PLUGIN_MANAGER_OF_CHOICE` for choice `$PLUGINS`. Then,
after `$NON_TRIVIAL_AMOUNT_OF_TIME`, something triggers
an upgrade somewhere: You deal with breaking changes,
confusing errors, and a wasted day of trying to stand
up your house of cards again.

Working on Go continues to improve my appreciation for
the simple solutions. I often find myself asking: *Am
I the only user?* If so, the simplest thing is often
to generate the minimally satisfying program and evolve
it with use.

For six years or so earlier in my life, I spent a lot
of time painting and occasionally sculpting. Nowadays,
I mainly shoot film photography as my creative outlet;
however, the feeling of crafting personalized software,
beholden to none other than yourself, does invoke the
similar sense of creativity and fulfillment I experienced
when painting or sculpting.

For years, I would go back-and-forth with my teacher,
debating the intersections of maths, programming,
and art. I wish that he were still around to school
me about what's what, especially in this day and age.

## Enter `patel.codes/mk.sh`

I never knew or cared about editing my dotfiles. I
just wanted something that provided the smallest
surface area for me to work productively. Over time,
it became obvious to me how primitive my configuration
was.

I finally took a step back: *Why is my setup not just
managed by an agent?* It exists for my productivity.
I know what I want and when I want it. So, I had an
agent completely nuke my previous setup. I spent a few
hours using the vanilla setup and materializing the
key programs that I needed for my workflows.

It was freeing to remove plugin managers for `nvim`
and `tmux` while also removing `oh-my-zsh`. I was
left with something pure: What I needed was the same
as what I wanted.

This culminated in the following library of bash scripts
that is invoked by a bootstrapping harness:

```text
05-tailscale.sh     # setup tailscaled, optionally with TS_AUTHKEY
10-zsh.sh           # no frameworks, just managed ~/.config/zsh
15-jj.sh            # with fsmonitor
20-tmux.sh          # with custom fzf-based window tui
30-go.sh            # download latest go release, bootstrap gotip
40-go-tools.sh      # everything i want `go install ...` in my system
50-neovim.sh        # built from stable: simple/bespoke plugins
```

During tight iteration, I was impressed that the
agent immediately reached for testing the setup
in Docker. I could not believe that I had never
thought about doing that before. After modifying
the testing harness, it became trivial to drop in
and iterate.

Now, I can open any machine and make it uniformly
mine by simply running:

```shell
curl -fsSL https://patel.codes/mk.sh | bash
exec zsh
```

Even better, any machine can indepedently checkout
the dotfiles2 repo, make changes, and push it to
the upstream. All other machines simply need to run
`d2sync && exec zsh` to pick up the changes.

## Closing thoughts

In the spirit of the linked exe.dev blog, this setup
is on my [GitHub](https://github.com/thatnealpatel/dotfiles2).
However, the beauty here is that you need not copy
or read my setup; you can simply sculpt your own by
pointing an agent at your dotfiles.
