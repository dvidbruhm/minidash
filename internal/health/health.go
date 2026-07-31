package health

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
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
		status:  map[string]string{},
		history: map[string][]string{},
		cap:     historyCap,
		client:  &http.Client{Timeout: timeout},
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
	if c.histPath != "" {
		_ = c.SaveHistory(c.histPath)
	}
}

// SetHistoryStore sets the persistence path and loads any existing history.
func (c *Checker) SetHistoryStore(path string) {
	c.histPath = path
	_ = c.LoadHistory(path)
}

// LoadHistory reads history from path. Missing or corrupt files are ignored
// (no error returned); the checker keeps whatever it already had.
func (c *Checker) LoadHistory(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var m map[string][]string
	if err := json.Unmarshal(data, &m); err != nil {
		return nil
	}
	c.mu.Lock()
	c.history = m
	if c.history == nil {
		c.history = map[string][]string{}
	}
	c.mu.Unlock()
	return nil
}

// SaveHistory writes history atomically, keeping only URLs currently in the
// status map (stale URLs are pruned). Best-effort; errors propagated.
func (c *Checker) SaveHistory(path string) error {
	c.mu.RLock()
	out := make(map[string][]string, len(c.status))
	for u := range c.status {
		if v, ok := c.history[u]; ok {
			cp := make([]string, len(v))
			copy(cp, v)
			out[u] = cp
		}
	}
	c.mu.RUnlock()
	data, err := json.Marshal(out)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".status-history.*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
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
