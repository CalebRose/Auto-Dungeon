package structs

// Obstacle -- entities blocking the party from traversing through a room
type Obstacle struct {
	ObstacleName        string
	ObstacleDescription string
	ObstacleType        string
	ObstacleRequirement int
	ObstaclePassed      bool
}
