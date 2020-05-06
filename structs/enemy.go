package structs

type Enemy struct {
	Name string
	Description string
	EnemyType string
	Condition string
	CombatRating int
	CombatAccuracy int
	Initiative int
	Loot Item
} 