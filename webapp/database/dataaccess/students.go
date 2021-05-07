package dataaccess

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/philippecery/maths/webapp/constant/user"
	"github.com/philippecery/maths/webapp/database"
	"github.com/philippecery/maths/webapp/database/model"
)

const queryGetStudentsByTeamID = `
SELECT
	u.UserID, u.FullName,
	g.GradeID, g.GradeName, g.GradeDescription
FROM Users u
	INNER JOIN TeamsUsers tu ON u.UserID = tu.UserID
	LEFT JOIN Students s ON u.UserID = s.UserID
	INNER JOIN Grades g ON s.GradeID = g.GradeID
WHERE tu.TeamID = ? AND u.RoleID = ? AND u.StatusID = ?
`

// GetStudentsByTeamID returns all the User documents in the Users collections where status is Enabled and role is Student
func GetStudentsByTeamID(teamID string) []*model.Student {
	var students []*model.Student
	if rows, err := database.Connection.QueryContext(context.TODO(), queryGetStudentsByTeamID, teamID, user.ChildOrStudent, user.Enabled); err == nil {
		defer rows.Close()
		for rows.Next() {
			student := new(model.Student)
			if err = rows.Scan(student.UserID, student.FullName, student.Grade.GradeID, student.Grade.Name, student.Grade.Description); err == nil {
				students = append(students, student)
			}
		}
	}
	return students
}

const queryGetStudentByID = `
SELECT
	u.UserID, u.FullName, u.LanguageCode, s.GradeID
FROM Users u
	LEFT JOIN Students s ON u.UserID = s.UserID
WHERE u.UserID = ? AND u.RoleID = ?
`

// GetStudentByID returns the User model from the Users collection where userid field is the provided id and the role is Student
func GetStudentByID(userID string) *model.Student {
	student := new(model.Student)
	var gradeID string
	err := database.Connection.QueryRowContext(context.TODO(), queryGetStudentByID, userID, user.ChildOrStudent).Scan(student.UserID, student.FullName, student.Language, gradeID)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("database: student %s doesn't exist", userID)
	case err != nil:
		log.Printf("database: unable to find student %s. Cause: %v", userID, err)
		return nil
	default:
		student.Grade = GetGradeByID(gradeID)
	}
	return student
}

const querySetGradeForStudents = "INSERT IGNORE Students (UserID, GradeID) VALUES "

// SetGradeForStudents updates the gradeid of selected students
func SetGradeForStudents(gradeID string, students []string) error {
	var err error
	if len(students) > 0 {
		valueStrings := make([]string, 0, len(students))
		valueArgs := make([]interface{}, 0, len(students)*2)
		for _, studentID := range students {
			valueStrings = append(valueStrings, "(?, ?)")
			valueArgs = append(valueArgs, studentID)
			valueArgs = append(valueArgs, gradeID)
		}
		query := querySetGradeForStudents + strings.Join(valueStrings, ", ")
		if _, err = database.Connection.ExecContext(context.TODO(), query, valueArgs...); err == nil {
			log.Printf("database: grade %s assigned to students", gradeID)
		} else {
			log.Printf("database: unable to assign grade %s to students. Cause: %v", gradeID, err)
		}
	}
	return err
}

const queryUnassignGradeForStudent = "DELETE FROM Students WHERE UserID = ? AND GradeID = ?"

// UnassignGradeForStudent removes the gradeid of selected student
func UnassignGradeForStudent(gradeID, studentID string) error {
	_, err := database.Connection.ExecContext(context.TODO(), queryUnassignGradeForStudent, studentID, gradeID)
	return err
}

const queryAssignGradeForStudent = "INSERT IGNORE Students (UserID, GradeID) VALUES (?, ?)"

// AssignGradeForStudent sets the gradeid for selected student
func AssignGradeForStudent(gradeID, studentID string) error {
	_, err := database.Connection.ExecContext(context.TODO(), queryAssignGradeForStudent, studentID, gradeID)
	return err
}
