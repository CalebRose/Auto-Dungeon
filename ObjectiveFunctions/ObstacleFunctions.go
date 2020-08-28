package objectivefunctions

import (
	// "firebase.google.com/go/auth"
	"fmt"
	"math/rand"

	"github.com/calebrose/Auto-Dungeon/structs"
)

// CheckObstacles - Check for any obstacles blocking the party from leaving the room
func CheckObstacles(party structs.Party, room *structs.Room) (structs.Party, *structs.Room) {
	//
	roll := rand.Intn(20) + 1
	fmt.Println("TEST")
	bonus := 0
	obstaclesPassed := 0

	for _, obstacle := range room.Obstacles {
		if obstacle.ObstaclePassed {
			obstaclesPassed++
			continue
		}
		// Loop through party
		for _, player := range party.Members {
			if player.Condition == "Dead" {
				continue
			}
			if obstacle.ObstacleType == "Door" {
				// Locked door -- lockpick it
				bonus = BonusAllocation(player.Proficiencies.Stealth.Level)
			} else if obstacle.ObstacleType == "Rubble" || obstacle.ObstacleType == "GiantMachine" || obstacle.ObstacleType == "Boulder" {
				// Physique check
				bonus = BonusAllocation(player.Proficiencies.Fisticuffs.Level)
			} else if obstacle.ObstacleType == "Inactive Machine" {
				bonus = BonusAllocation(player.Proficiencies.Engineering.Level)
			}

			if roll+bonus > obstacle.ObstacleRequirement || roll == 20 {
				// Obstacle passed
				obstacle.ObstaclePassed = true
				obstaclesPassed++
				player.Stats.ObstaclesOvercome++
				break
			}
		}
		if obstaclesPassed == len(room.Obstacles) {
			// Break Locks on Room
			room.Locked = false
			break
		}
	}

	return party, room
}
