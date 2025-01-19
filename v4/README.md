# s-expression parser, version 4

Re-implementation of  Dmitry Soshnikov's javascript s-expression parser in Go: https://gist.github.com/DmitrySoshnikov/2a434dda67019a4a7c37

To simplify reimplementation in Go, the parser signature is `parse(string) any`, returning atoms as `string` and lists as `[]any`

## Grammer

```
sexp -> atom | list

list -> '(' entries ')'

entries -> sexp entries | nil

atom -> string
```

## Examples

```
parse(``) -> nil
parse(`5`) -> 5
parse(`()`) -> []
parse(`(+ 5 10)`) ->  [+ 5 10]
parse(`(define (fact n) (if (= n 0) 1 (* n (fact (- n 1)))))`) -> [define [fact n][if [= n 0] 1 [* n [fact [- n 1]]]]]
```

