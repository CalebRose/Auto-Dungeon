package structs

// Objective - Goals for players to accomplish
type Objective struct {
	Name          string
	ObjectiveType string
	Description   string
	Fulfilled     bool
	Condition     string
	Target        string
}
