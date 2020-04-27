package structs

type Player struct {
	Name string
	Profession string
	Armor Armor
	Weapon Weapon
	HealthRating int
	Condition string
	Attributes Attribute
	Feats []Feat
}