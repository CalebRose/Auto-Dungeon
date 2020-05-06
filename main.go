package main

import (
	
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
  	// "firebase.google.com/go/auth"
	"fmt"
	"github.com/calebrose/Auto-Dungeon/configu"
	"github.com/calebrose/Auto-Dungeon/structs"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"math/rand"
	"math"
	"sort"
	"time"
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
	party := structs.Party {
		Name: "Main",
		Description: "Just a description",
		CurrentRoom: "TestRoom"
		Members: []*structs.Player{},
		Objectives: []*structs.Objective{},
		Value: 0,
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
	graph := structs.Graph {
		AdjList: map[string]*structs.Room{},
	}

	testRoom := structs.Room {
		Name: "TestRoom",
		RoomType: "StartingRoom",
		RoomConditions: "?",
		Visited: false,
		Locked: false,
		PlayerCover: 3,
		EnemyCover: 0,
		Enemies: []*structs.Enemy{},
		Discoveries: []structs.Discovery{},
		Edges: []string{"TestRoom2"},
		Key: "Room1",
	}

	testRoom2 := structs.Room {
		Name: "TestRoom2",
		RoomType: "Room",
		RoomConditions: "?",
		Visited: false,
		Locked: false,
		PlayerCover: 3,
		EnemyCover: 0,
		Enemies: []*structs.Enemy{},
		Discoveries: []structs.Discovery{},
		Edges: []string{},
		Key: "Room2",
	}
	// Map rooms to Graph
	graph.AddVertex(&testRoom)
	graph.AddVertex(&testRoom2)
	fmt.Println(graph)
	// Starting Room
	party.CurrentRoom = InitializeStartingPoint(graph.AdjList)
	path := []string {party.CurrentRoom}
	// Assign Objectives

	// The Loop
	for CheckObjectives(party.Objectives) == false {
		// Check Room
		currentRoom := graph.AdjList[party.CurrentRoom]
		currentRoom.Visited = true
		// Check for Enemies
		if(len(currentRoom.Enemies) > 0) {
			// Initiate Battle
			InitializeBattle(party.Members, currentRoom.Enemies)
		}

		// Check if at least one member of the party is alive
		// Otherwise, break loop
		
		// Check for Discoveries
		if(len(currentRoom.Discoveries) > 0) {
			// Roll for Discoveries
			// Loop through each player and run a check on discovering something
		}

		// Check Objectives
			// Check on whether party needed to be in room
			// Loop through player objectives 


		// Check for adjacent rooms
		if(len(currentRoom.Edges) > 0) {
			if(len(currentRoom.Edges) == 1) {
				// If there is only one adjacent room
				// Add current room to path taken
				path = append(path, party.CurrentRoom)
				// Traverse into room
				party.CurrentRoom = currentRoom.Edges[0]
			} else if (len(currentRoom.Edges) == 0) {
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
				safety_iterator := 0
				for unchosen == true {
					safety_iterator++
					random_pick := rand.Intn(len(currentRoom.Edges))
					chosenRoom := currentRoom.Edges[random_pick]
					if(graph.AdjList[chosenRoom].Visited == false) {
						unchosen = false
						path = append(path, party.CurrentRoom)
						party.CurrentRoom = chosenRoom
					}
					if(safety_iterator > len(currentRoom.Edges) * 2) {
						// Safety precaution to prevent infinite loop. Just go back one
						for i := 0; i < len(currentRoom.Edges); i++ {
							iteratedRoom = currentRoom.Edges[i]
							if(graph.AdjList[iteratedRoom].Visited == false) {
								path = append(path, party.CurrentRoom)
								party.CurrentRoom = chosenRoom
								unchosen = false
								break
							}
						}
						if(unchosen == false) {
							break
						} else {
							// Traverse backwards
							lastRoom := len(path) - 1
							party.CurrentRoom = path[lastRoom]
							// Remove last room from array
							path = path[:lastRoom]
							break
						}
					}
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



	enemyOne := structs.Enemy {
		Name: "Thief",
		Description: "A petty thief from the streets.",
		EnemyType: "Common",
		Condition: "Healthy",
		CombatRating: 4,
		CombatAccuracy: 40,
		Loot: structs.Item {
			Name: "Dollars",
			ItemType: "Currency",
			Description: "It's money.",
			Value: 30,
		},
	}

	testRoom.Enemies = append(testRoom.Enemies, &enemyOne)

}

func CheckObjectives(Objectives []structs.Objective) bool {
	totalObj := len(Objectives)
	completedObj := 0
	for _, obj := Objectives {
		if(obj.Fulfilled === true) {
			completedObj++
		}
	}
	if(completedObj == totalObj)
		return true
	else
		return false
}

func InitializeStartingPoint(rooms map[string]){
	for _, i := rooms {
		if(i.RoomType == "Start") {
			return i.Key
		}
	}
	return "NONE"
}

func InitializeBattle(players [] structs.Player, enemies []structs.Enemy) {
	//Take in slice of players & slice of enemies
	// Use Battle Queue Struct instead of map
	// create a slice to take in a struct with a name, type, and initiative. Generate initiative based on player's perception
	battleQueue := []structs.BattleQueue{}
	for i:= 0; i < len(players); i++ {
		playerQueue := structs.BattleQueue {
			Name: players[i].Name,
			CombatantType: "Player",
			EnemyType: "",
			Initiative: rand.Intn(players[i].Attributes.Perception) + 1,
		}
		battleQueue = append(battleQueue, playerQueue)
	}

	for i:= 0; i < len(enemies); i++ {
		enemyQueue := structs.BattleQueue {
			Name: enemies[i].Name,
			CombatantType: "Enemy",
			EnemyType: enemies[i].EnemyType,
			Initiative: enemies[i].Initiative,
		}
		battleQueue = append(battleQueue, enemyQueue)
	}
	// Sort the slice from highest initiative to lowest
	// FOR LATER -- do a feat check for players w/ sleight of hand
	battleSlice := battleQueue[:]
	sort.Slice(battleSlice, func(i,j int) bool {
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

		if(battleNode.CombatantType == "Player") {
			player := FindPlayer(players, battleNode)
			chosenEnemy := rand.Intn(len(enemies))
			SingularBattle(player, enemies[chosenEnemy])
		} else if (battleNode.CombatantType == "Enemy") {
			enemy := FindEnemy(enemies, battleNode)
			chosenPlayer := rand.Intn(len(players))
			SingularBattle(players[chosenPlayer], enemy)
		}

		battleSlice = append(battleSlice, battleNode)
		// Loop should break when either all enemies are defeated / mortally wounded, or all players are unable to fight
	}
	

	

}

func FindEnemy(enemies []structs.Enemy, battleNode structs.BattleQueue)  structs.Enemy {
	// Find the enemy for singular battle
	i := 0
	for i < len(enemies){
		if battleNode.Name == enemies[i].Name {
			break		
		}
		i++
	}
	return enemies[i]
}

func FindPlayer(players []structs.Player, battleNode structs.BattleQueue) structs.Player {
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

func SingularBattle(player structs.Player, enemy structs.Enemy) {
	// Revamp for multiple players
	
	// Establish Base Values & Variables
	playerBaseStrength := player.Attributes.Strength + player.Weapon.WeaponRating
	enemyBaseStrength := enemy.CombatRating

	fmt.Println("Player Base Strength:", playerBaseStrength)

	difference := 0

	currentWeaponCartridge := player.Weapon.WeaponCartridge
	// Since it's gun combat, makes sense to have one firing at a time.
	// Use playerTurn to signfiy whether it's the player's turn or the enemy's
	playerTurn := true

	// Change this loop
	// While the player is Healthy and the enemy ain't dead
	for player.Condition == "Healthy" && enemy.Condition != "Dead" {
		difference = 0
		currentPlayerStrength := 0
		currentEnemyStrength := 0

		// If it's the player's turn:
		if playerTurn == true {
			// If the player has ammo in his gun
			if currentWeaponCartridge > 0 {

				// Fires Shot -- if the shot is less than the Player's Accuracy rating
				shot := rand.Intn(100)

				currentWeaponCartridge--

				if shot <= player.Weapon.WeaponAccuracy {
					// Shot hits target
					currentPlayerStrength = rand.Intn(playerBaseStrength) + 1
				} else {
					fmt.Println("PLAYER MISSED!")
				}
			} else {
				currentWeaponCartridge = player.Weapon.WeaponCartridge
				fmt.Println("Reloading...")
			}
			// Set to false
			playerTurn = false
		} else {
			// Enemy's turn to fire
			enemyShot := rand.Intn(100)
			if enemyShot <= enemy.CombatAccuracy {
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
			if difference > enemyBaseStrength {
				enemy.Condition = "Dead"
			} else if float64(difference) > math.Floor(float64(enemyBaseStrength) * 0.5) {
				enemy.Condition = "Major Injury"
			} else {
				enemy.Condition = "'Tis but a scratch"
			}
			// if the enemy's strength is greater than 0
		} else if currentPlayerStrength < currentEnemyStrength {
			difference = currentEnemyStrength - currentPlayerStrength
			fmt.Println("Enemy was stronger. Difference:", difference)
			if difference > player.HealthRating {
				player.Condition = "Dead"
			} else if float64(difference) > math.Floor(float64(player.HealthRating) * 0.5) {
				player.Condition = "Major Injury"
			} else if float64(difference) > math.Floor(float64(player.HealthRating) * 0.2) {
				player.Condition = "Minor Injury"
			} else {
				// Tis but a scratch
			}
		} else {
			// It is a draw
		}
	}	
}
