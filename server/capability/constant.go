package capability

var (
	userID    = "1234"
	authCaps  = []string{"auth"}
	realmName = "Nova"
	capALL    = "ALL"
	capRoot   = "root"

	CAPLIST_REALMLEVEL = []string{"root", "config", "auth", "report"}
	CAPLIST_APPLEVEL   = []string{"schema", "rules"}
)
