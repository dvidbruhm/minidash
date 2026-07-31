package config

import "sync"

// Store holds the current config in memory behind a mutex; all disk writes
// serialize through Update.
type Store struct {
	mu   sync.RWMutex
	path string
	cur  *Config
}

func Open(path string) (*Store, error) {
	c, err := Load(path)
	if err != nil {
		return nil, err
	}
	return &Store{path: path, cur: c}, nil
}

// Snapshot returns a defensive copy of the current config.
func (s *Store) Snapshot() Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return *s.cur
}

// Update mutates the config under the lock and atomically writes it. If fn
// returns an error or the write fails, the in-memory copy is reverted.
func (s *Store) Update(fn func(c *Config) error) error {
	s.mu.Lock()
	prev := *s.cur
	if err := fn(s.cur); err != nil {
		s.cur = &prev
		s.mu.Unlock()
		return err
	}
	if err := Save(s.path, s.cur); err != nil {
		s.cur = &prev
		s.mu.Unlock()
		return err
	}
	s.mu.Unlock()
	return nil
}

// replace swaps the in-memory config (used by the fsnotify watcher).
func (s *Store) replace(c *Config) {
	s.mu.Lock()
	s.cur = c
	s.mu.Unlock()
}

func (s *Store) Path() string { return s.path }
