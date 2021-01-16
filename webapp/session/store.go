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

// GetSession returns the session of the current user.
// Retrieves the session ID from the cookie in the request.
// If session ID is not found, a new session is created.
func GetSession(w http.ResponseWriter, r *http.Request) *HTTPSession {
	httpSessionStore.lock.Lock()
	defer httpSessionStore.lock.Unlock()
	var session *HTTPSession
	now := time.Now()
	if cookie, err := r.Cookie(cookieName); err == nil && cookie.Value != "" {
		if element, ok := httpSessionStore.sessions[cookie.Value]; ok {
			element.Value.(*HTTPSession).lastAccessedTime = now
			httpSessionStore.list.MoveToFront(element)
			session = element.Value.(*HTTPSession)
		} else {
			log.Printf("Session %s not found\n", cookie.Value)
			http.SetCookie(w, &http.Cookie{Name: cookieName, Path: "/", HttpOnly: true, Secure: true, Expires: time.Time{}, MaxAge: -1})
		}
	} else {
		log.Printf("Session cookie not found\n")
	}
	if session == nil {
		log.Printf("Creating new HTTP session.")
		sessionID := util.GenerateRandomBytesToBase64(32)
		session = &HTTPSession{id: sessionID, creationTime: now, lastAccessedTime: now, attributes: make(map[string]interface{}, 0)}
		element := httpSessionStore.list.PushBack(session)
		httpSessionStore.sessions[sessionID] = element
		sessionCookie := &http.Cookie{Name: cookieName, Value: sessionID, Path: "/", HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode}
		http.SetCookie(w, sessionCookie)
		r.AddCookie(sessionCookie)
	}
	return session
}

// InvalidateSession deletes the session of the current user.
func InvalidateSession(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(cookieName); err == nil && cookie.Value != "" {
		httpSessionStore.lock.Lock()
		defer httpSessionStore.lock.Unlock()
		if element, ok := httpSessionStore.sessions[cookie.Value]; ok {
			delete(httpSessionStore.sessions, cookie.Value)
			httpSessionStore.list.Remove(element)
		}
		http.SetCookie(w, &http.Cookie{Name: cookieName, Path: "/", HttpOnly: true, Expires: time.Now(), MaxAge: -1})
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
