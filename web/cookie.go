package web

import (
	"github.com/gorilla/sessions"
	"net/http"
	"sync"
)

var (
	cookieKeyPairs                            = "update-me-please!"
	cookieStore         *sessions.CookieStore = nil
	cookieStoreOnceFlag *sync.Once            = &sync.Once{}
)

func SetCookieKeyPairs(keyPairs string) {
	cookieKeyPairs = keyPairs
	cookieStoreOnceFlag = &sync.Once{}
}

func GetCookieKeyPairs() (keyPairs string) {
	keyPairs = cookieKeyPairs
	return
}

func GetCookieStore() *sessions.CookieStore {
	cookieStoreOnceFlag.Do(func() {
		cookieStore = sessions.NewCookieStore([]byte(cookieKeyPairs))
	})
	return cookieStore
}

// Get returns a session for the given name after adding it to the registry.
//
// It returns a new session if the sessions doesn't exist. Access IsNew on
// the session to check if it is an existing session or a new one.
//
// It returns a new session and an error if the session exists but could
// not be decoded.
func GetCookie(r *http.Request, name string) (*sessions.Session, error) {
	return GetCookieStore().Get(r, name)
}

// New returns a session for the given name without adding it to the registry.
//
// The difference between New() and Get() is that calling New() twice will
// decode the session data twice, while Get() registers and reuses the same
// decoded session after the first call.
func NewCookie(r *http.Request, name string) (*sessions.Session, error) {
	return GetCookieStore().New(r, name)
}

// Save adds a single session to the response.
func SaveCookie(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	return GetCookieStore().Save(r, w, s)
}
