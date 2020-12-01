package main

import (
	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"

	// "firebase.google.com/go/auth"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	battle "github.com/calebrose/Auto-Dungeon/BattleFunctions"
	dis "github.com/calebrose/Auto-Dungeon/DiscoveryFunctions"
	mis "github.com/calebrose/Auto-Dungeon/Mission"
	obj "github.com/calebrose/Auto-Dungeon/ObjectiveFunctions"
	"github.com/calebrose/Auto-Dungeon/configu"
	"github.com/calebrose/Auto-Dungeon/structs"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

var client *db.Client

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	var sb strings.Builder
	sb.WriteString(exPath)
	sb.WriteString("/cred/serviceAccountKey.json")
	c := configu.Config()
	ctx := context.Background()
	opt := option.WithCredentialsFile(exPath + c["cred"])
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
	Mission := mis.Mission{
		MissionName:        "Test Mission",
		Location:           "The Walderlund",
		MissionEnemies:     []structs.Enemy{},
		MissionDiscoveries: []*structs.Discovery{},
	}
	partyObjective := structs.Objective{
		Name:           "Test",
		ObjectiveType:  "Location",
		Description:    "Party needs to reach this room",
		Fulfilled:      false,
		Condition:      "",
		TargetLocation: "TestRoom",
	}
	party := structs.Party{
		Name:          "Main",
		Description:   "Just a description",
		CurrentRoom:   "TestRoom",
		Members:       []*structs.Player{},
		Objectives:    []*structs.Objective{&partyObjective},
		Targets:       []*structs.Target{},
		IsStealth:     true,
		EscapeRoute:   []string{"TestRoom"},
		InRetreat:     false,
		MissionStatus: false,
		PartyBehavior: structs.PartyBehavior{
			PartyCuriosity:  0,
			PartyDiscovery:  0,
			PartyPersuasion: 0,
		},
	}
	// Players
	// docsnap := client.Collection("Players").Documents(ctx)
	// Grabs all players that are set to "Ready"
	docsnap := client.Collection("Players").Where("Ready", "==", true).Documents(ctx)

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
		player.StatAllocation("Mission", 0)
		party.Members = append(party.Members, &player)
	}
	// Set Party Behaviors
	party.SetBehaviors(party.Members)

	// Load Global Data
	Mission.LoadEnemies()
	Mission.LoadDiscoveries()

	// Rooms
	graph := structs.Graph{
		AdjList: map[string]*structs.Room{},
	}

	// Load Rooms Here

	testRoom := structs.Room{
		Name:           "TestRoom",
		RoomType:       "Start",
		RoomConditions: "?",
		Description:    "This is a test room",
		Visited:        false,
		Locked:         false,
		Continuous:     false,
		PlayerCover:    3,
		EnemyCover:     0,
		InitEnemyCount: 0,
		Enemies:        []*structs.Enemy{},
		Discoveries:    []*structs.Discovery{},
		Edges:          []string{"Room2"},
		Key:            "Room1",
	}

	testRoom2 := structs.Room{
		Name:           "TestRoom2",
		RoomType:       "Room",
		RoomConditions: "?",
		Description:    "This is a test room",
		Visited:        false,
		Locked:         false,
		Continuous:     false,
		PlayerCover:    3,
		EnemyCover:     0,
		InitEnemyCount: 3,
		Enemies:        []*structs.Enemy{},
		Discoveries:    []*structs.Discovery{},
		Edges:          []string{},
		Key:            "Room2",
	}
	fmt.Println(testRoom)
	// Map rooms to Graph
	graph.AddVertex(&testRoom)
	graph.AddVertex(&testRoom2)
	fmt.Println(graph)
	// Starting Room
	party.CurrentRoom = InitializeStartingPoint(graph.AdjList)
	// Assign Objectives

	// CheckConditions
	partyStatus := true    // The status of the party
	missionStatus := false // The status of the mission
	failure := false
	NoOtherPath := false

	// The Loop
	for {
		// Check Room
		currentRoom := graph.AdjList[party.CurrentRoom]
		currentRoom.Visited = true

		// LOAD DATA
		// If current room involves continuous enemies or the number of enemies is 0, load enemies into room
		if currentRoom.Continuous == true || len(currentRoom.Enemies) == 0 {
			currentRoom.LoadEnemies(Mission.MissionEnemies)
		}

		// LOAD DISCOVERIES
		if len(currentRoom.Discoveries) == 0 {
			currentRoom.LoadDiscoveries(Mission.MissionDiscoveries)
		}

		// Check for Enemies

		if len(currentRoom.Enemies) > 0 {
			// If party is under stealth

			if party.IsStealth {
				// Roll for detection?
				// How can I do a precision bonus?
				DetectionRoll := rand.Intn(100) + 1
				partyStealth := party.PartyBehavior.PartyStealth
				if DetectionRoll+partyStealth > 40 {
					// Keep Stealth, duh.
				} else {
					// Party is seen.
					party.KeepStealth(currentRoom.Enemies)
				}
			}
			// Initiate Battle
			if !party.IsStealth || rand.Intn(20) > 4 {
				party.Members = battle.InitializeBattle(party, currentRoom.Enemies, currentRoom)
			}

			// Assuming battle was made, if a player had a firearm, the party is revealed.
			if party.IsStealth {
				party.IsStealth = battle.Revealed(party.Members)
			}
		}

		// Check if at least one member of the party is alive
		// Otherwise, break loop
		partyStatus = CheckParty(party.Members)
		if !partyStatus {
			// If the Party is wiped, end the loop
			break
		}
		// Target Check -- for Rescue, Kidnap, Rescue, and Thieving
		if !party.MissionStatus || !failure {
			party, currentRoom = obj.TargetCheck(party, currentRoom)

			// Check for Discoveries
			if len(currentRoom.Discoveries) > 0 {
				// Roll for Discoveries
				// Loop through each player and run a check on discovering something
				party.Members, currentRoom.Discoveries = dis.MakeDiscovery(party, currentRoom.Discoveries)
			}

			// Check Objectives
			// With the room cleared, check if this is the room the party needs to be in
			party, failure = obj.CheckObjectiveCompletion(party, currentRoom)
		}
		// Check on whether party needed to be in room
		// Loop through player objectives

		// Check for Obstacles
		if len(currentRoom.Obstacles) > 0 {
			party, currentRoom = obj.CheckObstacles(party, currentRoom)
		}

		party.MissionStatus = obj.CheckAllObjectives(party.Objectives)
		if party.MissionStatus || failure {
			// Either all objectives are completed or one objective was a failure. Break the infinite loop
			for _, player := range party.Members {
				if player.Condition == "Dead" {
					continue
				}
				player.StatAllocation("ObjectivesCompleted", len(party.Objectives))
			}

			party.InRetreat = true
		}

		party.BreakLoop = obj.BreakLoop(party, failure)
		if party.BreakLoop {
			// Ends the Mission
			break
		}

		// If in escape
		if party.InRetreat == false {
			party, NoOtherPath = TraverseBackwards(party, NoOtherPath)
		} else if !currentRoom.Locked {
			// Check for adjacent rooms -- only when room is cleared? is movable?
			if len(currentRoom.Edges) > 0 {
				if len(currentRoom.Edges) == 1 {
					// If there is only one adjacent room
					// Add current room to path taken
					party.EscapeRoute = append(party.EscapeRoute, party.CurrentRoom)
					// Traverse into room
					party.CurrentRoom = currentRoom.Edges[0]
				} else if len(currentRoom.Edges) == 0 {
					// Else if there are no adjacent rooms
					// Traverse Backwards through Path
					lastRoom := len(party.EscapeRoute) - 1
					party.CurrentRoom = party.EscapeRoute[lastRoom]
					// Remove last room from array
					party.EscapeRoute = party.EscapeRoute[:lastRoom]
				} else {
					// Else if there is more than one adjacent room
					// Choose a Room
					unchosen := true
					safetyIterator := 0
					for unchosen == true {
						safetyIterator++
						randomPick := rand.Intn(len(currentRoom.Edges))
						fmt.Println(randomPick)
						chosenRoom := currentRoom.Edges[randomPick]
						if graph.AdjList[chosenRoom].Visited == false {
							unchosen = false
							party.EscapeRoute = append(party.EscapeRoute, party.CurrentRoom)
							party.CurrentRoom = chosenRoom
						}
						if safetyIterator > len(currentRoom.Edges)*2 {
							// Safety precaution to prevent infinite loop. Just go back one
							for i := 0; i < len(currentRoom.Edges); i++ {
								iteratedRoom := currentRoom.Edges[i]
								if graph.AdjList[iteratedRoom].Visited == false {
									party.EscapeRoute = append(party.EscapeRoute, party.CurrentRoom)
									party.CurrentRoom = chosenRoom
									unchosen = false
									break
								}
							}
							if unchosen == false {
								break
							} else {
								// Traverse backwards
								party, NoOtherPath = TraverseBackwards(party, NoOtherPath)
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

	}

	if partyStatus == false || failure {
		fmt.Println("The Party Has Wiped. You all have failed Lord Konderstahl.")
	} else if missionStatus == true {
		fmt.Println("Mission deemed a success. You all have pleased Lord Konderstahl greatly.")
	}
	// Log status of each player
	// Update all player objects into Firebase
	for _, player := range party.Members {
		playerRef := client.Collection("Players").Doc("Cayetano") // Replace with Discord ID
		// playerJSON, err := json.Marshal(party.Members[0])
		// if err != nil {
		// 	fmt.Println(err)
		// }
		playerRef.Set(ctx, player)
	}

}

// TraverseBackwards -- The Party retreats back
func TraverseBackwards(party structs.Party, NoOtherPath bool) (structs.Party, bool) {
	// Traverse backwards
	if len(party.EscapeRoute) < 1 {
		fmt.Println("THERE IS NO OTHER PATH TO GO")
		NoOtherPath = true
		return party, NoOtherPath
	}
	lastRoom := len(party.EscapeRoute) - 1
	party.CurrentRoom = party.EscapeRoute[lastRoom]
	// Remove last room from array
	party.EscapeRoute = party.EscapeRoute[:lastRoom]

	return party, NoOtherPath
}

// CheckParty - Check if the party is still alive
func CheckParty(party []*structs.Player) bool {
	for _, i := range party {
		if i.Condition != "Dead" {
			return true
		}
	}
	return false
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
