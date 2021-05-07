package model

import (
	"fmt"
	"time"

	"github.com/philippecery/maths/webapp/constant/homework"
	"github.com/philippecery/maths/webapp/i18n"
	"github.com/philippecery/maths/webapp/util"
)

// HomeworkSession represents a homwework session.
// Contains the homework assigned by the teacher, the operations generated, the answers submitted, and the results per operator.
type HomeworkSession struct {
	SessionID       string
	UserID          string
	StartTime       time.Time
	EndTime         time.Time
	Homework        *Homework
	Operations      []*Operation
	Additions       *Results
	Substractions   *Results
	Multiplications *Results
	Divisions       *Results
	Status          homework.SessionStatus
}

// Results contains the number of good and wrong answers submitted per operator during a session
type Results struct {
	NbGood  int
	NbWrong int
}

// NewHomeworkSession returns a new initialized homework session.
func NewHomeworkSession(userID string, homework Homework) *HomeworkSession {
	return &HomeworkSession{SessionID: util.GenerateUUID(), UserID: userID, StartTime: time.Now(), Homework: &homework, Operations: make([]*Operation, 0), Additions: &Results{}, Substractions: &Results{}, Multiplications: &Results{}, Divisions: &Results{}}
}

// GetCurrentOperation returns the latest operation added to this homework session.
func (s *HomeworkSession) GetCurrentOperation() *Operation {
	if len(s.Operations) > 0 {
		return s.Operations[len(s.Operations)-1]
	}
	return nil
}

// OperatorIDs returns the identifiers of operators where there are still operations remaining in this homework
func (s *HomeworkSession) OperatorIDs() []int {
	operationIDs := make([]int, 0)
	if s.Additions.NbGood < s.Homework.NbAdditions {
		operationIDs = append(operationIDs, 1)
	}
	if s.Substractions.NbGood < s.Homework.NbSubstractions {
		operationIDs = append(operationIDs, 2)
	}
	if s.Multiplications.NbGood < s.Homework.NbMultiplications {
		operationIDs = append(operationIDs, 3)
	}
	if s.Divisions.NbGood < s.Homework.NbDivisions {
		operationIDs = append(operationIDs, 4)
	}
	return operationIDs
}

// NbTotalGood returns the total number of good answers.
func (s *HomeworkSession) NbTotalGood() int {
	return s.Additions.NbGood + s.Substractions.NbGood + s.Multiplications.NbGood + s.Divisions.NbGood
}

// NbUpdate returns the number of answers, good or wrong, for the specified operator.
func (s *HomeworkSession) NbUpdate(isGood bool, operatorID int) int {
	var nbUpdate int
	var results *Results
	switch operatorID {
	case 1:
		results = s.Additions
	case 2:
		results = s.Substractions
	case 3:
		results = s.Multiplications
	case 4:
		results = s.Divisions
	}
	if isGood {
		results.NbGood++
		nbUpdate = results.NbGood
	} else {
		results.NbWrong++
		nbUpdate = results.NbWrong
	}
	return nbUpdate
}

const dateFormat = "Monday 02 January 2006 @ 15:04:05 GMT"

// FormattedDateTime returns the formatted and localized start datetime
func (s *HomeworkSession) FormattedDateTime(locale string) string {
	return i18n.FormatDateTime(s.StartTime, locale)
}

// FormattedDuration returns the formatted duration in minutes, seconds and milliseconds
func (s *HomeworkSession) FormattedDuration() string {
	var minutes, seconds, milliseconds time.Duration
	if !s.EndTime.IsZero() {
		duration := s.EndTime.Sub(s.StartTime).Round(time.Millisecond)
		minutes = duration / time.Minute
		duration -= minutes * time.Minute
		seconds = duration / time.Second
		duration -= seconds * time.Second
		milliseconds = duration / time.Millisecond
	}
	return fmt.Sprintf("%02d:%02d,%03d", minutes, seconds, milliseconds)
}
