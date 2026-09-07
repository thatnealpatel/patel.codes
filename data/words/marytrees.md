# Yet another "final" configuration spec

2025-08-05

## Update (2026)

I sunsetted the project that this was particularly central to;
however, from time of writing (Aug 2025) to present (May 2026),
this configuration spec remained unchanged.

## Motivations

I have a large personal project whose central dependency
is a configuration spec that dictactes events at runtime.
In this project, it is extremely important that the spec
is trivially evolvable to address unknown unknowns.

## Config v0 (2021)

The first iteration of my spec was simply JSON; everything
was completely adhoc. Arguably, this is where one should
start. There was no need to prematurely optimize.

This quickly become unmanagable; however, being naive
and without any practical design experience, I doubled
down and introduced mountains tech debt to compensate.

Two years of string parsing and esoteric semantics later,
it was time for a change.

## Config v1 (2023)

The second time around, I thought I had such a deep
understanding of my requirements that I could finally
write the "final version" of my configuration spec.

Initially, I was very excited by the simplicity in
using native structs and doing away with messy string
comparisons and untyped invarants baked into the existing
JSON schema.

For a while, this sufficed. When I needed new features
or invariants, I simply extended the configuration;
however, this config became complex, verbose, and
repetitive: 

```
type config struct {
    condition1 []struct {
        enabled bool
        type string
        minval float64
        val    float64
        maxval float64
    }
    option1a struct {
        enabled bool
        type string
        minval float64
        val    float64
        maxval float64
    }
    option1b struct {
        enabled bool
        type string
        minval float64
        val    float64
        maxval float64
    }
    condition2 []struct {
        enabled bool
        type string
        minval float64
        val    float64
        maxval float64
    }
}
```

There were many conditions which forced duplications of
application logic  this also made parsing out the config
in application logic extremely convoluted.

Parsing code eventually took the following form:

```
if config.option1a.enabled {
    var met bool
    for _, cond := range condition1 {
        if !cond.enabled {
            break
        }
        if cond.type == "qux" {
        } else if cond.type == "quux" {
        } else {
        }
    }
    if !met {
        return
    }
    // do something with option1a
    if options1b.enabled {
        // do something with option1b
        // in relation to option1
    }
}

if len(config.condition2) > 0 {
    for _, cond := range condition2 {
        if !cond.enabled {
            break
        }
        if cond.type == "qux" {
        } else if cond.type == "quux" {
        } else {
        }
    }
}
```

## Config v2 (2025)

I took a step back and started thinking more deeply
about my requirements; I had a configuration whose
purpose is to define the invariants over which dynamic
runtime inputs would be evaluated. It took me nearly 3
years to realize that I was writing a suboptimal SAT
parser.

Seeing as how I needed to involve the configuration
language to be more arbitrarily expressive anyways,
I realized that a great fit for my use case were m-ary
trees with a somewhat odd convention:


```
type config struct {
    enabled bool
    crit    *criteria
}

type criteria struct {
    type criteriaT
    val  any
    and  *criteria
    or   []*criteria
}

// example of a `val`
type lvh struct { lo, val, hi float64 }

type criteriaT string

const(
    // a criteria for which using
    // lvh is a natural choice of
    // expression.
    boundedCriteria criteriaT = "BOUNDED_CRITERIA"
)
```

Conceptually, the left-most node in the tree would be
the `and` node; the rest of the nodes, if any, would be
the `or` nodes. A root `criteria` (present in the `config`)
is said to be met iff its tree evalutes to true.

This lends itself to the very nice implementation:

```go
func walk(c *criteria, fn (*critiera) bool) (sat bool) {
    if c == nil {
        return true
    }
    sat = fn(c.val)
    for _, o := range c.or {
        sat = sat || fn(o)
    }
    return sat && fn(c.and)
}
```

With a little extra thinking, this implementation both reduced
the amount of esoteric code in my codebase and made reasoning
about new changes to my configuration spec extremely easy. The 
power I found in this design lies in how call-sites neatly
call `walk` in the following manner:

```
walk(config, func(c *criteria) bool {
    switch val := c.val.(type) {
    case foo:

    case bar:

    case baz:
        switch c.type {
        case "qux":
        case "quux":
        }
    }
})
```

This allowed for call-sites, regardless of their purpose,
to instrument the logic plainly without additional parsing
or sematic interpretation. Some call-sites actually desire
to walk the entire tree in which case the original `walk`
function is modified to contain no invariant tracking:

```
func walkall(c *criteria, fn (*critiera)) {
    if c == nil {
        return
    }
    fn(c.val)
    for _, o := range c.or {
        fn(o)
    }
    fn(c.and)
    return
}
```

Though I found myself saying this in the past, I am slightly
more convinced this time that this design will remain a
fixture for me. It's already proven to be as ubiquitous
and extensible as I had hoped. That being said, in writing this
it occurred to me that it may be better to simply use `val` 
with some repetitive, named types instead of shared types
distingushed by `criteriaT` enums.

```
type config struct {
    enabled bool
    crit    *criteria
}

type criteria struct {
    val  any
    and  *criteria
    or   []*criteria
}

type boundedLVH struct { lo, val, hi float64 }
type anotherBoundedLVH struct { lo, val, hi float64 }
```

## Final thoughts

As always, the code you wrote a year ago was written by a fool;
it's fun to be able to look back and laugh at the mistakes
you've made without realizing that you are only setting 
yourself up for a future punchline.
