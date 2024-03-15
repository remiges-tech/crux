/*
This file contains the functions that represent Cache tests for Load()/Purge()/Reload(). These functions are called
inside TestCache()) in do_matchest.go.

Some of the definitions of rulesets below deliberately use a lot of whitespace to keep the code consistent
and to make it easier to understand, add to, and edit these tests
*/

package crux

import (
	"context"
	"log"
	"testing"

	sqlc "github.com/remiges-tech/crux/matching-engine/db/sqlc-gen"

	"github.com/jackc/pgx/v5"
)

func testinit() (sqlc.DBQuerier, context.Context, error) {
	var ConnectionString = "host=localhost port=5432 user=postgres password=postgres dbname=crux sslmode=disable"
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, ConnectionString)
	if err != nil {
		log.Fatal("Failed to load data into cache:", err)
		return nil, nil, err
	}
	defer conn.Close(ctx)
	query := NewProvider(ConnectionString)
	return query, ctx, err

}

func testCache(tests *[]doMatchTest, t *testing.T) {

	//query, ctx,err := testinit()
	//if err != nil {
	//testLoadDB(tests, t,query,ctx)
	//}
	// Call the initializeRuleData function to populate ruleSchemas and ruleSets

	testLoad(tests, t)
	setSchemaRulesetCacheBuffer(t)

	//testPurge(tests, t)
	//testReload(tests, t, query, ctx)
}

func testLoadDB(tests *[]doMatchTest, t *testing.T, q sqlc.DBQuerier, c context.Context) {

	err := Load(q, c)
	if err != nil {
		t.Errorf("Error:%+v", err)
	}
}

func setSchemaRulesetCacheBuffer(t *testing.T) {

	err := loadInternal(mockSchemasets, mockRulesets)
	if err != nil {
		t.Errorf(" %v", err)
		return
	}

}

func testLoad(tests *[]doMatchTest, t *testing.T) {

	setSchemaRulesetCacheBuffer(t)

	//PrintAllSchemaCache()
	//PrintAllRuleSetCache()

}

func testPurge(tests *[]doMatchTest, t *testing.T) {

	err := Purge()
	if err != nil {
		t.Errorf("ERROR Purge %+v", err)
	}
}

func testReload(tests *[]doMatchTest, t *testing.T, q sqlc.DBQuerier, c context.Context) {

	/*err := Reload(q,c)
	if err != nil {
		t.Errorf("ERROR Reload %+v", err)
	}*/
	// Not needed its a combination of purge and load func
}
