package structs

// Party - structure for party of all player characters
type Party struct {
	Name        string
	Description string
	CurrentRoom string
	Members     []*Player
	Objectives  []*Objective
	Value       int
}
