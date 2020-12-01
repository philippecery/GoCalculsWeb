package session

import (
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
