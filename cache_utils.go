package main

import (
	"context"
	sqlc "crux/db/sqlc-gen"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func lockCache() {
	cacheLock.Lock()
}

func unlockCache() {
	cacheLock.Unlock()
}

func NewProvider(cfg string) sqlc.DBQuerier {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg)
	if err != nil {
		log.Fatal("error connecting db")
	}
	err = db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to the database")
	return sqlc.NewQuerierWithTX(db)
}

func AddReferencesToRuleSetCache() {
	for realmKey, perRealm := range rulesetCache {
		for appKey, perApp := range perRealm {
			for sliceKey, perSlice := range perApp {
				for _, rulesets := range perSlice.BRRulesets {
					for _, rule := range rulesets {
						if rule.RuleActions.ThenCall != "" {
							searchAndAddReferences(rule.RuleActions.ThenCall, rulesetCache, realmKey, appKey, sliceKey, rule)
						}
						if rule.RuleActions.ElseCall != "" {
							searchAndAddReferences(rule.RuleActions.ElseCall, rulesetCache, realmKey, appKey, sliceKey, rule)
						}
					}
				}
			}
		}
	}
}

func searchAndAddReferences(targetSetName string, cache map[realm_t]perRealm_t, realmKey realm_t, appKey app_t, sliceKey slice_t, sourceRule *Ruleset_t) {
	for _, perApp := range cache[realmKey] {
		for otherSliceKey, perSlice := range perApp {
			if otherSliceKey == sliceKey {
				continue
			}
			for _, existingRulesets := range perSlice.BRRulesets {
				for _, existingRule := range existingRulesets {
					if existingRule.SetName == targetSetName {
						sourceRule.RuleActions.References = append(sourceRule.RuleActions.References, existingRule)
					}
				}
			}
		}
	}
}

func PrintAllRuleSetCache() {
	for realmKey, perRealm := range rulesetCache {
		fmt.Println("Realm:", realmKey)
		for appKey, perApp := range perRealm {
			fmt.Println("\tApp:", appKey)
			for sliceKey, perSlice := range perApp {
				fmt.Println("\t\tSlice:", sliceKey)
				fmt.Println("\t\t\tLoadedAt:", perSlice.LoadedAt)

				// Print BRRulesets
				for className, rulesets := range perSlice.BRRulesets {
					fmt.Println("\t\t\tBRRulesets - Class:", className)
					for _, rule := range rulesets {
						fmt.Println("\t\t\t\tRulePatterns:", rule.RulePatterns)
						fmt.Println("\t\t\t\tRuleActions:", rule.RuleActions)
						fmt.Println("\t\t\t\tNMatched:", rule.NMatched)
						fmt.Println("\t\t\t\tNFailed:", rule.NFailed)

						// Print References if available
						for _, reference := range rule.RuleActions.References {
							fmt.Println("\t\t\t\t\tReferenced Rule:")
							fmt.Println("\t\t\t\t\t\tRulePatterns:", reference.RulePatterns)
							fmt.Println("\t\t\t\t\t\tRuleActions:", reference.RuleActions)
							fmt.Println("\t\t\t\t\t\tNMatched:", reference.NMatched)
							fmt.Println("\t\t\t\t\t\tNFailed:", reference.NFailed)
						}
					}
				}

				// Print Workflows
				for className, workflows := range perSlice.Workflows {
					fmt.Println("\t\t\tWorkflows - Class:", className)
					for _, workflow := range workflows {
						fmt.Println("\t\t\t\tRulePatterns:", workflow.RulePatterns)
						fmt.Println("\t\t\t\tRuleActions:", workflow.RuleActions)
						fmt.Println("\t\t\t\tNMatched:", workflow.NMatched)
						fmt.Println("\t\t\t\tNFailed:", workflow.NFailed)

						// Print References if available
						for _, reference := range workflow.RuleActions.References {
							fmt.Println("\t\t\t\t\tReferenced Rule:")
							fmt.Println("\t\t\t\t\t\tRulePatterns:", reference.RulePatterns)
							fmt.Println("\t\t\t\t\t\tRuleActions:", reference.RuleActions)
							fmt.Println("\t\t\t\t\t\tNMatched:", reference.NMatched)
							fmt.Println("\t\t\t\t\t\tNFailed:", reference.NFailed)
						}
					}
				}
			}
		}
	}
}
func PrintAllSchemaCache() {

	for realmKey, perRealm := range schemaCache {
		fmt.Println("Realm:", realmKey)
		for appKey, perApp := range perRealm {
			fmt.Println("\tApp:", appKey)
			for sliceKey, perSlice := range perApp {
				fmt.Println("\t\tSlice:", sliceKey)
				fmt.Println("\t\t\tLoadedAt:", perSlice.LoadedAt)
				for className, schemas := range perSlice.BRSchema {
					fmt.Println("\t\t\tBRSchema - Class:", className)
					for _, schema := range schemas {
						fmt.Println("\t\t\t\tPatternSchema:", schema.PatternSchema)
						fmt.Println("\t\t\t\tActionSchema:", schema.ActionSchema)
						fmt.Println("\t\t\t\tNChecked:", schema.NChecked)
					}
				}
				for className, schemas := range perSlice.WFSchema {
					fmt.Println("\t\t\tWFSchema - Class:", className)
					for _, schema := range schemas {
						fmt.Println("\t\t\t\tPatternSchema:", schema.PatternSchema)
						fmt.Println("\t\t\t\tActionSchema:", schema.ActionSchema)
						fmt.Println("\t\t\t\tNChecked:", schema.NChecked)
					}
				}

			}
		}
	}

}
