package main

import (
	"context"
	"fmt"
	"log"

	"go.neonxp.dev/lexpr"
)

func main() {
	ctx := context.Background()
	l := lexpr.New(lexpr.WithDefaults())

	// Simple math
	result1 := <-l.Eval(ctx, `2 + 2 * 2`)
	log.Println("Result 1:", result1.Value)

	// Helper for one result
	result2, err := l.OneResult(ctx, `len("test") + 10`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Result 2:", result2)

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
	log.Println("Result 3:", result3)

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
	result41, err := l.OneResult(ctx, `jsonData.key1name`) // = value1
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Result 4-1:", result41)
	result42, err := l.OneResult(ctx, `jsonData.rootKey2.childKey2`) // = value3
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Result 4-2:", result42)
	result43, err := l.OneResult(ctx, `jsonData.arrayKey.3`) // = array value 4
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Result 4-3:", result43)

	// Logic expressions
	result51, err := l.OneResult(ctx, `jsonData.key1name == "value1"`) // = 1
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Result 5-1:", result51)
	result52, err := l.OneResult(ctx, `10 >= 5 || 10 <= 5`) // = 1
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Result 5-2:", result52)
	result53, err := l.OneResult(ctx, `10 >= 5 && 10 <= 5`) // = 0
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Result 5-3:", result53)
}
