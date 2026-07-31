package config

import (
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watch hot-reloads the Store when its file changes externally. Returns a
// stop function. Debounces rapid events; ignores reload errors (keeps last
// good copy).
func Watch(st *Store) (func(), error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err := w.Add(st.path); err != nil {
		w.Close()
		return nil, err
	}
	done := make(chan struct{})
	go func() {
		var debounce *time.Timer
		for {
			select {
			case <-done:
				w.Close()
				return
			case ev, ok := <-w.Events:
				if !ok {
					return
				}
				if ev.Has(fsnotify.Write) || ev.Has(fsnotify.Create) {
					if debounce != nil {
						debounce.Stop()
					}
					debounce = time.AfterFunc(150*time.Millisecond, func() {
						if c, err := Load(st.path); err == nil {
							st.replace(c)
						}
					})
				}
			case _, ok := <-w.Errors:
				if !ok {
					return
				}
			}
		}
	}()
	return func() { close(done) }, nil
}
