package structs

// Stats - statistics for a player
type Stats struct {
	Missions            int
	BattlesEngaged      int
	ShotsFired          int
	ShotsMade           int
	DamageDone          int
	DamageTaken         int
	FeatsActivated      int
	PreferredWeapon     string
	PistolUse           int
	RifleUse            int
	ShotgunUse          int
	SniperRifleUse      int
	FistsUse            int
	MeleeWeaponUse      int
	MinorInjuries       int
	MajorInjuries       int
	SevereInjuries      int
	Recoveries          int
	DiscoveriesFound    int
	ItemsLooted         int
	EnemiesKilled       int
	ObstaclesOvercome   int
	BossesDefeated      int
	TrainingActions     int
	BarActions          int
	PersuasionChecks    int
	ItemsStolen         int
	PersonsRescued      int
	PersonsKidnapped    int
	AssassinationsMade  int
	ObjectivesCompleted int
}
