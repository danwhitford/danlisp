```
______            _     _           
|  _  \          | |   (_)          
| | | |__ _ _ __ | |    _ ___ _ __  
| | | / _  | '_ \| |   | / __| '_ \ 
| |/ / (_| | | | | |___| \__ \ |_) |
|___/ \__,_|_| |_\_____/_|___/ .__/ 
                             | |    
                             |_|    
```

## Running

To get a basic REPL

`go run cmd/danlisp/danlisp.go`

To run a file

`go run cmd/danlisp/danlisp.go <filename>`

## Examples

There are example programs [here](https://github.com/danwhitford/danlisp/tree/main/examples)

## Syntax

DanLisp uses paranthesized prefix notation, which will be familiar to anyone who has seen any other Lisp-like language.

### Literals

There are three types of literal so far; strings, numbers and nil.

eg
```
3
3.0
3.35
"foo"
"bar"
nil
```

### Operators

All the basic mathematical operators are present

```
(+ 2 2)
(- 12 7)
(* 5 2)
(/ 100 5)
```

Equality is a single `=`
```
(= 2 2)
```

Operators can be nested unambiguously

```
(+ 2 (- 10 5))
(+ (- 2 10) 5)
```

There are also logical operators for comparison, `lt` and `gt` for less than and greater than.

```
(lt 2 10)
(gt 1000 10)
```

There are boolean operators for combining conditions

```
(and (= 2 2) (= 4 (+ 2 2)))
(or (= 2 2) (= 5 (+ 2 2)))
```
### Built-ins

Printing to `stdout` is done with `prn`

```
(prn "Hello world")
(set name "Dan")
(prn "Hello" name)
```

If `prn` is given more than one argument it will print them all with a seperating space.

### Variables

Variables can be declared using `set` 

```
(set i 10)
(set name "Dan")
```

Variables are fully mutable. Example of incrememnting a variable.

```
(set count 1)
(set count (+ count 1))
```

It is possible to overwrite any builtin functions and keywords with user defined variables but don't do this.
### If

The `if` function will execute the first branch if the condition is true, and the second if it is false.  It has the form

```
(if cond
    true-branch
    false-branch)
```

For example

```
(if (gt 99999 n)
    (prn n "is a really big number")
    (prn n "is not that big"))
```

#### Conditionals

Any value that is not `nil` will evaluate to `true`. For conveniance the symbol `t` is provided and is set to true.

### Loops

There is one type of loop in DanLisp currently, the `while` loop of the form.

```
(while cond
    body)
```

For example to print "hello world" forever,

```
(while t
    (prn "hello world"))
```

You can simulate a `for` loop with a more verbose `while`,

```
(set i 1)
(while (lt i 10)
    (prn "Square of" i "is" (* i i))
    (set i (+ i 1)))
```

### Functions

Functions can be defined using `defn` in the form of 

```
(defn function-name (arglist)
    body)
```

For example a function that says hello

```
(defn say-hello (name)
    (prn "Hello" name))
```

or a function to add two numbers

```
(defn adder (a b)
    (+ a b))
```

The last expression in the function will implicitly be the return value. There is no way to return early.
