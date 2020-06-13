package discoveryfunctions

import (
	"math/rand"

	"github.com/calebrose/Auto-Dungeon/structs"
)

// MakeDiscovery -- Allows the party to make a discovery
func MakeDiscovery(party structs.Party, discoveries []*structs.Discovery) ([]*structs.Player, []*structs.Discovery) {
	for _, discovery := range discoveries {
		if discovery.Acquired == true {
			continue
		}
		pickPlayer := rand.Intn(len(party.Members))
		player := party.Members[pickPlayer]
		if discovery.Discovered == false {

			discoveryRoll := rand.Intn(20)
			if discoveryRoll+player.Attributes.Perception > discovery.DiscoveryRating {
				// Discovery Made
				discovery.Discovered = true
			}
		}
		if discovery.Discovered == true {
			if discovery.DiscoveryType == "Item" && (len(player.Inventory) < player.InventoryLimit) {
				item := structs.Item{
					Name:            discovery.Name,
					ItemType:        discovery.ItemType,
					ItemDescription: discovery.DiscoveryDescription,
					ItemOrigin:      discovery.ItemOrigin,
					ItemValue:       discovery.ItemValue,
					ItemRating:      discovery.ItemRating,
				}
				player.Inventory = append(player.Inventory, item)
				discovery.Acquired = true
			} else if discovery.DiscoveryType == "Weapon" && len(player.Holster) < player.HolsterLimit {
				// 1-10, worst; 11-20, second worst; 21-75 (55), normal; 76-90, uncommon (15); 91-100 (10), rare
				rarityRoll := rand.Intn(100)
				if rarityRoll < 11 {
					discovery.Name = "Worn-Down " + discovery.Name
					if discovery.Requirement > 1 {
						discovery.Requirement--
					}
					discovery.WeaponRating -= 2
					discovery.WeaponAccuracy -= 10
					discovery.WeaponValue -= 50
					discovery.WeaponReloadTime++
				} else if rarityRoll < 21 {
					discovery.Name = "Used " + discovery.Name
					if discovery.Requirement > 1 {
						discovery.Requirement--
					}
					discovery.WeaponRating--
					discovery.WeaponAccuracy -= 5
					discovery.WeaponValue -= 25
				} else if rarityRoll < 76 {
					// Do nothing. It basic
				} else if rarityRoll < 91 {
					discovery.Name = "Polished " + discovery.Name
					if discovery.Requirement < 10 {
						discovery.Requirement++
					}
					discovery.WeaponRating++
					discovery.WeaponAccuracy += 5
					discovery.WeaponValue += 25
				} else if rarityRoll < 101 {
					discovery.Name = "Mint " + discovery.Name
					if discovery.Requirement < 10 {
						discovery.Requirement++
					}
					discovery.WeaponRating += 2
					discovery.WeaponAccuracy += 10
					discovery.WeaponValue += 50
				}
				weapon := structs.Weapon{
					Name:              discovery.Name,
					WeaponDescription: discovery.DiscoveryDescription,
					WeaponType:        discovery.WeaponType,
					AttackType:        discovery.AttackType,
					WeaponValue:       discovery.WeaponValue,
					Requirement:       discovery.Requirement,
					WeaponRating:      discovery.WeaponRating,
					WeaponRange:       discovery.WeaponRange,
					WeaponAccuracy:    discovery.WeaponAccuracy,
					WeaponCartridge:   discovery.WeaponCartridge,
					CurrentCartridge:  discovery.WeaponCartridge,
					WeaponReloadTime:  discovery.WeaponReloadTime,
					CurrentReload:     discovery.WeaponReloadTime,
					WasFired:          false,
				}
				player.Holster = append(player.Holster, weapon)
				discovery.Acquired = true
			} else if discovery.DiscoveryType == "Vehicle" {

			} else if discovery.DiscoveryType == "Target" {
				target := structs.Target{
					TargetName:        discovery.TargetName,
					TargetDescription: discovery.DiscoveryDescription,
					TargetType:        discovery.TargetType,
					TargetHealth:      discovery.TargetHealth,
					CurrentHealth:     discovery.CurrentHealth,
					Condition:         discovery.Condition,
				}
				party.Targets = append(party.Targets, &target)
				discovery.Acquired = true
			}
		}

	}
	return (party.Members), (discoveries)
}
