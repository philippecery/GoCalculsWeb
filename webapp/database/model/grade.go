package model

// Grade model
type Grade struct {
	GradeID     string
	OwnerID     string
	Name        string
	Description string
	Homeworks   []*Homework
}

// ActionToken generates and returns a unique ID to pass as a query parameter for CSRF protection.
func (g *Grade) ActionToken() string {
	return actionToken(g.GradeID)
}

// VerifyGradeActionToken verifies the provided action token is valid for the provided grade.
func VerifyGradeActionToken(actionToken string, gradeID string) bool {
	return verifyActionToken(actionToken, gradeID)
}
