package dataaccess

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/philippecery/maths/webapp/constant/team"
	"github.com/philippecery/maths/webapp/constant/user"
	"github.com/philippecery/maths/webapp/database"
	"github.com/philippecery/maths/webapp/database/model"
)

const queryGetUserByID = `
SELECT
	u.UserID, u.UserPassword, u.StatusID, u.FailedAttempts, u.LastConnection,
	ou.TeamID, ou.RoleID
FROM Users u
	INNER JOIN TeamsUsers ou ON u.UserID = ou.UserID
WHERE u.UserID = ? AND ou.TeamDefault = ?
`

// GetUserByID returns the User model from the Users collection where userid field is the provided id
func GetUserByID(userID string) *model.User {
	user := new(model.User)
	err := database.Connection.QueryRowContext(context.TODO(), queryGetUserByID, userID, true).Scan(user.UserID, user.Password, user.Status, user.FailedAttempts, user.LastConnection)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("database: user %s doesn't exist", userID)
	case err != nil:
		log.Printf("database: unable to find user %s. Cause: %v", userID, err)
		return nil
	}
	return user
}

const queryUpdateLastConnection = "UPDATE Users SET LastConnection = NOW() WHERE UserID = ?"

// UpdateLastConnection updates the lastconnection field to current datetime for the User model where userid field is the provided id
func UpdateLastConnection(userID string) error {
	_, err := database.Connection.ExecContext(context.TODO(), queryUpdateLastConnection, userID)
	return err
}

const queryUpdateFailedAttempts = "UPDATE Users SET FailedAttempts = ? %s WHERE UserID = ?"

// UpdateFailedAttempts updates the number of failed login attempts for the User model where userid field is the provided id
func UpdateFailedAttempts(userID string, failedAttempts int) error {
	var statusUpdate string
	if failedAttempts > user.MaxFailedAttempts {
		statusUpdate = fmt.Sprintf(", StatusID = %d", user.Disabled)
	}
	query := fmt.Sprintf(queryUpdateFailedAttempts, statusUpdate)
	_, err := database.Connection.ExecContext(context.TODO(), query, failedAttempts, userID)
	return err
}

const queryGetUnregisteredUsersByTeamID = `
SELECT
	u.UserID, u.EmailAddress, u.RoleID, u.Token, u.Expires
FROM TeamsUsers ou
	INNER JOIN Users u ON TeamsUsers.UserID = Users.UserID
WHERE TeamsUsers.TeamID = ? AND Users.StatusID = ?
`

// GetUnregisteredUsersByTeamID returns all the User documents in the Users collections where status is Unregistered
func GetUnregisteredUsersByTeamID(teamID string) []*model.User {
	var users []*model.User
	if rows, err := database.Connection.QueryContext(context.TODO(), queryGetUnregisteredUsersByTeamID, teamID, user.Unregistered); err == nil {
		defer rows.Close()
		for rows.Next() {
			user := new(model.User)
			if err = rows.Scan(user.UserID, user.EmailAddress, user.Role, user.Token, user.Expires); err == nil {
				users = append(users, user)
			}
		}
	}
	return users
}

const queryGetRegisteredUsersByTeamID = `
SELECT
	u.UserID, u.EmailAddress, u.RoleID, u.Token, u.Expires
FROM TeamsUsers ou
	INNER JOIN Users u ON TeamsUsers.UserID = Users.UserID
WHERE TeamsUsers.TeamID = ? AND Users.StatusID <> ?
`

// GetRegisteredUsersByTeamID returns all the User documents in the Users collections where status is Unregistered
func GetRegisteredUsersByTeamID(teamID string) []*model.User {
	var users []*model.User
	if rows, err := database.Connection.QueryContext(context.TODO(), queryGetRegisteredUsersByTeamID, teamID, user.Unregistered); err == nil {
		defer rows.Close()
		for rows.Next() {
			user := new(model.User)
			if err = rows.Scan(user.UserID, user.EmailAddress, user.Role, user.Token, user.Expires); err == nil {
				users = append(users, user)
			}
		}
	}
	return users
}

// CreateNewUser creates a new User model in the Users collection
func CreateNewUser(teamID string, newUser *model.User, orgRoleID team.Role) error {
	var err error
	var tx *sql.Tx
	ctx := context.TODO()
	if tx, err = database.Connection.BeginTx(ctx, nil); err == nil {
		if err = createNewUser(tx, ctx, teamID, newUser, orgRoleID); err == nil {
			err = tx.Commit()
		} else {
			tx.Rollback()
		}
	}
	return err
}

const (
	queryCreateNewUser = "INSERT INTO Users (UserID, EmailAddress, Token, Expires, RoleID, StatusID) VALUES (?, ?, ?, ?, ?, ?)"
	queryAddUserToTeam = "INSERT INTO TeamsUsers (TeamID, UserID, RoleID, TeamDefault) VALUES (?, ?, ?, ?)"
)

