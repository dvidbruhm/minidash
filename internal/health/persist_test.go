package health

import (
	"encoding/json"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestHistoryPersistRoundTrip(t *testing.T) {
	up := httptest.NewServer(okHandler(200))
	defer up.Close()
	path := filepath.Join(t.TempDir(), "status-history.json")
	ch := New(500*time.Millisecond, 24)
	ch.CheckAll([]string{up.URL})
	if err := ch.SaveHistory(path); err != nil {
		t.Fatalf("SaveHistory: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not written: %v", err)
	}
	ch2 := New(500*time.Millisecond, 24)
	if err := ch2.LoadHistory(path); err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	h := ch2.History()
	if len(h[up.URL]) != 1 || h[up.URL][0] != "up" {
		t.Fatalf("loaded = %v", h[up.URL])
	}
}

func TestHistoryLoadMissingAndCorrupt(t *testing.T) {
	ch := New(500*time.Millisecond, 24)
	// missing
	if err := ch.LoadHistory(filepath.Join(t.TempDir(), "nope.json")); err != nil {
		t.Fatalf("missing should not error: %v", err)
	}
	// corrupt
	path := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(path, []byte("{not json"), 0o644)
	if err := ch.LoadHistory(path); err != nil {
		t.Fatalf("corrupt should not error: %v", err)
	}
	if len(ch.History()) != 0 {
		t.Fatal("expected empty history after corrupt load")
	}
}

func TestHistorySavePrunesStale(t *testing.T) {
	a := httptest.NewServer(okHandler(200))
	defer a.Close()
	b := httptest.NewServer(okHandler(200))
	defer b.Close()
	path := filepath.Join(t.TempDir(), "h.json")
	ch := New(500*time.Millisecond, 24)
	ch.CheckAll([]string{a.URL, b.URL})
	ch.SaveHistory(path)

	ch2 := New(500*time.Millisecond, 24)
	ch2.LoadHistory(path)
	ch2.CheckAll([]string{a.URL}) // only a this cycle
	ch2.SaveHistory(path)

	data, _ := os.ReadFile(path)
	var m map[string][]string
	json.Unmarshal(data, &m)
	if _, ok := m[b.URL]; ok {
		t.Fatal("stale url b not pruned")
	}
	if _, ok := m[a.URL]; !ok {
		t.Fatal("a missing after prune")
	}
}

func TestCheckAllAutoPersists(t *testing.T) {
	up := httptest.NewServer(okHandler(200))
	defer up.Close()
	path := filepath.Join(t.TempDir(), "auto.json")
	ch := New(500*time.Millisecond, 24)
	ch.SetHistoryStore(path)
	ch.CheckAll([]string{up.URL})
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("auto-save did not write file: %v", err)
	}
}
