package structs

import (
	"math"
	"math/rand"
)

// Room - the unique vacinities and areas the party traverses through
type Room struct {
	Name           string
	RoomType       string // Type of room. Interior, Exterior
	RoomConditions string
	Description    string
	Visited        bool // Indicates whether the party has visited the room previously
	Locked         bool // Indicates whether the room is locked
	Continuous     bool // Value indicates whether enemies will continuously load in room if the party stays in room
	PlayerCover    int  // Number of players that can get into cover in a room
	EnemyCover     int  // Number of enemies that can get into cover in a room
	InitEnemyCount int  // The total number of enemies that can appear in this room
	Enemies        []*Enemy
	Discoveries    []*Discovery
	Targets        []Target // Change to pointer?
	Edges          []string
	Obstacles      []Obstacle
	Key            string
}

// LoadEnemies - load enemies into a room
func (r *Room) LoadEnemies(enemiesJSON []Enemy) {
	// Load initial enemies
	for i := 0; i < r.InitEnemyCount; i++ {
		// Roll to generate enemy
		roll := rand.Intn(100) + 1
		// Roll is dependent on the 10 enemies being loaded
		num := math.Floor(float64(roll) / 10.0)
		chosenEnemy := enemiesJSON[int(num)]
		r.Enemies = append(r.Enemies, &chosenEnemy)

	}
}

// LoadDiscoveries - Load Discoveries into room
func (r *Room) LoadDiscoveries(discoveries []*Discovery) {
	// There can be 0, 1, and 2 discoveries in a room
	roll := rand.Intn(100) + 1
	count := 0
	if roll < 11 {
		count = 0
	} else if roll > 10 && roll < 86 {
		count = 2
	} else {
		count = 3
	}

	for i := 0; i < count; i++ {
		discoRoll := rand.Intn(len(discoveries))
		chosenDiscovery := discoveries[discoRoll]
		r.Discoveries = append(r.Discoveries, chosenDiscovery)
	}

}
