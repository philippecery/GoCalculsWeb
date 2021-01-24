package controller

import (
	"log"
	"net/http"

	"github.com/philippecery/maths/webapp/controller/app/student/wss"

	"github.com/philippecery/maths/webapp/config"
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

	handleFunc("/admin/user/list", noCache(accessControl(admin.UserList)))
	handleFunc("/admin/user/new", noCache(accessControl(admin.UserNew)))
	handleFunc("/admin/user/status", accessControl(admin.UserStatus))
	handleFunc("/admin/user/delete", accessControl(admin.UserDelete))

	handleFunc("/teacher/grade/list", noCache(accessControl(teacher.GradeList)))
	handleFunc("/teacher/grade/new", noCache(accessControl(teacher.GradeNew)))
	handleFunc("/teacher/grade/edit", noCache(accessControl(teacher.GradeEdit)))
	handleFunc("/teacher/grade/copy", noCache(accessControl(teacher.GradeCopy)))
	handleFunc("/teacher/grade/save", noCache(accessControl(teacher.GradeSave)))
	handleFunc("/teacher/grade/students", noCache(accessControl(teacher.GradeStudents)))
	handleFunc("/teacher/grade/unassign", accessControl(teacher.GradeUnassign))
	handleFunc("/teacher/grade/delete", accessControl(teacher.GradeDelete))
	handleFunc("/teacher/student/list", noCache(accessControl(teacher.StudentList)))
	handleFunc("/teacher/student/grade", noCache(accessControl(teacher.StudentGrade)))
	handleFunc("/teacher/student/assign", accessControl(teacher.GradeAssign))

	handleFunc("/student/dashboard", noCache(accessControl(student.Dashboard)))
	handleFunc("/student/operations", noCache(accessControl(student.Operations)))
	handleFunc("/student/results", noCache(accessControl(student.Results)))
	handleFunc("/student/websocket", accessControl(wss.Endpoints))

	handleFunc("/user/profile", noCache(accessControl(common.Profile)))

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
			w.Header().Set("Content-Security-Policy", "frame-ancestors 'none'; block-all-mixed-content; default-src 'none'; connect-src wss://"+config.Config.Hostname+" 'self'; font-src 'self'; img-src 'self'; style-src 'self' 'nonce-"+nonce+"'; form-action 'self'; base-uri 'self'; script-src 'nonce-"+nonce+"'")
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

func accessControl(h func(http.ResponseWriter, *http.Request, *session.HTTPSession, *session.UserInformation)) func(http.ResponseWriter, *http.Request, *session.HTTPSession) {
	return func(w http.ResponseWriter, r *http.Request, httpsession *session.HTTPSession) {
		if user := httpsession.GetAuthenticatedUser(); user != nil {
			if user.HasAccessTo(r.RequestURI) {
				h(w, r, httpsession, user)
				return
			}
			log.Printf("User %s does not have access to %s", user.UserID, r.RequestURI)
		} else {
			log.Println("User is not authenticated")
		}
		http.Redirect(w, r, "/logout", http.StatusFound)
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
