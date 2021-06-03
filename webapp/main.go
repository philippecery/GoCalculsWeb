package main // import "github.com/philippecery/maths/student"

import (
	"log"

	_ "github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/services/email"

	"github.com/philippecery/maths/webapp/controller"
	"github.com/philippecery/maths/webapp/database"
	"github.com/philippecery/maths/webapp/server"
)

func main() {
	var err error
	err = database.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	email.Setup()
	controller.SetupRoutes()
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
