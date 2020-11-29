package controller

import (
	"net/http"

	"github.com/philippecery/maths/webapp/controller/api"
	"github.com/philippecery/maths/webapp/controller/app"
)

// SetupRoutes defines the handlers for the request paths
func SetupRoutes() {
	handleStatic("css")
	handleStatic("fonts")
	handleStatic("js")

	http.HandleFunc("/", app.Home)
	http.HandleFunc("/login", app.Login)
	http.HandleFunc("/admin/users", app.Users)
	http.HandleFunc("/admin/newUser", app.NewUser)
	http.HandleFunc("/admin/status", app.Status)
	http.HandleFunc("/admin/delete", app.Delete)
	http.HandleFunc("/register", app.Register)
	http.HandleFunc("/operations", app.Todo)
	http.HandleFunc("/logout", app.Logout)

	http.HandleFunc("/websocket", api.Endpoints)
}

func handleStatic(path string) {
	http.Handle("/"+path+"/", http.StripPrefix("/"+path+"/", http.FileServer(http.Dir("./static/"+path))))
}
