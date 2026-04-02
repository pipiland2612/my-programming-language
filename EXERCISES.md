# MEPL Programming Exercises

A set of exercises for learning MEPL. Each exercise focuses on a different language feature. Solutions are provided at the end.

---

## Exercise 1: Arithmetic and Variables

Write a program that computes the area and perimeter of a rectangle with width 7 and height 3, and prints both values.

**Concepts:** variables, integers, arithmetic, printing

---

## Exercise 2: Boolean Logic and Conditionals

Write a function `clamp` that takes three integers: a value, a minimum, and a maximum. It should return the value if it is within the range, the minimum if the value is too low, or the maximum if the value is too high.

Then test it: `clamp 5 1 10` should return `5`, `clamp -3 0 100` should return `0`, and `clamp 50 0 10` should return `10`.

**Concepts:** functions, booleans, if-then-else, comparisons

---

## Exercise 3: Recursive List Processing

Write a recursive function `sumList` that computes the sum of all integers in a list. Then write a function `filterPositive` that returns a new list containing only the positive numbers from the input list.

Test with: `sumList [1, -2, 3, -4, 5]` should return `3`, and `filterPositive [1, -2, 3, -4, 5]` should return `[1, 3, 5]`.

**Concepts:** lists, recursion (fix), case matching, cons

---

## Exercise 4: Higher-Order Functions

Write a `map` function that applies a given function to every element of an integer list, producing a new list. Then use it to double every element in `[1, 2, 3, 4]`.

**Concepts:** functions as values, closures, higher-order functions, lists

---

## Exercise 5: Pairs and Sum Types

Model a simple result type using sums: a computation can return either a successful integer result (inl) or an error string (inr). Write a function `safeDivide` that takes two integers and returns `inl (a / b)` if b is not zero, or `inr "division by zero"` otherwise. Then write a function `showResult` that pattern-matches on the result and prints an appropriate message.

**Concepts:** pairs, sum types (inl/inr), case matching, functions

---

## Exercise 6: Records and Tuples

Create a record type representing a 2D point with `x` and `y` fields. Write a function `distance` that takes two point records and computes the squared Euclidean distance between them (since we only have integer arithmetic, return `(x2-x1)^2 + (y2-y1)^2`). Use a tuple to return both the two points and their squared distance.

**Concepts:** records, tuples, field access, functions

---

## Exercise 7: Strings and Loops

Write a program that uses a for loop to print each character of the string `"MEPL"` on a separate line, using `charAt` and `length`.

**Concepts:** strings, charAt, length, for loops

---

# Solutions

## Solution 1

```
let width : Int = 7;
let height : Int = 3;

let area : Int = width * height;
let perimeter : Int = 2 * (width + height);

println area;
println perimeter
```

Save as `exercises/ex1.mepl` and run with `./mepl exercises/ex1.mepl`.
Expected output:
```
21
20
```

## Solution 2

```
let clamp : Int -> Int -> Int -> Int =
  fn (val: Int) => fn (lo: Int) => fn (hi: Int) =>
    if val < lo then lo
    else if val > hi then hi
    else val;

println (clamp 5 1 10);
println (clamp (0 - 3) 0 100);
println (clamp 50 0 10)
```

Expected output:
```
5
0
10
```

## Solution 3

```
let sumList : [Int] -> Int = fix (fn (self: [Int] -> Int) =>
  fn (l: [Int]) => case l of
    | [] => 0
    | h :: t => h + (self t));

let filterPositive : [Int] -> [Int] = fix (fn (self: [Int] -> [Int]) =>
  fn (l: [Int]) => case l of
    | [] => [] : [Int]
    | h :: t => if h > 0 then h :: (self t) else self t);

println (sumList [1, -2, 3, -4, 5]);
println (filterPositive [1, -2, 3, -4, 5])
```

Expected output:
```
3
[1, 3, 5]
```

## Solution 4

```
let map : (Int -> Int) -> [Int] -> [Int] = fix (fn (self: (Int -> Int) -> [Int] -> [Int]) =>
  fn (f: Int -> Int) => fn (l: [Int]) => case l of
    | [] => [] : [Int]
    | h :: t => (f h) :: (self f t));

let double : Int -> Int = fn (x: Int) => x * 2;
println (map double [1, 2, 3, 4])
```

Expected output:
```
[2, 4, 6, 8]
```

## Solution 5

```
let safeDivide : Int -> Int -> (Int + String) =
  fn (a: Int) => fn (b: Int) =>
    if b == 0
    then inr "division by zero" as Int + String
    else inl (a / b) as Int + String;

let showResult : (Int + String) -> Unit =
  fn (r: Int + String) => case r of
    | inl n => println n
    | inr msg => println msg;

showResult (safeDivide 10 3);
showResult (safeDivide 10 0)
```

Expected output:
```
3
division by zero
```

## Solution 6

```
let p1 = {x = 1, y = 2};
let p2 = {x = 4, y = 6};

let distance : {x: Int, y: Int} -> {x: Int, y: Int} -> Int =
  fn (a: {x: Int, y: Int}) => fn (b: {x: Int, y: Int}) =>
    let dx = b.x - a.x in
    let dy = b.y - a.y in
    dx * dx + dy * dy;

let dist = distance p1 p2;
let result = (p1, p2, dist);
println result
```

Expected output:
```
({x = 1, y = 2}, {x = 4, y = 6}, 25)
```

## Solution 7

```
let s = "MEPL";
let len = length s;
for i = 0 to len do
  println (charAt s i)
end
```

Expected output:
```
M
E
P
L
```
