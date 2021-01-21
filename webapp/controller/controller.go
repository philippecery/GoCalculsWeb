package controller

import (
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/config"
	"github.com/philippecery/maths/webapp/controller/api"
	"github.com/philippecery/maths/webapp/controller/app"
	"github.com/philippecery/maths/webapp/controller/app/admin"
	"github.com/philippecery/maths/webapp/controller/app/common"
	"github.com/philippecery/maths/webapp/controller/app/student"
	"github.com/philippecery/maths/webapp/controller/app/teacher"
	"github.com/philippecery/maths/webapp/session"
)

// SetupRoutes defines the handlers for the request paths
func SetupRoutes() {
	handleStatic("css")
	handleStatic("fonts")
	handleStatic("js")
	handleStatic("img")

	handleFunc("/", noCache(app.Home))
	handleFunc("/register", noCache(app.Register))
	handleFunc("/login", noCache(app.Login))
	handleFunc("/logout", app.Logout)

	handleFunc("/admin/user/list", noCache(authenticated(admin.UserList)))
	handleFunc("/admin/user/new", noCache(authenticated(admin.UserNew)))
	handleFunc("/admin/user/status", authenticated(admin.UserStatus))
	handleFunc("/admin/user/delete", authenticated(admin.UserDelete))

	handleFunc("/teacher/grade/list", noCache(authenticated(teacher.GradeList)))
	handleFunc("/teacher/grade/new", noCache(authenticated(teacher.GradeNew)))
	handleFunc("/teacher/grade/edit", noCache(authenticated(teacher.GradeEdit)))
	handleFunc("/teacher/grade/copy", noCache(authenticated(teacher.GradeCopy)))
	handleFunc("/teacher/grade/save", noCache(authenticated(teacher.GradeSave)))
	handleFunc("/teacher/grade/students", noCache(authenticated(teacher.GradeStudents)))
	handleFunc("/teacher/grade/unassign", authenticated(teacher.GradeUnassign))
	handleFunc("/teacher/grade/delete", authenticated(teacher.GradeDelete))
	handleFunc("/teacher/student/list", noCache(authenticated(teacher.StudentList)))
	handleFunc("/teacher/student/grade", noCache(authenticated(teacher.StudentGrade)))
	handleFunc("/teacher/student/assign", authenticated(teacher.GradeAssign))

	handleFunc("/student/dashboard", noCache(authenticated(student.Dashboard)))
	handleFunc("/student/operations", noCache(authenticated(student.Operations)))
	handleFunc("/student/results", noCache(authenticated(student.Results)))

	handleFunc("/profile", noCache(authenticated(common.Profile)))

	handleFunc("/websocket", authenticated(api.Endpoints))
}

func handleStatic(path string) {
	http.Handle("/"+path+"/", http.StripPrefix("/"+path+"/", http.FileServer(http.Dir("./static/"+path))))
}

func handleFunc(pattern string, h func(http.ResponseWriter, *http.Request, *session.HTTPSession)) {
	http.HandleFunc(pattern, securityHeaders(h))
}

func securityHeaders(h func(http.ResponseWriter, *http.Request, *session.HTTPSession)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		if httpsession := session.GetSession(w, r); httpsession != nil {
			nonce := httpsession.SetCSPNonce()
			w.Header().Set("Strict-Transport-Security", "max-age=31536000 ; includeSubDomains")
			w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'; block-all-mixed-content; default-src 'none'; connect-src wss://"+config.Config.Hostname+"; font-src 'self'; img-src 'self'; style-src 'self' 'nonce-"+nonce+"'; form-action 'self'; base-uri 'self'; script-src 'nonce-"+nonce+"'")
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
			w.Header().Set("Referrer-Policy", "no-referrer")
			w.Header().Set("X-Frame-Options", "deny")
			w.Header().Set("X-XSS-Protection", "0")
			h(w, r, httpsession)
		} else {
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}
}

func authenticated(h func(http.ResponseWriter, *http.Request, *session.HTTPSession, *session.UserInformation)) func(http.ResponseWriter, *http.Request, *session.HTTPSession) {
	return func(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession) {
		if user := httpsession.GetAuthenticatedUser(); user != nil {
			h(w, r, httpsession, user)
		} else {
			log.Println("User is not authenticated")
			http.Redirect(w, r, "/logout", http.StatusFound)
		}
	}
}

func noCache(h func(http.ResponseWriter, *http.Request, *session.HTTPSession)) func(http.ResponseWriter, *http.Request, *session.HTTPSession) {
	return func(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession) {
		if r.Method == "GET" {
			httpsession.SetLastVisitedPage(r.RequestURI)
		}
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		h(w, r, httpsession)
	}
}
