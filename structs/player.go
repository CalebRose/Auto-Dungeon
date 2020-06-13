package structs

// Player - structure for player character
type Player struct {
	Name               string
	Profession         string
	Armor              Armor
	Weapon             Weapon
	Attributes         Attribute
	Proficiencies      Proficiencies
	Behavior           Behaviors
	Stats              Stats
	Inventory          []Item
	Holster            []Weapon
	InventoryLimit     int
	HolsterLimit       int
	Level              int
	Experience         int
	ExperienceRequired int
	HealthRating       int
	CurrentHealth      int
	Condition          string
	InCover            bool
	Ready              bool
	HasFought          bool
	Feats              []Feat
}
