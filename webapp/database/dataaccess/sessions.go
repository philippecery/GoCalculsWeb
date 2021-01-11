package dataaccess

import (
	"context"
	"errors"
	"log"

	"github.com/philippecery/maths/webapp/database/collection"
	"github.com/philippecery/maths/webapp/database/document"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StoreHomeworkSession stores the current homework session
func StoreHomeworkSession(newSession *document.HomeworkSession) error {
	if _, err := collection.Sessions.InsertOne(context.TODO(), newSession); err != nil {
		return errors.New("HomeworkSession creation failed")
	}
	log.Printf("Session stored.")
	return nil
}

const nbSessionsPerPage = 10

// GetSessionsByUserID returns the homework sessions for the specified user
func GetSessionsByUserID(userID string, page int) []*document.HomeworkSession {
	var err error
	var cursor *mongo.Cursor
	findOptions := options.Find().SetSort(bson.M{"startdate": -1}).SetLimit(nbSessionsPerPage).SetSkip(int64(page * nbSessionsPerPage))
	if cursor, err = collection.Sessions.Find(context.TODO(), bson.M{"userid": userID}, findOptions); err != nil {
		log.Printf("Unable to find HomeworkSession documents. Cause: %v", err)
		return nil
	}
	var homeworkSessions []*document.HomeworkSession
	if err = cursor.All(context.TODO(), &homeworkSessions); err != nil {
		log.Printf("Unable to decode HomeworkSession documents. Cause: %v", err)
		return nil
	}
	return homeworkSessions
}

// GetSessionByID returns the homework sessions for the specified user
func GetSessionByID(id string) *document.HomeworkSession {
	homeworkSession := new(document.HomeworkSession)
	if err := collection.Sessions.FindOne(context.TODO(), bson.M{"sessionid": id}).Decode(homeworkSession); err != nil {
		log.Printf("Unable to find HomeworkSession with id %s. Cause: %v", id, err)
		return nil
	}
	return homeworkSession
}
