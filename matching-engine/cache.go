package crux

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/remiges-tech/crux/db/sqlc-gen"
)

func loadInternalSchema(dbResponseSchema []sqlc.Schema) error {

	if len(dbResponseSchema) == 0 {
		return fmt.Errorf("dbResponseRuleSet is empty")
	}

	for _, row := range dbResponseSchema {
		realmKey := Realm_t(string(row.Realm))
		perRealm, exists := SchemaCache[realmKey]
		if !exists {
			perRealm = make(PerRealm_t)
			SchemaCache[realmKey] = perRealm
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
		if row.Brwf == "B" {
			perApp[sliceKey].BRSchema[classNameKey] = schemaData
		} else if row.Brwf == "W" {
			perApp[sliceKey].WFSchema[classNameKey] = schemaData
		}

	}
	return nil
}

func loadInternalRuleSet(dbResponseRuleSet []sqlc.Ruleset) error {

	if len(dbResponseRuleSet) == 0 {
		return fmt.Errorf("dbResponseRuleSet is empty")
	}
	for _, row := range dbResponseRuleSet {

		realmKey := Realm_t(string(row.Realm))
		perRealm, exists := RulesetCache[realmKey]
		if !exists {
			perRealm = make(PerRealm_t)
			RulesetCache[realmKey] = perRealm
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
				LoadedAt:   time.Now(),
				BRRulesets: make(map[ClassName_t][]*Ruleset_t),
				Workflows:  make(map[ClassName_t][]*Ruleset_t),
			}
		}

		var rules []Rule_t
		err := json.Unmarshal(row.Ruleset, &rules)
		if err != nil {
			fmt.Println("Error unmarshaling rules:", err)
			return nil
		}

		classNameKey := ClassName_t(row.Class)
		newRuleset := &Ruleset_t{
			Id:      row.ID,
			Class:   row.Class,
			SetName: row.Setname,
			Rules:   rules,
		}
		if row.Brwf == "B" {
			perApp[sliceKey].BRRulesets[classNameKey] = append(perApp[sliceKey].BRRulesets[classNameKey], newRuleset)

		} else if row.Brwf == "W" {
			fmt.Printf("className: %v", classNameKey)
			fmt.Printf("workflow name: %v", newRuleset.SetName)
			perApp[sliceKey].Workflows[classNameKey] = append(perApp[sliceKey].Workflows[classNameKey], newRuleset)
		}

	}
	AddReferencesToRuleSetCache()
	return nil
}
func loadInternal(dbResponseSchema []sqlc.Schema, dbResponseRuleSet []sqlc.Ruleset) error {
	RulesetCache = make(RulesetCache_t)
	SchemaCache = make(SchemaCache_t)

	err := loadInternalSchema(dbResponseSchema)
	if err != nil {
		return err
	}

	err = loadInternalRuleSet(dbResponseRuleSet)
	if err != nil {
		return err
	}
	return nil
}

func purgeInternal() error {
	RulesetCache = make(RulesetCache_t)
	SchemaCache = make(SchemaCache_t)
	return nil

}

func Load(query sqlc.Querier, ctx context.Context) error {

	lockCache()
	defer unlockCache()

	dbResponseSchema, err := query.AllSchemas(ctx)
	if err != nil {
		return err
	}
	if len(dbResponseSchema) == 0 {
		return fmt.Errorf("Didn't get schema")
	}

	dbResponseRuleSet, err := query.AllRuleset(ctx)
	if err != nil {
		return err
	}
	if len(dbResponseRuleSet) == 0 {
		return fmt.Errorf("Didn't get rule set")
	}
	err = loadInternal(dbResponseSchema, dbResponseRuleSet)
	if err != nil {
		return err
	}

	return nil
}

func Purge() error {
	lockCache()
	defer unlockCache()

	if err := purgeInternal(); err != nil {
		log.Fatal("Failed to purge data from cache:", err)
		return err
	}

	return nil
}

func Reload(query sqlc.DBQuerier, ctx context.Context) error {
	lockCache()
	defer unlockCache()

	if err := purgeInternal(); err != nil {
		log.Fatal("Failed to purge data from cache:", err)
		return err
	}

	dbResponseSchema, err := query.AllSchemas(ctx)
	if err != nil {
		return err
	}

	dbResponseRuleSet, err := query.AllRuleset(ctx)
	if err != nil {
		return err
	}
	err = loadInternal(dbResponseSchema, dbResponseRuleSet)
	if err != nil {
		return err
	}

	return nil
}
