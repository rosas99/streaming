package util

import "sync"

type Semaphore struct {
	sync.Locker
	ch chan struct{}
}

func NewSemaphore(capacity uint) *Semaphore {
	return &Semaphore{ch: make(chan struct{}, capacity)}
}

func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

func (s *Semaphore) AcquireMutex() {
	s.Lock()
	defer s.Unlock()
	<-s.ch
}
func (s *Semaphore) Release() {
	<-s.ch
}
