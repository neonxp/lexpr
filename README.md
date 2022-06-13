
# Lexpr - universal expression evaluator

This library can evaluate any types of expressions: math expression, logic expression, simple DSLs.

## Installation

`go get go.neonxp.dev/lexpr`

## Usage

Full example: [/example/main.go](/example/main.go)

```go
ctx := context.Background()
l := lexpr.New(lexpr.WithDefaults())

// Simple math
result1 := <-l.Eval(ctx, `2 + 2 * 2`) // Output channel can return many results
log.Println("Result 1:", result1.Value) // Output: 6

// Helper for exact one result
result2, err := l.OneResult(ctx, `len("test") + 10`)
if err != nil {
 log.Fatal(err)
}
log.Println("Result 2:", result2) // Output: 14

// Custom functions
l.SetFunction("add", func(ts *lexpr.TokenStack) error {
 a, okA := ts.Pop().Number() // first func argument
 b, okB := ts.Pop().Number() // second func argument
 if !okA || !okB {
  return fmt.Errorf("Both args must be number")
 }
 ts.Push(lexpr.TokenFromInt(a + b))
 return nil
})
result3, err := l.OneResult(ctx, `add(12, 24) * 2`)
if err != nil {
 log.Fatal(err)
}
log.Println("Result 3:", result3) // Output: 72

// JSON extraction via dots and variables
jsonString := `{
 "rootKey1": "value1",
 "rootKey2": {
  "childKey1": "value2",
  "childKey2": "value3"
 },
 "arrayKey": [
  "array value 1",
  "array value 2",
  "array value 3",
  "array value 4"
 ]
}`
key1name := "rootKey1"
l.SetVariable("jsonData", jsonString)
l.SetVariable("key1name", key1name)
result41, err := l.OneResult(ctx, `jsonData.key1name`)
if err != nil {
 log.Fatal(err)
}
log.Println("Result 4-1:", result41) // Output: "value1"
result42, err := l.OneResult(ctx, `jsonData.rootKey2.childKey2`)
if err != nil {
 log.Fatal(err)
}
log.Println("Result 4-2:", result42) // Output: "value3"
result43, err := l.OneResult(ctx, `jsonData.arrayKey.3`)
if err != nil {
 log.Fatal(err)
}
log.Println("Result 4-3:", result43) // Output: "array value 4"

// Logic expressions
result51, err := l.OneResult(ctx, `jsonData.key1name == "value1"`)
if err != nil {
 log.Fatal(err)
}
log.Println("Result 5-1:", result51) // Output: 1
result52, err := l.OneResult(ctx, `10 >= 5 || 10 <= 5`)
if err != nil {
 log.Fatal(err)
}
log.Println("Result 5-2:", result52) // Output: 0
result53, err := l.OneResult(ctx, `10 >= 5 && 10 <= 5`)
if err != nil {
 log.Fatal(err)
}
log.Println("Result 5-3:", result53) // Output: 0
```

## Default operators

|Operator|Description|Example|
|:------:|:---------:|:-----:|
||JSON operators||
|`.`|Extract field from json|`jsonData.key1.0.key2`|
||Math operators||
|`**`|Power number|`3 ** 3` = 27|
|`*`|Multiple numbers|`2 * 4` = 8|
|`/`|Divide number|`6 / 3` = 2|
|`%`|Rem of division|`5 % 3` = 2|
|`+`|Sum|`2 + 2` = 4|
|`-`|Substract|`6 - 2` = 4|
||Logic operators||
|`!`|Logic not|`!1` = 0|
|`>`|More|`3 > 2` = 1|
|`>=`|More or equal|`3 >= 3` = 1|
|`<`|Less|`3 < 2` = 0|
|`<=`|Less or equal|`3 <= 3` = 1|
|`==`|Equal|`1==1` = 1|
|`!=`|Not equal|`1!=1` = 0|
|`&&`|Logic and|`3 > 0 && 1 > 0` = 1|
|`||`|Logic or|`1 > 0 || 1 == 1` = 1|

## Default functions

|Function|Description|Example|
|:------:|:---------:|:-----:|
|max|returns max of two values|`max(1,2)` = 2|
|min|returns min of two values|`max(1,2)` = 1|
|len|returns length of string|`len("test")` = 4|
|atoi|converts string to number|`atoi("123")` = 123|
|itoa|converts number to string|`itoa(123)` = "123"|

## Contribution

PRs are welcome.

## Author

Alexander Kiryukhin <i@neonxp.dev>

## License

![GPL v3](https://www.gnu.org/graphics/gplv3-with-text-136x68.png)
