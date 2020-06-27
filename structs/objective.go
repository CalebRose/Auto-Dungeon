package structs

// Objective - Goals for players to accomplish
type Objective struct {
	Name           string
	ObjectiveType  string
	Description    string
	Fulfilled      bool
	Condition      string
	TargetPerson   string
	TargetLocation string
	TargetItem     string
}
