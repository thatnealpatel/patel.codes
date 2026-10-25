# Proving Kimberling's harmonic-floor conjecture

2026-10-25

[`A108211`](https://oeis.org/A108211) is the sequence
$a(n)=16n^2+1$, starting at $n=1$:

$$
17,65,145,257,401,577,\ldots.
$$

A conjecture of Clark Kimberling from September
2014 gives a different way to obtain these integers:

> For every $n\geq1$,
> $$
> \left\lfloor
> \frac{1}{\frac1{4n}-\log2+H_{2n}-H_n}
> \right\rfloor=16n^2+1.
> $$

Here $H_n=1+1/2+\cdots+1/n$ is the $n$th
[harmonic number](https://en.wikipedia.org/wiki/Harmonic_number), $\log$ is the natural logarithm,
and $\lfloor x\rfloor$ rounds $x$ down to an integer.
For example, at $n=1$ the denominator is
$3/4-\log2$, whose reciprocal lies between $17$
and $18$.

**This conjecture is true.** The proof compiles
using only `propext`, `Classical.choice`, and
`Quot.sound`.

## Proof

The definition and theorem statement below omit
the proof body:

:::gen

```lean
def a (n : ℕ) : ℕ := 16 * n ^ 2 + 1

theorem kimberling_conjecture (n : ℕ) (hn : 1 ≤ n) :
    ⌊(1 : ℝ) / (1 / (4 * (n : ℝ)) - Real.log 2 +
      ∑ k ∈ Icc (n + 1) (2 * n), (1 : ℝ) / k)⌋
      = (a n : ℤ)
```

:::

Here `Icc (n + 1) (2 * n)` includes both endpoints,
so the sum is $1/(n+1)+\cdots+1/(2n)$, exactly
$H_{2n}-H_n$. The index `k` simply runs through
those denominators, replacing OEIS's written-out
sum without changing its value. The divisions take place in the
real numbers, and the floor is an integer. The
hypothesis `hn` restricts the claim to $n\geq1$,
matching the sequence's first index.

Write the denominator as

$$
x_n=\frac1{4n}-\log2+H_{2n}-H_n.
$$

It is enough to prove

$$
\frac1{16n^2+2}<x_n<\frac1{16n^2+1}.
$$

Taking reciprocals will then put $1/x_n$ strictly
between the consecutive integers $16n^2+1$ and
$16n^2+2$.

### Comparing consecutive terms

:::gen

Subtracting consecutive terms cancels $\log2$
and almost all of the harmonic sums:

$$
\begin{aligned}
x_n-x_{n+1}
&=\frac1{4n}-\frac1{4(n+1)}
  +\frac1{n+1}-\frac1{2n+1}-\frac1{2n+2}\\
&=\frac1{4n(n+1)(2n+1)}.
\end{aligned}
$$

Also, $H_m-\log m$ converges to the
[Euler–Mascheroni constant](https://en.wikipedia.org/wiki/Euler%27s_constant). Consequently,

$$
x_n=\frac1{4n}
 +(H_{2n}-\log(2n))-(H_n-\log n)
 \longrightarrow0.
$$

Now introduce the two comparison sequences

$$
L_n=\frac1{16n^2+2},
\qquad
U_n=\frac1{16n^2+1}.
$$

Both tend to zero. Their consecutive differences
trap the difference of $x_n$:

$$
L_n-L_{n+1}<x_n-x_{n+1}<U_n-U_{n+1}
\qquad(n\geq1).
$$

These are exact rational inequalities. To check
them, put $D_n=4n(n+1)(2n+1)$. Direct subtraction
gives

$$
(U_n-U_{n+1})-\frac1{D_n}
=
\frac{32n^2+32n-17}
{D_n(16n^2+1)(16(n+1)^2+1)}>0
$$

and

$$
\frac1{D_n}-(L_n-L_{n+1})
=
\frac{36}
{D_n(16n^2+2)(16(n+1)^2+2)}>0.
$$

Every denominator is positive, and both numerators
are positive for $n\geq1$.

:::

### From differences to bounds

:::gen

The upper comparison says

$$
x_n-U_n<x_{n+1}-U_{n+1}.
$$

Thus $x_n-U_n$ is strictly increasing. Since it
converges to zero, every term must be negative:
A nonnegative term would force all later terms
to stay above some positive value. Hence

$$
x_n<U_n=\frac1{16n^2+1}.
$$

Similarly, the lower comparison says

$$
x_n-L_n>x_{n+1}-L_{n+1}.
$$

This difference is strictly decreasing to zero,
so every term is positive. Therefore

$$
\frac1{16n^2+2}=L_n<x_n.
$$

In particular, $x_n>0$. Taking reciprocals
reverses the inequalities and gives

$$
16n^2+1<\frac1{x_n}<16n^2+2.
$$

Rounding down proves

$$
\boxed{\left\lfloor\frac1{x_n}\right\rfloor=16n^2+1
\qquad(n\geq1).}
$$

:::

## Source

The source code for this proof can be found
[on my GitHub](https://github.com/thatnealpatel/proofs/blob/e0c9810910fd2faa896326257d0a7b95216756cc/Proofs/Analysis/A108211HarmonicFloor.lean).

## AI Disclosure

These results were obtained primarily using a mix
of Anthropic, OpenAI, and open weight models in a
private benchmark using a bespoke harness.
