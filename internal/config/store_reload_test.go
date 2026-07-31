package config

import (
	"path/filepath"
	"testing"
)

func TestStoreReload(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.yaml")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	// externally write a new config to the path
	c2 := Default()
	c2.Title = "Reloaded"
	if err := Save(path, &c2); err != nil {
		t.Fatal(err)
	}
	if err := st.Reload(); err != nil {
		t.Fatalf("Reload: %v", err)
	}
	if got := st.Snapshot().Title; got != "Reloaded" {
		t.Fatalf("after reload title = %q, want Reloaded", got)
	}
}
