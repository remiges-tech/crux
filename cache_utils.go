package main

import (
	"context"
	sqlc "crux/db/sqlc-gen"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"
	"time"

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
		for _, perApp := range perRealm {
			for sliceKey, perSlice := range perApp {
				for _, rulesets := range perSlice.BRRulesets {
					for _, rule := range rulesets {
						for _, subRule := range rule.Rules {
							if subRule.RuleActions.ThenCall != "" {
								searchAndAddReferences(subRule.RuleActions.ThenCall, rulesetCache, realmKey, sliceKey, rule, "thencall", subRule)
							}
							if subRule.RuleActions.ElseCall != "" {
								searchAndAddReferences(subRule.RuleActions.ElseCall, rulesetCache, realmKey, sliceKey, rule, "elsecall", subRule)
							}
						}
					}
				}
			}
		}
	}
}

func searchAndAddReferences(targetSetName string, cache map[realm_t]perRealm_t, realmKey realm_t,
	sliceKey slice_t, sourceRule *Ruleset_t, calltype string, subRule rule_t) {
	for _, perApp := range cache[realmKey] {
		for otherSliceKey, perSlice := range perApp {
			if otherSliceKey == sliceKey {
				continue
			}
			for _, existingRulesets := range perSlice.BRRulesets {
				for _, existingRule := range existingRulesets {
					if existingRule.SetName == targetSetName {
						existingRule.ReferenceType = calltype
						sourceRule.Rules[0].RuleActions.References = append(sourceRule.Rules[0].RuleActions.References, existingRule)

					}
				}
			}
			for _, existingRulesets := range perSlice.Workflows {
				for _, existingRule := range existingRulesets {
					if existingRule.SetName == targetSetName {
						existingRule.ReferenceType = calltype
						sourceRule.Rules[0].RuleActions.References = append(sourceRule.Rules[0].RuleActions.References, existingRule)

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
						for _, t := range rule.Rules {
							fmt.Println("\t\t\t\tRulePatterns:", t.RulePatterns)
							fmt.Println("\t\t\t\tRuleActions:", t.RuleActions)
							fmt.Println("\t\t\t\tNMatched:", t.NMatched)
							fmt.Println("\t\t\t\tNFailed:", rule.Rules[0].NFailed)

							// Print References if available
							for _, reference := range rule.Rules[0].RuleActions.References {
								fmt.Println("\t\t\t\t\tReferenced Rule:")
								fmt.Println("\t\t\t\t\t\tRulePatterns:", reference.Rules[0].RulePatterns)
								fmt.Println("\t\t\t\t\t\tRuleActions:", reference.Rules[0].RuleActions)
								fmt.Println("\t\t\t\t\t\tNMatched:", reference.Rules[0].NMatched)
								fmt.Println("\t\t\t\t\t\tNFailed:", reference.Rules[0].NFailed)
							}
						}
					}
				}

				// Print Workflows
				for className, workflows := range perSlice.Workflows {
					fmt.Println("\t\t\tWorkflows - Class:", className)
					for _, workflow := range workflows {
						fmt.Println("\t\t\t\tRulePatterns:", workflow.Rules[0].RulePatterns)
						fmt.Println("\t\t\t\tRuleActions:", workflow.Rules[0].RuleActions)
						fmt.Println("\t\t\t\tNMatched:", workflow.Rules[0].NMatched)
						fmt.Println("\t\t\t\tNFailed:", workflow.Rules[0].NFailed)

						// Print References if available
						for _, reference := range workflow.Rules[0].RuleActions.References {
							fmt.Println("\t\t\t\t\tReferenced Rule:")
							fmt.Println("\t\t\t\t\t\tRulePatterns:", reference.Rules[0].RulePatterns)
							fmt.Println("\t\t\t\t\t\tRuleActions:", reference.Rules[0].RuleActions)
							fmt.Println("\t\t\t\t\t\tNMatched:", reference.Rules[0].NMatched)
							fmt.Println("\t\t\t\t\t\tNFailed:", reference.Rules[0].NFailed)
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
func retrieveSchemasFromCacheByte(realm string, app string, class string, slice int, brwf string) ([]byte, []byte, string) {
	realmKey := realm_t(realm)
	perRealm, realmExists := schemaCache[realmKey]
	if !realmExists {
		return nil, nil, "Realmkey not match"
	}

	appKey := app_t(app)
	perApp, appExists := perRealm[appKey]
	if !appExists {
		return nil, nil, "AppKey not match"
	}

	sliceKey := slice_t(slice)
	perSlice, sliceExists := perApp[sliceKey]
	if !sliceExists {
		return nil, nil, "Slice key not match"
	}

	classNameKey := className_t(class)
	var schemas []*schema_t

	if brwf == "B" {
		schemas = perSlice.BRSchema[classNameKey]
	} else if brwf == "W" {
		schemas = perSlice.WFSchema[classNameKey]
	}

	if len(schemas) == 0 {
		return nil, nil, "No schemas found for the given class"
	}

	if len(schemas) == 0 {
		return nil, nil, "No schemas found for the given class"
	}

	patternSchemaJSON, err := json.Marshal(schemas[0].PatternSchema)
	if err != nil {
		return nil, nil, "JSON failed to marshal pattern"
	}

	actionSchemaJSON, err := json.Marshal(schemas[0].ActionSchema)
	if err != nil {
		return nil, nil, "JSON failed to marshal action"
	}

	return patternSchemaJSON, actionSchemaJSON, "success"
}

func retrieveRulesetFromCacheByte(realm string, app string, class string, slice int,
	brwf string) ([]byte, []byte, string, []*Ruleset_t) {
	realmKey := realm_t(realm)
	perRealm, exists := rulesetCache[realmKey]
	if !exists {
		return nil, nil, "Realmkey not match", nil
	}

	appKey := app_t(app)
	perApp, exists := perRealm[appKey]
	if !exists {
		return nil, nil, "AppKey not match", nil
	}

	sliceKey := slice_t(slice)
	perSlice, exists := perApp[sliceKey]
	if !exists {
		return nil, nil, "Slice key not match", nil
	}

	classNameKey := className_t(class)
	brwfKey := BrwfEnum(brwf)
	var rulesets []*Ruleset_t

	if brwfKey == "B" {
		rulesets = perSlice.BRRulesets[classNameKey]
	} else {
		rulesets = perSlice.Workflows[classNameKey]
	}

	if len(rulesets) == 0 {
		return nil, nil, "No rulesets found for the given class", nil
	}

	ruleset := rulesets[0]

	RuleActions, err := json.Marshal(ruleset.Rules[0].RuleActions)
	if err != nil {

		return nil, nil, "JSON failed to marshal rule actions", nil
	}

	RulePatterns, err := json.Marshal(ruleset.Rules[0].RulePatterns)
	if err != nil {

		return nil, nil, "JSON failed to marshal rule patterns", nil
	}

	return RulePatterns, RuleActions, "success", ruleset.Rules[0].RuleActions.References

}

func retrieveRuleSchemasFromCache(realm string, app string, class string, slice int) ([]*schema_t, error) {
	realmKey := realm_t(realm)

	perRealm, realmExists := schemaCache[realmKey]
	if !realmExists {

		return nil, errors.New("schema Realmkey not match")
	}

	appKey := app_t(app)
	perApp, appExists := perRealm[appKey]
	if !appExists {

		return nil, errors.New("schema AppKey not match")
	}

	sliceKey := slice_t(slice)

	perSlice, sliceExists := perApp[sliceKey]
	if !sliceExists {

		return nil, errors.New("schema Slice key not match")
	}

	var ruleSchemas []*schema_t

	brSchemas, brExists := perSlice.BRSchema[className_t(class)]
	if brExists {
		for _, schemas := range brSchemas {

			ruleSchemas = append(ruleSchemas, schemas)
		}
	}

	wfSchemas, wfExists := perSlice.WFSchema[className_t(class)]
	if wfExists {
		for _, schemas := range wfSchemas {
			ruleSchemas = append(ruleSchemas, schemas)
		}
	}

	return ruleSchemas, nil
}
func convertAttrValue(entityAttrVal string, valType valType_t) any {

	var entityAttrValConv any
	var err error
	switch valType {
	case valBool_t:
		entityAttrValConv, err = strconv.ParseBool(entityAttrVal)
	case valInt_t:
		entityAttrValConv, err = strconv.Atoi(entityAttrVal)
	case valFloat_t:
		entityAttrValConv, err = strconv.ParseFloat(entityAttrVal, 64)
	case valString_t, valEnum_t:
		entityAttrValConv = entityAttrVal
	case valTimestamp_t:
		entityAttrValConv, err = time.Parse(timeLayout, entityAttrVal)
	}
	if err != nil {
		return err
	}
	return entityAttrValConv
}

func retrieveRuleSetsFromCache(realm string, app string, class string, slice int) ([]*Ruleset_t, error) {
	realmKey := realm_t(realm)

	perRealm, realmExists := rulesetCache[realmKey]
	if !realmExists {
		return nil, errors.New("ruleset realmkey not match")
	}

	appKey := app_t(app)
	perApp, appExists := perRealm[appKey]
	if !appExists {
		return nil, errors.New("ruleset appKey not match")
	}

	sliceKey := slice_t(slice)
	perSlice, sliceExists := perApp[sliceKey]
	if !sliceExists {
		return nil, errors.New("ruleset slice key not match")
	}

	var ruleSets []*Ruleset_t

	for _, brRulesets := range perSlice.BRRulesets {
		ruleSets = append(ruleSets, brRulesets...)
	}
	for _, wfRulesets := range perSlice.Workflows {
		ruleSets = append(ruleSets, wfRulesets...)
	}

	return ruleSets, nil
}

func retriveRuleSchemasAndRuleSetsFromCache(realm string, app string, class string, slice string) ([]*schema_t, []*Ruleset_t) {
	s, _ := strconv.Atoi(slice)

	ruleSchemas, _ := retrieveRuleSchemasFromCache(realm, app, class, s)

	ruleSets, _ := retrieveRuleSetsFromCache(realm, app, class, s)
	return ruleSchemas, ruleSets

}
func printStats(statsData rulesetStats_t) {
	for realm, perRealm := range statsData {
		for app, perApp := range perRealm {
			for slice, perSlice := range perApp {
				fmt.Printf("Realm: %v, App: %v, Slice: %v\n", realm, app, slice)
				fmt.Printf("loadedAt: %v\n", perSlice.loadedAt)

				// Print stats for BRSchema
				for className, schema := range perSlice.BRSchema {
					fmt.Printf("Class: %v, nChecked: %v\n", className, schema.nChecked)
				}

				// Print stats for BRRulesets
				for className, rulesets := range perSlice.BRRulesets {
					for _, ruleset := range rulesets {
						fmt.Printf("Class: %v, nCalled: %v\n", className, ruleset.nCalled)
						for _, rule := range ruleset.rulesStats {
							fmt.Printf("nMatched: %v, nFailed: %v\n", rule.nMatched, rule.nFailed)
						}
					}
				}

				// Print stats for WFSchema
				for className, schema := range perSlice.WFSchema {
					fmt.Printf("Class: %v, nChecked: %v\n", className, schema.nChecked)
				}

				// Print stats for Workflows
				for className, workflows := range perSlice.Workflows {
					for _, workflow := range workflows {
						fmt.Printf("Class: %v, nCalled: %v\n", className, workflow.nCalled)
						for _, rule := range workflow.rulesStats {
							fmt.Printf("nMatched: %v, nFailed: %v\n", rule.nMatched, rule.nFailed)
						}
					}
				}
			}
		}
	}
}
