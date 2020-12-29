package controller

import (
	"net/http"

	"github.com/philippecery/maths/webapp/controller/app"
	"github.com/philippecery/maths/webapp/controller/app/admin"
	"github.com/philippecery/maths/webapp/controller/app/teacher"
	"github.com/philippecery/maths/webapp/session"
)

// SetupRoutes defines the handlers for the request paths
func SetupRoutes() {
	handleStatic("css")
	handleStatic("fonts")
	handleStatic("js")

	handleFunc("/", noCache(app.Home))
	handleFunc("/login", noCache(app.Login))
	handleFunc("/admin/user/list", noCache(admin.UserList))
	handleFunc("/admin/user/new", noCache(admin.UserNew))
	handleFunc("/admin/user/status", admin.UserStatus)
	handleFunc("/admin/user/delete", admin.UserDelete)

	handleFunc("/teacher/grade/list", noCache(teacher.GradeList))
	handleFunc("/teacher/grade/new", noCache(teacher.GradeNew))
	handleFunc("/teacher/grade/edit", noCache(teacher.GradeEdit))
	handleFunc("/teacher/grade/copy", noCache(teacher.GradeCopy))
	handleFunc("/teacher/grade/save", noCache(teacher.GradeSave))
	handleFunc("/teacher/grade/delete", teacher.GradeDelete)
	handleFunc("/teacher/student/list", noCache(teacher.StudentList))
	handleFunc("/teacher/student/edit", app.Todo)
	//handleFunc("/teacher/student/edit", noCache(teacher.StudentEdit))

	handleFunc("/register", noCache(app.Register))
	handleFunc("/operations", app.Todo)
	//handleFunc("/operations", noCache(student.Operations))
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
		w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'; block-all-mixed-content; default-src 'none'; connect-src 'self'; font-src 'self'; img-src 'self'; style-src 'self'; form-action 'self'; base-uri 'self'; script-src 'nonce-"+token+"' 'unsafe-inline'")
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
