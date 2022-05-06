package wait

/*
copy :https://github.com/HDT3213/godis/blob/master/src/lib/sync/wait/wait.go
*/

import (
	"sync"
	"time"
)

type Wait struct {
	wg sync.WaitGroup
}

func (w *Wait) Add(delta int) {
	w.wg.Add(delta)
}

func (w *Wait) Done() {
	w.wg.Done()
}

func (w *Wait) Wait() {
	w.wg.Wait()
}

func (w *Wait) WaitWithTimeOut(timeout time.Duration) bool {
	// wait for conn close
	c := make(chan bool, 1)
	go func() {
		defer close(c)
		w.wg.Wait()
		c <- true
	}()

	select {
	case <-time.After(timeout):
		return true // time out
	case <-c:
		return false // normally
	}
}
