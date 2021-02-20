package collection

import "go.mongodb.org/mongo-driver/mongo"

const (
	sessions = "sessions"
	users    = "users"
	grades   = "grades"
	homework = "homework"
)

// Sessions collection
var Sessions *mongo.Collection

// Users collection
var Users *mongo.Collection

// Grades collection
var Grades *mongo.Collection

// Homework collection
var Homework *mongo.Collection

// Register initiates the MongoDB collections
func Register(db *mongo.Database) {
	Sessions = db.Collection(sessions)
	Users = db.Collection(users)
	Grades = db.Collection(grades)
	Homework = db.Collection(homework)
}
