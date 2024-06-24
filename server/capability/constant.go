package capability

var (
	userID    = "admin"
	authCaps  = []string{"auth"}
	realmName = "Nova"
	capALL    = "ALL"
	capRoot   = "root"

	CAPLIST_REALMLEVEL = []string{"root", "config", "auth", "report"}
	CAPLIST_APPLEVEL   = []string{"schema", "rules"}
)
