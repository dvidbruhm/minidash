package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatchReloadsOnExternalEdit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	st, _ := Open(path)
	stop, err := Watch(st)
	if err != nil {
		t.Fatal(err)
	}
	defer stop()
	if err := os.WriteFile(path, []byte("title: FromDisk\ndefault_view: compact\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if st.Snapshot().Title == "FromDisk" {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("config not reloaded; title=%q", st.Snapshot().Title)
}
