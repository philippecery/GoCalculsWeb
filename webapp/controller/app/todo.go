package app

import (
	"log"
	"net/http"
)

// Todo controller is there only to display a WORK IN PROGRESS page when the feature is not available yet.
func Todo(w http.ResponseWriter, r *http.Request) {
	if err := Templates.ExecuteTemplate(w, "todo.html.tpl", nil); err != nil {
		log.Fatalf("Error while executing template 'todo': %v\n", err)
	}
}
