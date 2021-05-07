package model

import (
	"github.com/philippecery/maths/webapp/constant/operation"
)

// Operation represents a randomly generated operation to be done by the student, as well as the answer submitted.
type Operation struct {
	OperationID int
	OperatorID  int
	Operand1    int
	Operand2    int
	Answer      int
	Answer2     int
	Status      operation.Status
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
