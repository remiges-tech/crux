package main

import (
	"context"
	sqlc "crux/db/sqlc-gen"
	"encoding/json"
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
func retrieveSchemasFromCache(realm int, app string, class string, slice int, brwf string) ([]byte, []byte, string) {
	realmKey := realm_t(strconv.Itoa(realm))
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

func retrieveRulesetFromCache(realm int, app string, class string, slice int,
	brwf string) ([]byte, []byte, string, []*Ruleset_t) {
	realmKey := realm_t(strconv.Itoa(realm))
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

	RulePatterns, err := json.Marshal(ruleset.RulePatterns)
	if err != nil {
		return nil, nil, "JSON failed to marshal rule patterns", nil
	}

	RuleActions, err := json.Marshal(ruleset.RuleActions)
	if err != nil {
		return nil, nil, "JSON failed to marshal rule actions", nil
	}

	return RulePatterns, RuleActions, "success", ruleset.RuleActions.References
}
func initializeRuleSchemasFromCache(schemaCache schemaCache_t) {
	// Reset the global variable
	ruleSchemas = nil

	// Iterate over schemaCache
	for realmKey, perRealm := range schemaCache {
		for appKey, perApp := range perRealm {
			for sliceKey, perSlice := range perApp {
				// Process BRSchema
				for _, schemas := range perSlice.BRSchema {
					processRuleSchemas(realmKey, appKey, sliceKey, schemas)
				}

				// Process WFSchema
				for _, schemas := range perSlice.WFSchema {
					processRuleSchemas(realmKey, appKey, sliceKey, schemas)
				}
			}
		}
	}
}

func processRuleSchemas(realmKey realm_t, appKey app_t, sliceKey slice_t, schemas []*schema_t) {
	for _, schema := range schemas {
		ruleSchema := RuleSchema{
			class:         schema.Class,
			patternSchema: processPatternSchema(schema.PatternSchema),
			actionSchema:  processActionSchema(schema.ActionSchema),
		}

		ruleSchemas = append(ruleSchemas, ruleSchema)

	}
}

func processPatternSchema(patternSchema patternSchema_t) []AttrSchema {
	var attrSchemas []AttrSchema

	for _, attr := range patternSchema.Attr {
		attrSchema := AttrSchema{
			name:    attr.Name,
			valType: attr.ValType,
			vals:    make(map[string]bool),
			valMin:  attr.ValMin,
			valMax:  attr.ValMax,
			lenMin:  attr.LenMin,
			lenMax:  attr.LenMax,
		}

		for _, val := range attr.EnumVals {
			attrSchema.vals[val] = true
		}

		attrSchemas = append(attrSchemas, attrSchema)
	}
	return attrSchemas
}

func processActionSchema(actionSchema actionSchema_t) ActionSchema {
	return ActionSchema{
		tasks:      actionSchema.Tasks,
		properties: actionSchema.Properties,
	}
}

func initializeRuleSetsFromCache(rulesetCache rulesetCache_t) {

	ruleSets = make(map[string]RuleSet)

	for realmKey, perRealm := range rulesetCache {
		for appKey, perApp := range perRealm {
			for sliceKey, perSlice := range perApp {
				for className, brRulesets := range perSlice.BRRulesets {
					processRuleSets(realmKey, appKey, sliceKey, className, brRulesets)
				}
				for className, wfRulesets := range perSlice.Workflows {
					processRuleSets(realmKey, appKey, sliceKey, className, wfRulesets)
				}
			}
		}
	}
}

func processRuleSets(realmKey realm_t, appKey app_t, sliceKey slice_t, className className_t, rulesets []*Ruleset_t) {
	for _, rule := range rulesets {
		newRule := Rule{
			rulePattern: processRulePatterns(rule.RulePatterns),
			ruleActions: processRuleActions(rule.RuleActions),
		}

		ruleSets[rule.SetName] = RuleSet{
			class:   rule.Class,
			setName: rule.SetName,
			rules:   append(ruleSets[rule.SetName].rules, newRule),
		}
	}
}

func processRulePatterns(rulePatterns []rulePatternBlock_t) []RulePatternTerm {
	var patternTerms []RulePatternTerm

	for _, pattern := range rulePatterns {
		patternTerm := RulePatternTerm{
			attrName: pattern.Attr,
			op:       pattern.Op,
			attrVal:  convertAttrValue(pattern.Val, pattern.ValType),
		}
		patternTerms = append(patternTerms, patternTerm)
	}

	return patternTerms
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

func processRuleActions(ruleActions ruleActionBlock_t) RuleActions {
	return RuleActions{
		tasks:      ruleActions.Task,
		properties: ruleActions.Properties,
		thenCall:   ruleActions.ThenCall,
		elseCall:   ruleActions.ElseCall,
		willReturn: ruleActions.DoReturn,
		willExit:   ruleActions.DoExit,
	}
}
