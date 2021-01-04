package document

// Grade document
type Grade struct {
	GradeID     string
	Name        string
	Description string
	MentalMath  *Homework
	ColumnForm  *Homework
}

// Homework sub-document
type Homework struct {
	NbAdditions       int
	NbSubstractions   int
	NbMultiplications int
	NbDivisions       int
	Time              int
}

// NumberOfOperations returns the total number of operations for this homework
func (h *Homework) NumberOfOperations() int {
	return h.NbAdditions + h.NbSubstractions + h.NbMultiplications + h.NbDivisions
}

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
