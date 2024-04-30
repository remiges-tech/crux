package crux

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/remiges-tech/crux/db/sqlc-gen"
)

const (
	BRE = "B"
	WFE = "W"
)

type Cache struct {
	Ctx          context.Context
	Query        *sqlc.Queries
	RulesetCache RulesetCache_t
	SchemaCache  SchemaCache_t
}

func lockCache() {
	cacheLock.Lock()
}

func unlockCache() {
	cacheLock.Unlock()
}

/*
NewCache creates and returns an empty cache instance. The purpose of this function is to ensure that there is
always a single instance of the cache throughout the application's lifecycle.
*/
func NewCache(c context.Context, queries *sqlc.Queries) *Cache {
	return &Cache{
		Ctx:          c,
		Query:        queries,
		RulesetCache: make(RulesetCache_t),
		SchemaCache:  make(SchemaCache_t),
	}
}
func (c Cache) Load(realm, app, class, ruleSetName string, slice int32) error {
	err := c.loadSchema(realm, app, class, slice)
	if err != nil {
		return err
	}

	err = c.loadRuleSet(realm, app, class, ruleSetName, slice)
	if err != nil {
		return err
	}
	return nil
}
func (c Cache) loadSchema(realm, app, class string, slice int32) error {
	lockCache()
	defer unlockCache()

	dbResponseSchema, err := c.Query.LoadSchema(context.Background(), sqlc.LoadSchemaParams{
		RealmName: realm,
		Slice:     slice,
		Class:     class,
		App:       app,
	})
	if err != nil {
		return err
	}
	if len(dbResponseSchema) == 0 {
		return fmt.Errorf("schema not found")
	}

	err = c.loadInternalSchema(dbResponseSchema)
	if err != nil {
		return err
	}
	return nil
}

func (c Cache) loadRuleSet(realm, app, class, ruleSetName string, slice int32) error {

	lockCache()
	defer unlockCache()

	dbResponseRuleSet, err := c.Query.LoadRuleSet(c.Ctx,
		sqlc.LoadRuleSetParams{
			RealmName: realm,
			Slice:     slice,
			Class:     class,
			App:       app,
			Setname:   ruleSetName,
		})
	if err != nil {
		return err
	}
	if IsZeroOfUnderlyingType(dbResponseRuleSet) {
		return fmt.Errorf("rule set not found")
	}
	err = c.loadInternalRuleSet(dbResponseRuleSet)
	if err != nil {
		return err
	}
	return nil
}

func (c Cache) loadInternal(dbResponseSchema []sqlc.Schema, dbResponseRuleSet sqlc.Ruleset) error {

	err := c.loadInternalSchema(dbResponseSchema)
	if err != nil {
		return err
	}

	err = c.loadInternalRuleSet(dbResponseRuleSet)
	if err != nil {
		return err
	}
	return nil
}

func (c Cache) loadInternalSchema(dbResponseSchema []sqlc.Schema) error {

	if len(dbResponseSchema) == 0 {
		return fmt.Errorf("dbResponseRuleSet is empty")
	}

	for _, row := range dbResponseSchema {
		realmKey := Realm_t(string(row.Realm))
		perRealm, exists := c.SchemaCache[realmKey]
		if !exists {
			perRealm = make(PerRealm_t)
			c.SchemaCache[realmKey] = perRealm
		}

		appKey := App_t(row.App)
		perApp, exists := perRealm[appKey]
		if !exists {
			perApp = make(PerApp_t)
			perRealm[appKey] = perApp
		}

		sliceKey := Slice_t(row.Slice)
		_, exists = perApp[sliceKey]
		if !exists {
			perApp[sliceKey] = PerSlice_t{
				LoadedAt: time.Now(),
				BRSchema: make(map[ClassName_t]Schema_t),
				WFSchema: make(map[ClassName_t]Schema_t),
			}
		}
		var patterns []PatternSchema_t

		if err := json.Unmarshal(row.Patternschema, &patterns); err != nil {
			log.Println("Error unmarshaling Patternschema:", err)
			continue
		}

		var actions ActionSchema_t

		if err := json.Unmarshal(row.Actionschema, &actions); err != nil {
			log.Println("Error parsing ActionSchema JSON:", err)
			continue
		}

		schemaData := Schema_t{
			Class:         row.Class,
			PatternSchema: patterns,
			ActionSchema:  actions,
		}

		classNameKey := ClassName_t(row.Class)
		if row.Brwf == BRE {
			perApp[sliceKey].BRSchema[classNameKey] = schemaData
		} else if row.Brwf == WFE {
			perApp[sliceKey].WFSchema[classNameKey] = schemaData
		}

	}
	return nil
}

