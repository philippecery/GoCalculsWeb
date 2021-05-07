package homework

// Type type
type Type uint

const (
	MentalMath Type = iota + 1
	ColumnForm
)

// TypeDetails type
type TypeDetails struct {
	I18N                string
	Logo                string
	AdditionRange       *OperandRanges
	SubstractionRange   *OperandRanges
	MultiplicationRange *OperandRanges
	DivisionRange       *OperandRanges
}

// OperandRanges type
type OperandRanges struct {
	OperatorID int
	Operand1   *OperandRange
	Operand2   *OperandRange
}

// OperandRange type
type OperandRange struct {
	RangeMin int
	RangeMax int
	DecMax   int
}

func (t Type) Logo() string {
	return Types[t].Logo
}

func (t Type) I18N() string {
	return Types[t].I18N
}

// Types is the default configuration of operations.
// Todo: make it configurable
var Types = map[Type]*TypeDetails{
	MentalMath: &TypeDetails{
		I18N: "mentalmath",
		Logo: "hourglass",
		AdditionRange: &OperandRanges{
			OperatorID: 1,
			Operand1: &OperandRange{
				RangeMin: 1,
				RangeMax: 100,
				DecMax:   0,
			},
			Operand2: &OperandRange{
				RangeMin: 1,
				RangeMax: 10,
				DecMax:   0,
			},
		},
		SubstractionRange: &OperandRanges{
			OperatorID: 2,
			Operand1: &OperandRange{
				RangeMin: 10,
				RangeMax: 100,
				DecMax:   0,
			},
			Operand2: &OperandRange{
				RangeMin: 1,
				RangeMax: 10,
				DecMax:   0,
			},
		},
		MultiplicationRange: &OperandRanges{
			OperatorID: 3,
			Operand1: &OperandRange{
				RangeMin: 2,
				RangeMax: 10,
				DecMax:   0,
			},
			Operand2: &OperandRange{
				RangeMin: 2,
				RangeMax: 10,
				DecMax:   0,
			},
		},
		DivisionRange: &OperandRanges{
			OperatorID: 4,
			Operand1: &OperandRange{
				RangeMin: 10,
				RangeMax: 100,
				DecMax:   0,
			},
			Operand2: &OperandRange{
				RangeMin: 2,
				RangeMax: 10,
				DecMax:   0,
			},
		},
	},
	ColumnForm: &TypeDetails{
		I18N: "columnform",
		Logo: "pencil",
		AdditionRange: &OperandRanges{
			OperatorID: 1,
			Operand1: &OperandRange{
				RangeMin: 100,
				RangeMax: 1000000,
				DecMax:   2,
			},
			Operand2: &OperandRange{
				RangeMin: 100,
				RangeMax: 100000,
				DecMax:   2,
			},
		},
		SubstractionRange: &OperandRanges{
			OperatorID: 2,
			Operand1: &OperandRange{
				RangeMin: 100000,
				RangeMax: 1000000,
				DecMax:   2,
			},
			Operand2: &OperandRange{
				RangeMin: 1000,
				RangeMax: 100000,
				DecMax:   2,
			},
		},
		MultiplicationRange: &OperandRanges{
			OperatorID: 3,
			Operand1: &OperandRange{
				RangeMin: 100,
				RangeMax: 100000,
				DecMax:   2,
			},
			Operand2: &OperandRange{
				RangeMin: 100,
				RangeMax: 100000,
				DecMax:   2,
			},
		},
		DivisionRange: &OperandRanges{
			OperatorID: 4,
			Operand1: &OperandRange{
				RangeMin: 1000,
				RangeMax: 100000,
				DecMax:   0,
			},
			Operand2: &OperandRange{
				RangeMin: 10,
				RangeMax: 1000,
				DecMax:   0,
			},
		},
	},
}
