package crux

import (
	"fmt"
)

// func NewProvider(cfg string) sqlc.Querier {
// 	ctx := context.Background()
// 	db, err := pgxpool.New(ctx, cfg)
// 	if err != nil {
// 		log.Fatal("error connecting db")
// 	}
// 	err = db.Ping(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	log.Println("Successfully connected to the database")
// 	return sqlc.NewQuerierWithTX(db)
// }

// func AddReferencesToRuleSetCache() {
// 	for realmKey, perRealm := range RulesetCache {
// 		for _, perApp := range perRealm {
// 			for sliceKey, perSlice := range perApp {
// 				for _, rulesets := range perSlice.BRRulesets {
// 					for _, rule := range rulesets {
// 						for idx := range rule.Rules {
// 							processSubRule(&rule.Rules[idx], realmKey, sliceKey)
// 						}
// 					}
// 				}
// 				for _, rulesets := range perSlice.Workflows {
// 					for _, rule := range rulesets {
// 						for idx := range rule.Rules {
// 							processSubRule(&rule.Rules[idx], realmKey, sliceKey)
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }

// func processSubRule(subRule *Rule_t, realmKey Realm_t, sliceKey Slice_t) {
// 	if subRule.RuleActions.ThenCall != "" {
// 		referRuleset := searchAndAddReferences(subRule.RuleActions.ThenCall, RulesetCache, realmKey, sliceKey, "thencall")
// 		if referRuleset != nil {
// 			subRule.RuleActions.References = append(subRule.RuleActions.References, referRuleset)
// 		}
// 	}
// 	if subRule.RuleActions.ElseCall != "" {
// 		referRuleset := searchAndAddReferences(subRule.RuleActions.ElseCall, RulesetCache, realmKey, sliceKey, "elsecall")
// 		if referRuleset != nil {
// 			subRule.RuleActions.References = append(subRule.RuleActions.References, referRuleset)
// 		}
// 	}
// }

// func searchAndAddReferences(targetSetName string, cache map[Realm_t]PerRealm_t, realmKey Realm_t,
// 	sliceKey Slice_t, calltype string) *Ruleset_t {
// 	for _, perApp := range cache[realmKey] {
// 		for otherSliceKey, perSlice := range perApp {
// 			if otherSliceKey == sliceKey {
// 				continue
// 			}
// 			for _, existingRulesets := range perSlice.BRRulesets {
// 				for _, existingRule := range existingRulesets {
// 					if existingRule.SetName == targetSetName {
// 						existingRule.ReferenceType = calltype
// 						return existingRule
// 					}
// 				}
// 			}
// 			for _, existingRulesets := range perSlice.Workflows {
// 				for _, existingRule := range existingRulesets {
// 					if existingRule.SetName == targetSetName {
// 						existingRule.ReferenceType = calltype
// 						return existingRule
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }

// func PrintAllRuleSetCache() {

// 	for realmKey, perRealm := range RulesetCache {
// 		fmt.Println("Realm:", realmKey)
// 		for appKey, perApp := range perRealm {
// 			fmt.Println("\tApp:", appKey)
// 			for sliceKey, perSlice := range perApp {
// 				fmt.Println("\t\tSlice:", sliceKey)
// 				fmt.Println("\t\t\tLoadedAt:", perSlice.LoadedAt)

// 				// Print BRRulesets

// 				for className, BRRulesets := range perSlice.BRRulesets {
// 					fmt.Println("\t\t\tBRRulesets - Class:", className)
// 					for _, rule := range BRRulesets {
// 						for _, t := range rule.Rules {
// 							fmt.Println("\t\t\t\tRulePatterns:", t.RulePatterns)
// 							fmt.Println("\t\t\t\tRuleActions:", t.RuleActions)
// 							fmt.Println("\t\t\t\tNMatched:", t.NMatched)
// 							fmt.Println("\t\t\t\tNFailed:", t.NFailed)

// 							for _, refrule := range t.RuleActions.References {
// 								for _, z := range refrule.Rules {
// 									fmt.Println("\t\t\t\t\tReferenced Rule:")
// 									fmt.Println("\t\t\t\t\t\tRulePatterns:", z.RulePatterns)
// 									fmt.Println("\t\t\t\t\t\tRuleActions:", z.RuleActions)
// 									fmt.Println("\t\t\t\t\t\tNMatched:", z.NMatched)
// 									fmt.Println("\t\t\t\t\t\tNFailed:", z.NFailed)
// 								}
// 							}
// 						}
// 					}
// 				}

// 				// Print Workflows
// 				for className, workflows := range perSlice.Workflows {
// 					fmt.Println("\t\t\tWorkflows - Class:", className)
// 					for _, workflow := range workflows {
// 						for _, t := range workflow.Rules {
// 							fmt.Println("\t\t\t\tRulePatterns:", t.RulePatterns)
// 							fmt.Println("\t\t\t\tRuleActions:", t.RuleActions)
// 							fmt.Println("\t\t\t\tNMatched:", t.NMatched)
// 							fmt.Println("\t\t\t\tNFailed:", t.NFailed)

