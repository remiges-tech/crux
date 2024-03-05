package ruleset

var TRIGGER = false

func HasRootCapabilities() bool {
	return TRIGGER
}

func GeRuleSetsByRulesetRights() []string {
	return []string{"retailBANK", "nedbank"}
}
