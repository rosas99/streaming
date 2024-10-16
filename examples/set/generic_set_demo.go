package main

import (
	"fmt"
	set "github.com/rosas99/streaming/pkg/util/set/generics"
)

func main() {
	m := set.MakeSet[int]()
	m.Add(1)
	if m.Contains(2) {
		fmt.Println("contains 2")
	} else {
		fmt.Println("not contains 2")
	}
}
