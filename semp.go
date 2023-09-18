package semp

import (
	"errors"
	"sync"
)

type NotifyFunc func()

var (
	ErrCannotAcquire = errors.New("cannot acquire semaphore")
)

type Semp struct {
	weight uint32
	count  uint32
	ch     chan bool
	mu     *sync.Mutex
	notify NotifyFunc
}

func New(w uint, n NotifyFunc) *Semp {
	return &Semp{
		weight: uint32(w),
		count:  uint32(0),
		ch:     make(chan bool, 1),
		mu:     &sync.Mutex{},
		notify: n,
	}
}

// How to tackle this one
// should we have slots for this ?
func (s *Semp) Acquire(i int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.count+uint32(i) > s.weight {
		return ErrCannotAcquire
	}

	s.count += uint32(i)
	return nil
}

func (s *Semp) Release(i int) error {

	s.mu.Lock()
	defer s.mu.Unlock()
	s.count--
	go s.notify()
	return nil
}
