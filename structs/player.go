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

// StatAllocation -- Allocate, num int a stat to the player based on objectives
func (pl *Player) StatAllocation(stat string, num int) {
	if stat == "Rescue" {
		pl.Stats.PersonsRescued++
	} else if stat == "Kidnapping" {
		pl.Stats.PersonsKidnapped++
	} else if stat == "Theft" {
		pl.Stats.ItemsStolen++
	} else if stat == "Eliminate" {
		pl.Stats.AssassinationsMade++
	} else if stat == "BattleEngaged" {
		pl.Stats.BattlesEngaged++
	} else if stat == "EnemyKilled" {
		pl.Stats.EnemiesKilled++
	} else if stat == "Rifle" {
		pl.Stats.RifleUse++
	} else if stat == "Pistol" {
		pl.Stats.PistolUse++
	} else if stat == "Rifle" {
		pl.Stats.ShotgunUse++
	} else if stat == "Rifle" {
		pl.Stats.SniperRifleUse++
	} else if stat == "Melee" {
		pl.Stats.MeleeWeaponUse++
	} else if stat == "Fists" {
		pl.Stats.FistsUse++
	} else if stat == "ShotMade" {
		pl.Stats.ShotsMade++
	} else if stat == "ShotFired" {
		pl.Stats.ShotsFired++
	} else if stat == "Minorly Injured" {
		pl.Stats.MinorInjuries++
	} else if stat == "Majorly Injured" {
		pl.Stats.MajorInjuries++
	} else if stat == "Severely Injured" {
		pl.Stats.SevereInjuries++
	} else if stat == "DiscoveryMade" {
		pl.Stats.DiscoveriesFound++
	} else if stat == "BossDefeated" {
		pl.Stats.BossesDefeated++
	} else if stat == "DamageTaken" {
		pl.Stats.DamageTaken += num
	} else if stat == "DamageDone" {
		pl.Stats.DamageDone += num
	} else if stat == "ObjectivesCompleted" {
		pl.Stats.ObjectivesCompleted += num
	} else if stat == "Mission" {
		pl.Stats.Missions++
	}
}

// UseMedicalItem - Player Uses an item to heal themselves as opposed to attacking an enemy
func (pl *Player) UseMedicalItem() {
	length := len(pl.Inventory)
	for i := 0; i < length; i++ {
		item := pl.Inventory[i]
		if item.ItemType != "Medical" {
			continue
		}
		pl.CurrentHealth += item.ItemValue
		pl.Inventory[length-1], pl.Inventory[i] = pl.Inventory[i], pl.Inventory[length-1]
		pl.Inventory = pl.Inventory[:length-1]
	}
}
