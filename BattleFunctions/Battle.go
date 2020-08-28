package battlefunctions

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/calebrose/Auto-Dungeon/structs"
)

// InitializeBattle - Initialize Battle between Players and Enemies in a room
func InitializeBattle(party structs.Party, enemies []*structs.Enemy, room *structs.Room) []*structs.Player {
	//Take in slice of players & slice of enemies
	// Use Battle Queue Struct instead of map
	// create a slice to take in a struct with a name, type, and initiative. Generate initiative based on player's perception
	players := party.Members
	filteredPlayers := FilterPlayers(players)
	filteredEnemies := enemies
	battleQueue := []structs.BattleQueue{}
	stealthBonus := 0
	if party.IsStealth {
		stealthBonus += len(party.Members)
	}
	for i := 0; i < len(filteredPlayers); i++ {

		playerQueue := structs.BattleQueue{
			Name:          players[i].Name,
			CombatantType: "Player",
			EnemyType:     "",
			Initiative:    rand.Intn(players[i].Attributes.Perception) + 1 + stealthBonus,
		}
		players[i].StatAllocation("BattleEngaged", 0)
		battleQueue = append(battleQueue, playerQueue)

	}

	for i := 0; i < len(enemies); i++ {
		enemyQueue := structs.BattleQueue{
			Name:          enemies[i].Name,
			CombatantType: "Enemy",
			EnemyType:     enemies[i].EnemyType,
			Initiative:    enemies[i].Initiative,
		}
		battleQueue = append(battleQueue, enemyQueue)
	}
	// Sort the slice from highest initiative to lowest
	// FOR LATER -- do a feat check for players w/ sleight of hand
	battleSlice := battleQueue[:]
	sort.Slice(battleSlice, func(i, j int) bool {
		return battleSlice[i].Initiative > battleSlice[j].Initiative
	})

	playerCount := len(filteredPlayers)
	enemyCount := len(enemies)
	battleIterator := 0
	playerCoverCount := 0
	enemyCoverCount := 0

	// while all players are not dead || all enemies are not dead
	for playerCount > 0 && enemyCount > 0 {
		// pop from the beginning of the queue
		// battleNode, battleSlice := battleSlice[0], battleSlice[1:]
		battleNode := battleSlice[battleIterator]

		// Check whether the map received is a player or enemy
		// If player, randomly select an enemy from the enemies array. If enemy, randomly select a player from the players array

		if battleNode.CombatantType == "Player" {
			// Player Turn
			player := FindPlayer(filteredPlayers, battleNode)
			useMedicalItem := false
			if player == nil || player.Condition == "Dead" {
				battleIterator = Iterate(battleIterator, battleSlice)
				continue
			}
			// In Cover
			if player.InCover == false && playerCoverCount <= room.PlayerCover {
				player.InCover = true
				playerCoverCount++
			}

			if player.Behavior.MedicalUse {
				medRoll := rand.Intn(3)
				if medRoll > 1 {
					useMedicalItem = true
				}
			}
			if useMedicalItem {
				player.UseMedicalItem()
			} else {

				chosenEnemy := rand.Intn(enemyCount)
				player, filteredEnemies[chosenEnemy] = SingularBattle(player, filteredEnemies[chosenEnemy], true)
				if enemies[chosenEnemy].Condition == "Dead" {
					enemyCount--
					if enemies[chosenEnemy].EnemyType == "Boss" {
						player.StatAllocation("BossDefeated", 0)
					}
					if enemyCount > 0 {
						filteredEnemies = FilterEnemies(filteredEnemies)
					}
					player.StatAllocation("EnemyKilled", 0)
				}
			}
		} else if battleNode.CombatantType == "Enemy" {
			// Enemy Turn
			enemy := FindEnemy(enemies, battleNode)
			if enemy == nil || enemy.Condition == "Dead" {
				battleIterator = Iterate(battleIterator, battleSlice)
				continue
			}
			// In Cover
			if enemy.InCover == false && enemyCoverCount <= room.EnemyCover {
				enemy.InCover = true
				enemyCoverCount++
			}
			chosenPlayer := rand.Intn(playerCount)
			SingularBattle(filteredPlayers[chosenPlayer], enemy, false)
			if filteredPlayers[chosenPlayer].Condition == "Dead" {
				playerCount--
				if playerCount > 0 {
					filteredPlayers = FilterPlayers(filteredPlayers)
				}
			}
		}
		// battleSlice = append(battleSlice, battleNode)
		battleIterator = Iterate(battleIterator, battleSlice)

		// Loop should break when either all enemies are defeated / mortally wounded, or all players are unable to fight
	}
	fmt.Println(players)
	return players
}

