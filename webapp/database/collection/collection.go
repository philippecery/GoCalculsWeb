package collection

import "go.mongodb.org/mongo-driver/mongo"

const (
	users    = "users"
	grades   = "grades"
	sessions = "sessions"
)

// Users collection
var Users *mongo.Collection

// Grades collection
var Grades *mongo.Collection

// Sessions collection
var Sessions *mongo.Collection

// Register initiates the MongoDB collections
func Register(db *mongo.Database) {
	Users = db.Collection(users)
	Grades = db.Collection(grades)
	Sessions = db.Collection(sessions)
}
