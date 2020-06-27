package structs

// Enemy - Adversaries for the players
type Enemy struct {
	Name           string
	Description    string
	EnemyType      string
	Condition      string
	HitPoints      int
	CurrentHP      int
	CombatRating   int
	CombatAccuracy int
	Initiative     int
	InCover        bool
	Loot           Item
}