// 							for _, refrule := range t.RuleActions.References {
// 								for _, z := range refrule.Rules {
// 									fmt.Println("\t\t\t\t\tReferenced Rule:")
// 									fmt.Println("\t\t\t\t\t\tRulePatterns:", z.RulePatterns)
// 									fmt.Println("\t\t\t\t\t\tRuleActions:", z.RuleActions)
// 									fmt.Println("\t\t\t\t\t\tNMatched:", z.NMatched)
// 									fmt.Println("\t\t\t\t\t\tNFailed:", z.NFailed)
// 								}
// 							}
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }
// func PrintAllSchemaCache() {
// 	for realmKey, perRealm := range SchemaCache {
// 		fmt.Println("Realm:", realmKey)
// 		for appKey, perApp := range perRealm {
// 			fmt.Println("\tApp:", appKey)
// 			for sliceKey, perSlice := range perApp {
// 				fmt.Println("\t\tSlice:", sliceKey)
// 				fmt.Println("\t\t\tLoadedAt:", perSlice.LoadedAt)
// 				for className, schema := range perSlice.BRSchema {
// 					fmt.Println("\t\t\tBRSchema - Class:", className)
// 					//for _, schema := range schemas {
// 					fmt.Println("\t\t\t\tPatternSchema:", schema.PatternSchema)
// 					fmt.Println("\t\t\t\tActionSchema:", schema.ActionSchema)
// 					fmt.Println("\t\t\t\tNChecked:", schema.NChecked)
// 					//}
// 				}
// 				for className, schema := range perSlice.WFSchema {
// 					fmt.Println("\t\t\tWFSchema - Class:", className)
// 					//for _, schema := range schemas {
// 					fmt.Println("\t\t\t\tPatternSchema:", schema.PatternSchema)
// 					fmt.Println("\t\t\t\tActionSchema:", schema.ActionSchema)
// 					fmt.Println("\t\t\t\tNChecked:", schema.NChecked)
// 					//}
// 				}

// 			}
// 		}
// 	}
// }

// func containsField(value interface{}, fieldName string, t *testing.T) bool {

// 	switch v := value.(type) {

// 	case []byte:

// 		var raw json.RawMessage
// 		if err := json.Unmarshal(v, &raw); err != nil {
// 			fmt.Println("Error unmarshalling actual pattern:", err, v)
// 			return false
// 		}

// 		var data map[string]interface{}
// 		if err := json.Unmarshal(raw, &data); err != nil {
// 			var arrayData []interface{}
// 			if err := json.Unmarshal(raw, &arrayData); err != nil {
// 				fmt.Println("Error unmarshalling actual pattern:", err, v)
// 				return false
// 			}
// 			for _, element := range arrayData {
// 				if containsFieldName(element, fieldName) {
// 					return true
// 				}
// 			}
// 		}
// 		for _, value := range data {
// 			if containsFieldName(value, fieldName) {
// 				return true
// 			}
// 		}
// 	case map[string]interface{}:

// 		for key := range v {
// 			if key == fieldName {
// 				return true
// 			}
// 		}

// 	case []interface{}:
// 		for _, item := range v {
// 			if containsField(item, fieldName, t) {
// 				return true
// 			}
// 		}
// 	case string:
// 		return v == fieldName
// 	}
// 	return false
// }

// func containsFieldName(value interface{}, fieldName string) bool {

// 	v := reflect.ValueOf(value)

// 	switch v.Kind() {
// 	case reflect.Map:
// 		for _, key := range v.MapKeys() {
// 			if key.Kind() == reflect.String && key.String() == fieldName {
// 				return true
// 			}
// 			if nestedValue := v.MapIndex(key).Interface(); containsFieldName(nestedValue, fieldName) {
// 				return true
// 			}
// 		}

// 	case reflect.Slice:
// 		for i := 0; i < v.Len(); i++ {
// 			if nestedValue := v.Index(i).Interface(); containsFieldName(nestedValue, fieldName) {
// 				return true
// 			}
// 		}
// 	case reflect.String:
// 		return value.(string) == fieldName
// 	}
// 	return false
// }

// func convertAttrValue(entityAttrVal string, valType ValType_t) any {

// 	var entityAttrValConv any
// 	var err error
// 	switch valType {
// 	case ValBool_t:
// 		entityAttrValConv, err = strconv.ParseBool(entityAttrVal)
// 	case ValInt_t:
// 		entityAttrValConv, err = strconv.Atoi(entityAttrVal)
// 	case ValFloat_t:
// 		entityAttrValConv, err = strconv.ParseFloat(entityAttrVal, 64)
// 	case ValString_t, ValEnum_t:
// 		entityAttrValConv = entityAttrVal
// 	case ValTimestamp_t:
// 		entityAttrValConv, err = time.Parse(timeLayout, entityAttrVal)
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	return entityAttrValConv
// }

