package dataaccess

import (
	"context"
	"errors"
	"log"

	"github.com/philippecery/maths/webapp/constant"
	"github.com/philippecery/maths/webapp/database/collection"
	"github.com/philippecery/maths/webapp/database/document"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllStudents returns all the User documents in the Users collections where status is Enabled and role is Student
func GetAllStudents() []*document.Student {
	var err error
	var cursor *mongo.Cursor
	if cursor, err = collection.Users.Find(context.TODO(), bson.M{"status": constant.Enabled, "role": constant.Student}); err != nil {
		log.Printf("Unable to find User document. Cause: %v", err)
		return nil
	}
	defer cursor.Close(context.TODO())
	var students []*document.Student
	for cursor.Next(context.TODO()) {
		student := new(document.Student)
		if err = cursor.Decode(student); err != nil {
			log.Printf("Unable to decode User document. Cause: %v", err)
			return nil
		}
		if student.GradeID != "" {
			student.Grade = GetGradeByID(student.GradeID)
		}
		students = append(students, student)
	}
	return students
}

// GetStudentByID returns the User document from the Users collection where userid field is the provided id and the role is Student
func GetStudentByID(id string) *document.Student {
	student := new(document.Student)
	if err := collection.Users.FindOne(context.TODO(), bson.M{"userid": id, "role": constant.Student}).Decode(student); err != nil {
		log.Printf("Unable to find User with id %s and role Student. Cause: %v", id, err)
		return nil
	}
	if student.GradeID != "" {
		student.Grade = GetGradeByID(student.GradeID)
	}
	return student
}

// SetGradeForStudents updates the gradeid of selected students
func SetGradeForStudents(gradeID string, students []string) error {
	if _, err := collection.Users.UpdateMany(context.TODO(), bson.M{"role": constant.Student, "userid": bson.M{"$in": students}}, bson.M{"$set": bson.M{"gradeid": gradeID}}); err != nil {
		log.Printf("Unable to set grade %s for selected students. Cause: %v", gradeID, err)
		return errors.New("Unable to set grade %s for selected students")
	}
	return nil
}

// UnassignGradeForStudent removes the gradeid of selected student
func UnassignGradeForStudent(gradeID, studentID string) error {
	if _, err := collection.Users.UpdateOne(context.TODO(), bson.M{"role": constant.Student, "userid": studentID, "gradeid": gradeID}, bson.M{"$set": bson.M{"gradeid": ""}}); err != nil {
		log.Printf("Unable to reset grade %s for student %s. Cause: %v", gradeID, studentID, err)
		return errors.New("Unable to reset grade for selected student")
	}
	return nil
}

// AssignGradeForStudent sets the gradeid for selected student
func AssignGradeForStudent(gradeID, studentID string) error {
	if _, err := collection.Users.UpdateOne(context.TODO(), bson.M{"role": constant.Student, "userid": studentID}, bson.M{"$set": bson.M{"gradeid": gradeID}}); err != nil {
		log.Printf("Unable to set grade %s for student %s. Cause: %v", gradeID, studentID, err)
		return errors.New("Unable to reset grade for selected student")
	}
	return nil
}
