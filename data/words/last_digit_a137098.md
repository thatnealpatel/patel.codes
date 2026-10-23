# Proving last digit conjecture in `A137098`

2026-10-23

[`A137098`](https://oeis.org/A137098) is a sequence
consisting of numbers $k$ for which both $k$ and
$k^2$ use only the decimal digits

$$
D=\{2,4,5,7,9\}.
$$

A conjecture from Jan 2016 states that every term
beyond the first two ends in $7$. More precisely,

> If $k$ belongs to `A137098` and $k\ne2,5$, then
> $k\equiv7\pmod{10}$.

**This conjecture is true.** The proof machine-checked
and compiles with only `propext`, `Classical.choice`,
`Quot.sound`.

## Proof

:::gen

```lean
def allowedDigits : Finset ℕ := {2, 4, 5, 7, 9}

def UsesAllowed (n : ℕ) : Prop :=
  ∀ d ∈ Nat.digits 10 n, d ∈ allowedDigits

def IsTerm (k : ℕ) : Prop :=
  0 < k ∧ UsesAllowed k ∧ UsesAllowed (k ^ 2)

theorem last_digit_seven_except_first_two (k : ℕ) (h : IsTerm k)
    (hk2 : k ≠ 2) (hk5 : k ≠ 5) : k % 10 = 7
```

:::

We use the  membership hypothesis: To prove the statement,
first square the possible units digits modulo 10:

$$
2^2\equiv4,\quad4^2\equiv6,\quad5^2\equiv5,
\quad7^2\equiv9,\quad9^2\equiv1.
$$

Since 6 and 1 are forbidden, a term can end only in 2, 5, or 7.

Suppose a multidigit term ends in 2. Write $k=10a+2$ and let
$t=a\bmod10$, the tens digit of $k$. Expanding the square gives

$$
k^2=100a^2+40a+4.
$$

To find the tens digit of $k^2$, we need only its last two digits.
Adding $100a^2$, a multiple of 100, leaves both of those digits
unchanged. Writing $a=10b+t$ also gives $40a=400b+40t$, hence

$$
k^2\equiv40t+4\pmod{100}.
$$

The final $+4$ supplies the units digit without a
carry into the tens place. The tens digit is
therefore $4t\bmod10$. Because $k$ has at least two
digits, $t$ must belong to $D$. The five possibilities
are as follows:

:::gen

| Tens digit $t$ of $k$ | $4t$ | Tens digit of $k^2$ |
|---|---|---|
| 2 | 8 | 8 |
| 4 | 16 | 6 |
| 5 | 20 | 0 |
| 7 | 28 | 8 |
| 9 | 36 | 6 |

:::

Like before, these are forbidden. Thus, no term
beyond the first two can end in 2.

Similarly, if a multidigit term ends in 5, write
$k=10a+5$. Then

$$
k^2=100a(a+1)+25.
$$

Its hundreds digit is $t(t+1)\bmod10$. For $t\in D$,
that is only 0 or 6: Forbidden. Thus only the single
digit term 5 can end in 5.

**Every remaining term must end in 7.**

## Source

The source code for this proof can be found
[on my GitHub](https://github.com/thatnealpatel/proofs/blob/e0c9810910fd2faa896326257d0a7b95216756cc/Proofs/NumberComplexity/A137098BothUseDecimal.lean).

## AI Disclosure

These results were obtained primarily using a mix
of Anthropic, OpenAI, and open weight models in a
private benchmark using a bespoke harness.
