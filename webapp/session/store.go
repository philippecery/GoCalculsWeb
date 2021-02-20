package session

import (
	"container/list"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/philippecery/maths/webapp/util"
)

// HTTPSessionStore struct
type HTTPSessionStore struct {
	lock                sync.Mutex
	sessions            map[string]*list.Element
	list                *list.List
	maxInactiveInterval uint
}

var httpSessionStore = &HTTPSessionStore{
	list:                list.New(),
	sessions:            make(map[string]*list.Element, 0),
	maxInactiveInterval: defaultMaxInactiveInterval,
}

// NewSession creates a new session for the current user.
// If a session exists, it is invalidated.
func NewSession(w http.ResponseWriter, r *http.Request) *HTTPSession {
	httpSessionStore.lock.Lock()
	defer httpSessionStore.lock.Unlock()
	return newSession(w, r)
}

func newSession(w http.ResponseWriter, r *http.Request) *HTTPSession {
	invalidateSession(w, r)
	log.Printf("Creating new HTTP session.")
	now := time.Now()
	sessionID := util.GenerateRandomBytesToBase64(32)
	session := &HTTPSession{id: sessionID, creationTime: now, lastAccessedTime: now, attributes: make(map[string]interface{}, 0)}
	element := httpSessionStore.list.PushBack(session)
	httpSessionStore.sessions[sessionID] = element
	sessionCookie := &http.Cookie{Name: cookieName, Value: sessionID, Path: "/", HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode}
	http.SetCookie(w, sessionCookie)
	r.AddCookie(sessionCookie)
	return session
}

// GetSession returns the session of the current user.
// Returns nil if session does not exist.
func GetSession(w http.ResponseWriter, r *http.Request) *HTTPSession {
	httpSessionStore.lock.Lock()
	defer httpSessionStore.lock.Unlock()
	return getSession(w, r)
}

func getSession(w http.ResponseWriter, r *http.Request) *HTTPSession {
	var session *HTTPSession
	now := time.Now()
	if cookie, err := r.Cookie(cookieName); err == nil && cookie.Value != "" {
		if element, exists := httpSessionStore.sessions[cookie.Value]; exists {
			element.Value.(*HTTPSession).lastAccessedTime = now
			httpSessionStore.list.MoveToFront(element)
			session = element.Value.(*HTTPSession)
		} else {
			log.Printf("Session %s not found\n", cookie.Value)
			invalidateSession(w, r)
			return nil
		}
	} else {
		log.Printf("Session cookie not found\n")
	}
	if session == nil {
		session = newSession(w, r)
	}
	return session
}

// InvalidateSession deletes the session of the current user.
func InvalidateSession(w http.ResponseWriter, r *http.Request) {
	httpSessionStore.lock.Lock()
	defer httpSessionStore.lock.Unlock()
	invalidateSession(w, r)
}

func invalidateSession(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(cookieName); err == nil && cookie.Value != "" {
		if element, exists := httpSessionStore.sessions[cookie.Value]; exists {
			delete(httpSessionStore.sessions, cookie.Value)
			httpSessionStore.list.Remove(element)
		}
		sessionCookie := &http.Cookie{Name: cookieName, Path: "/", HttpOnly: true, Secure: true, Expires: time.Time{}, MaxAge: -1}
		http.SetCookie(w, sessionCookie)
		r.AddCookie(sessionCookie)
	} else {
		log.Printf("No session cookie found\n")
	}
}

// GetMaxInactiveInterval returns the maximum time interval, in seconds, that the server will keep a session open between client accesses.
func GetMaxInactiveInterval() uint {
	return httpSessionStore.maxInactiveInterval
}

// SetMaxInactiveInterval specifies the time, in seconds, between client requests before the server will invalidate a session.
// The provided interval must be greater than 0, otherwise the interval is not set.
func SetMaxInactiveInterval(interval uint) {
	if interval > 0 {
		httpSessionStore.maxInactiveInterval = interval
	}
}

func init() {
	invalidateExpiredSessions()
}

func invalidateExpiredSessions() {
	httpSessionStore.lock.Lock()
	defer httpSessionStore.lock.Unlock()
	interval := int64(httpSessionStore.maxInactiveInterval)
	now := time.Now().Unix()
	for {
		if element := httpSessionStore.list.Back(); element != nil {
			if session := element.Value.(*HTTPSession); (session.lastAccessedTime.Unix() + interval) < now {
				httpSessionStore.list.Remove(element)
				delete(httpSessionStore.sessions, session.id)
				continue
			}
		}
		break
	}
	time.AfterFunc(time.Duration(interval)*time.Second, func() { invalidateExpiredSessions() })
}
