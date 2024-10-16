package flow

import (
	"context"
	"github.com/rosas99/streaming/pkg/log"
	"github.com/rosas99/streaming/pkg/streams"
	"github.com/rosas99/streaming/pkg/streams/util"
	"golang.org/x/sync/semaphore"
)

// FlatMapFunction represents a FlatMap transformation function.
type FlatMapFunction[T, R any] func(T) []R

// FlatMap takes one element and produces zero, one, or more elements.
//
// in  -- 1 -- 2 ---- 3 -- 4 ------ 5 --
//
// [ -------- FlatMapFunction -------- ]
//
// out -- 1' - 2' -------- 4'- 4" - 5' -.
type FlatMap[T, R any] struct {
	flatMapFunction FlatMapFunction[T, R]
	in              chan any
	out             chan any
	parallelism     uint
}

// Verify FlatMap satisfies the Flow interface.
var _ streams.Flow = (*FlatMap[any, any])(nil)

// NewFlatMap returns a new FlatMap instance.
//
// flatMapFunction is the FlatMap transformation function.
// parallelism is the flow parallelism factor. In case the events order matters, use parallelism = 1.
func NewFlatMap[T, R any](flatMapFunction FlatMapFunction[T, R], parallelism uint) *FlatMap[T, R] {
	flatMap := &FlatMap[T, R]{
		flatMapFunction: flatMapFunction,
		in:              make(chan any),
		out:             make(chan any),
		parallelism:     parallelism,
	}
	go flatMap.doStream()

	return flatMap
}

// Via streams data through the given flow.
func (fm *FlatMap[T, R]) Via(flow streams.Flow) streams.Flow {
	go fm.transmit(flow)
	return flow
}

// To streams data to the given sink.
func (fm *FlatMap[T, R]) To(sink streams.Sink) {
	fm.transmit(sink)
}

// Out returns an output channel for sending data.
func (fm *FlatMap[T, R]) Out() <-chan any {
	return fm.out
}

// In returns an input channel for receiving data.
func (fm *FlatMap[T, R]) In() chan<- any {
	return fm.in
}

func (fm *FlatMap[T, R]) transmit(inlet streams.Inlet) {
	for element := range fm.Out() {
		inlet.In() <- element
	}
	close(inlet.In())
}

// doStream semaphore implement by buffer channel
func (fm *FlatMap[T, R]) doStream() {
	sem := make(chan struct{}, fm.parallelism)
	for elem := range fm.in {
		sem <- struct{}{}
		go func(element T) {
			defer func() { <-sem }()
			result := fm.flatMapFunction(element)
			for _, item := range result {
				fm.out <- item
			}
		}(elem.(T))
	}
	for i := 0; i < int(fm.parallelism); i++ {
		sem <- struct{}{}
	}
	close(fm.out)
}

// doStream2 semaphore implement by buffer channel with struct
func (fm *FlatMap[T, R]) doStream2() {
	sem := util.NewSemaphore(fm.parallelism)
	for elem := range fm.in {
		sem.Acquire()
		go func(element T) {
			defer func() { sem.Release() }()
			result := fm.flatMapFunction(element)
			for _, item := range result {
				fm.out <- item
			}
		}(elem.(T))
	}
	for i := 0; i < int(fm.parallelism); i++ {
		sem.Acquire()
	}
	close(fm.out)
}

// doStream3 semaphore implement by golang.org/x/sync/semaphore
func (fm *FlatMap[T, R]) doStream3() {
	sem := semaphore.NewWeighted(int64(fm.parallelism))

	for elem := range fm.in {
		if err := sem.Acquire(context.Background(), 1); err != nil {
			log.Errorf("...")
		}
		go func(element T) {
			defer sem.Release(1)
			result := fm.flatMapFunction(element)
			for _, item := range result {
				fm.out <- item
			}
		}(elem.(T))
	}

	for i := 0; i < int(fm.parallelism); i++ {
		if err := sem.Acquire(context.Background(), 1); err != nil {
			log.Errorf("...")
		}
	}
	close(fm.out)
}
