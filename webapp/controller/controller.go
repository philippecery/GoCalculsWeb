package controller

import (
	"github.com/philippecery/maths/webapp/services/auth"
	"github.com/philippecery/maths/webapp/services/public"
	"github.com/philippecery/maths/webapp/services/student/wss"

	"github.com/philippecery/maths/webapp/services/admin"
	"github.com/philippecery/maths/webapp/services/common"
	"github.com/philippecery/maths/webapp/services/student"
	"github.com/philippecery/maths/webapp/services/teacher"
)

// SetupRoutes defines the handlers for the request paths
func SetupRoutes() {
	handleStatic("css")
	handleStatic("fonts")
	handleStatic("js")
	handleStatic("img")

	handleFunc("/", noCache(public.Home))
	handleFunc("/signup", noCache(public.SignUp))
	handleFunc("/register", noCache(public.Register))
	handleFunc("/login", noCache(auth.Login))
	handleFunc("/logout", auth.Logout)

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
	handleFunc("/user/changePassword", noCache(accessControl(common.ChangePassword)))
}
