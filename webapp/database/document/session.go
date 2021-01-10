package document

import (
	"time"

	"github.com/philippecery/maths/webapp/constant"
)

// HomeworkSession represents a homwework session.
// Contains the homework assigned by the teacher, the operations generated, the answers submitted, and the results per operator.
type HomeworkSession struct {
	UserID          string
	StartTime       time.Time
	EndTime         time.Time
	TypeID          int
	Homework        *Homework
	Operations      []*Operation
	Additions       *Results
	Substractions   *Results
	Multiplications *Results
	Divisions       *Results
}

// Results contains the number of good and wrong answers submitted per operator during a session
type Results struct {
	NbGood  int
	NbWrong int
}

// Operation represents a randomly generated operation to be done by the student, as well as the answer submitted.
type Operation struct {
	OperatorID int
	Operand1   int
	Operand2   int
	Status     constant.Status
	Answer     int
	Answer2    int
}

// VerifyResult returns true if answer equals expected result
func (o *Operation) VerifyResult(answer, answer2 int) bool {
	o.Answer = answer
	o.Answer2 = answer2
	var result, result2 int
	switch o.OperatorID {
	case 1:
		result = o.Operand1 + o.Operand2
	case 2:
		result = o.Operand1 - o.Operand2
	case 3:
		result = o.Operand1 * o.Operand2
	case 4:
		result = o.Operand1 / o.Operand2
		result2 = o.Operand1 % o.Operand2
	}
	return answer == result && answer2 == result2
}

// GetResult returns the expected result for this operation, only if an answer was already submitted.
func (o *Operation) GetResult() (int, int) {
	if o.Answer > 0 {
		switch o.OperatorID {
		case 1:
			return o.Operand1 + o.Operand2, 0
		case 2:
			return o.Operand1 - o.Operand2, 0
		case 3:
			return o.Operand1 * o.Operand2, 0
		case 4:
			return o.Operand1 / o.Operand2, o.Operand1 % o.Operand2
		}
	}
	return 0, 0
}

// GetAnswer returns the submitted answer for this operation.
func (o *Operation) GetAnswer() (int, int) {
	return o.Answer, o.Answer2
}

// NewHomeworkSession returns a new initialized homework session.
func NewHomeworkSession(userID string, typeID int, homework Homework) *HomeworkSession {
	return &HomeworkSession{UserID: userID, StartTime: time.Now(), TypeID: typeID, Homework: &homework, Operations: make([]*Operation, 0), Additions: &Results{}, Substractions: &Results{}, Multiplications: &Results{}, Divisions: &Results{}}
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
