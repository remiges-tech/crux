package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	sqlc "crux/db/sqlc-gen"
)

func loadInternalSchema(dbResponseSchema []sqlc.Schema) error {

	if len(dbResponseSchema) == 0 {
		return fmt.Errorf("dbResponseRuleSet is empty")
	}

	for _, row := range dbResponseSchema {
		realmKey := realm_t(strconv.Itoa(int(row.Realm)))
		perRealm, exists := schemaCache[realmKey]
		if !exists {
			perRealm = make(perRealm_t)
			schemaCache[realmKey] = perRealm
		}

		appKey := app_t(row.App)
		perApp, exists := perRealm[appKey]
		if !exists {
			perApp = make(perApp_t)
			perRealm[appKey] = perApp
		}

		sliceKey := slice_t(row.Slice)
		_, exists = perApp[sliceKey]
		if !exists {
			perApp[sliceKey] = perSlice_t{
				LoadedAt: time.Now(),
				BRSchema: make(map[className_t][]*schema_t),
				WFSchema: make(map[className_t][]*schema_t),
			}

			var patterns []patternSchema_t

			if err := json.Unmarshal(row.Patternschema, &patterns); err != nil {
				log.Println("Error unmarshaling Patternschema:", err)
				continue
			}

			var actions actionSchema_t

			if err := json.Unmarshal(row.Actionschema, &actions); err != nil {
				log.Println("Error parsing ActionSchema JSON:", err)
				continue
			}

			schemaData := &schema_t{
				Class:         row.Class,
				PatternSchema: patterns,
				ActionSchema:  actions,
			}

			classNameKey := className_t(row.Class)
			if row.Brwf == "B" {
				perApp[sliceKey].BRSchema[classNameKey] = append(perApp[sliceKey].BRSchema[classNameKey], schemaData)
			} else if row.Brwf == "W" {
				perApp[sliceKey].WFSchema[classNameKey] = append(perApp[sliceKey].WFSchema[classNameKey], schemaData)
			}

		}
	}
	return nil
}

func loadInternalRuleSet(dbResponseRuleSet []sqlc.Ruleset) error {

	if len(dbResponseRuleSet) == 0 {
		return fmt.Errorf("dbResponseRuleSet is empty")
	}
	for _, row := range dbResponseRuleSet {

		realmKey := realm_t(strconv.Itoa(int(row.Realm)))
		perRealm, exists := rulesetCache[realmKey]
		if !exists {
			perRealm = make(perRealm_t)
			rulesetCache[realmKey] = perRealm
		}

		appKey := app_t(row.App)
		perApp, exists := perRealm[appKey]
		if !exists {
			perApp = make(perApp_t)
			perRealm[appKey] = perApp
		}

		sliceKey := slice_t(row.Slice)
		_, exists = perApp[sliceKey]
		if !exists {
			perApp[sliceKey] = perSlice_t{
				LoadedAt:   time.Now(),
				BRRulesets: make(map[className_t][]*Ruleset_t),
				Workflows:  make(map[className_t][]*Ruleset_t),
			}

			var rules []rule_t

			err := json.Unmarshal(row.Ruleset, &rules)
			if err != nil {
				fmt.Println("Error unmarshaling rules:", err)
				return nil
			}

			classNameKey := className_t(row.Setname)
			newRuleset := &Ruleset_t{
				Id:      row.ID,
				Class:   row.Class,
				SetName: row.Setname,
				Rules:   rules,
			}
			if row.Brwf == "B" {
				perApp[sliceKey].BRRulesets[classNameKey] = append(perApp[sliceKey].BRRulesets[classNameKey], newRuleset)

			} else if row.Brwf == "W" {
				perApp[sliceKey].Workflows[classNameKey] = append(perApp[sliceKey].Workflows[classNameKey], newRuleset)
			}

		}
	}
	AddReferencesToRuleSetCache()
	return nil
}
func loadInternal(dbResponseSchema []sqlc.Schema, dbResponseRuleSet []sqlc.Ruleset) error {
	rulesetCache = make(rulesetCache_t)
	schemaCache = make(schemaCache_t)

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
	rulesetCache = make(rulesetCache_t)
	schemaCache = make(schemaCache_t)
	return nil

}

func Load(query sqlc.DBQuerier, ctx context.Context) error {

	lockCache()
	defer unlockCache()

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
