# Refuting two prime-counting conjectures

2026-10-23

There exist two conjectures about the number of ways
you may write an integer as a sum of prime counts.
The first predicts existence and the second conjectures
about its uniqueness.

Here, I present two counterexamples: For [A262995](https://oeis.org/A262995),
we have $a(241540)=0$. For [A263100](https://oeis.org/A263100), we have
$a(15293)=1$, outside the conjectured list of exceptions.

Both have machine-checked proofs using only
`propext`, `Classical.choice`, and `Quot.sound`.

Here $\pi(x)$ counts the primes less than or equal
to $x$. For example, $\pi(10)=4$, counting $2,3,5,7$.

## Representation doesn't always exist

[`A262995`](https://oeis.org/A262995) counts positive
ordered pairs $(k,m)$ for which

$$
n=\pi(T_k)+\pi(1+T_m),
\qquad T_j=\frac{j(j+1)}2.
$$

Here $T_j=1+2+\cdots+j$ is the $j$th triangular
number: $1,3,6,10,\ldots$.

A conjecture from October 2015 states:

> $a(n)>0$ for every $n>0$.

The entry reports verification through $10^5=100000$.

**This conjecture is false.** Take $n=241540$,
beyond the reported search range of $100000$.
The conjecture predicts $a(241540)>0$, but we
establish $a(241540)=0$: There is no positive
pair $(k,m)$ representing it.

The argument below certifies this particular
counterexample.

### Proof

:::gen

```lean
def pairSet (n : ℕ) : Set (ℕ × ℕ) :=
  {p | 0 < p.1 ∧ 0 < p.2 ∧
    n = Nat.primeCounting (p.1 * (p.1 + 1) / 2) +
      Nat.primeCounting (1 + p.2 * (p.2 + 1) / 2)}

noncomputable def a (n : ℕ) : ℕ := (pairSet n).ncard

theorem pairSet_241540 : pairSet 241540 = ∅

theorem a_241540 : a 241540 = 0
```

:::

Here, `p.1` and `p.2` are the two coordinates of the
pair. `Nat.primeCounting` is the aforementioned $\pi$.
The first theorem states the set of representing pairs
is empty; the second says its cardinality is zero.

First, the proof certifies these two prime counts:

:::gen

$$
\pi(T_{2596})=241602,
\qquad
\pi(1+T_{2596})=241603.
$$

Both are larger than $241540$. If $k\geq2596$,
then $T_k\geq T_{2596}$, so

$$
\pi(T_k)+\pi(1+T_m)
\geq\pi(T_k)
\geq241602
>241540.
$$

Such a pair cannot work. The same argument excludes
$m\geq2596$. Every possible representation must
therefore satisfy

$$
1\leq k\leq2595,
\qquad
1\leq m\leq2595.
$$

:::

Now the problem is finite, but there are still two
things to establish: The prime counts in the tables
are correct, and no two relevant counts sum to the
target.

### Certifying the prime counts

The source uses a [sieve of Eratosthenes](https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes).
Every composite integer has a prime divisor no
larger than its square root, which is why the
following finite sieve suffices.

:::gen

The largest argument needed here is

$$
1+T_{2596}=3370907<1837^2.
$$

Thus primes below $1837$ suffice to sieve the
required range. The implementation works in
finite windows and proves that a surviving bit
corresponds exactly to a prime. Counting those
bits gives the number of primes in each interval.
Adding the interval counts certifies the tables

$$
A_k=\pi(T_k),\qquad B_m=\pi(1+T_m).
$$

The listed table values are not assumptions:
The source proves that its table checker returns
correct prime counts and then evaluates the
checks in Lean's kernel. The tables include indices
through $2596$, covering every possible solution
and the cutoff values used above.

:::

### Excluding every sum

:::gen

Write $N=241540$. For a fixed value $b=B_m$, the
equation we want to rule out is

$$
A_k+b=N.
$$

If $b>N$, it is impossible because $A_k\geq0$.
Otherwise it is equivalent to

$$
A_k=N-b.
$$

So instead of repeatedly adding every pair of
table entries, we can ask whether the complement
$N-b$ occurs anywhere in the first table.

The source records that first table as a bitmask:
An integer whose bit at position $r$ is 1 exactly
when $r$ occurs as an $A_k$ for a positive table
index $k$. For example, a table containing
$\{0,2,5\}$ would set bits 0, 2, and 5, giving

$$
2^0+2^2+2^5=37.
$$

Testing bit 2 would report a match; testing bit 3
would not. Repeated table values simply set the
same bit again. That is enough here because we
are testing existence, not counting how many
indices produce each value.

For every positive-index entry $b$ in the second
table, the check verifies either $b>N$ or that
bit $N-b$ is unset. Every entry passes: No
complement occurs. The source proves what the
bitmask lookup means and uses `decide +kernel`
to certify that the finite check returns `true`.

The bitmask makes the computation efficient; the
bounds make it exhaustive. This is not merely a
search that failed to find a pair. The bounds
exclude every larger pair, and the certified
calculation excludes every remaining pair. Thus

$$
\boxed{a(241540)=0.}
$$

:::

I am making no claim that this is the first
such zero.

## A unique representation in `A263100`

[`A263100`](https://oeis.org/A263100) counts positive
ordered pairs $(k,m)$ for which

$$
n=\pi(k^2)+\pi(\lfloor m^2/2\rfloor).
$$

Its October 2015 conjecture states:

> $a(n)>0$ for every $n>0$, and $a(n)=1$ only for
> $n=1,3,28$.

**The uniqueness clause is false.** At $n=15293$,
the only representing pair is $(339,291)$.
This does not disprove positivity: This input
has a representation.

### Proof

This sequence has its own definitions of `pairSet`
and `a`, in a separate Lean namespace:

:::gen

```lean
def pairSet (n : ℕ) : Set (ℕ × ℕ) :=
  {p | 0 < p.1 ∧ 0 < p.2 ∧
    n = Nat.primeCounting (p.1 ^ 2) +
      Nat.primeCounting (p.2 ^ 2 / 2)}

noncomputable def a (n : ℕ) : ℕ := (pairSet n).ncard

theorem mem_pairSet_15293 (k m : ℕ) :
    (k, m) ∈ pairSet 15293 ↔ k = 339 ∧ m = 291

theorem a_15293 : a 15293 = 1
```

Division of natural numbers in Lean rounds down,
so `p.2 ^ 2 / 2` is $\lfloor m^2/2\rfloor$.
The symbol `↔` means “if and only if”: A pair
represents $15293$ exactly when its coordinates
are $339$ and $291$.

First, that pair works:

$$
\pi(339^2)=\pi(114921)=10867,
$$

$$
\pi(\lfloor291^2/2\rfloor)=\pi(42340)=4426,
$$

and $10867+4426=15293$.

But finding one pair is not enough to prove
uniqueness. We must rule out every other pair.
As before, two certified prime counts give bounds:

$$
\pi(410^2)=15346>15293,
$$

$$
\pi(\lfloor579^2/2\rfloor)=15308>15293.
$$

If $k\geq410$ or $m\geq579$, one summand alone
already exceeds the target. Every solution must
therefore lie in

$$
1\leq k\leq409,
\qquad
1\leq m\leq578.
$$

The kernel-checked table calculation excludes
every $k$ except $339$. For that value, the second
summand must be $15293-10867=4426$. The table
then shows that only $m=291$ gives that count.

There is exactly one representing pair, so

$$
\boxed{a(15293)=1.}
$$

:::

Since $15293\notin\{1,3,28\}$, part (i) of the
stated conjecture is **false**. This refutes its
uniqueness clause, not its positivity clause.
The separate conjectures (ii) and (iii) are not
addressed here.

## Related counterexample

A related entry [`A263001`](https://oeis.org/A263001)
counts representations of a different form:

$$
n=\pi(k(k+1))+\pi\!\left(\frac{m(m+1)}2\right),
\qquad k,m>0.
$$

Its conjecture includes the claim that $a(n)=1$
only for $n=1,4,6$. OEIS also records a counterexample
to this proposed uniqueness list. Like the `A263100`
example above, it refutes the uniqueness clause,
not positivity.

## Source

The source code for both proofs is on my GitHub: [`A262995Refutation.lean`](https://github.com/thatnealpatel/proofs/blob/e0c9810910fd2faa896326257d0a7b95216756cc/Proofs/Sorry/NumberComplexity/A262995Refutation.lean), [`A263100Refutation.lean`](https://github.com/thatnealpatel/proofs/blob/e0c9810910fd2faa896326257d0a7b95216756cc/Proofs/NumberComplexity/A263100Refutation.lean).

The former file lives in a `Proofs/Sorry` *directory*
in the project; however, its contents machine check
in the manner stated at the top of this post. I will
eventually lift it to the `Proofs/` module and update
this post.

## AI Disclosure

These results were obtained primarily using a mix
of Anthropic, OpenAI, and open weight models in a
private benchmark using a bespoke harness.
