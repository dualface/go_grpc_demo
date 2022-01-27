package server

import (
	"sync"
	"sync/atomic"
)

type (
	RunSingleton struct {
		running uint32
		mutex   sync.Mutex
	}
)

func (rs *RunSingleton) Run(f func()) bool {
	if atomic.LoadUint32(&rs.running) != 0 {
		return false
	}

	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	if rs.running != 0 {
		return false
	}

	atomic.StoreUint32(&rs.running, 1)
	defer atomic.StoreUint32(&rs.running, 0)
	f()
	return true
}
