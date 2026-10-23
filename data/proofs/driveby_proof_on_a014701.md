# Drive by proof on `A014701`

2026-08-23

## Motivations

Two of my most persistent obsessions remain
maths and programming; one finds the former
often flowing into the latter, but only in
the recent century or so have we seen the
latter flowing into the former.

Naturally, formal methods and "auto-research"
are even more recent developments in this
space. While my disdain for large language
models largely lives on, there is virtually
no sound argument in denying the utility that
these tools create.

## The conjecture

On 15 May 2025, Jean-Marc Rebert wrote a comment on
[`A014701`](https://oeis.org/A014701) which conjectures:

```
Conjecture: a(n+1) is the minimal number of steps to
go from 0 to n, by choosing before each step, after
the first step, whether to keep the same step length
or double it. The initial step length is 1.
```

I came across this entry while doing some basic
conjecture mining (tooling in my harness) over
OEIS. The goal was to find open conjectures that
an LLM could formalize into theorems. My primary
motivation was to get familiar with the pitfalls and
where the human-in-the-loop is required. It later
evolved into a cohesive workflow that I am still
refining... a topic for a future post.


`A014701` counts the number of multiplications in the
Chandah-sutra method, a form of
[exponentiation by squaring](https://en.wikipedia.org/wiki/Exponentiation_by_squaring).

For positive `m`, the formula recorded on the entry is

```
A014701(m) = floor(log m) + popCount(m) - 1
           = bitLength(m) + popCount(m) - 2
```

where `log` is the binary logarithm.

This distinction matters: Bit length alone is
[`A070939`](https://oeis.org/A070939), while `popCount` is
[`A000120`](https://oeis.org/A000120). Neither is `A014701` outright.

[`A056792`](https://oeis.org/A056792) is numerically
related by `A014701(m) = A056792(m) - 1`, with one
key difference: It minimizes the steps that either
add one to the current position or double the current
position.

Rebert's walk always advances and chooses whether to
keep or double the **step length**. The numerical identity
is useful context, but it is not a rigorous proof that
the two walks are one in the same.

## The model and theorem

In
[`StepWalk.lean`](https://github.com/thatnealpatel/proofs/blob/fa83e94/Proofs/NumberComplexity/StepWalk.lean),
`Reach k p s` states: After `k` steps the walk is at position `p`
and its current step length is `s`.

This captures Rebert's formulation directly:

```lean
Reach 1 1 1

/- keep -/
Reach k p s → Reach (k + 1) (p + s) s

/- double -/
Reach k p s → Reach (k + 1) (p + 2 * s) (2 * s)
```

`Reachable n k` means that some current step length `s`
makes `Reach k n s` true. The theorem signature is

```lean
∀ (n k : ℕ), 1 ≤ n → (IsLeast {j | Reachable n j} k ↔ k + 2 = (n + 1).bits.length + popCount (n + 1))
```

so it proves the conjecture for positive destinations `n >= 1`.

It deliberately does not cover `n = 0`: Every modeled walk has
taken its initial step and is at a positive position, at least
as far as I understand it.

`A014701(1) = 0` corresponds to the empty walk.

`IsLeast {j | Reachable n j} k` give us both halves of the
"minimality argument:" a walk reaches `n` in `k` steps, and
every other achievable step count is at least `k`.

The Lean formulation deliberately avoids truncated subtraction
in `ℕ`; by the formula above, its RHS states exactly `k = A014701(n + 1)`.

## The geometry of the decision tree

I'll be honest: I looked at all that confused and thought
to myself *Well, isn't the geometry of this encoded
in the decision tree? Can't I just see it?*

Thus, I asked my stochastic parrot to produce some nice
ASCII visualizations with supporting prose that brings the
thought to life.

Each node in the tree below is `(position, step length)`
and `K` = keep, `D` = double.

:::gen

```
                                  (1,1)
                         K       /     \       D
                                /       \
                             (2,1)     (3,2)
                       K    /   \ D   K /   \ D
                           /     \     /     \
                        (3,1)  (4,2) (5,2)  (7,4)
```

:::

Now, the useful geometry appears when taking the decision
tree and folding it by `(p,s)`; consider the following

:::gen

```
                 K             K              D
(p,s) ------> (p+s,s) ----> (p+2s,s) ----> (p+4s,2s)
  |                                                   ^
  | D                                                 |
  v                                                   |
(p+2s,2s) -------------------- K ---------------------+
```

:::

The upper route is `KKD` and the lower route is `DK`. Both
finish at the same full state, but the lower route uses two
transitions instead of three. This is binary carrying made
visible:

```
KKD  -->  DK          two keeps of size s become one keep of size 2s
```

If two keeps occur at the very end (`...KK`), it can be
replaced by a single `D`.

Let's call any arbitrary string of `K` and `D` choices
a **decision word.** Since doubles only move down, every
such word has one block of `K` at each scale:

:::gen

```
scale 2^0:    K ... K    D
scale 2^1:    K ... K    D
   ...
scale 2^t:    K ... K    stop
              < c_t >
```

Writing `c_i` for the number of keeps on row `i` gives

```
p + 1 = 2^(t+1) + c_0 2^0 + c_1 2^1 + ... + c_t 2^t

k     = 1 + t + c_0 + c_1 + ... + c_t.
```

Whenever some `c_i >= 2`, the carry cell supplies a shorter
route to the same endpoint. A shortest walk therefore has only
zero or one keep on each row:

```
scale 2^0:    [K if bit 0 is 1]    D
scale 2^1:    [K if bit 1 is 1]    D
                 ...
scale 2^t:    [K if bit t is 1]    stop
```

These choices are exactly the lower binary digits of `p + 1`;
its leading `1` is the baseline `2^(t+1)`. The number of rows
supplies bit length, and the number of horizontal moves
supplies population count:

```
t + 2              = bitLength(p + 1)
1 + sum of the c_i = popCount(p + 1)

k = bitLength(p + 1) + popCount(p + 1) - 2.
```

:::

So, there you have it: An arbitrary walk acts like a
binary expansion with uncarried digits, while a shortest
walk is the unique fully carried word, with at most one
keep at each scale.

You may look at these visualizations and and similarly
be tempted to make the following conjecture:

```
for each positive target, there is exactly one shortest decision word.
```

So, I handed this off to my harness, and about 20
minutes later, it produced a machine-checked proof
for this conjecture. I am not convinced it is that
interesting or useful, but the experience of being
able to go from writing this post, seeing something
interesting, and backgrounding the proof whilst not
giving up focus on writing this post was notable.

The proof exists nearby at
[`List WalkDecision`](https://github.com/thatnealpatel/proofs/blob/b3386c2e038e4b8f0ae236f79b897d6580f50733/Proofs/NumberComplexity/StepWalk.lean#L690-L708):

:::gen

```lean
theorem existsUnique_shortest_decisionWord (n : ℕ) (hn : 1 ≤ n) :
    ∃! w : List WalkDecision,
      decisionPosition w = n ∧
        IsLeast {j : ℕ | Reachable n j} (w.length + 1)
```

:::

## Proof kernel

For the boring details, the proof contains
three claims, only two of which are required
to prove the conjecture as Rebert wrote:

1. No walk can beat the binary cost

2. A walk built from binary digits attains that cost

3. (Extra credit) the attained decicion word is unique

We get (1) and (2) by packaging

```lean
binCost m = m.bits.length + popCount m
```
and

```lean
/- q ≠ 0, r < 2^u -/
binCost (q * 2^u + r) = u + binCost q + popCount r
```

An induction over a `Reach k p s` derivation supplies
the lower bound by maintaining

```lean
s = 2^t ∧ 2^(t+1) ≤ p + 1 ∧ binCost (p + 1) ≤ k + 2
```

The keep and double cases both reduce to one exchange lemma:

```lean
2^u ≤ P → binCost (P + 2^u) ≤ binCost P + 1
```

Here `P = p + 1`. The last part of the invariant says that
every walk of length `k` reaching `n` satisfies

```lean
binCost (n + 1) ≤ k + 2
```

(I have no idea why the agent decided to style it this way.)

Finally, we show the above step count `k` satisfies

```lean
k + 2 = (n + 1).bits.length + popCount (n + 1)
```

:::gen

Together, the lower bound and construction prove the formula for the
minimum.

The uniqueness theorem needs one further layer because `Reach` is a proposition,
not path data. The proof therefore uses `List WalkDecision` and establishes
three facts corresponding directly to the geometry:

1. Every decision word gives a `Reach` derivation, and every `Reach` derivation
   has a decision word;
2. Every shortest word is `DecisionNormal`, meaning it contains no `KK`;
3. Two normal words ending at the same position are equal.

The second fact is the carry reduction, and the third says that a fully carried
word is determined by its endpoint. This proves actual uniqueness of the
keep-or-double choices, rather than the vacuous uniqueness of proofs of a
proposition.

:::

`StepWalk.lean` compiles without `sorry`, `admit`, or `native_decide`.

`#print axioms` reports exactly `[propext, Classical.choice, Quot.sound]`
for both `rebert_conjecture` and `rebert_conjecture_iInf`.

The file also checks 14 selected OEIS terms in the kernel:
indices `1..8`, `15`, `16`, `31`, `32`, `64`, and `86`,
through `a(86) = 9`.

## Literature check

I have little in the way of understanding
if this is an obvious exercise that any
practitioner can otherwise prove trivially.
After actually reading further into the
problem, I am fairly confident this is not
a notable result.

:::gen

1. The live `A014701` entry and its links and cross-references one hop out. The
   entry contains several proved formulas and neighboring characterizations,
   including the different `A056792` walk, the
   [Gruber–Holzer 2021](https://doi.org/10.4230/LIPIcs.MFCS.2021.52)
   max formula, and Cunningham's base-2 digit-sum comment, but not a proof of
   Rebert's walk;

2. The [RSOS assembly-theory paper published in 2026](https://doi.org/10.1098/rsos.260082)
   and now linked from the entry. Its accessible publisher text identifies
   `A014701` with the classical depth index and contains no Rebert attribution
   or keep-or-double walk;

3. An exact-phrase and citation probe for `A014701`, Rebert, and variants of
   the walk description;

4. SeqFan: The old pipermail host was unreachable, its available Wayback index
   predates the conjecture, and site-scoped probes of the current Google Groups
   archive returned nothing, although that indexing is spotty;

5. A full clone and text search of
   [`sequencelib`](https://github.com/provables/sequencelib), associated with
   [arXiv:2601.11757](https://arxiv.org/abs/2601.11757), which returned no
   `A014701` hit among 25,905 Lean files.

:::

## Closing thoughts

I spent **way more time** than I care to admit trying
to understand the generated proof and its formulation.
This blog post is me vicariously lifting my sunk-cost
fallacy to you.

It very well could be that I am the first person to
waste time on this exercise; in any case, it was a
nearly-free drive-by that was serving other purposes.
