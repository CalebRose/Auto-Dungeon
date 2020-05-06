package structs

type Player struct {
	Name string
	Profession string
	Armor Armor
	Weapon Weapon
	HealthRating int
	CurrentHealth int
	Condition string
	Attributes Attribute
	Proficiencies Proficiencies
	Feats []Feat
}