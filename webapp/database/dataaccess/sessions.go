package dataaccess

import (
	"context"
	"errors"
	"log"

	"github.com/philippecery/maths/webapp/database/collection"
	"github.com/philippecery/maths/webapp/database/document"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func StoreHomeworkSession(newSession *document.HomeworkSession) error {
	if _, err := collection.Sessions.InsertOne(context.TODO(), newSession); err != nil {
		return errors.New("HomeworkSession creation failed")
	}
	log.Printf("Session stored.")
	return nil
}

func GetSessionsByUserID(userID string) []*document.HomeworkSession {
	var err error
	var cursor *mongo.Cursor
	if cursor, err = collection.Sessions.Find(context.TODO(), bson.M{"userid": userID}); err != nil {
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