func (c Cache) loadInternalRuleSet(dbResponseRuleSet sqlc.Ruleset) error {

	if IsZeroOfUnderlyingType(dbResponseRuleSet) {
		return fmt.Errorf("didn't get rule set")
	}

	realmKey := Realm_t(string(dbResponseRuleSet.Realm))
	perRealm, exists := c.RulesetCache[realmKey]
	if !exists {
		perRealm = make(PerRealm_t)
		c.RulesetCache[realmKey] = perRealm
	}

	appKey := App_t(dbResponseRuleSet.App)
	perApp, exists := perRealm[appKey]
	if !exists {
		perApp = make(PerApp_t)
		perRealm[appKey] = perApp
	}

	sliceKey := Slice_t(dbResponseRuleSet.Slice)
	_, exists = perApp[sliceKey]
	if !exists {
		perApp[sliceKey] = PerSlice_t{
			LoadedAt:   time.Now(),
			BRRulesets: make(map[ClassName_t][]*Ruleset_t),
			Workflows:  make(map[ClassName_t][]*Ruleset_t),
		}
	}

	var rules []Rule_t
	err := json.Unmarshal(dbResponseRuleSet.Ruleset, &rules)
	if err != nil {
		fmt.Println("Error unmarshaling rules:", err)
		return nil
	}

	classNameKey := ClassName_t(dbResponseRuleSet.Class)
	newRuleset := &Ruleset_t{
		Id:       dbResponseRuleSet.ID,
		Class:    dbResponseRuleSet.Class,
		SetName:  dbResponseRuleSet.Setname,
		Rules:    rules,
		IsActive: dbResponseRuleSet.IsActive.Bool,
	}
	if dbResponseRuleSet.Brwf == BRE {
		perApp[sliceKey].BRRulesets[classNameKey] = append(perApp[sliceKey].BRRulesets[classNameKey], newRuleset)

	} else if dbResponseRuleSet.Brwf == WFE {
		fmt.Printf("className: %v", classNameKey)
		fmt.Printf("workflow name: %v", newRuleset.SetName)
		perApp[sliceKey].Workflows[classNameKey] = append(perApp[sliceKey].Workflows[classNameKey], newRuleset)
	}

	// AddReferencesToRuleSetCache()
	return nil
}

func (c Cache) Purge(brwf, app, realm, class ,rulesetName ,field string, slice int32) {
	lockCache()
	defer unlockCache()

	if brwf == "B" && field == "schema" {
		c.SchemaCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].BRSchema[ClassName_t(class)] = Schema_t{}
		c.RulesetCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].BRRulesets[ClassName_t(class)] = nil
	} else if brwf == "B" && field == "rule" {
		c.RulesetCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].BRRulesets[ClassName_t(class)] = nil
	} else if brwf == "w" && field == "schema" {
		c.SchemaCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].WFSchema[ClassName_t(class)] = Schema_t{}
		c.RulesetCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].Workflows[ClassName_t(class)] = nil
	} else if brwf == "W" && field == "rule" {
		c.RulesetCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].Workflows[ClassName_t(class)] = nil
	}else if brwf == "B" && field == "ruleset"{
		ruleset := c.RulesetCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].BRRulesets[ClassName_t(class)]
		 for _,r := range ruleset{
			r.SetName = rulesetName
			c.RulesetCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].BRRulesets[ClassName_t(class)]= nil
		 }
	} 
}

