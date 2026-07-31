package health

import (
	"net/http"
	"sync"
	"time"
)

type Checker struct {
	mu       sync.RWMutex
	status   map[string]string
	history  map[string][]string
	cap      int
	histPath string
	client   *http.Client
}

// New creates a Checker. timeout bounds each request; historyCap is the number
// of recent samples retained per URL (<=0 disables history).
func New(timeout time.Duration, historyCap int) *Checker {
	return &Checker{
		status:   map[string]string{},
		history:  map[string][]string{},
		cap:      historyCap,
		client:   &http.Client{Timeout: timeout},
	}
}

// CheckAll checks the given URLs concurrently and stores results.
func (c *Checker) CheckAll(urls []string) {
	var wg sync.WaitGroup
	for _, u := range urls {
		wg.Add(1)
		go func(u string) { defer wg.Done(); c.checkOne(u) }(u)
	}
	wg.Wait()
}

func (c *Checker) checkOne(u string) {
	status := "down"
	defer func() {
		if r := recover(); r != nil {
			status = "down"
		}
		c.mu.Lock()
		c.status[u] = status
		if c.cap > 0 {
			h := append(c.history[u], status)
			if len(h) > c.cap {
				h = h[len(h)-c.cap:]
			}
			c.history[u] = h
		}
		c.mu.Unlock()
	}()
	resp, err := c.client.Get(u)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 400 {
		status = "up"
	}
}

func (c *Checker) Snapshot() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string]string, len(c.status))
	for k, v := range c.status {
		out[k] = v
	}
	return out
}

// History returns a defensive copy of recent samples per URL.
func (c *Checker) History() map[string][]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string][]string, len(c.history))
	for k, v := range c.history {
		cp := make([]string, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}

// Start runs a background loop calling checkFn every interval until stop is
// closed. Performs an immediate first check.
func (c *Checker) Start(interval time.Duration, checkFn func() []string, stop <-chan struct{}) {
	run := func() { c.CheckAll(checkFn()) }
	run()
	t := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-stop:
				t.Stop()
				return
			case <-t.C:
				run()
			}
		}
	}()
}
