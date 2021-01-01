package app

import (
	"log"
	"net/http"
)

// Home redirects to the login page
func Home(w http.ResponseWriter, r *http.Request) {
	log.Println("/: Redirecting to Login page")
	http.Redirect(w, r, "/login", http.StatusFound)
}
