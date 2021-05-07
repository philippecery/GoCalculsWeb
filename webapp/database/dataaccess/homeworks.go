package dataaccess

import (
	"context"
	"database/sql"
	"log"

	"github.com/philippecery/maths/webapp/constant/homework"

	"github.com/philippecery/maths/webapp/database"
	"github.com/philippecery/maths/webapp/database/model"
)

const queryGetHomeworkByID = `
SELECT
	h.HomeworkID, h.HomeworkName, h.HomeworkTypeID, h.GradeID,
	h.NbAdditions, h.NbSubstractions, h.NbMultiplications, h.NbDivisions, h.Time, h.StatusID,
FROM Homeworks h
WHERE h.HomeworkID = ?
`

func GetHomeworkByID(homeworkID string) (*model.Homework, string) {
	homework := new(model.Homework)
	var gradeID string
	err := database.Connection.QueryRowContext(context.TODO(), queryGetHomeworkByID, homeworkID).Scan(homework.HomeworkID, homework.Name, homework.Type, gradeID, homework.NbAdditions, homework.NbSubstractions, homework.NbMultiplications, homework.NbDivisions, homework.Time, homework.Status)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("database: homework %s doesn't exist", homeworkID)
	case err != nil:
		log.Printf("database: unable to find homework %s. Cause: %v", homeworkID, err)
	}
	if err != nil {
		return nil, ""
	}
	return homework, gradeID
}

const queryGetHomeworkByIDAndStatus = `
SELECT
	h.HomeworkID, h.HomeworkName, h.HomeworkTypeID,
	h.NbAdditions, h.NbSubstractions, h.NbMultiplications, h.NbDivisions, h.Time, h.StatusID
FROM Homeworks h
WHERE h.HomeworkID = ? AND h.StatusID = ?
`

func GetHomeworkByIDAndStatus(homeworkID string, status homework.Status) *model.Homework {
	homework := new(model.Homework)
	err := database.Connection.QueryRowContext(context.TODO(), queryGetHomeworkByIDAndStatus, homeworkID, status).Scan(homework.HomeworkID, homework.Name, homework.Type, homework.NbAdditions, homework.NbSubstractions, homework.NbMultiplications, homework.NbDivisions, homework.Time, homework.Status)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("database: homework %s doesn't exist", homeworkID)
	case err != nil:
		log.Printf("database: unable to find homework %s. Cause: %v", homeworkID, err)
		return nil
	}
	return homework
}

const queryGetHomeworksByGradeID = `
SELECT
	h.HomeworkID, h.HomeworkName, h.HomeworkTypeID,
	h.NbAdditions, h.NbSubstractions, h.NbMultiplications, h.NbDivisions, h.Time, h.StatusID
FROM Homeworks h
WHERE h.GradeID = ?
`

func GetHomeworksByGradeID(gradeID string) []*model.Homework {
	var homeworks []*model.Homework
	if rows, err := database.Connection.QueryContext(context.TODO(), queryGetHomeworksByGradeID, gradeID); err == nil {
		defer rows.Close()
		for rows.Next() {
			homework := new(model.Homework)
			if err = rows.Scan(homework.HomeworkID, homework.Name, homework.Type, homework.NbAdditions, homework.NbSubstractions, homework.NbMultiplications, homework.NbDivisions, homework.Time, homework.Status); err == nil {
				homeworks = append(homeworks, homework)
			}
		}
	}
	return homeworks
}

const queryAddHomework = `
INSERT INTO Homeworks (
	GradeID, HomeworkID, HomeworkName, HomeworkTypeID, StatusID,
	NbAdditions, NbSubstractions, NbMultiplications, NbDivisions, HomeworkTime)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`

// AddHomework creates a new Grade model in the Grades collection
func AddHomework(gradeID string, homework *model.Homework) error {
	var err error
	if homework != nil {
		_, err = database.Connection.ExecContext(context.TODO(), queryAddHomework, gradeID, homework.HomeworkID, homework.Name, homework.Type, homework.Status, homework.NbAdditions, homework.NbSubstractions, homework.NbMultiplications, homework.NbDivisions, homework.Time)
	}
	return err
}

const queryUpdateHomeworkStatus = "UPDATE Homeworks SET StatusID = ? WHERE HomeworkID = ? AND StatusID = ?"

func PublishHomework(homeworkID string) error {
	_, err := database.Connection.ExecContext(context.TODO(), queryUpdateHomeworkStatus, homework.Online, homeworkID, homework.Draft)
	return err
}

func ArchiveHomework(homeworkID string) error {
	_, err := database.Connection.ExecContext(context.TODO(), queryUpdateHomeworkStatus, homework.Archived, homeworkID, homework.Online)
	return err
}

const queryUpdateHomework = "UPDATE Homeworks SET HomeworkName = ? WHERE HomeworkID = ?"

// UpdateHomework
func UpdateHomework(homework *model.Homework) error {
	var err error
	if homework != nil {
		_, err = database.Connection.ExecContext(context.TODO(), queryUpdateHomework, homework.Name, homework.HomeworkID)
	}
	return err
}

const (
	queryGetHomeworkInUse    = "SELECT COUNT(*) FROM HomeworkSessions WHERE HomeworkID = ? AND StatusID <> ?" // StatusID <> 0
	queryDeleteHomeworkInUse = "DELETE FROM HomeworkSessions WHERE HomeworkID = ?"
	queryDeleteHomework      = "DELETE FROM Homeworks WHERE HomeworkID = ?"
)

// DeleteHomework deletes the Grade model from the Grades collection where the gradeid field is the provided id
func DeleteHomework(homeworkID string) error {
	var err error
	var tx *sql.Tx
	ctx := context.TODO()
	if tx, err = database.Connection.BeginTx(ctx, nil); err == nil {
		var usedInSessions uint
		err = tx.QueryRowContext(ctx, queryGetHomeworkInUse, homeworkID, homework.Draft).Scan(usedInSessions)
		switch {
		case err == sql.ErrNoRows:
			log.Printf("database: homework %s doesn't exist", homeworkID)
		case err != nil:
			log.Printf("database: unable to find homework %s. Cause: %v", homeworkID, err)
		}
		if err != nil {
			tx.Rollback()
			return err
		}
		if usedInSessions == 0 {
			if _, err = tx.ExecContext(ctx, queryDeleteHomeworkInUse, homeworkID); err == nil {
				_, err = tx.ExecContext(ctx, queryDeleteHomework, homeworkID)
			}
		} else {
			_, err = tx.ExecContext(ctx, queryUpdateHomeworkStatus, homework.Archived, homeworkID, homework.Online)
		}
		if err == nil {
			err = tx.Commit()
		} else {
			tx.Rollback()
		}
	}
	return err
}

const queryHasAccessToHomework = `
SELECT
	COUNT(*)
FROM Homeworks h
	INNER JOIN Grades g ON h.GradeID = g.GradeID
	INNER JOIN Students s ON g.GradeID = s.GradeID
WHERE s.UserID = ? AND h.HomeworkID = ?
`

func HasAccessToHomework(userID, homeworkID string) bool {
	var hasAccess int
	database.Connection.QueryRowContext(context.TODO(), queryHasAccessToHomework, userID, homeworkID).Scan(hasAccess)
	return hasAccess == 1
}
