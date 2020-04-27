package structs

type Room struct {
	Name string
	RoomType string
	RoomConditions string
	Visited bool
	Locked bool
	PlayerCover int
	EnemyCover int
	Enemies []Enemy
	Discoveries []Discovery
	Edges []string
	Key string
} 