package app

import (
	"sync"
)

// Pool is a []byte processing worker pool
type Pool struct {
	name string
	mu   sync.RWMutex
	cnt  int
	fn   func([]byte)
	in   chan []byte
}

// NewPool returns a worker pool
func NewPool(name string, fn func([]byte), workers int) *Pool {
	p := &Pool{
		name: name,
		fn:   fn,
		in:   make(chan []byte),
	}
	for i := 0; i < workers; i++ {
		p.Start()
	}
	return p
}

// Start starts a worker, which will run until either
// it receives a nil message or the channel is closed
func (p *Pool) Start() {
	p.mu.Lock()
	p.cnt++
	p.mu.Unlock()
	go func() {
		for b := range p.in {
			if b == nil {
				break
			}
			p.fn(b)
		}
		p.mu.Lock()
		p.cnt--
		p.mu.Unlock()
	}()
}

// Process give the bytes to be processed by the worker pool
func (p *Pool) Process(in []byte) {
	go func() {
		p.in <- in
	}()
}

// Drop removes a worker from the pool
func (p *Pool) Drop() {
	go func() {
		p.in <- nil
	}()
}

// Shutdown cancels all workers
func (p *Pool) Shutdown() {
	close(p.in)
}

// Count returns the number of workers in the pool
func (p *Pool) Count() int {
	p.mu.RLock()
	cnt := p.cnt
	p.mu.RUnlock()
	return cnt
}
