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




	// Configure to Firebase

	// Load Assets
	party := structs.Party {
		Name: "Main",
		Description: "Just a description",
		Members: []structs.Player{},
		Value: 0,
	}
	// Players
	// LoadPlayers()
		// Rooms

		// Map rooms to Graph

		// Starting Room

		// Assign Objectives

	// Loop


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
	fmt.Println("Player One: " + playerOne.Condition)
	fmt.Println("Enemy: " + enemyOne.Condition)


	// fmt.Println(party)
	// fmt.Println(playerOne)
	// fmt.Println(playerOne.Weapon.Value > playerTwo.Weapon.Value)
	// fmt.Println(playerTwo)

}

func InitializeBattle (players [] struc.Player, enemies []structs.Enemy) {
	//Take in slice of players & slice of enemies
	// Use Battle Queue Struct instead of map
	// create a slice to take in a map with a name, type, and initiative. Generate initiative based on player's perception

	// Sort the slice from highest initiative to lowest. Confirm whether it's possible to sort a number in string form
	// FOR LATER -- do a feat check for players w/ sleight of hand

	// Utilize slice as a 'queue'

	// while all players are not dead || all enemies are not dead

	// pop from the beginning of the queue

	// Check whether the map received is a player or enemy

	// If player, randomly select an enemy from the enemies array. If enemy, randomly select a player from the players array

	// Utilize SingularBattle function. Include boolean for player/enemy turn as parameter to signify whether an enemy or player is attacking

	// Loop should break when either all enemies are defeated / mortally wounded, or all players are unable to fight
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
				enemyOne.Condition = "Dead"
			} else if float64(difference) > math.Floor(float64(enemyBaseStrength) * 0.5) {
				enemyOne.Condition = "Major Injury"
			} else {
				enemyOne.Condition = "'Tis but a scratch"
			}
			// if the enemy's strength is greater than 0
		} else if currentPlayerStrength < currentEnemyStrength {
			difference = currentEnemyStrength - currentPlayerStrength
			fmt.Println("Enemy was stronger. Difference:", difference)
			if difference > player.HealthRating {
				playerOne.Condition = "Dead"
			} else if float64(difference) > math.Floor(float64(playerOne.HealthRating) * 0.5) {
				playerOne.Condition = "Major Injury"
			} else if float64(difference) > math.Floor(float64(playerOne.HealthRating) * 0.2) {
				playerOne.Condition = "Minor Injury"
			} else {
				// Tis but a scratch
			}
		} else {
			// It is a draw
		}
	}	
}

func LoadPlayers(){
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

		party.Members = append(party.Members, player)
	}
}
