# MEPL Project TODO List

> CS-C2170 Modern and Emerging Programming Languages - Statically-typed programming language project
> Total achievable: 5000p (capped) out of 6450 theoretical max

## Recommended Build Order

Lexer → Parser → Type Checker → Evaluator → REPL → Features (integers/booleans/variables → functions → pairs/lists/sums → recursion → declarations/imports) → Docs/Examples last

---

## Core Infrastructure

- [ ] **Build lexer/tokenizer** — Tokenize source code: keywords, operators, literals, identifiers, punctuation
- [ ] **Build parser (AST generation)** — Produce AST from tokens: let bindings, integers, booleans, if-then-else, lambdas, application, pairs, lists, sums, case matching, fix, declarations, type annotations
- [ ] **Build static type checker** — Verify types before evaluation: Int, Bool, arrow types (A -> B), pair types, list types, sum types. Must reject type errors before evaluation
- [ ] **Build evaluator/interpreter** — Tree-walking interpreter: variable lookup with scoping, arithmetic, boolean ops, closures, pairs/projections, list ops, sum matching, recursion via fix

---

## Section 1: Getting Started (700p) — MANDATORY

- [ ] **1.1 Running the project (300p)** — README.md with clear instructions for compiling/running on Linux, how to run example files, optional Dockerfile. *Prerequisite for all other points*
- [ ] **1.2 Language introduction (200p)** — High-level overview: syntax for terms/expressions, types, declarations, and examples
- [ ] **1.3 Exercises (200p)** — EXERCISES.md with 5+ programming exercises (problem statements + solutions), each focusing on a different feature, distinct from example files

---

## Section 2: Language Features (3100p core + 1150p bonus)

Each feature needs: a working example file + a type error example file in `examples/`

### Core Features

- [ ] **2.1 Variables (250p)**
  - [ ] Can create variables using `let` (125p)
  - [ ] Can access values in the right scope (125p)
  - [ ] Example: `examples/variables.mepl`

- [ ] **2.2 Integers (200p)**
  - [ ] Integer literals (50p)
  - [ ] Arithmetic: +, -, *, / (100p)
  - [ ] Type error example (50p)
  - [ ] Examples: `examples/integers.mepl`, `examples/integers-errors.mepl`

- [ ] **2.3 Booleans (250p)**
  - [ ] Boolean literals true/false (50p)
  - [ ] Operations: if-then-else, and, or, not (100p)
  - [ ] Type error example, e.g. if branches with different types (100p)
  - [ ] Examples: `examples/booleans.mepl`, `examples/booleans-errors.mepl`

- [ ] **2.4 Functions (600p)**
  - [ ] Anonymous functions / lambda abstractions (150p)
  - [ ] Function application (150p)
  - [ ] Closures capture environments correctly (100p)
  - [ ] Arrow type syntax A -> B (150p)
  - [ ] Type error example (50p)
  - [ ] Examples: `examples/functions.mepl`, `examples/functions-errors.mepl`

- [ ] **2.5 Pairs (200p)**
  - [ ] Construct pairs (50p)
  - [ ] Destructure with fst/snd (50p)
  - [ ] Type error example (100p)
  - [ ] Examples: `examples/pairs.mepl`, `examples/pairs-errors.mepl`

- [ ] **2.6 Lists (400p)**
  - [ ] Construct and combine lists (100p)
  - [ ] Destructure using case matching (150p)
  - [ ] Recursive operations: map, fold (100p)
  - [ ] Type error example (50p)
  - [ ] Examples: `examples/lists.mepl`, `examples/lists-errors.mepl`

- [ ] **2.7 Sums (200p)**
  - [ ] Construct sums with inl/inr (50p)
  - [ ] Destructure using case matching (50p)
  - [ ] Type error example (100p)
  - [ ] Examples: `examples/sums.mepl`, `examples/sums-errors.mepl`

- [ ] **2.8 Recursion (400p)**
  - [ ] Fixed-point operator (fix) (150p)
  - [ ] Evaluation of recursive calls (100p)
  - [ ] Type error example, e.g. fix on wrong type (150p)
  - [ ] Examples: `examples/recursion.mepl`, `examples/recursion-errors.mepl`

- [ ] **2.9 Declarations & Imports (600p)**
  - [ ] Named function/constant declarations (150p)
  - [ ] Import code from another file (150p)
  - [ ] Type-checked declarations with error examples (200p)
  - [ ] Examples: `examples/declarations.mepl`, `examples/declarations-errors.mepl`, `examples/imports.mepl`

### Bonus Features

- [ ] **2.10 Comments (100p B)** — Single/multi-line comments → `examples/comments.mepl`
- [ ] **2.11 Output/Printing (100p B)** — Print basic types (50p) + all values including functions (50p) → `examples/printing.mepl`
- [ ] **2.12 Tuples (150p B)** — Construct (50p), access/destructure (50p), type errors (50p) → `examples/tuples.mepl`, `examples/tuples-errors.mepl`
- [ ] **2.13 Strings & Characters (150p B)** — Construct/combine (50p), char access (50p), type errors (50p) → `examples/strings.mepl`, `examples/strings-errors.mepl`
- [ ] **2.14 Records (150p B)** — Construct (50p), field access (50p), type errors (50p) → `examples/records.mepl`, `examples/records-errors.mepl`
- [ ] **2.15 Algebraic Data Types (300p B)** — N-ary sums/custom types (100p), pattern matching (200p) → `examples/adt.mepl`
- [ ] **2.16 Loops (200p B)** — For or while loops (200p) → `examples/loops.mepl`

---

## Section 3: Working Interpreter (1200p + 300p bonus)

- [ ] **3.1 REPL (700p)** — Interactive evaluation, basic computational features, type-check before eval
- [ ] **3.2 REPL from file (500p)** — Start REPL with a file loaded so all its declarations are available
- [ ] **3.3 Informative error messages (300p B)** — Useful error info (100p) + source line/column numbers (200p)

---

## Points Summary

| Section | Core | Bonus |
|---------|------|-------|
| 1. Getting Started | 700 | 0 |
| 2. Language Features | 3100 | 1150 |
| 3. Working Interpreter | 1200 | 300 |
| **Total** | **5000** | **1450** |

> Cap: 5000p. Bonus points help reach the cap if core features are missing.
