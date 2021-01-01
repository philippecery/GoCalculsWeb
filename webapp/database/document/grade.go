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
