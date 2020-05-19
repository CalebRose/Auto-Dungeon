package main

import (
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"

	// "firebase.google.com/go/auth"
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/calebrose/Auto-Dungeon/configu"
	"github.com/calebrose/Auto-Dungeon/structs"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

var client *db.Client

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	c := configu.Config()
	ctx := context.Background()
	opt := option.WithCredentialsFile(c["cred"])
	// config := &firebase.Config{ProjectID: c["projectId"]}
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	// Load Assets
	partyObjective := structs.Objective{
		Name:          "Test",
		ObjectiveType: "Location",
		Description:   "Party needs to reach this room",
		Fulfilled:     false,
		Condition:     "",
		Target:        "TestRoom2",
	}
	party := structs.Party{
		Name:        "Main",
		Description: "Just a description",
		CurrentRoom: "TestRoom",
		Members:     []*structs.Player{},
		Objectives:  []*structs.Objective{&partyObjective},
		Value:       0,
	}
	// Players
	docsnap := client.Collection("Players").Documents(ctx)

	defer docsnap.Stop()

	// Map players to party
	for {
		doc, err := docsnap.Next()

		if doc == nil {
			break
		}
		if err != nil {
			log.Fatalf("Something went wrong %v", err)
		}

		var player structs.Player
		doc.DataTo(&player)

		party.Members = append(party.Members, &player)
	}
	fmt.Println(party)
	// Rooms
	graph := structs.Graph{
		AdjList: map[string]*structs.Room{},
	}

	testRoom := structs.Room{
		Name:           "TestRoom",
		RoomType:       "Start",
		RoomConditions: "?",
		Visited:        false,
		Locked:         false,
		PlayerCover:    3,
		EnemyCover:     0,
		Enemies:        []*structs.Enemy{},
		Discoveries:    []structs.Discovery{},
		Edges:          []string{"Room2", "Room3"},
		Key:            "Room1",
	}

	testRoom2 := structs.Room{
		Name:           "TestRoom2",
		RoomType:       "Room",
		RoomConditions: "?",
		Visited:        false,
		Locked:         false,
		PlayerCover:    3,
		EnemyCover:     0,
		Enemies:        []*structs.Enemy{},
		Discoveries:    []structs.Discovery{},
		Edges:          []string{},
		Key:            "Room2",
	}
	enemyOne := structs.Enemy{
		Name:           "Thief",
		Description:    "A petty thief from the streets.",
		EnemyType:      "Common",
		Condition:      "Healthy",
		HitPoints:      10,
		CurrentHP:      10,
		CombatRating:   4,
		CombatAccuracy: 40,
		Initiative:     0,
		InCover:        false,
		Loot: structs.Item{
			Name:        "Dollars",
			ItemType:    "Currency",
			Description: "It's money.",
			Value:       30,
		},
	}

	testRoom.Enemies = append(testRoom.Enemies, &enemyOne)
	fmt.Println(testRoom)
	// Map rooms to Graph
	graph.AddVertex(&testRoom)
	graph.AddVertex(&testRoom2)
	fmt.Println(graph)
	// Starting Room
	party.CurrentRoom = InitializeStartingPoint(graph.AdjList)
	path := []string{party.CurrentRoom}
	// Assign Objectives

	// The Loop
	for {
		// Check Room
		currentRoom := graph.AdjList[party.CurrentRoom]
		currentRoom.Visited = true
		fmt.Println(currentRoom)
		// Check for Enemies

		if len(currentRoom.Enemies) > 0 {
			// Initiate Battle
			party.Members = InitializeBattle(party.Members, currentRoom.Enemies)
			party.Objectives[0].Fulfilled = CheckObjective(party.Objectives[0], currentRoom, "Eliminate")
			party.Objectives[0].Fulfilled = CheckObjective(party.Objectives[0], currentRoom, "Item")
		}

		// Check if at least one member of the party is alive
		// Otherwise, break loop

		// Check for Discoveries
		if len(currentRoom.Discoveries) > 0 {
			// Roll for Discoveries
			// Loop through each player and run a check on discovering something
		}
		// With the room cleared, check if this is the room the party needs to be in
		party.Objectives[0].Fulfilled = CheckObjective(party.Objectives[0], currentRoom, "Location")

		// Check for Obstacles

		// Check Objectives
		// Check on whether party needed to be in room
		// Loop through player objectives
		if CheckAllObjectives(party.Objectives) == true {
			// If All Objectives are complete, break the infinite loop
			break
		}

		// Check for adjacent rooms -- only when room is cleared? is movable?
		if len(currentRoom.Edges) > 0 {
			if len(currentRoom.Edges) == 1 {
				// If there is only one adjacent room
				// Add current room to path taken
				path = append(path, party.CurrentRoom)
				// Traverse into room
				party.CurrentRoom = currentRoom.Edges[0]
			} else if len(currentRoom.Edges) == 0 {
				// Else if there are no adjacent rooms
				// Traverse Backwards through Path
				lastRoom := len(path) - 1
				party.CurrentRoom = path[lastRoom]
				// Remove last room from array
				path = path[:lastRoom]
			} else {
				// Else if there is more than one adjacent room
				// Choose a Room
				unchosen := true
				NoOtherPath := false
				safetyIterator := 0
				for unchosen == true {
					safetyIterator++
					randomPick := rand.Intn(len(currentRoom.Edges))
					fmt.Println(randomPick)
					chosenRoom := currentRoom.Edges[randomPick]
					if graph.AdjList[chosenRoom].Visited == false {
						unchosen = false
						path = append(path, party.CurrentRoom)
						party.CurrentRoom = chosenRoom
					}
					if safetyIterator > len(currentRoom.Edges)*2 {
						// Safety precaution to prevent infinite loop. Just go back one
						for i := 0; i < len(currentRoom.Edges); i++ {
							iteratedRoom := currentRoom.Edges[i]
							if graph.AdjList[iteratedRoom].Visited == false {
								path = append(path, party.CurrentRoom)
								party.CurrentRoom = chosenRoom
								unchosen = false
								break
							}
						}
						if unchosen == false {
							break
						} else {
							// Traverse backwards
							if len(path) < 1 {
								fmt.Println("THERE IS NO OTHER PATH TO GO")
								NoOtherPath = true
								break
							}
							lastRoom := len(path) - 1
							party.CurrentRoom = path[lastRoom]
							// Remove last room from array
							path = path[:lastRoom]
							break
						}
					}
				}
				if NoOtherPath == true {
					// Party fails
					fmt.Println("For some reason the party ended up back on the starting square. Y'all lost.")
					break
				}
			}
		}
	}

	// So at this point all main objectives were either completed or the party is dead
	// If victory achieved
	// Congrats! Mission completed! You all did it!
	// Log the victory
	// Else
	// Log defeat

	// Log status of each player
	// Update all player objects into Firebase

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

// CheckObjective - Check if an objective is complete
func CheckObjective(obj *structs.Objective, room *structs.Room, objectiveType string) bool {
	// Current Logic implemented for One Objective Only, but structured for multiple just in case
	result := false
	// If the ObjectiveType matches that of the section this function is placed, do the check, if not, ignore
	if obj.ObjectiveType == objectiveType {
		if obj.ObjectiveType == "Eliminate" {
			// Eliminate an Enemy Unit
			// Loop through dead enemy units
			enemies := room.Enemies
			for i := 0; i < len(enemies); i++ {
				// If the enemy in the objective was in the room & dead
				if obj.Target == enemies[i].Name && enemies[i].Condition == "Dead" {
					result = true
					break
				}
			}
		} else if obj.ObjectiveType == "Location" {
			if obj.Target == room.Name {
				result = true
			}
		} else if obj.ObjectiveType == "Item" {
			enemies := room.Enemies
			for i := 0; i < len(enemies); i++ {
				// If the enemy in the objective was in the room & dead
				if obj.Target == enemies[i].Loot.Name && enemies[i].Condition == "Dead" {
					result = true
					break
				}
			}
		}
	}
	return result
}

// InitializeStartingPoint - find the starting room.
func InitializeStartingPoint(rooms map[string]*structs.Room) string {
	for _, i := range rooms {
		if i.RoomType == "Start" {
			return i.Key
		}
	}
	return "NONE"
}

// InitializeBattle - Initialize Battle between Players and Enemies in a room
func InitializeBattle(players []*structs.Player, enemies []*structs.Enemy) []*structs.Player {
	//Take in slice of players & slice of enemies
	// Use Battle Queue Struct instead of map
	// create a slice to take in a struct with a name, type, and initiative. Generate initiative based on player's perception
	battleQueue := []structs.BattleQueue{}
	for i := 0; i < len(players); i++ {
		playerQueue := structs.BattleQueue{
			Name:          players[i].Name,
			CombatantType: "Player",
			EnemyType:     "",
			Initiative:    rand.Intn(players[i].Attributes.Perception) + 1,
		}
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

	playerCount := len(players)
	enemyCount := len(enemies)

	// while all players are not dead || all enemies are not dead
	for playerCount > 0 && enemyCount > 0 {
		// pop from the beginning of the queue
		battleNode, battleSlice := battleSlice[0], battleSlice[1:]

		// Check whether the map received is a player or enemy
		// If player, randomly select an enemy from the enemies array. If enemy, randomly select a player from the players array

		if battleNode.CombatantType == "Player" {
			player := FindPlayer(players, battleNode)
			chosenEnemy := rand.Intn(len(enemies))
			SingularBattle(player, enemies[chosenEnemy], true)
			if enemies[chosenEnemy].Condition == "Dead" {
				enemyCount--
			}
		} else if battleNode.CombatantType == "Enemy" {
			enemy := FindEnemy(enemies, battleNode)
			chosenPlayer := rand.Intn(len(players))
			SingularBattle(players[chosenPlayer], enemy, false)
			if players[chosenPlayer].Condition == "Dead" {
				playerCount--
			}
		}

		battleSlice = append(battleSlice, battleNode)
		// Loop should break when either all enemies are defeated / mortally wounded, or all players are unable to fight
	}
	return players
}

// FindEnemy - find Enemy from the battle queue
func FindEnemy(enemies []*structs.Enemy, battleNode structs.BattleQueue) *structs.Enemy {
	// Find the enemy for singular battle
	i := 0
	for i < len(enemies) {
		if battleNode.Name == enemies[i].Name {
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
			break
		}
		i++
	}
	return players[i]
}

// SingularBattle : A turn of Battle between a player and an enemy.
func SingularBattle(player *structs.Player, enemy *structs.Enemy, playerTurn bool) (*structs.Player, *structs.Enemy) {
	// Establish Base Values & Variables
	playerBaseStrength := player.Attributes.Strength + player.Weapon.WeaponRating
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

			// Fires Shot -- if the shot is less than the Player's Accuracy rating
			shot := rand.Intn(100)

			player.Weapon.CurrentCartridge--

			if shot <= player.Weapon.WeaponAccuracy-enemyBonus {
				// Shot hits target
				currentPlayerStrength = rand.Intn(playerBaseStrength) + 1
			} else {
				fmt.Println("PLAYER MISSED!")
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
		fmt.Println("Player was stronger. Difference:", difference)
		enemy.CurrentHP, enemy.Condition = DamageCalculation(difference, enemy.CurrentHP, enemy.HitPoints)

	} else if currentPlayerStrength < currentEnemyStrength {
		difference = currentEnemyStrength - currentPlayerStrength
		fmt.Println("Enemy was stronger. Difference:", difference)
		player.CurrentHealth, player.Condition = DamageCalculation(difference, player.CurrentHealth, player.HealthRating)

	} else {
		// It is a draw
		fmt.Println("Ended up in a draw")
	}
	return player, enemy
}

// DamageCalculation - calculate damage dealt between two adversaries
func DamageCalculation(Damage int, currentHitpoints int, Hitpoints int) (int, string) {
	//
	if Damage > currentHitpoints {
		return 0, "Dead"
	}
	currHP := currentHitpoints - Damage
	condition := ""
	if float64(currHP) > math.Floor(float64(Hitpoints)*0.8) {
		condition = "Barely Scratched"
	} else if float64(currHP) > math.Floor(float64(Hitpoints)*0.6) {
		condition = "Minorly Injured"
	} else if float64(currHP) > math.Floor(float64(Hitpoints)*0.4) {
		condition = "Major Injury"
	} else if float64(currHP) > math.Floor(float64(Hitpoints)*0.2) {
		condition = "Severely Injured"
	} else {
		condition = "Mortally Wounded"
	}

	return currHP, condition
}
