package structs

// Stats - statistics for a player
type Stats struct {
	Missions          int
	BattlesEngaged    int
	ShotsFired        int
	ShotsMade         int
	DamageDone        int
	DamageTaken       int
	FeatsActivated    int
	PreferredWeapon   string
	MinorInjuries     int
	MajorInjuries     int
	SevereInjuries    int
	Recoveries        int
	DiscoveriesFound  int
	EnemiesLooted     int
	ObstaclesOvercome int
	BossesDefeated    int
	TrainingActions   int
	BarActions        int
}
