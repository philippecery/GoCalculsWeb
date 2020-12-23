package controller

import (
	"net/http"

	"github.com/philippecery/maths/webapp/controller/app"
	"github.com/philippecery/maths/webapp/session"
)

// SetupRoutes defines the handlers for the request paths
func SetupRoutes() {
	handleStatic("css")
	handleStatic("fonts")
	handleStatic("js")

	handleFunc("/", noCache(app.Home))
	handleFunc("/login", noCache(app.Login))
	handleFunc("/admin/users", noCache(app.Users))
	handleFunc("/admin/newUser", noCache(app.NewUser))
	handleFunc("/admin/status", app.Status)
	handleFunc("/admin/delete", app.Delete)
	handleFunc("/register", noCache(app.Register))
	handleFunc("/operations", app.Todo)
	handleFunc("/logout", app.Logout)

	//handleFunc("/websocket", api.Endpoints)
}

func handleStatic(path string) {
	http.Handle("/"+path+"/", http.StripPrefix("/"+path+"/", http.FileServer(http.Dir("./static/"+path))))
}

func handleFunc(pattern string, h func(http.ResponseWriter, *http.Request)) {
	http.HandleFunc(pattern, securityHeaders(h))
}

func securityHeaders(h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token := session.GetSession(w, r).SetCSPNonce()
		w.Header().Set("Strict-Transport-Security", "max-age=31536000 ; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'; block-all-mixed-content; default-src 'none'; connect-src 'self'; font-src 'self'; style-src 'self'; base-uri 'self'; script-src 'nonce-"+token+"' 'unsafe-inline'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		h(w, r)
	}
}

func noCache(h func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		h(w, r)
	}
}
