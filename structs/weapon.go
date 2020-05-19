package structs

// Weapon - objects players use to take on adversaries
type Weapon struct {
	Name             string
	WeaponType       string
	Value            int
	Description      string
	Requirement      int
	WeaponRating     int
	WeaponRange      int
	WeaponAccuracy   int
	WeaponCartridge  int
	CurrentCartridge int
	WeaponReloadTime int
	CurrentReload    int
}
