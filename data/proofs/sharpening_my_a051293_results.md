# What I actually proved about `A051293`

2026-08-21

## Motivations

In May 2026, I published
[An accidentally novel combinatorics proof](/words/an_accidentally_novel_combinatorics_proof_1.html).
The post described a machine-checked proof of the same
[`A051293`](https://oeis.org/A051293) asymptotic that
AlphaProof had formalized.

I claimed a general version of the result, and I also
wrote that I did not have the time to rigorously vet the
proof to present it in my own voice.

While that distinction was intentionally unsatisfying,
it left some questions on the table:

1. Did my Lean definition actually describe `A051293`?

2. How did my theorem compare with AlphaProof's theorem?

3. Was the general `M` result necessary to prove the conjecture?

4. Which novelty claims were substantiated?

Let's re-visit the results in this post; I intentionally
elected to keep the original unedited.

## What I said in May

For those who have not read the original post, it
roughly reduces to:

1. The generated Lean proof appeared to be non-vacuous

2. It proved Cloitre's conjecture for every `M`

3. Its route appeared to be different from AlphaProof's route

4. One intermediate combinatorial identity appeared to be proven for the first time

While none of these claims were/are false, I wanted
to put them to bed once and for all by elucidating
the topic to the best of my ability.

## The sequence and the conjecture

`A051293(n)` counts the nonempty subsets of `{1,...,n}`
whose arithmetic mean is an integer. Cloitre recorded an
asymptotic expansion whose coefficients are the ordered
Bell, or Fubini, numbers:

$$1, 1, 3, 13, 75, 541, \ldots$$

Cloitre's general statement on [`A051293`](https://oeis.org/A051293)
reads as follows:

> More precisely, I conjecture for any `m > 0`, `a(n) = 2^(n+1)/n *
> Sum_{k=0..m} A000670(k)/n^k + o(1/n^(m+1))` (`A000670` =
> preferential arrangements of n labeled elements).

He then gives the fixed statement:

> In fact I conjecture that `a(n) = 2^(n+1)/n * (1 + 1/n + 3/n^2 + 13/n^3 + 75/n^4 + 541/n^5 + o(1/n^5))`.

:::gen

This explicit sentence ends at `M = 5` and fixes the intended indexing:

$$
a(n) = \frac{2^{n+1}}{n}
\left(
1 + \frac{1}{n} + \frac{3}{n^2} + \frac{13}{n^3}
+ \frac{75}{n^4} + \frac{541}{n^5} + o\left(\frac{1}{n^5}\right)
\right).
$$

:::

AlphaProof proves this fixed truncation. My development
proves a corrected asymptotic expansion for arbitrary `M`,
and obtains the displayed `M = 5` statement as a corollary.

Formally, general `M` contains the fixed instance. Practically,
general `M` was not necessary to settle the explicit conjecture
AlphaProof proved. It is cool and (maybe) useful because it
names the coefficient pattern.

## Did I formalize the right sequence?

Yes, but not rigorously.

The original development used `a_comb`, a count over subsets
of `{0,...,n-1}` whose elements are shifted by one when their
mean is tested. That is a convenient Lean representation of
subsets of `{1,...,n}`.

`Proofs/Enumerative/A051293/Cloitre.lean` now adds a literal
version of the OEIS definition:

```lean
def a_oeis (n : ℕ) : ℕ :=
  ((Finset.Icc 1 n).powerset.filter (fun S : Finset ℕ =>
    S.Nonempty ∧ S.card ∣ S.sum id)).card
```

I now also check (by kernel `decide`) the first ten terms
of `a_comb` (original) and `a_oeis` (new) in my proof.
Additionally, `a_comb_eq_a_oeis` proves the two counts
equal for every `n` by shifting each element by one.

This bijection proves that the previous unsubstantiated
use of the convenient internal definition is in fact
the literal definition given.

## What AlphaProof proved

AlphaProof presents `A051293 n` directly as the number
of nonempty subsets of `Finset.Icc 1 n` whose cardinality
divides their sum. Its final theorem, `target_theorem_0`,
is the limit

:::gen

$$
\frac{a(n)-\frac{2^{n+1}}{n}
  \left(1+\frac{1}{n}+\frac{3}{n^2}+\frac{13}{n^3}
  +\frac{75}{n^4}+\frac{541}{n^5}\right)}
 {\frac{2^{n+1}}{n^6}}
\rightarrow 0.
$$

This is exactly the fixed `M = 5` statement, rather than
merely a similar asymptotic. `cloitre_explicit_tendsto`
has the same normalized limit expression over `a_oeis`,
and follows from `cloitre_conjecture 5`. The two declarations
live in separate repositories, but their sequence definitions
are the same literal `Finset.Icc 1 n` count.

:::

## How my proof differs

The proofs overlap more than my original post
suggested; both use a roots-of-unity filter
and both eventually reduce the dominant term to

$$
S(n)=\sum_{j=1}^n \frac{2^j}{j}.
$$

:::gen

AlphaProof gets there by grouping subsets by their cardinality `k`.
For each `k`, its roots-of-unity argument separates the principal term
`choose n k / k` from the nontrivial roots and bounds the latter
exponentially. Summing the principal terms gives

$$
\sum_{k=1}^n \frac{1}{k} \binom{n}{k}=S(n)-H(n).
$$

:::

My proof takes a different combinatorial bridge. It groups the
integer-mean subsets by their maximum and uses the Zumkeller
identity to reduce the count to a sum involving `b_comb(k)`.

A separate roots-of-unity argument identifies `b_comb(k)` with a
divisor-sum formula. The divisor `d = 1` contributes `2^k/k`; after
summing the remaining odd-divisor terms, the proof obtains a polynomial
times `2^(n/3)` bound. This is exponentially negligible relative to
the main `2^n/n` scale, again leaving `S(n)` as the dominant term.

There is also a real difference in how much of the coefficient pattern
is formalized. AlphaProof proves six exact finite geometric-moment
identities for

$$
\sum_{j<n}\frac{j^m}{2^j}, \qquad 0\leq m\leq5.
$$

Their explicit correction terms imply the limiting values
`2, 2, 6, 26, 150, 1082`. Those are twice
`1, 1, 3, 13, 75, 541`, so the coefficients are not arbitrary
constants that merely make the final algebra work. My proof
packages the same phenomenon uniformly as

$$
\sum_{j\geq0}\frac{j^m}{2^j}=2F(m),
$$

then carries the expansion through for arbitrary `M`. The
distinction is therefore not “One proof explains the coefficients
and the other does not.” It is that AlphaProof verifies the first
six moment formulas individually, while my development proves the
Fubini pattern uniformly and uses a different combinatorial route
to reach the same dominant sum.

And, experimentally, this is what I sought to achieve:
without deep background, I wanted to test my experimental
harness to see if it could produce a proof using a
different route.

## The indexing ambiguity

:::gen

Cloitre's general and explicit sentences do not give the little-`o`
term the same clear scope. The general sentence writes
`+ o(1/n^(m+1))` after an unparenthesized product, while the explicit
sentence puts `o(1/n^5)` inside the parenthesized expansion.

Rather than treat the informal general sentence as a separate stronger
claim, this development follows the unambiguous explicit statement.
`cloitre_conjecture M` gives an error of

$$o\left(\frac{2^n}{n^{M+1}}\right).$$

After dividing by the prefactor `2^(n+1)/n`, this is

$$o\left(\frac{1}{n^M}\right).$$

At `M = 5`, that is exactly the parenthesized `o(1/n^5)` remainder in
Cloitre's explicit sentence. The Lean theorem follows that convention
without taking a position on how the general OEIS sentence should be
repunctuated.

:::

## What survived

Broadly speaking, the original claims are
true; however, they are more faithfully
sharpened:

- Generated a machine-checked proof of the corrected general-`M` expansion;
- The convenient combinatorial count is proved equal to a literal `A051293` count for every `n`;
- Both definitions are checked against the first ten published terms;
- The explicit `M = 5` result proved by AlphaProof follows from the general theorem;
- The proof route explains the Fubini coefficients uniformly;

The general `M` theorem is cool; it expands six convenient
coefficients into a pattern and lifts a story as to why they
appear; however, it is not strictly necessary to prove the
conjecture.

## Literature check

(I did write some of the prose below, but I have left
the shadowing to indicate that I merely adopted the
literature check provided by my research harness.)

:::gen

The exact per-`k` equality is recorded in [A082550](https://oeis.org/A082550)
and [A063776](https://oeis.org/A063776) through observations by Papadopoulos in
2016 and Wiseman in 2019. Neither entry supplies a proof.

A literature sweep found published  neighboring results on
zero-sum subsets and necklaces, but no published proof of this
exact integer-mean equality and no independent formalization of it.
The Lean file therefore supplies a formal proof of an OEIS-observed
identity.

Additionally, one OEIS cross-reference is shifted: [A082550](https://oeis.org/A082550)
prints `A051293(n+1) - A051293(n)`, while the definitions and terms
give `A051293(n) - A051293(n-1)`.

The `k+1` in the Lean summation is intentional: It converts `Finset.range n`
from zero-based indices to maxima `1,...,n`. The proof derives the relation
from the underlying counts, so this does not affect its results.

:::

## Closing thoughts

In general, most of my original hedging was warranted. Though,
in this case, the missing work was small. It has become more
obvious to me that writing these posts is going to be the
bottleneck.

I understand how to audit my proofs more rigorously now, and I
have built some substantial tooling in order to continue this
type of research... however, that is a topic for another post.

## References

- [arXiv:2605.22763](https://arxiv.org/abs/2605.22763): Exact theorem and the paper's description of the conjecture;
- [`google-deepmind/AlphaProof-nexus-results`](https://github.com/google-deepmind/alphaproof-nexus-results/blob/main/APNOutputs/OEIS/oeis_51293_conjecture_0.lean): Exact `target_theorem_0` signature and proof structure;
- [`Proofs/Enumerative/A051293/Counting.lean`](https://github.com/thatnealpatel/proofs/blob/main/Proofs/Enumerative/A051293/Counting.lean): `cloitre_conjecture`;
- [`Proofs/Enumerative/A051293/Cloitre.lean`](https://github.com/thatnealpatel/proofs/blob/main/Proofs/Enumerative/A051293/Cloitre.lean): `a_oeis`, `a_comb_eq_a_oeis`, and `cloitre_explicit_tendsto`;
