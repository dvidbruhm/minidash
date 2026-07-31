package health

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestChecker(t *testing.T) {
	up := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }))
	defer up.Close()
	down := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(500) }))
	defer down.Close()
	ch := New(time.Second, 500*time.Millisecond)
	ch.CheckAll([]string{up.URL, down.URL, "http://127.0.0.1:0/nope"})
	st := ch.Snapshot()
	if st[up.URL] != "up" {
		t.Fatalf("up got %q", st[up.URL])
	}
	if st[down.URL] != "down" {
		t.Fatalf("down got %q", st[down.URL])
	}
	if st["http://127.0.0.1:0/nope"] != "down" {
		t.Fatalf("err got %q", st["http://127.0.0.1:0/nope"])
	}
}
