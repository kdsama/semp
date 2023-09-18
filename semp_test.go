package semp

import (
	"fmt"
	"sync"
	"testing"
)

func TestAcquire(t *testing.T) {

	type obj struct {
		name      string
		totalReqs int
		semInput  int
		want      int
		f         func()
	}
	tests := []obj{
		{
			name:      "Try acquiring more locks than the Maximum weight ",
			totalReqs: 8,
			semInput:  5,
			want:      3,
			f:         func() {},
		},
		{
			name:      "Try acquiring less locks than the Maximum weight ",
			totalReqs: 3,
			semInput:  5,
			want:      0,
			f:         func() {},
		},
	}
	for _, o := range tests {
		t.Run(o.name, func(t *testing.T) {
			sem := New(uint(o.semInput), o.f)

			errCount := []error{}
			wg := sync.WaitGroup{}
			wg.Add(o.totalReqs)
			mu := sync.Mutex{}
			for i := 0; i < o.totalReqs; i++ {
				go func() {
					defer wg.Done()
					e := sem.Acquire(1)
					if e == nil {
						fmt.Println("Acquired")
					} else {
						mu.Lock()
						errCount = append(errCount, e)
						mu.Unlock()
					}

				}()
			}
			wg.Wait()
			want := o.want
			got := len(errCount)
			if want != got {
				t.Errorf("wanted %v but got %v", want, got)
			}
		})
	}

}
