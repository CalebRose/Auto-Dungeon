package objectivefunctions

import (
	// "firebase.google.com/go/auth"
	"fmt"
	"math/rand"
	"github.com/calebrose/Auto-Dungeon/structs"
)

// BreakLoop - Breaks Loop to end the game
func BreakLoop(party structs.Party, failure bool) bool {
	// Future implementation -- Getaway vehicles?
	obj := party.Objectives[0]
	if party.MissionStatus && obj.TargetLocation == party.CurrentRoom {
		return true
	}
	if obj.ObjectiveType == "Escort" || obj.ObjectiveType == "Kidnapping" || obj.ObjectiveType == "Rescue" || obj.ObjectiveType == "Location" {
		if failure && obj.TargetLocation == party.CurrentRoom {
			return true
		}
	}
	return false
}

// CheckAllObjectives - Confirm that all objectives are complete
func CheckAllObjectives(Objectives []*structs.Objective) bool {
	totalObj := len(Objectives)
	completedObj := 0
	for i := 0; i < totalObj; i++ {
		if Objectives[i].Fulfilled == true {
			completedObj++
		}
	}
	if completedObj == totalObj {
		return true
	}
	return false
}

// CheckObjectiveCompletion - Check if an objective is complete
func CheckObjectiveCompletion(party structs.Party, room *structs.Room) (structs.Party, bool) {
	// Current Logic implemented for One Objective Only, but structured for multiple just in case

	obj := party.Objectives[0]
	if obj.ObjectiveType == "Eliminate" {
		// Eliminate an Enemy Unit
		// Loop through dead enemy units
		enemies := room.Enemies
		for i := 0; i < len(enemies); i++ {
			// If the enemy in the objective was in the room & dead
			if obj.TargetPerson == enemies[i].Name && enemies[i].Condition == "Dead" {
				for _, player := range party.Members {
					if player.Condition == "Dead" {
						continue
					}
					player.StatAllocation(obj.ObjectiveType, 0)
				}
				obj.Fulfilled = true
				break
			}
		}
	} else if obj.ObjectiveType == "Location" {
		if obj.TargetLocation == room.Name {
			obj.Fulfilled = true
		}
	} else if obj.ObjectiveType == "Item" {
		enemies := room.Enemies
		for i := 0; i < len(enemies); i++ {
			// If the enemy in the objective was in the room & dead
			if obj.TargetItem == enemies[i].Loot.Name && enemies[i].Condition == "Dead" {
				obj.Fulfilled = true
				break
			}
		}
	} else if obj.ObjectiveType == "Escort" {
		// Check if Target is alive && party has reached a location.
		for _, i := range party.Targets {
			if i.CurrentHealth == 0 {
				// Mission Failure
				obj.Fulfilled = false
				party.InRetreat = true
				return party, true
			}
		}
		if obj.TargetLocation == room.Name {
			obj.Fulfilled = true
		}
	} else if obj.ObjectiveType == "Espionage" {
		// Party must remain in stealth
		if obj.TargetLocation != room.Name || party.IsStealth == false {
			// Mission Failure -- it might make sense to have a Starting Point
			party.InRetreat = true
			return party, true
		}
		obj.Fulfilled = true

	} else if obj.ObjectiveType == "Kidnapping" || obj.ObjectiveType == "Rescue" {
		// Party must acquire a specific target and reach a room (starting room, escape room, vehicle)
		hasTarget := false
		if len(party.Targets) == 0 {
			return party, false
		}
		for _, target := range party.Targets {
			if target.CurrentHealth == 0 {
				// Mission Failure
				party.InRetreat = true
				return party, true
			}
			if target.TargetName == obj.TargetPerson {
				hasTarget = true
				party.InRetreat = true
				break
			}
		}
		if hasTarget && obj.TargetLocation != room.Name {
			return party, false
		}
		obj.Fulfilled = true

	} else if obj.ObjectiveType == "Theft" {
		// Party must have a specific item
		for _, target := range party.Targets {
			if target.TargetName == obj.TargetItem {
				obj.Fulfilled = true
				party.InRetreat = true
				break
			}
		}
	} else if obj.ObjectiveType == "Informant" {

	}
	return party, false
}

