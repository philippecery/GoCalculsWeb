package app

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/philippecery/maths/webapp/constant"
	"github.com/philippecery/maths/webapp/util"
)

var validUserID = regexp.MustCompile("^[a-z]{2,}(\\.?[a-z]{2,})*$")

var validEmailAddress = regexp.MustCompile("^[a-z0-9._%+-]+@[a-z0-9.-]+\\.[a-z]{2,}$")
var validName = regexp.MustCompile("^[A-ZÀ-ÿa-z][-,. 'A-ZÀ-ÿa-z]*$")

var validPassword = []*regexp.Regexp{
	regexp.MustCompile("^.*[0-9]+.*$"),
	regexp.MustCompile("^.*[a-z]+.*$"),
	regexp.MustCompile("^.*[A-Z]+.*$"),
	regexp.MustCompile("^.*[^0-9a-zA-Z]+.*$"),
}

// ValidateUserID check the submitted user ID is valid.
func ValidateUserID(userID string) (string, error) {
	if len(userID) <= 32 && validUserID.MatchString(userID) {
		return userID, nil
	}
	return "", fmt.Errorf("errorInvalidUserID")
}

// ValidateRoleID checks the submitted role ID is valid.
func ValidateRoleID(role string) (int, error) {
	if roleID, _ := strconv.Atoi(role); roleID > 0 && constant.UserRole(roleID).IsValid() {
		return roleID, nil
	}
	return 0, fmt.Errorf("errorInvalidRoleID")
}

// ValidateEmailAddress checks the submitted email address is valid.
func ValidateEmailAddress(emailAddress string) (string, error) {
	if len(emailAddress) <= 254 && validEmailAddress.MatchString(emailAddress) {
		return emailAddress, nil
	}
	return "", fmt.Errorf("errorInvalidEmailAddress")
}

// ValidatePassword checks the submitted password is valid.
func ValidatePassword(password, passwordConfirm string) (string, error) {
	if password == passwordConfirm {
		if validatePasswordStrength(password) {
			return util.ProtectPassword(password), nil
		}
	}
	return "", fmt.Errorf("errorPassword")
}

func validatePasswordStrength(password string) bool {
	if len(password) < 8 {
		return false
	}
	for _, regex := range validPassword {
		if !regex.MatchString(password) {
			return false
		}
	}
	return true
}

// ValidateName checks the submitted name is valid.
func ValidateName(name string) (string, error) {
	if validName.MatchString(name) {
		return name, nil
	}
	return "", fmt.Errorf("errorInvalidName")
}

// ValidateDate checks the submitted date is valid.
func ValidateDate(date string) (*time.Time, error) {
	if date, err := time.Parse("2000-12-31", date); err == nil {
		return &date, nil
	}
	return nil, fmt.Errorf("errorInvalidDate")
}
