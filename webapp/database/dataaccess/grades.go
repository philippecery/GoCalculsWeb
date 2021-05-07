package dataaccess

import (
	"context"
	"database/sql"
	"log"

	"github.com/philippecery/maths/webapp/database"
	"github.com/philippecery/maths/webapp/database/model"
)

const queryGetGradesByOwnerID = "SELECT g.GradeID, g.GradeName, g.GradeDescription FROM Grades g WHERE g.OwnerID = ?"

// GetGradesByOwnerID returns all the Grade documents in the Grades collection
func GetGradesByOwnerID(ownerID string) []*model.Grade {
	var grades []*model.Grade
	if rows, err := database.Connection.QueryContext(context.TODO(), queryGetGradesByOwnerID, ownerID); err == nil {
		defer rows.Close()
		for rows.Next() {
			grade := new(model.Grade)
			if err = rows.Scan(grade.GradeID, grade.Name, grade.Description); err == nil {
				grades = append(grades, grade)
			}
		}
	}
	return grades
}

const queryGetGradeByStudentID = `
SELECT
	g.GradeID, g.GradeName, g.GradeDescription
FROM Students s
	INNER JOIN Grades g ON s.GradeID = g.GradeID
WHERE s.UserID = ?
`

// GetGradeByStudentID returns the User model from the Users collection where userid field is the provided id and the role is Student
func GetGradeByStudentID(userID string) *model.Grade {
	grade := new(model.Grade)
	err := database.Connection.QueryRowContext(context.TODO(), queryGetGradeByStudentID, userID).Scan(grade.GradeID, grade.Name, grade.Description)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("database: student %s doesn't exist or does not have a grade assigned", userID)
	case err != nil:
		log.Printf("database: unable to find student %s. Cause: %v", userID, err)
		return nil
	}
	grade.Homeworks = GetHomeworksByGradeID(grade.GradeID)
	return grade
}

const queryGetGradesByTeamID = `
SELECT
	g.GradeID, g.GradeName, g.GradeDescription
FROM Grade g
	INNER JOIN Students s ON s.GradeID = g.GradeID
	INNER JOIN Users u ON s.UserID = u.UserID
WHERE u.TeamID = ?
`

// GetGradesByTeamID returns the User model from the Users collection where userid field is the provided id and the role is Student
func GetGradesByTeamID(teamID string) []*model.Grade {
	var grades []*model.Grade
	if rows, err := database.Connection.QueryContext(context.TODO(), queryGetGradesByTeamID, teamID); err == nil {
		defer rows.Close()
		for rows.Next() {
			grade := new(model.Grade)
			if err = rows.Scan(grade.GradeID, grade.Name, grade.Description); err == nil {
				grades = append(grades, grade)
			}
		}
	}
	return grades
}

const queryGetGradeByID = "SELECT g.GradeID, g.GradeName, g.GradeDescription FROM Grades g WHERE g.GradeID = ?"

// GetGradeByID returns the Grade model from the Grades collection where gradeid field is the provided id
func GetGradeByID(gradeID string) *model.Grade {
	grade := new(model.Grade)
	err := database.Connection.QueryRowContext(context.TODO(), queryGetGradeByID, gradeID).Scan(grade.GradeID, grade.Name, grade.Description)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("database: grade %s doesn't exist", gradeID)
	case err != nil:
		log.Printf("database: unable to find grade %s. Cause: %v", gradeID, err)
		return nil
	}
	grade.Homeworks = GetHomeworksByGradeID(gradeID)
	return grade
}

const queryCreateNewGrade = "INSERT INTO Grades (GradeID, OwnerID, GradeName, GradeDescription) VALUES (?, ?, ?, ?)"

// CreateNewGrade creates a new Grade model in the Grades collection
func CreateNewGrade(newGrade *model.Grade) error {
	var err error
	if newGrade != nil {
		_, err = database.Connection.ExecContext(context.TODO(), queryCreateNewGrade, newGrade.GradeID, newGrade.OwnerID, newGrade.Name, newGrade.Description)
	}
	return err
}

const queryUpdateGrade = "UPDATE Grades SET GradeName = ?, GradeDescription = ? WHERE GradeID = ?"

// UpdateGrade retrieves and replaces the Grade model where gradeid field equals the one in the new Grade model
func UpdateGrade(grade *model.Grade) error {
	var err error
	if grade != nil {
		_, err = database.Connection.ExecContext(context.TODO(), queryUpdateGrade, grade.Name, grade.Description, grade.GradeID)
	}
	return err
}

const queryHasAccessToGrade = `
SELECT
	COUNT(*)
FROM Students s
WHERE s.UserID = ? AND s.GradeID = ?
`

func HasAccessToGrade(userID, gradeID string) bool {
	var hasAccess int
	database.Connection.QueryRowContext(context.TODO(), queryHasAccessToGrade, userID, gradeID).Scan(hasAccess)
	return hasAccess == 1
}

const queryIsGradeOwner = `
SELECT
	COUNT(*)
FROM Grade g
WHERE g.OwnerID = ? AND g.GradeID = ?
`

func IsGradeOwner(userID, gradeID string) bool {
	var hasAccess int
	database.Connection.QueryRowContext(context.TODO(), queryIsGradeOwner, userID, gradeID).Scan(hasAccess)
	return hasAccess == 1
}
