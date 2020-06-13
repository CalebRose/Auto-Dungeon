package structs

// Target - individuals whom the party may target in terms of rescue, kidnap, or escort
type Target struct {
	TargetName        string
	TargetDescription string
	TargetType        string
	TargetHealth      int
	CurrentHealth     int
	Condition         string
	ObjRequirement    int
	Acquired          bool
	ObjRequired       bool
}
