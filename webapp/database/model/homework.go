package model

import (
	"github.com/philippecery/maths/webapp/constant/homework"
)

// Homework sub-model
type Homework struct {
	HomeworkID        string
	Name              string
	Type              homework.Type
	NbAdditions       int
	NbSubstractions   int
	NbMultiplications int
	NbDivisions       int
	Time              int
	Status            homework.Status
}

// NumberOfOperations returns the total number of operations for this homework
func (h *Homework) NumberOfOperations() int {
	return h.NbAdditions + h.NbSubstractions + h.NbMultiplications + h.NbDivisions
}

// NumberOfOperationsByOperator returns the number of operations in this homework for the sepcified operator
func (h *Homework) NumberOfOperationsByOperator(operatorID int) int {
	var nbOperations int
	switch operatorID {
	case 1:
		nbOperations = h.NbAdditions
	case 2:
		nbOperations = h.NbSubstractions
	case 3:
		nbOperations = h.NbMultiplications
	case 4:
		nbOperations = h.NbDivisions
	}
	return nbOperations
}

// ActionToken generates and returns a unique ID to pass as a query parameter for CSRF protection.
func (h *Homework) ActionToken() string {
	return actionToken(h.HomeworkID)
}

// VerifyHomeworkActionToken verifies the provided action token is valid for the provided grade.
func VerifyHomeworkActionToken(actionToken string, homeworkID string) bool {
	return verifyActionToken(actionToken, homeworkID)
}
