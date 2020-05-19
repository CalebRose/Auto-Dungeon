package structs

// Room - the unique vacinities and areas the party traverses through
type Room struct {
	Name           string
	RoomType       string
	Description    string
	RoomConditions string
	Visited        bool
	Locked         bool
	PlayerCover    int
	EnemyCover     int
	Enemies        []*Enemy
	Discoveries    []Discovery
	Edges          []string
	Key            string
}
