package api

import (
	"fmt"
	"log"

	"github.com/philippecery/maths/webapp/constant"
	"github.com/philippecery/maths/webapp/database/document"
	"github.com/philippecery/maths/webapp/util"
)

func (r *request) operation() map[string]interface{} {
	if session := r.getHomeworkSession(); session != nil {
		if homeworkType, exists := constant.HomeworkTypes[session.TypeID]; exists {
			operatorIDs := session.OperatorIDs()
			if len(operatorIDs) > 0 {
				fmt.Printf("Remaining operators: %v\n", operatorIDs)
				rndIdx, _ := util.GetNumber(len(operatorIDs))
				nextOperation := &document.Operation{OperatorID: operatorIDs[rndIdx], Status: constant.Todo}
				var operandRange *constant.OperandRanges
				switch nextOperation.OperatorID {
				case 1:
					operandRange = homeworkType.AdditionRange
				case 2:
					operandRange = homeworkType.SubstractionRange
				case 3:
					operandRange = homeworkType.MultiplicationRange
				case 4:
					operandRange = homeworkType.DivisionRange
				}
				nextOperation.Operand1, _ = util.GetNumberInRange(operandRange.Operand1.RangeMin, operandRange.Operand1.RangeMax)
				nextOperation.Operand2, _ = util.GetNumberInRange(operandRange.Operand2.RangeMin, operandRange.Operand2.RangeMax)
				r.addOperation(nextOperation)
				fmt.Printf("Nb operations in session = %d\n", len(r.getHomeworkSession().Operations))
				operator := constant.Operators[nextOperation.OperatorID]
				return map[string]interface{}{
					"response": "operation", "operationName": r.getLocalizedMessage(operator.I18N), "operand1": nextOperation.Operand1, "operand2": nextOperation.Operand2, "operator": operator.Symbol,
				}
			}
			log.Printf("/websocket[operation]: session completed")
		} else {
			log.Printf("/websocket[operation]: Invalid type ID")
		}
	} else {
		log.Printf("/websocket[operation]: No HomeworkSession found in session")
	}
	return map[string]interface{}{
		"response": "error", "message": r.getLocalizedMessage("errorGenericMessage"),
	}
}

func (r *request) answer() map[string]interface{} {
	if session := r.getHomeworkSession(); session != nil {
		operation := session.GetCurrentOperation()
		if result, result2 := operation.Result(); result > 0 {
			var answer, answer2 int
			answer, _ = r.toInt("answer")
			if session.TypeID == 4 {
				answer2, _ = r.toInt("answer2")
			}
			good := answer == result && answer2 == result2
			nbUpdate := session.NbUpdate(good, operation.OperatorID)
			percentUpdate := (nbUpdate * 100) / session.Homework.NumberOfOperationsByOperator(operation.OperatorID)
			var percentAll int
			if good {
				operation.Status = constant.Good
				percentAll = (session.NbTotalGood() * 100) / session.Homework.NumberOfOperations()
			} else {
				operation.Status = constant.Wrong
			}
			nbTotalRemaining := session.Homework.NumberOfOperations() - session.NbTotalGood()
			r.saveHomeworkSession(session)
			return map[string]interface{}{
				"response": "answer", "good": good, "nbUpdate": nbUpdate, "percentUpdate": percentUpdate, "percentAll": percentAll, "nbTotalRemaining": nbTotalRemaining,
			}
		}
	} else {
		log.Printf("/websocket[answer]: No HomeworkSession found in session")
	}
	return map[string]interface{}{
		"response": "error", "message": r.getLocalizedMessage("errorGenericMessage"),
	}
}

func (r *request) results() map[string]interface{} {
	if timeout, err := r.getBool("timeout"); err == nil {
		return map[string]interface{}{
			"response": "results", "timeout": timeout, "additions": map[string]interface{}{"nbTotal": 10, "nbGood": 10, "nbWrong": 0}, "substractions": map[string]interface{}{"nbTotal": 10, "nbGood": 10, "nbWrong": 0}, "multiplications": map[string]interface{}{"nbTotal": 10, "nbGood": 10, "nbWrong": 0}, "divisions": map[string]interface{}{"nbTotal": 10, "nbGood": 10, "nbWrong": 0},
		}
	}
	return map[string]interface{}{
		"response": "error", "message": r.getLocalizedMessage("errorGenericMessage"),
	}
}

func (r *request) toggle() map[string]interface{} {
	if show, err := r.getBool("show"); err == nil {
		if show {
			return map[string]interface{}{
				"type": "toggle", "result": 3, "result2": 0,
			}
		}
		return map[string]interface{}{
			"type": "toggle", "answer": 3, "answer2": 0,
		}
	}
	return map[string]interface{}{
		"response": "error", "message": r.getLocalizedMessage("errorGenericMessage"),
	}
}
