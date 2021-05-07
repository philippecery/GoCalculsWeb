package dataaccess

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/philippecery/maths/webapp/constant/homework"
	"github.com/philippecery/maths/webapp/constant/operation"
	"github.com/philippecery/maths/webapp/database"
	"github.com/philippecery/maths/webapp/database/model"
)

const queryNewHomeworkSession = "INSERT INTO HomeworkSession (SessionID, UserID, StartTime, HomeworkID, StatusID) VALUES (?, ?, ?, ?, ?)"

// NewHomeworkSession stores the current homework session
func NewHomeworkSession(newSession *model.HomeworkSession) error {
	var err error
	if newSession != nil {
		_, err = database.Connection.ExecContext(context.TODO(), queryNewHomeworkSession, newSession.SessionID, newSession.UserID, newSession.StartTime, newSession.Homework.HomeworkID, newSession.Status)
	}
	return err
}

const queryAddOperation = "INSERT INTO SessionOperations (SessionID, OperationID, OperatorID, Operand1, Operand2, Status) VALUES (?, ?, ?, ?, ?, ?)"

// AddOperation
func AddOperation(sessionID string, operation *model.Operation) error {
	var err error
	if operation != nil {
		_, err = database.Connection.ExecContext(context.TODO(), queryAddOperation, sessionID, operation.OperationID, operation.OperatorID, operation.Operand1, operation.Operand2, operation.Status)
	}
	return err
}

const querySaveAnswer = "UPDATE SessionOperations SET Answer = ?, Answer2 = ?, StatusID = ? WHERE SessionID = ? AND OperationID = ?"

// SaveAnswer
func SaveAnswer(sessionID string, operationID, answer, answer2 int, status operation.Status) error {
	_, err := database.Connection.ExecContext(context.TODO(), querySaveAnswer, answer, answer2, status, sessionID, operationID)
	return err
}

const queryEndSession = "UPDATE HomeworkSessions SET EndTime = ?, StatusID = ? WHERE SessionID = ? AND EndTime IS NULL"

// EndSession
func EndHomeworkSession(sessionID string, endTime time.Time, status homework.SessionStatus) error {
	_, err := database.Connection.ExecContext(context.TODO(), queryEndSession, endTime, status, sessionID)
	return err
}

const nbPerPage = 10

const queryGetSessionsByUserID = `
SELECT
	hs.SessionID, hs.StartTime, hs.EndTime, hs.StatusID,
	h.HomeworkID, h.HomeworkName, h.HomeworkTypeID,
	h.NbAdditions, h.NbSubstractions, h.NbMultiplications, h.NbDivisions, h.Time, h.StatusID
FROM HomeworkSessions hs
	INNER JOIN Homeworks h ON hs.HomeworkID = h.HomeworkID
WHERE hs.UserID = ?
LIMIT ?, ?
`

const queryCountSessionsByUserID = `
SELECT
	COUNT(*)
FROM HomeworkSessions hs
WHERE hs.UserID = ?
`

// GetSessionsByUserID returns the paginated homework sessions for the specified user, along with the total number of sessions
func GetSessionsByUserID(userID string, homeworkType, status, page int) ([]*model.HomeworkSession, int) {
	var homeworkSessions []*model.HomeworkSession
	var count int
	var err error
	var rows *sql.Rows
	if rows, err = database.Connection.QueryContext(context.TODO(), queryGetSessionsByUserID, userID, page*nbPerPage, nbPerPage); err == nil {
		defer rows.Close()
		for rows.Next() {
			homeworkSession := new(model.HomeworkSession)
			if err = rows.Scan(homeworkSession.SessionID, homeworkSession.StartTime, homeworkSession.EndTime, homeworkSession.Status, homeworkSession.Homework.HomeworkID, homeworkSession.Homework.Name, homeworkSession.Homework.Type, homeworkSession.Homework.NbAdditions, homeworkSession.Homework.NbSubstractions, homeworkSession.Homework.NbMultiplications, homeworkSession.Homework.NbDivisions, homeworkSession.Homework.Time, homeworkSession.Homework.Status); err != nil {
				break
			}
			homeworkSessions = append(homeworkSessions, homeworkSession)
		}
		if err == nil {
			err = database.Connection.QueryRowContext(context.TODO(), queryCountSessionsByUserID, userID).Scan(count)
			switch {
			case err == sql.ErrNoRows:
				log.Printf("database: user %s doesn't exist", userID)
			case err != nil:
				log.Printf("database: unable to find user %s. Cause: %v", userID, err)
			}
		}
	}
	if err != nil {
		return nil, 0
	}
	return homeworkSessions, count
}

const queryGetSessionByID = `
SELECT
	hs.SessionID, hs.UserID, hs.StartTime, hs.EndTime, hs.StatusID,
	h.HomeworkID, h.HomeworkName, h.HomeworkTypeID,
	h.NbAdditions, h.NbSubstractions, h.NbMultiplications, h.NbDivisions, h.Time, h.StatusID
FROM HomeworkSessions hs
	INNER JOIN Homeworks h ON hs.HomeworkID = h.HomeworkID
WHERE hs.SessionID = ?
`

// GetSessionByID returns the homework sessions for the specified user
func GetSessionByID(sessionID string) *model.HomeworkSession {
	homeworkSession := new(model.HomeworkSession)
	err := database.Connection.QueryRowContext(context.TODO(), queryGetSessionByID, sessionID).Scan(homeworkSession.SessionID, homeworkSession.UserID, homeworkSession.StartTime, homeworkSession.EndTime, homeworkSession.Status, homeworkSession.Homework.HomeworkID, homeworkSession.Homework.Name, homeworkSession.Homework.Type, homeworkSession.Homework.NbAdditions, homeworkSession.Homework.NbSubstractions, homeworkSession.Homework.NbMultiplications, homeworkSession.Homework.NbDivisions, homeworkSession.Homework.Time, homeworkSession.Homework.Status)
	switch {
	case err == sql.ErrNoRows:
		log.Printf("database: session %s doesn't exist", sessionID)
	case err != nil:
		log.Printf("database: unable to find session %s. Cause: %v", sessionID, err)
	}
	if err != nil {
		return nil
	}
	homeworkSession.Operations = GetOperationsBySessionID(sessionID)
	return homeworkSession
}

const queryGetOperationsBySessionID = "SELECT OperationID, OperatorID, Operand1, Operand2, Answer, Answer2, StatusID FROM SessionOperations WHERE SessionID = ?"

// GetOperationsBySessionID
func GetOperationsBySessionID(sessionID string) []*model.Operation {
	var operations []*model.Operation
	if rows, err := database.Connection.QueryContext(context.TODO(), queryGetOperationsBySessionID, sessionID); err == nil {
		defer rows.Close()
		for rows.Next() {
			operation := new(model.Operation)
			if err = rows.Scan(operation.OperationID, operation.OperatorID, operation.Operand1, operation.Operand2, operation.Answer, operation.Answer2, operation.Status); err == nil {
				operations = append(operations, operation)
			}
		}
	}
	return operations
}
