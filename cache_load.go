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
	"github.com/jackc/pgx/v5/pgxpool"

	sqlc "remiges/crux/db/sqlc-gen"
)

type realm_t string
type app_t string
type slice_t int
type className_t string
type BrwfEnum string

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
	Class         string            `json:"class"`
	PatternSchema []patternSchema_t `json:"patternschema"`
	ActionSchema  actionSchema_t    `json:"actionschema"`
	NChecked      int32             `json:"n_checked"`
}
type patternSchema_t struct {
	Attr  []attr_t `json:"attr"`
	Class string   `json:"class"`
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
	Tasks      []string `json:"tasks"`
	Properties []string `json:"properties"`
}

type rulePatternBlock_t struct {
	Attr    string    `json:"attr"`
	Op      string    `json:"op"`
	Val     string    `json:"val"`
	ValType valType_t `json:"valtype"`
}

type propertyBlock_t struct {
	Val  string `json:"val"`
	Name string `json:"name"`
}

type ruleActionBlock_t struct {
	Task               []string          `json:"tasks"`
	Properties         []propertyBlock_t `json:"properties"`
	ThenCall, ElseCall string            `json:"thencall,omitempty" elsecall,omitempty"`
	DoReturn, DoExit   bool              `json:"doreturn,omitempty" doexit,omitempty"`
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
	RulePatterns []rulePatternBlock_t `json:"rulepattern"`
	RuleActions  ruleActionBlock_t    `json:"ruleactions"`
	NMatched     int                  `json:"nMatched"`
	NFailed      int                  `json:"nFailed"`
}

type perSlice_t struct {
	LoadedAt   time.Time
	BRSchema   map[className_t][]schema_t
	BRRulesets map[className_t][]Ruleset_t
	WFSchema   map[className_t][]schema_t
	Workflows  map[className_t][]Ruleset_t
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

func init() {
	rulesetCacheStatsSince = time.Now()
}

func loadInternal() error {
	ctx := context.Background()
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=crux sslmode=disable"
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatal("Failed to load data into cache:", err)
		return err
	}
	defer conn.Close(ctx)
	query := NewProvider(connStr)
	rulesetCache = make(rulesetCache_t)
	schemaCache = make(schemaCache_t)

	dbResponseSchema, err := query.AllSchemas(ctx)
	if err != nil {
		return err
	}
	var patternCacheSchema []patternSchema_t
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
				BRSchema: make(map[className_t][]schema_t),
				WFSchema: make(map[className_t][]schema_t),
			}

			var patterns patternSchema_t
			if err := json.Unmarshal(row.Patternschema, &patterns); err != nil {
				log.Println("Error unmarshaling Patternschema:", err)
				continue
			}
			patternCacheSchema = append(patternCacheSchema, patterns)

			var actions actionSchema_t
			if err := json.Unmarshal(row.Actionschema, &actions); err != nil {
				log.Println("Error parsing ActionSchema JSON:", err)
				continue
			}
			schemaData := schema_t{
				Class:         row.Class,
				PatternSchema: patternCacheSchema,
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

	dbResponseRuleSet, err := query.AllRuleset(ctx)
	if err != nil {
		return err
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
				BRRulesets: make(map[className_t][]Ruleset_t),
				Workflows:  make(map[className_t][]Ruleset_t),
			}

			var rulesets []Ruleset_t
			if err := json.Unmarshal(row.Ruleset, &rulesets); err != nil {
				log.Println("Error unmarshaling Ruleset:", err)
				continue
			}

			for _, rule := range rulesets {
				classNameKey := className_t(row.Class)
				if row.Brwf == "B" {
					perApp[sliceKey].BRRulesets[classNameKey] = append(perApp[sliceKey].BRRulesets[classNameKey], rule)
				} else if row.Brwf == "W" {
					perApp[sliceKey].Workflows[classNameKey] = append(perApp[sliceKey].Workflows[classNameKey], rule)
				}
			}
		}
	}

	return nil
}

func purgeInternal() error {
	rulesetCache = make(rulesetCache_t)
	schemaCache = make(schemaCache_t)
	return nil

}

func Load() error {
	lockCache()
	defer unlockCache()

	err := loadInternal()
	if err != nil {
		log.Fatal("Failed to load data into cache:", err)
		// Handle the error as needed
		return err
	}
	fmt.Println("RULESET CACHE DATA ", rulesetCache)

	fmt.Println("SHCEMA CACHE DATA ", schemaCache)
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

	err := loadInternal()
	if err != nil {
		log.Fatal("Failed to load data into cache:", err)
		// Handle the error as needed
		return err
	}
	fmt.Println("RULESET CACHE DATA ", rulesetCache)
	return nil
}