// Iterate - Function for iterating through the BattleQueue
func Iterate(i int, queue []structs.BattleQueue) int {
	i++
	if i >= len(queue) {
		i = 0
	}
	return i
}

// FilterEnemies - array for filtering enemies based on the Dead condition. Can be altered later to cover more scenarios
func FilterEnemies(enemies []*structs.Enemy) []*structs.Enemy {
	temp := []*structs.Enemy{}
	for _, i := range enemies {
		if i.Condition == "Dead" {
			temp = append(temp, i)
		}
	}
	fmt.Println(temp)
	return temp
}

// FilterPlayers - array for filtering players based on the Dead condition. Can be altered later to cover more scenarios
func FilterPlayers(players []*structs.Player) []*structs.Player {
	temp := []*structs.Player{}
	injuredPlayers := []*structs.Player{}
	for _, player := range players {
		if player.Condition != "Dead" {
			if (player.Condition == "Minorly Injured" && player.Behavior.EngagementThreshold <= 4) ||
				(player.Condition == "Majorly Injured" && player.Behavior.EngagementThreshold <= 3) ||
				(player.Condition == "Severely Injured" && player.Behavior.EngagementThreshold <= 2) ||
				(player.Behavior.EngagementThreshold == 1) {
				injuredPlayers = append(injuredPlayers, player)
			} else {
				temp = append(temp, player)
			}
		}
	}
	if len(temp) <= (len(injuredPlayers) / 2) {
		for i := 0; i < len(injuredPlayers)/2; i++ {
			temp = append(temp, injuredPlayers[i])
		}
	}
	return temp

}

// FindEnemy - find Enemy from the battle queue
func FindEnemy(enemies []*structs.Enemy, battleNode structs.BattleQueue) *structs.Enemy {
	// Find the enemy for singular battle
	i := 0
	for i < len(enemies) {
		if battleNode.Name == enemies[i].Name {
			if enemies[i].Condition == "Dead" {
				return nil
			}
			break
		}
		i++
	}
	return enemies[i]
}

// FindPlayer - Find a player from the battle queue
func FindPlayer(players []*structs.Player, battleNode structs.BattleQueue) *structs.Player {
	// Find the player for Singular Battle
	i := 0
	for i < len(players) {
		if battleNode.Name == players[i].Name {
			if players[i].Condition == "Dead" {
				return nil
			}
			break
		}
		i++
	}
	return players[i]
}

