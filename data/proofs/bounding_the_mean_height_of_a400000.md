# A logarithmic bound for `A400000`

2026-11-25

## Introduction

In Eric Fox's [city-building sequence (`A400000`)](https://oeis.org/A400000),
the first building has height zero. Each new building takes the
*least* nonnegative integer height $h$ from which it sees exactly
$h$ earlier rooftops. A rooftop blocks the view only if it lies
strictly above the line of sight. Collinear rooftops all count.

The sequence starts $0,1,1,2,1,2,1,3$. Write $a(n)$ for the height of the
$n$th building and let $S_N=\sum_{n=1}^{N}a(n)$. Then $S_N/N$ is the
mean height of the first $N$ buildings. The partial sums are listed in
[`A397740`](https://oeis.org/A397740), along with the following conjecture:

$$
\limsup_{N\to\infty}\frac{S_N}{N}\le 3.
$$

In this post, I present a logarithmic upper bound that I
obtained on Nov 7, 2026 and attempted to improve before
running out of time to do so. My
[recent efforts](https://github.com/thatnealpatel/grh)
to help formalize and repair
[Masullo's *The Music of the Primes*](https://math216.wordpress.com/wp-content/uploads/2026/11/the-music-of-the-primes.pdf)
is likely going to backlog this further. As such, I figured
writing up the bound is interesting enough.

For every $N\ge 1$, we obtain the bound

$$
S_N\le 100N\bigl(\lfloor\log_3 N\rfloor+2\bigr).
$$

The previous statement is machine-checked in Lean 4 and uses
only `propext`, `Classical.choice`, and `Quot.sound`:

```lean
theorem term_partialSum_log_bound (N : ℕ) :
    (∑ i ∈ Finset.range N, term i) ≤ 100 * N * (Nat.log 3 N + 2)
```

The full conjecture remains **open** for $\leq 3$.

On a personal note, I hope to revise this more in my own voice.

## Looking back from a lower roof

:::gen

Picture the first three rooftops one unit apart, at heights zero, one,
and one. Stand at the fourth site and look left. From ground level you
see only the nearest roof. At height one you see both height-one roofs:
The nearer roof lies on the horizontal sightline, not above it.

At height two you still see just those two roofs, so two is the first
height that matches the count. The proof uses the same idea in reverse:
Given a tall building, what would you see from a lower roof at its site?

Number the positions from zero, as in Lean, and write $x_n$ for
the height at position $n$. Suppose $i<n$ are both of height at
least $B\ge4$, but their distance $d=n-i$ satisfies $4d\le B$.
At $n$, try a roof of height $h=\lfloor B/2\rfloor$ instead of
$x_n$. Two things happen:

1. **The trial sees too many roofs.** Raising the trial roof can
   only add views. Thus the visible count minus the trial height
   drops by at most one per unit step. It starts positive (the
   immediate predecessor is visible even from the ground) and
   cannot become negative without first hitting zero. Since $h$
   is below the *least* fixed point, the trial sees at least $h+1$
   roofs. At most $d$ are at or after $i$, so at least $h+1-d$
   lie strictly before $i$.
2. **Those earlier roofs are tall.** The sightline to any such
   roof $p$ must clear building $i$, which stands at least
   $B-h\ge 2d$ above the trial. Going left from $n$, the line
   must therefore rise at least two units per position:
   $$
   x_p\ge h+2(n-p).
   $$

The second observation gives a *lower* bound on each visible
height. For the upper bound, ask what that building could see
when it was built.

:::

## Passing a sightline from one building to another

:::gen

List all the roofs visible from the trial that are strictly taller
than $h$, from left to right. Take two successive ones, $p<q$,
with $r$ roofs before $p$ in the list. Then

$$
x_q\le r+(q-p).
$$

Here's why. Building $q$ sees exactly $x_q$ earlier roofs at its
actual height. At most $q-p$ of them lie at or after $p$. I claim
that every roof $t<p$ seen by $q$ is also in the trial's list.

Draw the two known sightlines: From $t$ to $q$, and from $p$ to
the trial point at $n$. The first passes on or above $p$, and the
second passes on or above $q$. In slope notation this gives

$$
s(t,q)\le s(p,q)\le s(q,n),
$$

where $n$ denotes the trial point. Hence the chord from $t$ to $n$
lies on or above $q$. It dominates the clear chord $tq$ from $t$
to $q$, and, since it also clears $p$, it dominates the clear
chord $pn$ from $p$ to $n$. These intervals cover every possible
blocker. So the trial sees $t$. Moreover, $x_p>x_q$: The chord
from $p$ to the lower trial roof clears $q$, whose height exceeds
$h$. Thus $s(t,q)\le s(p,q)<0$, and $x_t>x_q>h$. This puts $t$
among the $r$ listed roofs before $p$, as claimed.

:::

## Where the logarithm comes from

:::gen

Put $k=h-d$. By the trial count, the first $k+1$ roofs in the
list all lie before $i$. Take consecutive roofs $p,q$ among them,
so that the rank $r$ of $p$ satisfies $r<k\le h$. Our two bounds
on $x_q$ give

$$
h+2(n-q)\le x_q\le r+(q-p).
$$

In particular $q-p\ge 2(n-q)$, or $n-p\ge 3(n-q)$.
Read the list backwards: **Each step at least triples the distance
to the trial roof.** There are $k$ steps, the nearest distance is
at least one, and the farthest is at most $n$. Consequently

$$
3^{\lfloor B/2\rfloor-d}\le n.
$$

Now let $L=\lfloor\log_3 N\rfloor$. If $n<N$ and $B>4(L+1)$,
then $4d\le B$ would give $\lfloor B/2\rfloor-d\ge L+1$.
There isn't room: $3^{L+1}>N$. This proves the spacing result.

To finish, count the sum one floor at a time. Write $T_B$ for
the number of buildings of height at least $B$ in the prefix.
Then $S_N=\sum_{B=1}^{N}T_B$, since a building cannot see more
roofs than there are predecessors. Up to $b=4(L+1)$ use $T_B\le N$;
above $b$ the spacing gives $T_B\le 1+4N/B$. Therefore

$$
S_N\le bN+N+4N\sum_{b<B\le N}\frac1B
    \le N\bigl(b+1+4\ln N\bigr).
$$

This implies the stated bound. The Lean proof counts the same
tail in geometric bands rather than summing the harmonic series.

:::

## Why this doesn't give a bounded mean

:::gen

For comparison, write down the largest power of two dividing each
positive integer: $1,2,1,4,1,2,1,8,\ldots$. Terms at least $B$
are spaced at least $B$ apart, *better* than the spacing above.
Yet among positions $1$ through $2^m-1$, each height $2^j$ for
$0\le j<m$ occurs $2^{m-j-1}$ times. Each of the $m$ height
classes contributes the same $2^{m-1}$ to the sum. The mean grows
like $m/2$, not like a constant. This isn't a city-building
sequence; it shows exactly what is lost by keeping only the
spacing information.

:::

## Source

The [Lean proof](https://github.com/thatnealpatel/proofs/blob/6a52b1f276f81a150e279e03cbad2e9c1daf89f6/Proofs/Sorry/Combinatorics/A400000LogBound.lean) was checked on November 7, 2026.
I still need to lift this out of the `Proofs/Sorry` directory.

## AI Disclosure

These results were obtained primarily using a mix
of Anthropic, OpenAI, and open weight models in a
private benchmark using a bespoke harness.
