package server

import (
	"path/filepath"
	"testing"
	"time"

	"minidash/internal/auth"
	"minidash/internal/config"
	"minidash/internal/health"
)

type stubIcons struct{}

func (stubIcons) Search(q, prefix string, limit int) []IconResult { return nil }
func (stubIcons) SVG(prefix, name string) (string, bool)          { return "", false }
func (stubIcons) Collections() []string                           { return nil }

func NewTestDeps(t *testing.T) Deps {
	t.Helper()
	st, err := config.Open(filepath.Join(t.TempDir(), "config.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	a, err := auth.New(filepath.Join(t.TempDir(), ".secret"), "", "")
	if err != nil {
		t.Fatal(err)
	}
	return Deps{Store: st, Auth: a, Health: health.New(time.Second, 500*time.Millisecond), Icons: stubIcons{}}
}