// SingularBattle : A turn of Battle between a player and an enemy.
func SingularBattle(player *structs.Player, enemy *structs.Enemy, playerTurn bool) (*structs.Player, *structs.Enemy) {
	// Establish Base Values & Variables
	coreStrength := 0
	if player.Weapon.WeaponType == "Rifle" ||
		player.Weapon.WeaponType == "Pistol" ||
		player.Weapon.WeaponType == "Shotgun" ||
		player.Weapon.WeaponType == "SniperRifle" {
		coreStrength = player.Attributes.Dexterity

	} else {
		coreStrength = player.Attributes.Strength
	}
	player.StatAllocation(player.Weapon.WeaponType, 0)
	playerBaseStrength := coreStrength + player.Weapon.WeaponRating
	enemyBaseStrength := enemy.CombatRating
	currentPlayerStrength := 0
	currentEnemyStrength := 0
	playerBonus := 0 // Used to detect if player is in cover
	enemyBonus := 0  // Used to detect if enemy is in cover

	difference := 0

	// If it's the player's turn:
	if playerTurn == true {
		// If the player has ammo in his gun
		currentWeaponCartridge := player.Weapon.CurrentCartridge

		if enemy.InCover == true {
			enemyBonus = 20
		}

		if currentWeaponCartridge > 0 {
			accuracy := player.Weapon.WeaponAccuracy
			fireRate := player.Weapon.FireRate
			if player.Behavior.BattleAggression {
				fireRate = int(math.Ceil(float64(fireRate) * 1.5))
				accuracy -= 15
			}
			// Fires Shot -- if the shot is less than the Player's Accuracy rating
			for i := 0; i < fireRate && currentWeaponCartridge > 0; i++ {
				shot := rand.Intn(100)
				player.Weapon.CurrentCartridge--
				player.StatAllocation("ShotFired", 0)
				if shot <= accuracy-enemyBonus {
					// Shot hits target
					player.StatAllocation("ShotMade", 0)
					currentPlayerStrength = currentPlayerStrength + rand.Intn(playerBaseStrength) + 1
				} else {
					fmt.Println("PLAYER MISSED!")
				}
			}

		} else {
			player.Weapon.CurrentReload++
			if player.Weapon.CurrentReload == player.Weapon.WeaponReloadTime {
				currentWeaponCartridge = player.Weapon.WeaponCartridge
				player.Weapon.CurrentReload = 0
			}

			fmt.Println("Reloading...")
		}
	} else {
		// Enemy's turn to fire
		if player.InCover == true {
			playerBonus = 33
		}

		enemyShot := rand.Intn(100)
		if enemyShot <= enemy.CombatAccuracy-playerBonus {
			currentEnemyStrength = rand.Intn(enemyBaseStrength)
		} else {
			fmt.Println("ENEMY MISSED!")
		}
		playerTurn = true
	}
	fmt.Println("Player Strength:", currentPlayerStrength)
	fmt.Println("Enemy Strength:", currentEnemyStrength)
	// If the player's strength is greater than 0
	if currentPlayerStrength > currentEnemyStrength {
		difference = currentPlayerStrength - currentEnemyStrength
		player.StatAllocation("DamageDone", difference)
		fmt.Println("Player was stronger. Difference:", difference)
		enemy.CurrentHP, enemy.Condition = DamageCalculation(difference, enemy.CurrentHP, enemy.HitPoints)

	} else if currentPlayerStrength < currentEnemyStrength {
		difference = currentEnemyStrength - currentPlayerStrength
		fmt.Println("Enemy was stronger. Difference:", difference)
		player.CurrentHealth, player.Condition = DamageCalculation(difference, player.CurrentHealth, player.HealthRating)
		player.StatAllocation("DamageTaken", difference)
		player.StatAllocation(player.Condition, 0)
	} else {
		// It is a draw
		fmt.Println("Ended up in a draw")
	}
	return player, enemy
}

// DamageCalculation - calculate damage dealt between two adversaries
func DamageCalculation(Damage int, currentHitpoints int, Hitpoints int) (int, string) {
	if Damage > currentHitpoints {
		return 0, "Dead"
	}
	currHP := currentHitpoints - Damage
	condition := ""
	if float64(currHP) > math.Floor(float64(Hitpoints)*0.8) && currHP < Hitpoints {
		condition = "Barely Scratched"
		currHP = Hitpoints
	} else if float64(currHP) > math.Floor(float64(Hitpoints)*0.6) {
		condition = "Minorly Injured"
		currHP = int(float64(Hitpoints) * 0.8)
	} else if float64(currHP) > math.Floor(float64(Hitpoints)*0.4) {
		condition = "Major Injury"
		currHP = int(float64(Hitpoints) * 0.6)
	} else if float64(currHP) > math.Floor(float64(Hitpoints)*0.2) {
		condition = "Severely Injured"
		currHP = int(float64(Hitpoints) * 0.4)
	} else {
		condition = "Mortally Wounded"
		currHP = int(float64(Hitpoints) * 0.2)
	}

	return currHP, condition
}

// Revealed - If a party member with a firearm engaged in battle, roll to determine if party is stealthed or not
func Revealed(players []*structs.Player) bool {
	for _, player := range players {
		if player.HasFought == true {
			if player.Weapon.WeaponType == "Pistol" || player.Weapon.WeaponType == "Rifle" || player.Weapon.WeaponType == "Shotgun" || player.Weapon.WeaponType == "Sniper Rifle" {
				if player.Weapon.WasFired == true {
					// Roll
					roll := rand.Intn(100) + 1
					// Modify roll with player's attributes
					if roll > (30 + player.Proficiencies.Stealth.Level) {
						// Reveal
						return false
					}
				}
			}
		}
	}
	return true
}
