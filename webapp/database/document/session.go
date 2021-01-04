package document

import "github.com/philippecery/maths/webapp/constant"

type HomeworkSession struct {
	TypeID          int
	Homework        *Homework
	Operations      []*Operation
	Additions       *Results
	Substractions   *Results
	Multiplications *Results
	Divisions       *Results
}

type Results struct {
	NbGood  int
	NbWrong int
}

type Operation struct {
	OperatorID int
	Operand1   int
	Operand2   int
	Status     constant.Status
}

func (o *Operation) Result() (int, int) {
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
	return 0, 0
}

func NewHomeworkSession(typeID int, homework Homework) *HomeworkSession {
	return &HomeworkSession{TypeID: typeID, Homework: &homework, Operations: make([]*Operation, 0), Additions: &Results{}, Substractions: &Results{}, Multiplications: &Results{}, Divisions: &Results{}}
}

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

func (s *HomeworkSession) NbTotalGood() int {
	return s.Additions.NbGood + s.Substractions.NbGood + s.Multiplications.NbGood + s.Divisions.NbGood
}

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
