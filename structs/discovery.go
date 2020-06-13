package structs

// Discovery - What the player can discover in a room
type Discovery struct {
	Name                 string
	DiscoveryDescription string
	DiscoveryType        string
	DiscoveryRating      int
	DiscoveryValue       int
	Discovered           bool
	Acquired             bool
	Weapon
	Item
	Target
	Vehicle
}
