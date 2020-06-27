package structs

// PartyBehavior - Party Set Behavior based on player input
type PartyBehavior struct {
	PartyCuriosity  int // Will the Party explore additional rooms along the main path?
	PartyDiscovery  int // Will the party take time to make discoveries
	PartyPersuasion int // Will the Party attempt to be peaceful or aggressive when caught under stealth?
	PartyStealth    int
}
