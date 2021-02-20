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

// NewHomeworkSession stores the current homework session
func NewHomeworkSession(newSession *document.HomeworkSession) error {
	if _, err := collection.Homework.InsertOne(context.TODO(), newSession); err != nil {
		return errors.New("HomeworkSession creation failed")
	}
	log.Printf("Session stored.")
	return nil
}

// UpdateHomeworkSession updates the current homework session
func UpdateHomeworkSession(newSession *document.HomeworkSession) error {
	var err error
	var result *mongo.UpdateResult
	if result, err = collection.Homework.ReplaceOne(context.TODO(), bson.M{"sessionid": newSession.SessionID}, newSession); err != nil {
		return errors.New("HomeworkSession update failed")
	}
	if result.MatchedCount == 1 && result.ModifiedCount == 1 {
		log.Printf("Session updated.")
	} else {
		log.Printf("MatchedCount = %d ; ModifedCount = %d", result.MatchedCount, result.ModifiedCount)
	}
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
	if cursor, err = collection.Homework.Find(context.TODO(), filters, findOptions); err != nil {
		log.Printf("Unable to find HomeworkSession documents for user %s. Cause: %v", userID, err)
		return nil, 0
	}
	var homeworkSessions []*document.HomeworkSession
	if err = cursor.All(context.TODO(), &homeworkSessions); err != nil {
		log.Printf("Unable to decode HomeworkSession documents. Cause: %v", err)
		return nil, 0
	}
	var nbTotal int64
	if nbTotal, err = collection.Homework.CountDocuments(context.TODO(), bson.M{"userid": userID}); err != nil {
		log.Printf("Unable to count HomeworkSession documents. Cause: %v", err)
		return nil, 0
	}
	return homeworkSessions, nbTotal
}

// GetSessionByID returns the homework sessions for the specified user
func GetSessionByID(id string) *document.HomeworkSession {
	homeworkSession := new(document.HomeworkSession)
	if err := collection.Homework.FindOne(context.TODO(), bson.M{"sessionid": id}).Decode(homeworkSession); err != nil {
		log.Printf("Unable to find HomeworkSession with id %s. Cause: %v", id, err)
		return nil
	}
	return homeworkSession
}