func createNewUser(tx *sql.Tx, ctx context.Context, teamID string, newUser *model.User, orgRoleID team.Role) error {
	var err error
	if _, err = tx.ExecContext(ctx, queryCreateNewUser, newUser.UserID, newUser.EmailAddress, newUser.Token, newUser.Expires, newUser.Role, newUser.Status); err == nil {
		if _, err = tx.ExecContext(ctx, queryAddUserToTeam, teamID, newUser.UserID, orgRoleID, true); err == nil {
			log.Printf("database: registration token %s created for user %s. Token expires on %v.", newUser.Token, newUser.UserID, newUser.Expires)
		}
	} else {
		log.Printf("database: unable to create user %s. Cause: %v", newUser.UserID, err)
	}
	return err
}

// IsUserIDAvailable returns true if there is no User model in the Users collection where userid is the provided id
func IsUserIDAvailable(userID string) bool {
	return GetUserByID(userID) == nil
}

const queryDeleteUser = "DELETE FROM Users WHERE UserID = ?"

// DeleteUser deletes the User model from the Users collection where the userid field is the provided id
func DeleteUser(userID string) error {
	_, err := database.Connection.ExecContext(context.TODO(), queryDeleteUser, userID)
	return err
}

const queryToggleUserStatus = "UPDATE Users SET StatusID = ? WHERE UserID = ?"

// ToggleUserStatus retrieves the User model where the userid field is the provided id, then update the status field to Disabled if current status is Enabled, or Enabled otherwise.
func ToggleUserStatus(userID string) error {
	var err error
	if u := GetUserByID(userID); u != nil {
		var newStatus user.Status
		if u.Status == user.Enabled {
			newStatus = user.Disabled
		} else {
			newStatus = user.Enabled
		}
		_, err = database.Connection.ExecContext(context.TODO(), queryToggleUserStatus, newStatus, userID)
	}
	return err
}

const queryGetUserByToken = "SELECT UserID, EmailAddress, RoleID, Token, Expires FROM Users WHERE Token = ? AND StatusID = ?"

// GetUserByToken returns the User model from the Users collection where token field is the provided token
func GetUserByToken(token string) *model.User {
	userToken := new(model.User)
	err := database.Connection.QueryRowContext(context.TODO(), queryGetUserByToken, token, user.Unregistered).Scan(userToken.UserID, userToken.EmailAddress, userToken.Role, userToken.Token, userToken.Expires)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("database: user with token %s doesn't exist", token)
	case err != nil:
		log.Printf("database: unable to find user with token %s. Cause: %v", token, err)
		return nil
	}
	return userToken
}

const queryRegisterUser = `
UPDATE Users
SET
	Password = ?, FullName = ?, LanguageCode = ?, StatusID = ?,
	PasswordDate = NOW(), FailedAttempts = 0,
	Token = NULL, Expires = NULL
WHERE UserID = ? AND Token = ?
`

// RegisterUser retrieves and replaces the User model where userid field equals the one in the new User model
func RegisterUser(newUser *model.RegisteredUser, token string) error {
	_, err := database.Connection.ExecContext(context.TODO(), queryRegisterUser, newUser.Password, newUser.FullName, newUser.Language, newUser.Status, newUser.UserID, token)
	return err
}

const queryGetUserProfileByID = "SELECT UserID, EmailAddress, FullName, LanguageCode FROM Users WHERE UserID = ?"

// GetUserProfileByID returns the User model from the Users collection where userid field is the provided id
func GetUserProfileByID(userID string) *model.UserProfile {
	userProfile := new(model.UserProfile)
	err := database.Connection.QueryRowContext(context.TODO(), queryGetUserProfileByID, userID).Scan(userProfile.UserID, userProfile.EmailAddress, userProfile.FullName, userProfile.Language)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("database: user %s doesn't exist", userID)
	case err != nil:
		log.Printf("database: unable to find user %s. Cause: %v", userID, err)
		return nil
	default:
		userProfile.Teams = GetTeamsByUserID(userID)
	}
	return userProfile
}

const queryUpdateUserProfile = "UPDATE Users SET FullName = ?, EmailAddressTmp = ? WHERE UserID = ?"

// UpdateUserProfile updates the User model where userid field equals the one in the new User model
func UpdateUserProfile(userProfile *model.User) error {
	_, err := database.Connection.ExecContext(context.TODO(), queryUpdateUserProfile, userProfile.FullName, userProfile.EmailAddressTmp, userProfile.UserID)
	return err
}

const queryUpdateUserPassword = "UPDATE Users SET Password = ?, PasswordDate = NOW() WHERE UserID = ?"

// UpdateUserPassword updates the User model where userid field equals the one in the new User model
func UpdateUserPassword(userPassword *model.User) error {
	if _, err := database.Connection.ExecContext(context.TODO(), queryUpdateUserPassword, userPassword.Password, userPassword.UserID); err == nil {
		log.Printf("database: password updated for user %s", userPassword.UserID)
		return nil
	} else {
		log.Printf("database: unable to update password for user %s. Cause: %v", userPassword.UserID, err)
	}
	return errors.New("database: user password update failed")
}
