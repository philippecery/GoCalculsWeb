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

// GetSessionsByUserID returns the paginated homework sessions for the specified user, along with the total number of sessions
func GetSessionsByUserID(userID string, homeworkType, status, page int) ([]*document.HomeworkSession, int64) {
	var err error
	var cursor *mongo.Cursor
	filters := bson.M{"userid": userID}
	if homeworkType >= 0 {
		filters["typeid"] = homeworkType
	}
	if status >= 0 {
		filters["status"] = status
	}
	findOptions := options.Find().SetSort(bson.M{"starttime": -1}).SetLimit(nbSessionsPerPage).SetSkip(int64((page - 1) * nbSessionsPerPage))
	if cursor, err = collection.Sessions.Find(context.TODO(), filters, findOptions); err != nil {
		log.Printf("Unable to find HomeworkSession documents for user %s. Cause: %v", userID, err)
		return nil, 0
	}
	var homeworkSessions []*document.HomeworkSession
	if err = cursor.All(context.TODO(), &homeworkSessions); err != nil {
		log.Printf("Unable to decode HomeworkSession documents. Cause: %v", err)
		return nil, 0
	}
	var nbTotal int64
	if nbTotal, err = collection.Sessions.CountDocuments(context.TODO(), bson.M{"userid": userID}); err != nil {
		log.Printf("Unable to count HomeworkSession documents. Cause: %v", err)
		return nil, 0
	}
	return homeworkSessions, nbTotal
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
