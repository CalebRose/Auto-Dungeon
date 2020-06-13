package structs

// Party - structure for party of all player characters
type Party struct {
	Name          string
	Description   string
	CurrentRoom   string
	Members       []*Player
	Objectives    []*Objective
	Targets       []*Target
	IsStealth     bool
	InRetreat     bool // The party is retreating from the mission -- either due to failure or due to mission objective
	EscapeRoute   []string
	MissionStatus bool
	BreakLoop     bool
	PartyBehavior PartyBehavior
}

// SetBehaviors - Set the behaviors of the party
func (p *Party) SetBehaviors(players []*Player) {
	playerCount := len(players)
	initCuriosity := 0
	initDiscovery := 0
	initPersuasion := 0
	if 20%playerCount == 0 {
		// 2, 4, 5, 10, 20 -- set base to 10
		initCuriosity = 10
		initDiscovery = 10
		initPersuasion = 10
	}
	for _, player := range players {
		behavior := player.Behavior
		if behavior.Curiosity {
			if initCuriosity+2 > 20 {
				initCuriosity = 20
			} else {
				initCuriosity += 2
			}
		} else {
			if initCuriosity-2 < 0 {
				initCuriosity = 0
			} else {
				initCuriosity -= 2
			}
		}
		if behavior.Discovery {
			if initDiscovery+2 > 20 {
				initDiscovery = 20
			} else {
				initDiscovery += 2
			}
		} else {
			if initDiscovery-2 < 0 {
				initDiscovery = 0
			} else {
				initDiscovery -= 2
			}
		}
		if behavior.Persuasion {
			if initPersuasion+2 > 20 {
				initPersuasion = 20
			} else {
				initPersuasion += 2
			}
		} else {
			if initPersuasion-2 < 0 {
				initPersuasion = 0
			} else {
				initPersuasion -= 2
			}
		}
	}
	p.PartyBehavior.PartyCuriosity = initCuriosity
	p.PartyBehavior.PartyDiscovery = initDiscovery
	p.PartyBehavior.PartyPersuasion = initPersuasion
}
