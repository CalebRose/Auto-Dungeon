package structs

import (
	"fmt"
	"math/rand"
)

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
	initStealth := 0
	if 20%playerCount == 0 {
		// 2, 4, 5, 10, 20 -- set base to 10
		initCuriosity = 10
		initDiscovery = 10
		initPersuasion = 10
		initStealth = 10
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
		if player.Proficiencies.Stealth.Level > 5 {
			if initStealth+2 > 20 {
				initStealth = 20
			} else {
				initStealth += 2
			}
		} else if player.Proficiencies.Stealth.Level > 1 {
			initStealth += 0
		} else {
			initStealth -= 2
		}
	}
	p.PartyBehavior.PartyCuriosity = initCuriosity
	p.PartyBehavior.PartyDiscovery = initDiscovery
	p.PartyBehavior.PartyPersuasion = initPersuasion
	p.PartyBehavior.PartyStealth = initStealth
}

// KeepStealth - The Party is seen by enemies, run a check if the party can persuade the group on their presence
func (p *Party) KeepStealth(enemies []*Enemy) {
	count := 0 // Count for how many enemies believe the group
	roll20 := false
	for _, player := range p.Members {
		persuasionRoll := rand.Intn(20) + 1
		persuasion := player.Proficiencies.Persuasion
		bonus := 0
		if persuasion.Level > 15 {
			bonus = 4
		} else if persuasion.Level > 12 {
			bonus = 3
		} else if persuasion.Level > 9 {
			bonus = 2
		} else if persuasion.Level > 6 {
			bonus++
		} else if persuasion.Level < 3 {
			bonus--
		} else if persuasion.Level <= 1 {
			bonus -= 2
		}
		if persuasionRoll == 1 {
			roll20 = true
			break
		} else if persuasionRoll+bonus > 14 {
			count++
			player.Stats.PersuasionChecks++
		}
	}
	if roll20 || float64(count) >= (float64(len(enemies))*0.66) {
		fmt.Println("The Party passed all persuasion checks. Stealth was retained.")
		p.IsStealth = true
	} else {
		fmt.Println("The Party could not talk their way out.")
		p.IsStealth = false
	}
}
