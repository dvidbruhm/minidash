package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func okHandler(code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(code) })
}

func TestHistoryRingBuffer(t *testing.T) {
	up := httptest.NewServer(okHandler(200))
	defer up.Close()
	ch := New(500*time.Millisecond, 3) // cap 3
	for i := 0; i < 4; i++ { // 4 checks -> oldest evicted
		ch.CheckAll([]string{up.URL})
	}
	h := ch.History()
	if len(h[up.URL]) != 3 {
		t.Fatalf("len = %d, want 3", len(h[up.URL]))
	}
	for _, s := range h[up.URL] {
		if s != "up" {
			t.Fatalf("sample = %q, want up", s)
		}
	}
}

func TestHistoryRecordsDownAndUp(t *testing.T) {
	down := httptest.NewServer(okHandler(500))
	defer down.Close()
	up := httptest.NewServer(okHandler(200))
	defer up.Close()
	ch := New(500*time.Millisecond, 24)
	ch.CheckAll([]string{down.URL}) // down
	ch.CheckAll([]string{up.URL, down.URL})
	h := ch.History()
	if got := h[down.URL]; len(got) != 2 || got[0] != "down" || got[1] != "down" {
		t.Fatalf("down history = %v", got)
	}
	if got := h[up.URL]; len(got) != 1 || got[0] != "up" {
		t.Fatalf("up history = %v", got)
	}
}
