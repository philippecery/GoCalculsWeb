package dataaccess

import (
	"context"

	"github.com/philippecery/maths/webapp/constant/user"
	"github.com/philippecery/maths/webapp/database"
	"github.com/philippecery/maths/webapp/database/model"
)

const queryGetAllUnregisteredUsers = "SELECT UserID, EmailAddress, RoleID, Token, Expires FROM Users WHERE StatusID = ?"

// GetAllUnregisteredUsers returns all the User documents in the Users collections where status is Unregistered
func GetAllUnregisteredUsers() []*model.User {
	var users []*model.User
	if rows, err := database.Connection.QueryContext(context.TODO(), queryGetAllUnregisteredUsers, user.Unregistered); err == nil {
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

const queryGetAllRegisteredUsers = "SELECT UserID, EmailAddress, RoleID, Token, Expires FROM Users WHERE StatusID <> ?"

// GetAllRegisteredUsers returns all the User documents in the Users collections where status is NOT Unregistered
func GetAllRegisteredUsers() []*model.User {
	var users []*model.User
	if rows, err := database.Connection.QueryContext(context.TODO(), queryGetAllRegisteredUsers, user.Unregistered); err == nil {
		defer rows.Close()
		for rows.Next() {
			user := new(model.User)
			if err = rows.Scan(user.UserID, user.FullName, user.EmailAddress, user.Role, user.LastConnection); err == nil {
				users = append(users, user)
			}
		}
	}
	return users
}
