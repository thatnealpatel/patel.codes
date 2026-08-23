# proving `A092482`'s closed form 

2026-08-23

## motivations

i am currently engaged in various side
research projects, one of which requires
me to build an evaluation set for grading
the efficacy of an auto-research program
built natively into a custom harness i
have been rolling from scratch in Go.

i originally intended the eval set to be
a set of formalization tasks and a few
proof tasks over known results. i suppose
i should not have been surprised that i
picked off some low-hanging fruit in the
process.

speaking of motivation, i was unfortunately
not particularly motivated to spend more
time on distilling and reducing the core
kernels of the proof when  writing this post;
as such, you might notice that a larger portion
of it is generated than usual. while i did spend
hours in exploring and writing this post, i
hope that it does not detract from your reading
experience.

## the sequence and the typo

a three-term arithmetic progression, or 3-AP, is a triple $x<y<z$
satisfying $x+z=2y$; once some terms have been chosen, we say a
candidate $m$ is **blocked** if two earlier terms $x<y<m$ satisfy

$$
x+m=2y
$$

the greedy rule chooses the least candidate that is not blocked, with
`(1,2,3)` as the sole permitted exception.

[`A092482`](https://oeis.org/A092482) is the greedy
increasing sequence that contains no 3-AP other than
its initial terms `1, 2, 3`

```
1, 2, 3, 6, 7, 14, 15, 17, 18, 36, 37, 39, 40, 45, 46, 48, 49, 98, ...
```

in prose, it is defined as

```
a(1)=1, a(2)=2, a(3)=3; a(n) is least k such that no three terms of
a(1), a(2), ..., a(n-1), k form an arithmetic progression, except
for the first triple (1,2,3).
```

and the sequence contains a comment with the conjectured closed form

```
For n > 2, a(n+2) = 1 + 2^floor(log_2(n)) + Sum_{k=1..n}
(3^A007814(n) + 1)/2 = 1 + A053644(n) + A005836(n)
(conjectured and checked up to n=512).
```

the summand contains a typo: `A007814(n)` should be `A007814(k)`.
correcting it gives

$$
a(n+2)=1+2^{\lfloor\log_2 n\rfloor}
 +\sum_{k=1}^{n}\frac{3^{A007814(k)}+1}{2}
$$

let's define $\tau(n)$ as *reading the binary digits of $n$ as a base-3 numeral*

$$
0,1,2,3,4,5,6,7
\quad\xmapsto{\tau}\quad
0,1,3,4,9,10,12,13
$$

these are the values listed by `A005836`. in terms of its published
one-based indexing,

$$
\tau(n)=A005836(n+1)
$$

so the final equality in the quoted comment should read

$$
a(n+2)=1+A053644(n)+A005836(n+1)
$$

equivalently, the repaired theorem presents

$$
\boxed{
a(n+2)=1+2^{\lfloor\log_2 n\rfloor}+\tau(n)
}
\qquad(n\geq1)
$$

which is ever so slightly stronger than the given `n >2`
stated on the entry. in the zero-indexed Lean definition,
`greedySeq r = a(r+1)` and the precise formalization as

:::gen

```lean
theorem greedySeq_add_two (m : ℕ) :
    greedySeq (m + 2) =
      1 + 2 ^ Nat.log 2 (m + 1) + binToTernary (m + 1)
```

:::

here `Nat.log 2 n` is $\lfloor\log_2n\rfloor$ and `binToTernary` is
$\tau$. at $n=1,2,3$, the formula gives `3,6,7`, respectively. only
$n=1,2$ extend the stated range $n>2$: they give the final seed value `3`
(block $L=0$) and the first value after the seed, `6` (the start of block
$L=1$). the $n=3$ case, giving `7`, was already in the stated range and is
the second value of block $L=1$.

## why the sum is `tau`

incrementing a binary number flips its trailing ones
to zeroes and carries a new one. if `k` has $\nu_2(k)$
trailing zeroes, the step from `k-1` to `k`, read in
base-3, is

$$
\tau(k)-\tau(k-1)=\frac{3^{\nu_2(k)}+1}{2}
$$

here $\nu_2(k)$ is the
[largest $j$ such that $2^j$ divides $k$](https://en.wikipedia.org/wiki/P-adic_valuation);
it is also the number of *trailing zeroes* in the binary expansion of
positive $k$. put $j=\nu_2(k)$. then $k-1$ ends in exactly $j$ binary ones.
incrementing replaces those ones by zeroes (read: forces a carry) and changes
the preceding zero to one. reading the same change in base 3 gives

$$
3^j-\sum_{i=0}^{j-1}3^i
 =3^j-\frac{3^j-1}{2}
 =\frac{3^j+1}{2}
$$

for example,

```text
k-1  = 3   011_2   tau(011_2) => 011_3 = 4
k    = 4   100_2   tau(100_2) => 100_3 = 9  # diff of 5=(3^2 + 1)/2
```

where `_{n}` is short-hand for the base and
the absence of `_{n}` implies base-10.

telescoping from $\tau(0)=0$ gives

$$
\tau(n)=\sum_{k=1}^{n}\frac{3^{\nu_2(k)}+1}{2}
$$

Lean records the same identity without natural-number division:

:::gen

```lean
theorem two_mul_binToTernary_eq_sum (n : ℕ) :
    2 * binToTernary n =
      ∑ k ∈ Finset.Icc 1 n, (3 ^ padicValNat 2 k + 1)
```

:::

here `padicValNat 2 k` is $\nu_2(k)$, and `Finset.Icc 1 n` is the
integer interval $\{1,\ldots,n\}$.


## the block geometry

similarly to my previous post, i just needed to *see* things:
not an uncommon feeling in these contexts.

an explanation for the proof centers around decomposing and
partitioning the sequence terms into blocks; with that, you
can make arguments about which subsets belong in which blocks
and how moving within and between blocks gives the greedy
sequence by induction.

so, starting with how we construct the blocks

:::gen

```text
{1, 2} | 3 | 6 7 | 14 15 17 18 | 36 37 39 40 45 46 48 49 | 98 ...
        L=0  L=1       L=2                 L=3                L=4
```

:::

let

$$
B_L=2^L+3^L+1
$$

and let $T_L$ be the set of *offsets* ($T_L\subseteq\mathbb{N}$)
below $3^L$ whose ternary digits are all `0` or `1`.
equivalently,

$$
T_L=\{\tau(r):0\leq r<2^L\}
$$

in the set notation below,

$$
B_L+T_L=\{B_L+t:t\in T_L\}
$$

viewing offsets as length-$L$ ternary strings padded with leading zeroes
(with the empty string representing `0` when $L=0$), each digit may
independently be `0` or `1`, so

$$
 | T_L |=2^L,
\qquad
\max T_L=\frac{3^L-1}{2},
\qquad
2t<3^L\quad(t\in T_L)
$$

splitting by the leading ternary digit also gives

$$
T_{L+1}=T_L\cup(3^L+T_L)
$$

we then propose the closed form $V$

$$
V=\{1,2\}\cup\bigcup_{L\geq0}(B_L+T_L)
$$

for the first few blocks, we observe

:::gen

| $L$ | $B_L$ | $T_L$ | $B_L+T_L$ |
|---:|---:|---|---|
| 0 | 3 | $\{0\}$ | $\{3\}$ |
| 1 | 6 | $\{0,1\}$ | $\{6,7\}$ |
| 2 | 14 | $\{0,1,3,4\}$ | $\{14,15,17,18\}$ |
| 3 | 36 | $\{0,1,3,4,9,10,12,13\}$ | $\{36,37,39,40,45,46,48,49\}$ |

:::

this rewrites the closed form as a block plus an offset within that block.
for example, take the formula parameter $n=6$, which corresponds to the
sequence term $a(8)$. then $n=2^2+2$, so $L=2$, $r=2$, and
$\tau(r)=10_3=3$. the term is

$$
a(8)=B_2+\tau(2)=14+3=17.
$$

every positive `n` has a unique decomposition

$$
n=2^{\lfloor\log_2n\rfloor}+r,
\qquad
0\leq r<2^{\lfloor\log_2n\rfloor}.
$$

splitting off the leading binary one and reading in base 3 gives

$$
\tau(n)=3^{\lfloor\log_2n\rfloor}+\tau(r).
$$

therefore

$$
1+2^{\lfloor\log_2n\rfloor}+\tau(n)
 =1+2^{\lfloor\log_2n\rfloor}+3^{\lfloor\log_2n\rfloor}+\tau(r)
 =B_{\lfloor\log_2n\rfloor}+\tau(r).
$$

as i understand it, this is just the fancy way of saying:
in base-2, $n$ has a leading `1` followed by exactly
$\lfloor\log_2n\rfloor$ bits: this resolves which block. the values of these
bits form $r$. when read in base-2, $r$ is the intra-block
index; when read in base-3, $\tau(r)$, it is the actual
offset from $B_{\lfloor\log_2n\rfloor}$.

we capture this as `closedForm_eq_block` in the Lean proof.

## where the power-of-two term comes from

[`A005836`](https://oeis.org/A005836) sequences a set of
nonnegative integers such that all base-3 expansions contain
only `0` and `1`. adding one to each term gives [`A003278`](https://oeis.org/A003278),
the ordinary greedy 3-AP-free sequence beginning at `1`:

```
1, 2, 4, 5, 10, 11, 13, 14, ...
```

notably,
[lemma 6.4 of Moy and Rolnick's *Novel structures in Stanley sequences*](https://arxiv.org/abs/1502.06013)
proves that these Stanley numbers (`A005836`) are **3-AP-free and greedy**:
no three distinct terms form an arithmetic progression, and every
omitted number completes one with two earlier terms. (also, yes i linked
the pre-print on purpose.)

`A092482` differs in that it  explicitly permits `(1,2,3)`.

for an exact comparison, let

$$
S(r)=A003278(r+1)=1+\tau(r)
$$

$$
G(r)=A092482(r+1)
$$

both with zero-based arguments. in the Lean source these functions are
`stanleyGreedy` and `greedySeq`, respectively. block $L$ begins at
index $2^L$ in $S$ but at index $2^L+1$ in $G$. thus the table below
compares corresponding block starts, not equal sequence indices.
subtracting the two block formulas below shows the displacement created
by accepting the seed rule.

:::gen

```text
L:                  0       1       2       3       4
A092482 start:      3       6      14      36      98
zero-based index:   2       3       5       9      17
A003278 start:      2       4      10      28      82
zero-based index:   1       2       4       8      16
difference:         1       2       4       8      16
```

:::

the difference doubles from one block to the next. formally,

$$
\operatorname{greedySeq}(2^L+1)
 =\operatorname{stanleyGreedy}(2^L)+2^L,
$$

:::gen

this is `greedySeq_defect`. it is the $r=0$ case of the stronger blockwise
identity: whenever $n=2^L+r$ with $0\leq r<2^L$,

$$
S(n) =1+3^L+\tau(r)
$$

$$
G(n+1) =1+3^L+\tau(r)+2^L
$$

therefore the entire `A092482` block is the corresponding ordinary Stanley
block translated by $2^L$. since $L=\lfloor\log_2n\rfloor$, this translation
is exactly the $2^{\lfloor\log_2n\rfloor}$ term in the closed form.

accepting `3` explains how the displacement begins: `4` is then blocked by
`(2,3,4)` and `5` by `(1,3,5)`, so the next term is `6`. the later
admissibility and covering arguments prove that the translated blocks continue
to obey the greedy rule.

:::

## the proof

earlier, we constructed $V$ such that

$$
V=\{1,2\}\cup\bigcup_{L\geq0}(B_L+T_L).
$$

and, in the proof, we claim that enumeration
of this set (in increasing order) is exactly
the greedy process that gives us the closed
form.

in order to prove this, we need further show:

:::gen

1. **admissibility:** no three distinct terms of $V$ form an arithmetic
   progression except `(1,2,3)`;
   
2. **minimality:** every integer greater than `2` outside $V$ is blocked by two
   smaller terms of $V$.

induction on the prefix length then forces the greedy process to
select exactly the increasing enumeration of $V$.

:::

## admissibility

:::gen

suppose $a<b<c$ are terms of $V$ and $a+c=2b$. if $c=3$, the only
possibility is the permitted progression `(1,2,3)`. now suppose that $c$ lies
in block $L>0$.

first note that twice any term before block $L$ is at most $B_L$.

indeed, the largest such term lies in block $L-1$, and

$$
2\left(B_{L-1}+\max T_{L-1}\right)
 =2\left(2^{L-1}+3^{L-1}+1+\frac{3^{L-1}-1}{2}\right)
 =B_L.
$$

if $b$ were in an earlier block, then $2b\leq B_L$, whereas

$$
2b=a+c>c\geq B_L,
$$

a contradiction. hence $b$ and $c$ lie in the same block. write

$$
b=B_L+s,
\qquad
c=B_L+t,
\qquad
s<t,
\qquad
s,t\in T_L.
$$

then

$$
a=B_L+2s-t
$$

if $a$ were earlier than block $L$, the same bound would
give $2a\leq B_L$. but $2t<3^L$, so

$$
2a=2B_L+4s-2t
$$

$$
>2B_L-3^L
$$

$$
=B_L+2^L+1
$$

$$
>B_L
$$

again a contradiction.

:::

all three terms must lie in block $L$. subtracting
$B_L$ produces a nontrivial 3-AP in $T_L$. but $T_L$
consists of the first $2^L$ values of the zero-reindexed
`A005836`.

we saw above that Moy and Rolnick's `Lemma 6.4` proves
that this ternary-digit sequence is 3-AP-free. this proves
admissibility: the Lean theorem is `noThreeAPExceptSeed_Vset`.

## minimality: within a block

:::gen

now take a candidate in the $L$-th digit window,

$$
m=B_L+u,
\qquad
0\leq u<3^L,
$$

but suppose $m\notin V$. then $u\notin T_L$, so at least one ternary digit of
$u$ is `2`. replace every `2` by `0` to obtain
$\operatorname{keepOnes}(u)$, and by `1` to obtain
$\operatorname{capDigits}(u)$. digit by digit,

$$
\operatorname{keepOnes}(u)+u
 =2\operatorname{capDigits}(u),
$$

and the presence of a `2` gives

$$
\operatorname{keepOnes}(u)
 <\operatorname{capDigits}(u)<u.
$$

both replacements lie in $T_L$. translating by $B_L$ gives

$$
\bigl(B_L+\operatorname{keepOnes}(u)\bigr)+(B_L+u)
 =2\bigl(B_L+\operatorname{capDigits}(u)\bigr).
$$

thus the two smaller block terms block $m$.

for example, $B_2=14$ and the missing offset $u=5=12_3$ gives

```text
keepOnes(12_3) = 10_3 = 3
capDigits(12_3) = 11_3 = 4
```

hence

$$
14+3=17,
\qquad
14+4=18,
\qquad
14+5=19,
$$

with $17+19=2\cdot18$. this is the internal branch of `exists_blocking`.

:::

## minimality: between the blocks

:::gen

the digit argument covers every omitted value in the ambient digit window

$$
[B_L,B_L+3^L).
$$

the actual block $B_L+T_L$ is a sparse subset of this window. it remains to
cover

$$
[B_L+3^L,B_{L+1}),
$$

whose length is

$$
B_{L+1}-(B_L+3^L)=3^L+2^L.
$$

for example, the block-$3$ digit window is `[36,63)`, while the next block
starts at `98`; the inter-block gap is `[63,98)`.

let

$$
\operatorname{Pre}(L)
 =\{1,2\}\cup\bigcup_{j<L}(B_j+T_j),
$$

the terms before block $L$. the recursive invariant is

$$
Q_L:\quad
0\leq\mu<3^L+2^L
\Longrightarrow
\exists\,t\in T_L, p\in\operatorname{Pre}(L),\quad
\mu+p=2t+2^L+1
$$

for $m=B_L+3^L+\mu$, this gives

$$
m+p=2(B_L+t).
$$

both witnesses are earlier: $p<B_L<m$, and $t<3^L$ gives
$B_L+t<B_L+3^L\leq m$.

at $L=0$, $T_0=\{0\}$ and $\operatorname{Pre}(0)=\{1,2\}$. the values
$\mu=0,1$ use $(t,p)=(0,2)$ and $(0,1)$:

$$
4+2=2\cdot3,
\qquad
5+1=2\cdot3.
$$

for the induction step, abbreviate $a=2^L$, $b=3^L$, and
$B_L=a+b+1$. the level-$(L+1)$ range is $[0,2a+3b)$. the following
inclusions will be used:

$$
T_L\subseteq T_{L+1},
\qquad
b+T_L\subseteq T_{L+1},
\qquad
\operatorname{Pre}(L)\subseteq\operatorname{Pre}(L+1).
$$

for any $0\leq u<b$, there are also $s,s'\in T_L$ with $u+s'=2s$. if
$u\in T_L$, take $s=s'=u$; otherwise the `capDigits`/`keepOnes` construction
from the preceding section supplies the witnesses. the four cases are:

| range for $\mu$ | reduced parameter | witnesses at level $L+1$ |
|---|---|---|
| $[0,a)$ | $u=\mu+b-a$, with $0\leq u<b$ | from $u+s'=2s$, take $t=s$, $p=B_L+s'$ |
| $[a,a+b)$ | $\mu'=\mu-a$, with $0\leq\mu'<a+b$ | from $Q_L(\mu')$, keep $t=s$, $p$ |
| $[a+b,a+2b)$ | $u=\mu-a-b$, with $0\leq u<b$ | from $u+s'=2s$, take $t=b+s$, $p=B_L+s'$ |
| $[a+2b,2a+3b)$ | $\mu'=\mu-a-2b$, with $0\leq\mu'<a+b$ | from $Q_L(\mu')$, take $t=b+s$, keep $p$ |

here $t\mapsto b+t$ prefixes a leading ternary `1`, which explains the
second embedding. the four arithmetic checks are

$$
\mu+(B_L+s')
$$

$$
=u+2a+1+s'=2s+2a+1
$$

$$
\mu+p
$$

$$
=\mu'+a+p=2s+2a+1
$$

$$
\mu+(B_L+s')
$$

$$
=u+s'+2a+2b+1=2(b+s)+2a+1
$$

$$
\mu+p
$$

$$
=\mu'+p+a+2b=2(b+s)+2a+1
$$

in every case the final expression is $2t+2^{L+1}+1$, proving
$Q_{L+1}$. this four-case induction is `q_covering`. together with internal
digit blocking, it proves that every nonterm above `2` is blocked by two
smaller terms of $V$.

:::

## some assembly required

:::gen

admissibility is `noThreeAPExceptSeed_Vset`, and minimality is
`exists_blocking`. the formal prefix induction starts from `{1}`.
admissibility makes each next closed-form value legal. in the first two
applications of `nextGreedy_key`, there is no integer strictly between `1` and
`2`, or between `2` and `3`. thereafter every smaller unselected candidate is
greater than `2`, so minimality makes it illegal. the least legal next value
is therefore the next closed-form value, and the induction gives

```lean
theorem greedySeq_eq_closedForm : greedySeq = closedForm
```

the induction derives the seed prefix `{1,2,3}`; the formal corollary recording
this prefix is `prefixSet_two`.

:::

[`Proofs/Enumerative/No3APGreedy.lean`](https://github.com/thatnealpatel/proofs/blob/main/Proofs/Enumerative/No3APGreedy.lean)
contains the complete formal proof: it compiles without
`sorry`, `admit`, or non-standard axioms. the proved
formula also derives the first 57 terms displayed on
the OEIS entry.

## literature check

a cursory literature check done with my harness did
not yield anything adjacent or anything that would
imply a previous formalization exists.

## closing thoughts

i spent *hours* trying to understand this proof; what
frustrated me the most, especially of all the proofs
in my backlog, is that this one seemed to be the one
that is closest to my understanding of maths as a
programmer. yet, it took much longer than expected to
figure out what the argument being made was.

nonetheless, i found the blockwise visualization
to be quite coherent in trying to wrap my head around
the proof: even more, i was surprised by the fidelity
of the "conversations" i had with my research harness
when trying to build intuition.

i will admit: there is something _unsatisfying_ about
formalizing and proving things in this manner; in part,
i think it is because these results do not excite me,
and i would rather be doing other things with my time
than trying to elucidate a result i do not care about.
though, that is not to say that this exercise is not
without material benefit.