// func Reload(ctx context.Context, query sqlc.Querier, slice int32, app, class, realm, workflowName string) error {
// 	lockCache()
// 	defer unlockCache()
// 	if err := purgeInternal(); err != nil {
// 		log.Fatal("Failed to purge data from cache:", err)
// 		return err
// 	}
// 	// dbResponseSchema, err := query.AllSchemas(ctx)
// 	dbResponseSchema, err := query.LoadSchema(ctx, sqlc.LoadSchemaParams{
// 		RealmName: realm,
// 		Slice:     slice,
// 		Class:     class,
// 		App:       app,
// 	})
// 	if err != nil {
// 		return err
// 	}
// 	if len(dbResponseSchema) == 0 {
// 		return fmt.Errorf("Didn't get schema")
// 	}
// 	dbResponseRuleSet, err := query.LoadRuleSet(ctx,
// 		sqlc.LoadRuleSetParams{
// 			Slice:     slice,
// 			App:       app,
// 			Class:     class,
// 			Setname:   workflowName,
// 			RealmName: realm,
// 		})
// 	if err != nil {
// 		return err
// 	}
// 	// if IsZeroOfUnderlyingType(dbResponseRuleSet) {
// 	// 	return fmt.Errorf("Didn't get rule set")
// 	// }
// 	err = loadInternal(dbResponseSchema, dbResponseRuleSet)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (c Cache) RetriveRuleSchemasAndRuleSetsFromCache(brwf, app, realm, class, ruleSetName string, slice int32) (*Schema_t, *Ruleset_t, error) {

	ruleSchemas, err := c.RetrieveRuleSchemasFromCache(brwf, app, realm, class, slice)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to retrieveRuleSchemasFromCache: %v", err)
	}

	ruleSets, err := c.RetrieveWorkflowRuleSetFromCache(brwf, app, realm, class, ruleSetName, slice)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to RetrieveRuleSetsFromCache: %v", err)
	}

	return ruleSchemas, ruleSets, nil
}

func (c Cache) RetrieveRuleSchemasFromCache(brwf, app, realm, class string, slice int32) (*Schema_t, error) {
	if brwf == BRE {
		brSchemas, brExists := c.SchemaCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].BRSchema[ClassName_t(class)]
		if brExists {
			return &brSchemas, nil
		}
		if !brExists {
			if err := c.loadSchema(realm, app, class, slice); err != nil {
				return nil, fmt.Errorf("error while loading cache in RetrieveWorkflowRulesetFromCache: %v", err)
			} else {
				brSchemas, brExists := c.SchemaCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].BRSchema[ClassName_t(class)]
				if brExists {
					return &brSchemas, nil
				} else {
					return nil, fmt.Errorf("no brschema found")
				}
			}
		}
	} else if brwf == WFE {
		wfSchemas, wfExists := c.SchemaCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].WFSchema[ClassName_t(class)]
		if wfExists {
			return &wfSchemas, nil
		}
		if !wfExists {
			if err := c.loadSchema(realm, app, class, slice); err != nil {
				return nil, fmt.Errorf("error while loading cache in RetrieveWorkflowRulesetFromCache: %v", err)
			} else {
				wfSchemas, wfExists := c.SchemaCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].WFSchema[ClassName_t(class)]
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

func (c Cache) RetrieveWorkflowRuleSetFromCache(brwf, app, realm, class, ruleSetName string, slice int32) (*Ruleset_t, error) {

	ruleSets, exists := c.getRulesetsFromCacheWithName(brwf, app, realm, class, ruleSetName, slice)
	if exists {
		return ruleSets, nil
	} else {
		if err := c.loadRuleSet(realm, app, class, ruleSetName, slice); err != nil {
			return nil, fmt.Errorf("error while loading cache in RetrieveWorkflowRulesetFromCache: %v", err)
		} else {
			ruleSets, exists := c.getRulesetsFromCacheWithName(brwf, app, realm, class, ruleSetName, slice)
			if exists {
				return ruleSets, nil
			} else {
				return nil, fmt.Errorf("rule set not exist for given specification")
			}
		}
	}
}

func (c Cache) GetRulesetName(brwf, app, realm, class, ruleSetName string, slice int32) (*Ruleset_t, bool, error) {
	if brwf == BRE {
		brRulesets, exist := c.RulesetCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].Workflows[ClassName_t(class)]
		if exist {
			for _, r := range brRulesets {
				if r.SetName == ruleSetName {
					return r, true, nil
				}
			}
		} else {
			if err := c.loadRuleSet(realm, app, class, ruleSetName, slice); err != nil {
				return nil, false, fmt.Errorf("error while loading cache in GetRulesetName: %v", err)
			} else {
				ruleSets, exists := c.getRulesetsFromCacheWithName(brwf, app, realm, class, ruleSetName, slice)
				if exists {
					return ruleSets, true, nil
				} else {
					return nil, false, fmt.Errorf("rule set not exist for given specification")
				}
			}
		}
	} else if brwf == WFE {
		workflows, exist := c.RulesetCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].Workflows[ClassName_t(class)]
		if exist {
			for _, w := range workflows {
				if w.SetName == ruleSetName {
					return w, true, nil
				}
			}
		} else {
			if err := c.loadRuleSet(realm, app, class, ruleSetName, slice); err != nil {
				return nil, false, fmt.Errorf("error while loading cache in RetrieveWorkflowRulesetFromCache: %v", err)
			} else {
				ruleSets, exists := c.getRulesetsFromCacheWithName(brwf, app, realm, class, ruleSetName, slice)
				if exists {
					return ruleSets, true, nil
				} else {
					return nil, false, fmt.Errorf("rule set not exist for given specification")
				}
			}
		}
	}
	return nil, false, fmt.Errorf("no ruleset found ")
}

