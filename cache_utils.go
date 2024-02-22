package main

import (
	"context"
	sqlc "crux/db/sqlc-gen"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"

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
							searchAndAddReferences(rule.RuleActions.ThenCall, rulesetCache, realmKey, appKey, sliceKey, rule, "thencall")
						}
						if rule.RuleActions.ElseCall != "" {
							searchAndAddReferences(rule.RuleActions.ElseCall, rulesetCache, realmKey, appKey, sliceKey, rule, "elsecall")
						}
					}
				}
			}
		}
	}
}

func searchAndAddReferences(targetSetName string, cache map[realm_t]perRealm_t, realmKey realm_t, appKey app_t,
	sliceKey slice_t, sourceRule *Ruleset_t, calltype string) {
	for _, perApp := range cache[realmKey] {
		for otherSliceKey, perSlice := range perApp {
			if otherSliceKey == sliceKey {
				continue
			}
			for _, existingRulesets := range perSlice.BRRulesets {
				for _, existingRule := range existingRulesets {
					if existingRule.SetName == targetSetName {
						existingRule.ReferenceType = calltype
						sourceRule.RuleActions.References = append(sourceRule.RuleActions.References, existingRule)

					}
				}
			}
			for _, existingRulesets := range perSlice.Workflows {
				for _, existingRule := range existingRulesets {
					if existingRule.SetName == targetSetName {
						existingRule.ReferenceType = calltype
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

func containsField(value interface{}, fieldName string, t *testing.T) bool {

	switch v := value.(type) {

	case []byte:

		var raw json.RawMessage
		if err := json.Unmarshal(v, &raw); err != nil {
			fmt.Println("Error unmarshalling actual pattern:", err, v)
			return false
		}

		var data map[string]interface{}
		if err := json.Unmarshal(raw, &data); err != nil {
			var arrayData []interface{}
			if err := json.Unmarshal(raw, &arrayData); err != nil {
				fmt.Println("Error unmarshalling actual pattern:", err, v)
				return false
			}
			for _, element := range arrayData {
				if containsFieldName(element, fieldName) {
					return true
				}
			}
		}
		for _, value := range data {
			if containsFieldName(value, fieldName) {
				return true
			}
		}
	case map[string]interface{}:

		for key := range v {
			if key == fieldName {
				return true
			}
		}

	case []interface{}:
		for _, item := range v {
			if containsField(item, fieldName, t) {
				return true
			}
		}
	case string:
		return v == fieldName
	}
	return false
}

func containsFieldName(value interface{}, fieldName string) bool {

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Map:
		for _, key := range v.MapKeys() {
			if key.Kind() == reflect.String && key.String() == fieldName {
				return true
			}
			if nestedValue := v.MapIndex(key).Interface(); containsFieldName(nestedValue, fieldName) {
				return true
			}
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			if nestedValue := v.Index(i).Interface(); containsFieldName(nestedValue, fieldName) {
				return true
			}
		}
	case reflect.String:
		return value.(string) == fieldName
	}
	return false
}