// func printStats(statsData rulesetStats_t) {
// 	for realm, perRealm := range statsData {
// 		for app, perApp := range perRealm {
// 			for slice, perSlice := range perApp {
// 				fmt.Printf("Realm: %v, App: %v, Slice: %v\n", realm, app, slice)
// 				fmt.Printf("loadedAt: %v\n", perSlice.loadedAt)
// 				// Print stats for BRSchema
// 				for className, schema := range perSlice.BRSchema {
// 					fmt.Printf("Class: %v, nChecked: %v\n", className, schema.nChecked)
// 				}
// 				// Print stats for BRRulesets
// 				for className, rulesets := range perSlice.BRRulesets {
// 					for _, ruleset := range rulesets {
// 						fmt.Printf("Class: %v, nCalled: %v\n", className, ruleset.nCalled)
// 						for _, rule := range ruleset.rulesStats {
// 							fmt.Printf("nMatched: %v, nFailed: %v\n", rule.nMatched, rule.nFailed)
// 						}
// 					}
// 				}
// 				// Print stats for WFSchema
// 				for className, schema := range perSlice.WFSchema {
// 					fmt.Printf("Class: %v, nChecked: %v\n", className, schema.nChecked)
// 				}
// 				// Print stats for Workflows
// 				for className, workflows := range perSlice.Workflows {
// 					for _, workflow := range workflows {
// 						fmt.Printf("Class: %v, nCalled: %v\n", className, workflow.nCalled)
// 						for _, rule := range workflow.rulesStats {
// 							fmt.Printf("nMatched: %v, nFailed: %v\n", rule.nMatched, rule.nFailed)
// 						}
// 					}
// 				}
// 			}
// 		}
// 	}
// }

func lockCache() {
	cacheLock.Lock()
}

func unlockCache() {
	cacheLock.Unlock()
}

func (c Cache) RetriveRuleSchemasAndRuleSetsFromCache(brwf string) (*Schema_t, *Ruleset_t, error) {

	ruleSchemas, err := c.RetrieveRuleSchemasFromCache(brwf)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieveRuleSchemasFromCache: %v", err)
	}

	ruleSets, err := c.RetrieveWorkflowRuleSetFromCache(brwf)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to RetrieveRuleSetsFromCache: %v", err)
	}

	return ruleSchemas, ruleSets, nil
}

func (c Cache) RetrieveRuleSchemasFromCache(brwf string) (*Schema_t, error) {
	if brwf == "B" {
		brSchemas, brExists := SchemaCache[c.Realm][c.App][c.Slice].BRSchema[ClassName_t(c.Class)]
		if brExists {
			return &brSchemas, nil
		}
		if !brExists {
			if err := c.Load(); err != nil {
				return nil, fmt.Errorf("error while loading cache in RetrieveWorkflowRulesetFromCache: %v", err)
			} else {
				brSchemas, brExists := SchemaCache[c.Realm][c.App][c.Slice].BRSchema[ClassName_t(c.Class)]
				if brExists {
					return &brSchemas, nil
				} else {
					return nil, fmt.Errorf("no brschema found")
				}
			}
		}
	} else if brwf == "W" {
		wfSchemas, wfExists := SchemaCache[c.Realm][c.App][c.Slice].WFSchema[ClassName_t(c.Class)]
		if wfExists {
			return &wfSchemas, nil
		}
		if !wfExists {
			if err := c.Load(); err != nil {
				return nil, fmt.Errorf("error while loading cache in RetrieveWorkflowRulesetFromCache: %v", err)
			} else {
				wfSchemas, wfExists := SchemaCache[c.Realm][c.App][c.Slice].WFSchema[ClassName_t(c.Class)]
				if wfExists {
					return &wfSchemas, nil
				} else {
					return nil, fmt.Errorf("no wfschema found")
				}
			}
		}
	}
	return nil, fmt.Errorf("no schema found")
}

func (c Cache) RetrieveWorkflowRuleSetFromCache(brwf string) (*Ruleset_t, error) {

	ruleSets, exists := c.GetRulesetsFromCacheWithName(brwf)
	if exists {
		return ruleSets, nil
	} else {
		if err := c.Load(); err != nil {
			return nil, fmt.Errorf("error while loading cache in RetrieveWorkflowRulesetFromCache: %v", err)
		} else {
			ruleSets, exists := c.GetRulesetsFromCacheWithName(brwf)
			if exists {
				return ruleSets, nil
			} else {
				return nil, fmt.Errorf("rule set not exist for given specification")
			}
		}
	}
}

func (c Cache) GetRulesetsFromCacheWithName(brwf string) (*Ruleset_t, bool) {
	if brwf == "B" {
		brRulesets, exist := RulesetCache[c.Realm][c.App][c.Slice].Workflows[c.Class]
		if exist {
			for _, r := range brRulesets {
				if r.SetName == c.WorkflowName {
					return r, true
				}
			}
		}
	} else if brwf == "W" {
		workflows, exist := RulesetCache[c.Realm][c.App][c.Slice].Workflows[c.Class]
		if exist {
			for _, w := range workflows {
				if w.SetName == c.WorkflowName {
					return w, true
				}
			}
		}
	}
	return nil, false
}
