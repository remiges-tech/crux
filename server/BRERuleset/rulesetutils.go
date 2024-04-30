package breruleset

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/remiges-tech/crux/db/sqlc-gen"
)

func hasRulesetCapability(app string, query *sqlc.Queries, c *gin.Context) (bool, error) {
	count, err := query.GetRuleSetCapabilityForApp(c, sqlc.GetRuleSetCapabilityForAppParams{
		Userid: userID,
		Realm:  realmName,
		App:    pgtype.Text{String: app, Valid: true},
		Cap:    RULESET,
	})
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	} else {
		return false, fmt.Errorf("caller does not have ruleset capability for the specified app : %v", app)
	}
}
