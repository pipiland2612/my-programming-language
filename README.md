# MEPL — Modern and Emerging Programming Language

MEPL is a statically-typed functional programming language with an interactive REPL, implemented in Go.

## Getting Started

### Running with Docker

If you do not have Go installed locally, you can still build and run MEPL using Docker alone.

```bash

docker build -t mepl .
docker run --rm mepl examples/integers.mepl
```

```bash
# Using REPL mode
docker run --rm -it mepl   # REPL mode
```

### Running with Go
Required **Go 1.25** or later

### Quick Start
```bash
# Build and run directly
go build -o mepl . && ./mepl examples/integers.mepl
```

### Building

```bash
go build -o mepl .
```

### Running Example Files

```bash
# Run any example file
./mepl examples/integers.mepl
./mepl examples/functions.mepl
./mepl examples/lists.mepl

# Or without building first
go run . examples/integers.mepl
```

### Running the REPL

```bash
# Start an interactive REPL
./mepl

# Load a file and then enter the REPL (all declarations available)
./mepl examples/declarations.mepl -i
```

### Running Tests

```bash
go test ./...
```

---

## Language Introduction

MEPL is an expression-based language where everything evaluates to a value. It features static type checking, first-class functions, algebraic data types, and pattern matching.

### Types

| Type | Syntax | Description |
|------|--------|-------------|
| Integer | `Int` | Whole numbers |
| Boolean | `Bool` | `true` or `false` |
| String | `String` | Text in double quotes |
| Unit | `Unit` | The unit value `()` |
| Function | `A -> B` | Function from `A` to `B` |
| Pair | `(A, B)` | Two-element pair |
| Tuple | `(A, B, C, ...)` | N-element tuple (3+) |
| List | `[A]` | Homogeneous list |
| Sum | `A + B` | Tagged union (left or right) |
| Record | `{x: Int, y: Bool}` | Named fields |
| ADT | `type Option = None | Some of Int` | Named algebraic data type |

### Variables and Declarations

```
// Top-level declaration
let x : Int = 42

// Let-in expression (local binding)
let y = 10 in y + 1
```

### Integers and Arithmetic

```
let a = 10 + 3 * 2     // 16
let b = 10 % 3          // 1
let c = -5              // unary minus
```

### Booleans and Conditionals

```
let x = true && (false || true)
let y = if x then 1 else 0
let z = 3 == 3          // true
let w = 5 < 10          // true
```

### Functions

Functions are defined with `fn`, applied by juxtaposition (curried by default):

```
// Single-argument function
let inc : Int -> Int = fn (x: Int) => x + 1

// Multi-argument (curried)
let add : Int -> Int -> Int = fn (x: Int) => fn (y: Int) => x + y
let result = add 3 4    // 7
```

### Pairs and Tuples

```
// Pairs (2 elements)
let p : (Int, Bool) = (1, true)
let a = fst p           // 1
let b = snd p           // true

// Tuples (3+ elements) with index access
let t : (Int, Bool, String) = (1, true, "hello")
let x = t.0             // 1
let y = t.1             // true
let z = t.2             // "hello"
```

### Lists

```
let xs : [Int] = [1, 2, 3]
let ys = 0 :: xs                    // [0, 1, 2, 3]
let empty = [] : [Int]              // empty list with type annotation

// Case matching on lists
let sum : [Int] -> Int = fix (fn (self: [Int] -> Int) =>
  fn (l: [Int]) => case l of
    | [] => 0
    | h :: t => h + (self t))
```

### Sums (Tagged Unions)

```
let x : Int + Bool = inl 42 as Int + Bool
let y : Int + Bool = inr true as Int + Bool

let describe = fn (v: Int + Bool) => case v of
  | inl n => n + 1
  | inr b => if b then 1 else 0
```

### Records

```
let person : {name: String, age: Int} = {name = "Alice", age = 30}
let n = person.name     // "Alice"
let a = person.age      // 30
```

### Algebraic Data Types

```
type Option = None | Some of Int

let getValue : Option -> Int = fn (opt : Option) =>
  case opt of
    | None => 0
    | Some n => n
```

### Recursion

Recursion uses the fixed-point operator `fix`:

```
let factorial : Int -> Int = fix (fn (f: Int -> Int) =>
  fn (n: Int) => if n == 0 then 1 else n * (f (n - 1)))
```

### Loops

```
for i = 0 to 10 do
  println i
end
```

### Strings

```
let s = "Hello" + " " + "World"     // concatenation
let len = length s                    // 11
let ch = charAt s 0                   // "H"
```

### Printing

```
print "no newline"
println "with newline"
println 42
println [1, 2, 3]
```

### Comments

```
// Single-line comment

/* Multi-line
   comment */
```

### Declarations and Imports

```
// In mathlib.mepl
let square : Int -> Int = fn (x: Int) => x * x

// In main.mepl
import "mathlib.mepl"
println (square 5)       // 25
```

---

## Example Files

All examples are in the `examples/` directory. For each feature, there is a working example and a type-error example:

| Feature | Working Example | Type Error Example |
|---------|----------------|-------------------|
| Variables | `variables.mepl` | `variables-errors.mepl` |
| Integers | `integers.mepl` | `integers-errors.mepl` |
| Booleans | `booleans.mepl` | `booleans-errors.mepl` |
| Functions | `functions.mepl` | `functions-errors.mepl` |
| Pairs | `pairs.mepl` | `pairs-errors.mepl` |
| Lists | `lists.mepl` | `lists-errors.mepl` |
| Sums | `sums.mepl` | `sums-errors.mepl` |
| Recursion | `recursion.mepl` | `recursion-errors.mepl` |
| Declarations | `declarations.mepl` | `declarations-errors.mepl` |
| Imports | `imports.mepl` | `imports-errors.mepl` |
| Comments | `comments.mepl` | `comments-errors.mepl` |
| Printing | `printing.mepl` | `printing-errors.mepl` |
| Tuples | `tuples.mepl` | `tuples-errors.mepl` |
| Strings | `strings.mepl` | `strings-errors.mepl` |
| Records | `records.mepl` | `records-errors.mepl` |
| Loops | `loops.mepl` | `loops-errors.mepl` |
| ADTs | `adt.mepl` | `adt-errors.mepl` |
