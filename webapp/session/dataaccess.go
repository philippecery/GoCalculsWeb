package session

import (
	"strings"

	"github.com/philippecery/maths/webapp/database/document"
	"github.com/philippecery/maths/webapp/util"
)

// SetAuthenticatedUser stores the non-sensitive data of the authenticated user in this session.
func (s *HTTPSession) SetAuthenticatedUser(user *document.User) {
	s.SetAttribute("user", &UserInformation{
		UserID:    user.UserID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
	})
}

// GetAuthenticatedUser returns the UserInformation struct retrieved from the user session.
// Returns nil if there is no session created for the current user or if it does not contain a UserInformation struct.
func (s *HTTPSession) GetAuthenticatedUser() *UserInformation {
	if userInfo, isUserInformation := s.GetAttribute("user").(*UserInformation); isUserInformation {
		return userInfo
	}
	return nil
}

// SetErrorMessageID stores an error message ID in the user session.
func (s *HTTPSession) SetErrorMessageID(messageID string) {
	s.SetAttribute("errorMessageID", messageID)
}

// GetErrorMessageID returns the error message ID from the session, if any.
func (s *HTTPSession) GetErrorMessageID() string {
	if messageID, isString := s.GetAttribute("errorMessageID").(string); isString {
		s.RemoveAttribute("errorMessageID")
		return messageID
	}
	return ""
}

// SetCSRFToken generates a random token and stores it in this session.
func (s *HTTPSession) SetCSRFToken() string {
	token := util.GenerateRandomBytesToBase64(32)
	s.SetAttribute("csrf", token)
	return token
}

// GetCSRFToken returns the CSRF token from the session.
func (s *HTTPSession) GetCSRFToken() string {
	if token, isString := s.GetAttribute("csrf").(string); isString {
		return token
	}
	return ""
}

// NewCSWHToken generates a random token and stores it in this session.
func (s *HTTPSession) NewCSWHToken() string {
	token := util.GenerateRandomBytesToBase64(32)
	s.SetAttribute("cswh", token)
	return token
}

// GetCSWHToken returns the CSRF token from the session.
func (s *HTTPSession) GetCSWHToken() string {
	if token, isString := s.GetAttribute("cswh").(string); isString {
		return token
	}
	return ""
}

// SetCSPNonce generates a random nonce for strict CSP and stores it in this session.
func (s *HTTPSession) SetCSPNonce() string {
	nonce := util.GenerateRandomBytesToBase64(32)
	s.SetAttribute("nonce", nonce)
	return nonce
}

// GetCSPNonce returns the nonce for strict CSP from the session.
func (s *HTTPSession) GetCSPNonce() string {
	if nonce, isString := s.GetAttribute("nonce").(string); isString {
		s.RemoveAttribute("nonce")
		return nonce
	}
	return ""
}

var pagesToIgnore = []string{"/register", "/login", "/profile", "/student/operations"}

// SetLastVisitedPage stores the URI of the last visited page
func (s *HTTPSession) SetLastVisitedPage(uri string) {
	var toIgnore bool
	for _, pageToIgnore := range pagesToIgnore {
		if strings.HasPrefix(uri, pageToIgnore) {
			toIgnore = true
			break
		}
	}
	if !toIgnore {
		s.SetAttribute("lastVisitedPage", uri)
	}
}

// GetLastVisitedPage returns the URI of the last visited page
func (s *HTTPSession) GetLastVisitedPage() string {
	return s.GetStringAttribute("lastVisitedPage")
}
