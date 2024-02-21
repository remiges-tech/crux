package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	pgx "github.com/jackc/pgx/v5"

	sqlc "crux/db/sqlc-gen"
)

type realm_t string
type app_t string
type slice_t int
type className_t string
type BrwfEnum string

var ConnectionString = "host=localhost port=5432 user=postgres password=postgres dbname=crux sslmode=disable"

type statsSchema_t struct {
	NChecked int
}

type statsRuleset_t struct {
	NCalled    int
	RulesStats []map[realm_t]map[app_t]map[slice_t][]statsSchema_t
}

type statsPerSlice_t struct {
	LoadedAt   time.Time
	BRSchema   map[className_t]statsSchema_t
	BRRulesets map[className_t][]statsRuleset_t
	WFSchema   map[className_t]statsSchema_t
	Workflows  map[className_t][]statsRuleset_t
}

type statsPerApp_t map[slice_t]statsPerSlice_t

type statsPerRealm_t map[app_t]statsPerApp_t

type rulesetStats_t map[realm_t]statsPerRealm_t

type schema_t struct {
	Class         string          `json:"class"`
	PatternSchema patternSchema_t `json:"patternschema"`
	ActionSchema  actionSchema_t  `json:"actionschema"`
	NChecked      int32           `json:"n_checked"`
}
type patternSchema_t struct {
	Attr []attr_t `json:"attr"`
}
type attr_t struct {
	Name     string   `json:"name"`
	ValType  string   `json:"valtype"`
	EnumVals []string `json:"vals,omitempty"`
	ValMin   float64  `json:"valmin,omitempty"`
	ValMax   float64  `json:"valmax,omitempty"`
	LenMin   int      `json:"lenmin,omitempty"`
	LenMax   int      `json:"lenmax,omitempty"`
}
type actionSchema_t struct {
	Tasks      []string          `json:"tasks"`
	Properties map[string]string `json:"properties"`
}

type rulePatternBlock_t struct {
	Attr    string    `json:"attr"`
	Op      string    `json:"op"`
	Val     string    `json:"val"`
	ValType valType_t `json:"valtype,omitempty"`
}

/*
type propertyBlock_t struct {
	Val  string `json:"val"`
	Name string `json:"name"`
}*/

type ruleActionBlock_t struct {
	Task       []string          `json:"tasks"`
	Properties map[string]string `json:"properties"`
	ThenCall   string            `json:"thencall,omitempty"`
	ElseCall   string            `json: "elsecall,omitempty"`
	DoReturn   bool              `json:"doreturn,omitempty"`
	DoExit     bool              `json:"doexit,omitempty"`
	References []*Ruleset_t      `json:"references,omitempty"`
}

type valType_t int

const (
	valInt_t valType_t = iota
	valFloat_t
	valString_t
	valBool_t
	valTimestamp_t
	valEnum_t
)

type Ruleset_t struct {
	Ver          int                  `json:"ver,omitempty""`
	Class        string               `json:"class,omitempty""`
	SetName      string               `json:"setname,omitempty""`
	RulePatterns []rulePatternBlock_t `json:"rulepattern"`
	RuleActions  ruleActionBlock_t    `json:"ruleactions"`
	NMatched     int                  `json:"nMatched"`
	NFailed      int                  `json:"nFailed"`

	NextRuleset *Ruleset_t
}

type perSlice_t struct {
	LoadedAt   time.Time
	BRSchema   map[className_t][]*schema_t
	BRRulesets map[className_t][]*Ruleset_t
	WFSchema   map[className_t][]*schema_t
	Workflows  map[className_t][]*Ruleset_t
}

type perApp_t map[slice_t]perSlice_t

type perRealm_t map[app_t]perApp_t

type rulesetCache_t map[realm_t]perRealm_t

type schemaCache_t map[realm_t]perRealm_t

var (
	rulesetCache           rulesetCache_t
	schemaCache            schemaCache_t
	cacheLock              sync.RWMutex
	nCheckedcounter        int32
	rulesetCacheStatsSince time.Time
)

func init() {
	rulesetCacheStatsSince = time.Now()
}
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

			var patterns patternSchema_t
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
				NChecked:      nCheckedcounter,
			}
			nCheckedcounter++

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

			var rulesets []Ruleset_t

			if err := json.Unmarshal(row.Ruleset, &rulesets); err != nil {
				log.Println("Error unmarshaling Ruleset:", err, string(row.Ruleset))
				continue
			}

			for _, rule := range rulesets {
				classNameKey := className_t(row.Class)
				newRuleset := &Ruleset_t{
					Class:        row.Class,
					SetName:      row.Setname,
					RulePatterns: rule.RulePatterns,
					RuleActions:  rule.RuleActions,
					NMatched:     rule.NMatched,
					NFailed:      rule.NFailed,
				}
				if row.Brwf == "B" {

					perApp[sliceKey].BRRulesets[classNameKey] = append(perApp[sliceKey].BRRulesets[classNameKey], newRuleset)

				} else if row.Brwf == "W" {
					perApp[sliceKey].Workflows[classNameKey] = append(perApp[sliceKey].Workflows[classNameKey], newRuleset)
				}
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

func retrieveRulesetFromCache(realm int, app string, class string, slice int, brwf string) ([]byte, []byte, string, []*Ruleset_t) {
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

func purgeInternal() error {
	rulesetCache = make(rulesetCache_t)
	schemaCache = make(schemaCache_t)
	return nil

}

func Load() error {
	lockCache()
	defer unlockCache()
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, ConnectionString)
	if err != nil {
		log.Fatal("Failed to load data into cache:", err)
		return err
	}
	defer conn.Close(ctx)
	query := NewProvider(ConnectionString)

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

func Reload() error {
	lockCache()
	defer unlockCache()

	if err := purgeInternal(); err != nil {
		log.Fatal("Failed to purge data from cache:", err)
		return err
	}
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, ConnectionString)
	if err != nil {
		log.Fatal("Failed to load data into cache:", err)
		return err
	}
	defer conn.Close(ctx)
	query := NewProvider(ConnectionString)

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
