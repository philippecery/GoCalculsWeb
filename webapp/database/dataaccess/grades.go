package dataaccess

import (
	"context"
	"errors"
	"log"

	"github.com/philippecery/maths/webapp/constant/user"
	"github.com/philippecery/maths/webapp/database/collection"
	"github.com/philippecery/maths/webapp/database/document"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type gradeCache map[string]*document.Grade

var cache = make(gradeCache)

func (c *gradeCache) put(grade *document.Grade) {
	cache[grade.GradeID] = grade
}

func (c *gradeCache) remove(id string) {
	delete(cache, id)
}

func (c *gradeCache) get(id string) *document.Grade {
	if grade, found := cache[id]; found {
		return grade
	}
	return nil
}

// GetAllGrades returns all the Grade documents in the Grades collection
func GetAllGrades() []*document.Grade {
	var err error
	var cursor *mongo.Cursor
	if cursor, err = collection.Grades.Find(context.TODO(), bson.M{}); err != nil {
		log.Printf("Unable to find Grade documents. Cause: %v", err)
		return nil
	}
	defer cursor.Close(context.TODO())
	var grades []*document.Grade
	for cursor.Next(context.TODO()) {
		grade := new(document.Grade)
		if err = cursor.Decode(grade); err != nil {
			log.Printf("Unable to decode Grade document. Cause: %v", err)
			return nil
		}
		cache.put(grade)
		grades = append(grades, grade)
	}
	return grades
}

// GetGradeByID returns the Grade document from the Grades collection where gradeid field is the provided id
func GetGradeByID(id string) *document.Grade {
	var grade *document.Grade
	if grade = cache.get(id); grade == nil {
		grade = new(document.Grade)
		if err := collection.Grades.FindOne(context.TODO(), bson.M{"gradeid": id}).Decode(grade); err != nil {
			log.Printf("Unable to find Grade with id %s. Cause: %v", id, err)
			return nil
		}
		cache.put(grade)
	}
	return grade
}

// CreateNewGrade creates a new Grade document in the Grades collection
func CreateNewGrade(newGrade *document.Grade) error {
	if _, err := collection.Grades.InsertOne(context.TODO(), newGrade); err != nil {
		return errors.New("Grade creation failed")
	}
	cache.put(newGrade)
	log.Printf("Grade %s created.", newGrade.GradeID)
	return nil
}

// UpdateGrade retrieves and replaces the Grade document where gradeid field equals the one in the new Grade document
func UpdateGrade(grade *document.Grade) error {
	if _, err := collection.Grades.ReplaceOne(context.TODO(), bson.M{"gradeid": grade.GradeID}, grade); err != nil {
		log.Printf("Unable to replace grade %s. Cause: %v", grade.GradeID, err)
		return errors.New("Grade update failed")
	}
	cache.put(grade)
	log.Printf("Grade %s is update", grade.GradeID)
	return nil
}

// DeleteGrade deletes the Grade document from the Grades collection where the gradeid field is the provided id
func DeleteGrade(id string) error {
	if _, err := collection.Users.UpdateMany(context.TODO(), bson.M{"role": user.Student, "gradeid": id}, bson.M{"$set": bson.M{"gradeid": ""}}); err == nil {
		if _, err := collection.Grades.DeleteOne(context.TODO(), bson.M{"gradeid": id}); err == nil {
			log.Printf("Grade %s is deleted.", id)
			cache.remove(id)
			return nil
		}
		log.Printf("Unable to delete grade %s. Cause: %v", id, err)
	} else {
		log.Printf("Unable to unassign grade %s to students. Cause: %v", id, err)
	}
	return errors.New("Grade deletion failed")
}
