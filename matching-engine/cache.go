package crux

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/remiges-tech/crux/db/sqlc-gen"
)

type Cache struct {
	Ctx          context.Context
	Query        sqlc.Querier
	Slice        Slice_t
	App          App_t
	Class        ClassName_t
	Realm        Realm_t
	WorkflowName string
}

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

func loadInternalRuleSet(dbResponseRuleSet sqlc.Ruleset) error {

	if IsZeroOfUnderlyingType(dbResponseRuleSet) {
		return fmt.Errorf("didn't get rule set")
	}

	realmKey := Realm_t(string(dbResponseRuleSet.Realm))
	perRealm, exists := RulesetCache[realmKey]
	if !exists {
		perRealm = make(PerRealm_t)
		RulesetCache[realmKey] = perRealm
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
		Id:      dbResponseRuleSet.ID,
		Class:   dbResponseRuleSet.Class,
		SetName: dbResponseRuleSet.Setname,
		Rules:   rules,
	}
	if dbResponseRuleSet.Brwf == "B" {
		perApp[sliceKey].BRRulesets[classNameKey] = append(perApp[sliceKey].BRRulesets[classNameKey], newRuleset)

	} else if dbResponseRuleSet.Brwf == "W" {
		fmt.Printf("className: %v", classNameKey)
		fmt.Printf("workflow name: %v", newRuleset.SetName)
		perApp[sliceKey].Workflows[classNameKey] = append(perApp[sliceKey].Workflows[classNameKey], newRuleset)
	}

	// AddReferencesToRuleSetCache()
	return nil
}

func loadInternal(dbResponseSchema []sqlc.Schema, dbResponseRuleSet sqlc.Ruleset) error {
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

func (c Cache) Load() error {

	lockCache()
	defer unlockCache()

	// dbResponseSchema, err := query.AllSchemas(ctx)
	dbResponseSchema, err := c.Query.LoadSchema(c.Ctx, sqlc.LoadSchemaParams{
		RealmName: string(c.Realm),
		Slice:     int32(c.Slice),
		Class:     string(c.Class),
		App:       string(c.App),
	})
	if err != nil {
		return err
	}
	if len(dbResponseSchema) == 0 {
		return fmt.Errorf("schema not found")
	}

	dbResponseRuleSet, err := c.Query.LoadRuleSet(c.Ctx,
		sqlc.LoadRuleSetParams{
			RealmName: string(c.Realm),
			Slice:     int32(c.Slice),
			Class:     string(c.Class),
			App:       string(c.App),
			Setname:   c.WorkflowName,
		})
	if err != nil {
		return err
	}
	if IsZeroOfUnderlyingType(dbResponseRuleSet) {
		return fmt.Errorf("rule set not found")
	}
	err = loadInternal(dbResponseSchema, dbResponseRuleSet)
	if err != nil {
		return err
	}
	return nil
}

func (c Cache) Purge(brwf, field string) {
	lockCache()
	defer unlockCache()

	if brwf == "B" && field == "schema" {
		SchemaCache[c.Realm][c.App][c.Slice].BRSchema[ClassName_t(c.Class)] = Schema_t{}
		RulesetCache[c.Realm][c.App][c.Slice].BRRulesets[c.Class] = nil
	} else if brwf == "B" && field == "rule" {
		RulesetCache[c.Realm][c.App][c.Slice].BRRulesets[c.Class] = nil
	} else if brwf == "w" && field == "schema" {
		SchemaCache[c.Realm][c.App][c.Slice].WFSchema[ClassName_t(c.Class)] = Schema_t{}
		RulesetCache[c.Realm][c.App][c.Slice].Workflows[c.Class] = nil
	} else if brwf == "W" && field == "rule" {
		RulesetCache[c.Realm][c.App][c.Slice].Workflows[c.Class] = nil
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
