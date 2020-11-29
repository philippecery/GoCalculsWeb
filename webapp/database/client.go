package database

import (
	"context"
	"log"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/database/collection"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

// Connect creates the mongodb client and connects to the mongodb server
func Connect() error {
	var err error
	clientOptions := options.Client().SetAuth(
		options.Credential{
			AuthMechanism: config.Config.DB.AuthMechanism,
			AuthSource:    config.Config.DB.AuthSource,
			Username:      config.Config.DB.UserName,
			Password:      config.Config.DB.Password,
		},
	).ApplyURI(config.Config.DB.URL)
	if client, err = mongo.Connect(context.TODO(), clientOptions); err == nil {
		if err = client.Ping(context.TODO(), nil); err == nil {
			log.Println("Connected to MongoDB!")
			database := client.Database("maths")
			collection.Users = database.Collection("users")
		}
	}
	return err
}

// Disconnect disconnects the client from the server
func Disconnect() error {
	err := client.Disconnect(context.TODO())
	if err == nil {
		log.Println("Connection to MongoDB closed.")
	}
	return err
}
