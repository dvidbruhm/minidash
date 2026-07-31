package auth

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/alexedwards/argon2id"
)

func TestPasswordPlain(t *testing.T) {
	a, _ := New(filepath.Join(t.TempDir(), ".secret"), "plaintext-pw", "")
	if !a.Verify("plaintext-pw") {
		t.Fatal("plain verify failed")
	}
	if a.Verify("wrong") {
		t.Fatal("wrong accepted")
	}
}

func TestPasswordHash(t *testing.T) {
	hash, _ := argon2id.CreateHash("hunter2", argon2id.DefaultParams)
	a, _ := New(filepath.Join(t.TempDir(), ".secret"), "", hash)
	if !a.Verify("hunter2") {
		t.Fatal("hash verify failed")
	}
	if a.Verify("nope") {
		t.Fatal("wrong hash accepted")
	}
}

func TestDisabledWhenNoPassword(t *testing.T) {
	a, _ := New(filepath.Join(t.TempDir(), ".secret"), "", "")
	if a.Enabled() {
		t.Fatal("should be disabled")
	}
}

func TestCookieRoundTrip(t *testing.T) {
	a, _ := New(filepath.Join(t.TempDir(), ".secret"), "pw", "")
	rec := httptest.NewRecorder()
	a.SetSession(rec, true)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(rec.Result().Cookies()[0])
	if !a.Valid(req) {
		t.Fatal("cookie not valid")
	}
	a.Clear(rec)
}
