package session

import (
	"time"
)

// HTTPSession struct
type HTTPSession struct {
	id               string
	creationTime     *time.Time
	lastAccessedTime *time.Time
	attributes       map[string]interface{}
}

// GetAttribute returns the object bound with the specified name in this session, or nil if no object is bound under the name.
func (s *HTTPSession) GetAttribute(key string) interface{} {
	if v, ok := s.attributes[key]; ok {
		return v
	}
	return nil
}

// GetStringAttribute returns the object bound with the specified name in this session, or nil if no object is bound under the name.
func (s *HTTPSession) GetStringAttribute(key string) string {
	if v, isString := s.GetAttribute(key).(string); isString {
		return v
	}
	return ""
}

// GetAttributeNames returns a slice containing the names of all the objects bound to this session.
func (s *HTTPSession) GetAttributeNames() []string {
	names := make([]string, 0, len(s.attributes))
	for name := range s.attributes {
		names = append(names, name)
	}
	return names
}

// GetID returns a string containing the unique identifier assigned to this session.
func (s *HTTPSession) GetID() string {
	return s.id
}

// GetLastAccessedTime returns the last time the server accessed this session.
func (s *HTTPSession) GetLastAccessedTime() *time.Time {
	return s.lastAccessedTime
}

// SetAttribute binds an object to this session, using the name specified.
func (s *HTTPSession) SetAttribute(key string, value interface{}) {
	s.attributes[key] = value
}

// RemoveAttribute removes the object bound with the specified name from this session.
func (s *HTTPSession) RemoveAttribute(key string) {
	delete(s.attributes, key)
}
