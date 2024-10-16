package main

import (
	"context"
	"fmt"
	"github.com/rosas99/streaming/pkg/id"
)

func main() {
	o := func(*id.SonyflakeOptions) {
		id.WithSonyflakeMachineId(1) // 自定义机器ID，默认为自动检测
	}

	snowIns := id.NewSonyflake(o)
	id := snowIns.Id(context.Background())
	fmt.Print("id is :", id)
}
