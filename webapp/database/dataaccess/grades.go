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

// GetAllGrades returns all the Grade documents in the Grades collection
func GetAllGrades() []*document.Grade {
	var err error
	var cursor *mongo.Cursor
	if cursor, err = collection.Grades.Find(context.TODO(), bson.M{}); err != nil {
		log.Printf("Unable to find Grade documents. Cause: %v", err)
		return nil
	}
	var grades []*document.Grade
	for cursor.Next(context.TODO()) {
		var grade document.Grade
		if err = cursor.Decode(&grade); err != nil {
			log.Printf("Unable to decode Grade document. Cause: %v", err)
			return nil
		}
		grades = append(grades, &grade)
	}
	return grades
}

// GetGradeByID returns the Grade document from the Grades collection where gradeid field is the provided id
func GetGradeByID(id string) *document.Grade {
	var grade document.Grade
	if err := collection.Grades.FindOne(context.TODO(), bson.M{"gradeid": id}).Decode(&grade); err != nil {
		log.Printf("Unable to find Grade with id %s. Cause: %v", id, err)
		return nil
	}
	return &grade
}

// CreateNewGrade creates a new Grade document in the Grades collection
func CreateNewGrade(newGrade *document.Grade) error {
	if _, err := collection.Grades.InsertOne(context.TODO(), newGrade); err != nil {
		return errors.New("Grade creation failed")
	}
	log.Printf("Grade %s created.", newGrade.GradeID)
	return nil
}

// UpdateGrade retrieves and replaces the Grade document where gradeid field equals the one in the new Grade document
func UpdateGrade(newGrade *document.Grade) error {
	if _, err := collection.Grades.ReplaceOne(context.TODO(), bson.M{"gradeid": newGrade.GradeID}, newGrade); err != nil {
		log.Printf("Unable to replace grade %s. Cause: %v", newGrade.GradeID, err)
		return errors.New("Grade creation failed")
	}
	log.Printf("Grade %s is created", newGrade.GradeID)
	return nil
}

// DeleteGrade deletes the Grade document from the Grades collection where the gradeid field is the provided id
func DeleteGrade(id string) error {
	if _, err := collection.Grades.DeleteOne(context.TODO(), bson.M{"gradeid": id}); err != nil {
		log.Printf("Unable to delete grade %s. Cause: %v", id, err)
		return errors.New("Grade deletion failed")
	}
	log.Printf("Grade %s is deleted.", id)
	return nil
}
