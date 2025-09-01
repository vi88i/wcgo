package service

import "sync"

type ConcurrentCounter struct {
	Store map[string]uint64
	Mutex sync.Mutex
}

func NewConcurrentCounter() *ConcurrentCounter {
	return &ConcurrentCounter{
		Store: make(map[string]uint64),
		Mutex: sync.Mutex{},
	}
}

func (c *ConcurrentCounter) Increment(key string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Store[key]++
}
