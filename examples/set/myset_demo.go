package main

import (
	"fmt"
	"github.com/rosas99/streaming/pkg/util/set"
)

func main() {

	mySet := set.NewSet(1, 2, 3)

	mySet.Add("apple")

	fmt.Println("Contains 'apple':", mySet.Contains("apple")) // true

	mySet.Remove("apple")
	fmt.Println("Contains 'apple' after removal:", mySet.Contains("apple")) // false

	fmt.Println("Set:", mySet.String())

	fmt.Println("Size of Set:", mySet.Size())

	mySet.Clear()
	fmt.Println("Set after clear:", mySet.String())
}