// TargetCheck -- mission objective variants regarding objectives
func TargetCheck(party structs.Party, room *structs.Room) (structs.Party, *structs.Room) {
	//
	obj := party.Objectives[0]
	enemiesAlive := EnemiesAlive(room.Enemies)
	bonus := 0
	attempts := 0

	// Loop through room's targets

	for _, target := range room.Targets {

		req := target.ObjRequirement
		if party.IsStealth && enemiesAlive {
			req += 2
		}

		// If Target is a vehicle
		// FOR FUTURE IMPLEMENTATION
		//

		for _, player := range party.Members {
			if player.Condition == "Dead" {
				continue
			}
			roll := rand.Intn(20) + 1
			if obj.ObjectiveType == "Theft" {

				if target.TargetType != "Item" {
					// Break the loop. This ain't the right target
					break
				}
				if enemiesAlive {
					bonus = BonusAllocation(player.Proficiencies.Pickpocket.Level)
				} else {
					target.ObjRequirement = 1
				}

				if roll+bonus > target.ObjRequirement || roll == 20 {
					fmt.Println("Player " + player.Name + " was able to pocket a " + target.TargetName)
					party.Targets = append(party.Targets, &target)
					target.Acquired = true
					if room.Locked {
						room.Locked = false
					}
				}

			} else if obj.ObjectiveType == "Kidnapping" {

				if target.TargetType != "Person" {
					// Break the loop. This ain't the right target
					break
				}

				if attempts < 3 || enemiesAlive {
					bonus = BonusAllocation(player.Proficiencies.Persuasion.Level)
				} else if player.Profession == "Conscript" {
					bonus = BonusAllocation(player.Proficiencies.LongRangeWeapons.Level)
				} else {
					bonus = BonusAllocation(player.Proficiencies.Fisticuffs.Level)
				}

				if roll+bonus > target.ObjRequirement || roll == 20 {
					if enemiesAlive {
						fmt.Println("Player " + player.Name + " has kidnapped " + target.TargetName + " right from under their enemies' noses!")
					} else {
						fmt.Println("Player " + player.Name + " has kidnapped their target! Time to retreat!")
					}
					party.Targets = append(party.Targets, &target)
					target.Acquired = true
					if room.Locked {
						room.Locked = false
					}
				}

			} else if obj.ObjectiveType == "Rescue" {

				bonus = BonusAllocation(player.Proficiencies.Persuasion.Level)
				if roll+bonus > target.ObjRequirement || roll == 20 {
					if enemiesAlive {
						fmt.Println("The party managed to rescue " + target.TargetName + " right from under their enemies' noses!")
					} else {
						fmt.Println("The party has rescued their target! Time to retreat!")
					}
					party.Targets = append(party.Targets, &target)
					target.Acquired = true
					if room.Locked {
						room.Locked = false
					}
				}

			}
			attempts++
		}
		if !target.Acquired && target.ObjRequired {
			room.Locked = true
		} else if target.Acquired && target.ObjRequired {
			for _, player := range party.Members {
				if player.Condition == "Dead" {
					continue
				}
				player.StatAllocation(obj.ObjectiveType, 0)
			}
		}
	}
	return party, room
}

// EnemiesAlive -- Detect Whether Enemies are still alive
func EnemiesAlive(enemies []*structs.Enemy) bool {
	for _, enemy := range enemies {
		if enemy.Condition != "Dead" {
			return true
		}
	}
	return false
}

// BonusAllocation -- Allocate Bonus based on proficiency
func BonusAllocation(skill int) int {
	if skill > 9 {
		return 3
	} else if skill > 6 {
		return 2
	} else if skill > 3 {
		return 1
	}
	return 0
}
