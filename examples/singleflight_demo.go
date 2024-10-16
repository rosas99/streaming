package main

import (
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
)

func main() {
	var sg singleflight.Group
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err, shared := sg.Do("key", func() (interface{}, error) {
				// 这里是你的耗时操作，比如数据库查询
				fmt.Println("get one")
				return "value", nil
			})
			if err != nil {
				panic(err)
			}
			fmt.Printf("v: %v, shared: %v\n", v, shared)
		}()
	}
	wg.Wait()
}
