# Monkey Lango

Because we type like monkeys

### Functional

Purely functional language, no loops, no mutation, no side effects, everything is immutable.

### Simple

Simple language, no classes, no modules, no namespaces, everything is a function or a macro.
If you need something different you can build it with functions or macros.

### Grammar

```js
// Assignment
let name = value; // Optional semicolon

// Functions
let func = fn(arg1, arg2, ..., argN) {
	// Function body
}

"string" // Strings
123 // Integers
true // Booleans
null // Null

// Indexing
arr[index] // Indexing arrays, for now only integers
obj["key"] // Indexing objects, anything Hashable like strings, integers, booleans, etc
obj["missing"] // returns null

// Calling
func(arg1, arg2, ..., argN)

// If
if condition { // Everything is true but 0, false and null. Should I add empty strings, hashes, and arrays?
	// If body
}

"string" + "string" // String concatenation
"abcde" - "abc" // String substraction returns "de"
1 + 1 - (5 - 2) * 3 / 2 // Integer operations
```

### Builtins

Some builtins utilities:

```js
len(arr) // Returns the length of the array
first(arr) // Returns the first element of the array
last(arr) // Returns the last element of the array
rest(arr) // Returns the rest of the array
push(arr, value) // Pushes a value to the end of the array
string(value, value, ..., value) // Converts any value to a string
echo(value, value, ..., value) // Echos any value to the console
read(file) // Reads a file and returns its content
eval(file) // Evaluates a string as code and returns its content
```

### Macros

Here is the "unique" and truly powerful feature of Monkey Lango, macros.

```js
let func = macro(name, lp, arg, type, rp, lb, body, rb) {
    // Macro body
}
```

### Contributing

Contributions are welcome, just open a PR.
