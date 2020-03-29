package strToolkit

import "sync"

type AutoIncrementID struct {
	mutex   sync.Mutex
	counter int64
}

func (a *AutoIncrementID) Generate() int64 {
	a.mutex.Lock()
	a.counter++
	i := a.counter
	a.mutex.Unlock()
	return i
}