func (c Cache) getRulesetsFromCacheWithName(brwf, app, realm, class, ruleSetName string, slice int32) (*Ruleset_t, bool) {

	if brwf == BRE {
		brRulesets, exist := c.RulesetCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].Workflows[ClassName_t(class)]
		if exist {
			for _, r := range brRulesets {
				if r.SetName == ruleSetName {
					return r, true
				}
			}
		}
	} else if brwf == WFE {
		workflows, exist := c.RulesetCache[Realm_t(realm)][App_t(app)][Slice_t(slice)].Workflows[ClassName_t(class)]
		if exist {
			for _, w := range workflows {
				if w.SetName == ruleSetName {
					return w, true
				}
			}
		}
	}
	return nil, false
}

// This fuction is used to retrieve a specific ruleset from cache and checks whether it is active if not then make it active.
// It also checks whether it has thencall then it recursively retrive all child rulesets and make all of them active.
func (c Cache) RetrieveAndCheckIsActiveRuleSet(brwf, app, realm, class, ruleSetName string, slice int32) (ntraversed int, nactivated int, err error) {

	// Retrieve the rule set from cache
	currentRuleset, exists, err := c.GetRulesetName(brwf, app, realm, class, ruleSetName, slice)
	if err != nil {
		return 0, 0, err
	}
	if !exists {
		return 0, 0, fmt.Errorf("rule set not found in cache: %s", ruleSetName)
	}

	// Check if the current rule set is active
	if !currentRuleset.IsActive {
		currentRuleset.IsActive = true
		nactivated++
	}

	// Iterate over the rules in the current rule set
	for _, rule := range currentRuleset.Rules {
		// Check if the rule has a then call
		if rule.RuleActions.ThenCall != "" {
			// Recursively retrieve and check the next rule set
			if _, _, err := c.RetrieveAndCheckIsActiveRuleSet(brwf, app, realm, class, rule.RuleActions.ThenCall, slice); err != nil {
				return 0, 0, err
			} else {
				ntraversed++
			}
		}
	}

	return ntraversed, nactivated, nil
}

// This function is used to retrieve a specific ruleset from cache and make it inactive
func (c Cache) RetrieveAndDeActiveRuleSet(brwf, app, realm, class, ruleSetName string, slice int32) error {

	// Retrieve the rule set from cache
	currentRuleset, exists, err := c.GetRulesetName(brwf, app, realm, class, ruleSetName, slice)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("rule set not found in cache: %s", ruleSetName)
	}

	// Check if the current rule set is active
	if currentRuleset.IsActive {
		currentRuleset.IsActive = false
	}

	return nil
}
