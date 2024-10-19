package main

import (
	"context"
	"fmt"
	r9 "github.com/redis/go-redis/v9"
	"github.com/rosas99/streaming/pkg/streams/connector/extension"
	"github.com/rosas99/streaming/pkg/streams/connector/redis"
	"github.com/rosas99/streaming/pkg/streams/flow"
	"strconv"
	"time"
)

func main() {
	_ = Run()
}

// addUTC appends a UTC timestamp to the beginning of the message value.
var addUTC = func(msg *message) *message {
	timestamp := time.Now().Format(time.DateTime)

	// Concatenate the UTC timestamp with msg.Value
	msg.msg = timestamp + " " + (msg.msg)
	return msg
}

var addUTC2 = func(msg *r9.Message) *r9.Message {
	timestamp := time.Now().Format(time.DateTime)

	// Concatenate the UTC timestamp with msg.Value
	msg.Payload = timestamp + " " + (msg.Payload)
	return msg
}

func Run() error {
	source := extension.NewChanSource(tickerChan(time.Second * 1))

	filter := flow.NewMap(addUTC, 1)

	sink := extension.NewStdoutSink()

	source.Via(filter).To(sink)
	return nil
}

type message struct {
	msg string
}

func Run2() error {
	options := &r9.Options{
		Addr:         "127.0.0.1:6379",
		Username:     "",
		Password:     "onex(#)666",
		DB:           0,
		MaxRetries:   3,
		MinIdleConns: 0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	}

	rdb := r9.NewClient(options)

	// check redis if is ok
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		fmt.Print("eer")
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	source, err := redis.NewRedisSource(ctx, options, "aa")

	time.AfterFunc(time.Second*10, cancelFunc)

	if err != nil {
		fmt.Print("ccc")
	}
	filter := flow.NewMap(addUTC2, 1)

	sink := redis.NewRedisSink(options, "bb")

	source.Via(filter).To(sink)
	return nil
}

func tickerChan(repeat time.Duration) chan any {
	//ticker := time.NewTicker(repeat)
	//oc := ticker.C
	nc := make(chan any)
	go func() {
		//for range oc {
		//	nc <- &message{strconv.FormatInt(time.Now().UTC().UnixNano(), 10)}
		//}

		// 发送数据到数据源
		for i := 1; i <= 5; i++ {
			nc <- &message{strconv.FormatInt(time.Now().UTC().UnixNano(), 10)}
			//time.Sleep(time.Second)
		}
		// 关闭数据源的输出通道，表示数据发送完毕
		close(nc)
	}()

	return nc
}
