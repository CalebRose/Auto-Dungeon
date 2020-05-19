package structs

type Player struct {
	Name string
	Profession string
	Armor Armor
	Weapon Weapon
	HealthRating int
	CurrentHealth int
	Condition string
	Level int
	Attributes Attribute
	Proficiencies Proficiencies
	Experience int
	ExperienceRequired int
	InCover bool
	Ready bool

	Feats []Feat
}