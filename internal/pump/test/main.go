package main

import (
	"github.com/rosas99/streaming/pkg/streams/connector/extension"
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

func tickerChan(repeat time.Duration) chan any {
	//ticker := time.NewTicker(repeat)
	//oc := ticker.C
	nc := make(chan any)
	go func() {
		//for range oc {
		//nc <- &message{strconv.FormatInt(time.Now().UTC().UnixNano(), 10)}
		//}

		for i := 0; i < 300; i++ {
			nc <- &message{strconv.FormatInt(time.Now().UTC().UnixNano(), 10)}
		}
		close(nc)

	}()
	return nc
}
