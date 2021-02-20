package main // import "github.com/philippecery/maths/student"

import (
	"log"

	_ "github.com/philippecery/maths/webapp/config"

	"github.com/philippecery/maths/webapp/controller"
	"github.com/philippecery/maths/webapp/database"
	"github.com/philippecery/maths/webapp/server"
)

func main() {
	var err error
	err = database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer database.Disconnect()
	controller.SetupRoutes()
	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
