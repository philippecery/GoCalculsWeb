package homework

// Operator type
type Operator struct {
	I18N   string
	Symbol string
}

// Operators list
var Operators = map[int]*Operator{
	1: &Operator{I18N: "addition", Symbol: "+"},
	2: &Operator{I18N: "substraction", Symbol: "-"},
	3: &Operator{I18N: "multiplication", Symbol: "*"},
	4: &Operator{I18N: "division", Symbol: "/"},
}
