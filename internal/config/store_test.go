package config

import (
	"path/filepath"
	"sync"
	"testing"
)

func TestStoreUpdateConcurrent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = st.Update(func(c *Config) error {
				c.Links = append(c.Links, Link{Name: "x", URL: "http://x"})
				return nil
			})
		}()
	}
	wg.Wait()
	if got := len(st.Snapshot().Links); got != 51 { // 50 + default seed link
		t.Fatalf("len(links) = %d", got)
	}
}
