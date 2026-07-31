package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
)

const cookieName = "minidash_session"

type Auth struct {
	secret       []byte
	plain        string
	hash         string
	rememberDays int
}

// New creates an Auth. secretPath holds the HMAC key (auto-generated on first
// run). Either plain or hash may be empty; both empty = disabled (open).
func New(secretPath, plain, hash string) (*Auth, error) {
	secret, err := loadOrCreateSecret(secretPath)
	if err != nil {
		return nil, err
	}
	return &Auth{secret: secret, plain: plain, hash: hash, rememberDays: 30}, nil
}

func loadOrCreateSecret(path string) ([]byte, error) {
	if path == "" {
		path = ".secret"
	}
	if b, err := os.ReadFile(path); err == nil {
		return b, nil
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return nil, err
		}
	}
	if err := os.WriteFile(path, b, 0o600); err != nil {
		return nil, err
	}
	return b, nil
}

func (a *Auth) Enabled() bool { return a.plain != "" || a.hash != "" }

func (a *Auth) Verify(pw string) bool {
	if a.plain != "" {
		return hmac.Equal([]byte(a.plain), []byte(pw))
	}
	if a.hash != "" {
		match, err := argon2id.ComparePasswordAndHash(pw, a.hash)
		return err == nil && match
	}
	return true
}

func (a *Auth) SetSession(w http.ResponseWriter, remember bool) {
	days := 1
	if remember {
		days = a.rememberDays
	}
	exp := time.Now().Add(time.Duration(days) * 24 * time.Hour).Unix()
	payload := strconv.FormatInt(exp, 10)
	http.SetCookie(w, &http.Cookie{
		Name: cookieName, Value: payload + "." + a.sign(payload), Path: "/",
		HttpOnly: true, SameSite: http.SameSiteLaxMode, Expires: time.Unix(exp, 0),
	})
}

func (a *Auth) Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: cookieName, Value: "", Path: "/", MaxAge: -1})
}

func (a *Auth) Valid(r *http.Request) bool {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return false
	}
	payload, sig, ok := strings.Cut(c.Value, ".")
	if !ok || !hmac.Equal([]byte(a.sign(payload)), []byte(sig)) {
		return false
	}
	exp, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		return false
	}
	return time.Now().Unix() < exp
}

func (a *Auth) sign(payload string) string {
	mac := hmac.New(sha256.New, a.secret)
	mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

// Require allows the request through when auth is disabled or the session is
// valid; otherwise redirects to /login.
func (a *Auth) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.Enabled() || a.Valid(r) {
			next.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})
}
