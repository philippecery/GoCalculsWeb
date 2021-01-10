package dataaccess

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/philippecery/maths/webapp/constant"
	"github.com/philippecery/maths/webapp/database/collection"
	"github.com/philippecery/maths/webapp/database/document"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllUsers returns all the User documents in the Users collection
func GetAllUsers() []*document.User {
	var err error
	var cursor *mongo.Cursor
	if cursor, err = collection.Users.Find(context.TODO(), bson.M{}); err != nil {
		log.Printf("Unable to find User documents. Cause: %v", err)
		return nil
	}
	var users []*document.User
	if err = cursor.All(context.TODO(), &users); err != nil {
		log.Printf("Unable to decode User documents. Cause: %v", err)
		return nil
	}
	return users
}

// GetUserByID returns the User document from the Users collection where userid field is the provided id
func GetUserByID(id string) *document.User {
	user := new(document.User)
	if err := collection.Users.FindOne(context.TODO(), bson.M{"userid": id}).Decode(user); err != nil {
		log.Printf("Unable to find User with id %s. Cause: %v", id, err)
		return nil
	}
	return user
}

// UpdateLastConnection updates the lastconnection field to current datetime for the User document where userid field is the provided id
func UpdateLastConnection(id string) {
	if _, err := collection.Users.UpdateOne(context.TODO(), bson.M{"userid": id}, bson.M{"$set": bson.M{"lastconnection": time.Now()}}); err != nil {
		log.Printf("Unable to update %s's last connection time. Cause: %v", id, err)
	}
}

// GetAllUnregisteredUsers returns all the User documents in the Users collections where status is Unregistered
func GetAllUnregisteredUsers() []*document.User {
	var err error
	var cursor *mongo.Cursor
	if cursor, err = collection.Users.Find(context.TODO(), bson.M{"status": constant.Unregistered}); err != nil {
		log.Printf("Unable to find User document. Cause: %v", err)
		return nil
	}
	var unregisteredUsers []*document.User
	if err = cursor.All(context.TODO(), &unregisteredUsers); err != nil {
		log.Printf("Unable to decode User document. Cause: %v", err)
		return nil
	}
	return unregisteredUsers
}

// GetAllRegisteredUsers returns all the User documents in the Users collections where status is NOT Unregistered
func GetAllRegisteredUsers() []*document.User {
	var err error
	var cursor *mongo.Cursor
	if cursor, err = collection.Users.Find(context.TODO(), bson.M{"status": bson.M{"$ne": constant.Unregistered}}); err != nil {
		log.Printf("Unable to find User document. Cause: %v", err)
		return nil
	}
	var registeredUsers []*document.User
	if err = cursor.All(context.TODO(), &registeredUsers); err != nil {
		log.Printf("Unable to decode User document. Cause: %v", err)
		return nil
	}
	return registeredUsers
}

// CreateNewUser creates a new User document in the Users collection
func CreateNewUser(newUser *document.UnregisteredUser) error {
	if _, err := collection.Users.InsertOne(context.TODO(), newUser); err != nil {
		return errors.New("Registration token creation failed")
	}
	log.Printf("Registration token %s created for user %s. Token expires on %v.", newUser.Token, newUser.UserID, newUser.Expires)
	return nil
}

// IsUserIDAvailable returns true if there is no User document in the Users collection where userid is the provided id
func IsUserIDAvailable(id string) bool {
	if user := GetUserByID(id); user == nil {
		return true
	}
	return false
}

// DeleteUser deletes the User document from the Users collection where the userid field is the provided id
func DeleteUser(id string) error {
	if _, err := collection.Users.DeleteOne(context.TODO(), bson.M{"userid": id}); err != nil {
		log.Printf("Unable to delete user %s. Cause: %v", id, err)
		return errors.New("User deletion failed")
	}
	log.Printf("User %s is deleted.", id)
	return nil
}

// ToggleUserStatus retrieves the User document where the userid field is the provided id, then update the status field to Disabled if current status is Enabled, or Enabled otherwise.
func ToggleUserStatus(id string) error {
	if user := GetUserByID(id); user != nil {
		var newStatus constant.UserStatus
		if user.Status == constant.Enabled {
			newStatus = constant.Disabled
		} else {
			newStatus = constant.Enabled
		}
		if _, err := collection.Users.UpdateOne(context.TODO(), bson.M{"userid": id}, bson.M{"$set": bson.M{"status": newStatus}}); err != nil {
			log.Printf("Unable to update status for user %s. Cause: %v", id, err)
			return errors.New("User status update failed")
		}
		log.Printf("Status is updated for user %s.", id)
		return nil
	}
	return fmt.Errorf("User %s not found", id)
}

// GetUserByToken returns the User document from the Users collection where token field is the provided token
func GetUserByToken(token string) *document.User {
	unregisteredUser := new(document.User)
	if err := collection.Users.FindOne(context.TODO(), bson.M{"token": token}).Decode(unregisteredUser); err != nil {
		log.Printf("Unable to find unregistered user with token %s. Cause: %v", token, err)
		return nil
	}
	return unregisteredUser
}

// RegisterUser retrieves and replaces the User document where userid field equals the one in the new User document
func RegisterUser(newUser *document.RegisteredUser, token string) error {
	if _, err := collection.Users.ReplaceOne(context.TODO(), bson.M{"userid": newUser.UserID}, newUser); err != nil {
		log.Printf("Unable to replace user %s. Cause: %v", newUser.UserID, err)
		return errors.New("User creation failed")
	}
	log.Printf("User %s is created", newUser.UserID)
	return nil
}
