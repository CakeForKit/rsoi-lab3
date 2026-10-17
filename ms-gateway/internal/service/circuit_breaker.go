package service

import (
	"sync"
	"time"
)

const (
	circuitFailureThreshold = 3
	circuitFailureWindow    = 3 * time.Second
	circuitRetryTimeout     = 3 * time.Second
)

type circuitBreaker struct {
	mu          sync.Mutex
	failures    [circuitFailureThreshold]time.Time
	next        int
	count       int
	openedAt    time.Time
	halfOpenRun bool
}

func (c *circuitBreaker) allow() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.openedAt.IsZero() {
		return true
	}
	if time.Since(c.openedAt) < circuitRetryTimeout || c.halfOpenRun {
		return false
	}

	c.halfOpenRun = true
	return true
}

func (c *circuitBreaker) success() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.halfOpenRun {
		return
	}

	c.failures = [circuitFailureThreshold]time.Time{}
	c.next = 0
	c.count = 0
	c.openedAt = time.Time{}
	c.halfOpenRun = false
}

func (c *circuitBreaker) failure() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.halfOpenRun {
		c.openedAt = time.Now()
		c.halfOpenRun = false
		return
	}

	fallenAt := time.Now()
	c.failures[c.next] = fallenAt
	if c.count < circuitFailureThreshold {
		c.count++
	}
	c.next = (c.next + 1) % circuitFailureThreshold

	if c.count < circuitFailureThreshold {
		return
	}

	head := c.failures[c.next]
	tail := c.failures[(c.next+circuitFailureThreshold-1)%circuitFailureThreshold]
	if tail.Sub(head) <= circuitFailureWindow {
		c.openedAt = fallenAt
	}
}
