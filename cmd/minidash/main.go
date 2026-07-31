package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"minidash/internal/auth"
	"minidash/internal/config"
	"minidash/internal/health"
	"minidash/internal/icons"
	"minidash/internal/server"
)

type iconSvc struct{}

func (iconSvc) Search(q, prefix string, limit int) []server.IconResult {
	rs := icons.Search(q, prefix, limit)
	out := make([]server.IconResult, len(rs))
	for i, r := range rs {
		out[i] = server.IconResult{Prefix: r.Prefix, Name: r.Name, SVG: r.SVG}
	}
	return out
}
func (iconSvc) SVG(prefix, name string) (string, bool) { return icons.SVG(prefix, name) }
func (iconSvc) Collections() []string                  { return icons.Collections() }

func main() {
	cfgPath := envOr("MINIDASH_CONFIG", "config.yaml")
	addr := envOr("MINIDASH_ADDR", ":8080")
	pw := os.Getenv("MINIDASH_PASSWORD")

	store, err := config.Open(cfgPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	stopWatch, err := config.Watch(store)
	if err != nil {
		log.Printf("warn: fsnotify unavailable: %v", err)
	} else {
		defer stopWatch()
	}

	secretPath := filepath.Join(filepath.Dir(cfgPath), ".secret")
	a, err := auth.New(secretPath, pw, store.Snapshot().PasswordHash)
	if err != nil {
		log.Fatalf("auth: %v", err)
	}
	if !a.Enabled() {
		log.Printf("warn: no password set (MINIDASH_PASSWORD / password_hash) - Settings is open")
	}

	c := store.Snapshot()
	ch := health.New(time.Duration(c.Health.TimeoutSeconds)*time.Second, 24)
	ch.SetHistoryStore(filepath.Join(filepath.Dir(cfgPath), "status-history.json"))
	stop := make(chan struct{})
	defer close(stop)
	if c.Health.Enabled {
		ch.Start(time.Duration(c.Health.IntervalSeconds)*time.Second, func() []string {
			cfg := store.Snapshot()
			var urls []string
			for _, l := range cfg.Links {
				on := cfg.Health.Enabled
				if l.Health != nil {
					on = *l.Health
				}
				if on {
					urls = append(urls, l.URL)
				}
			}
			return urls
		}, stop)
	}

	srv, err := server.New(server.Deps{Store: store, Auth: a, Health: ch, Icons: iconSvc{}})
	if err != nil {
		log.Fatalf("server: %v", err)
	}

	hs := &http.Server{Addr: addr, Handler: srv.Handler()}
	go func() {
		log.Printf("minidash listening on %s", addr)
		if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = hs.Shutdown(ctx)
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
