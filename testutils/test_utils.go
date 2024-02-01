package testutils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	"github.com/jackc/pgx/v5"
// 	"github.com/jackc/tern/v2/migrate"
// 	"github.com/remiges-tech/alya/router"
// 	"github.com/remiges-tech/alya/service"
// 	"github.com/remiges-tech/alya/wscutils"
// 	"github.com/remiges-tech/logharbour/logharbour"
// )

// type queriesFunc func(*sqlc.Queries) error

// // setupDependencies setups dependencies passed from array and returns
// // initialized service
// func setupDependencies(r *gin.Engine, dependent []string, db *pgx.Conn, queries sqlc.DBQuerier) *service.Service {
// 	service := service.NewService(r)
// 	for _, value := range dependent {
// 		switch value {
// 		case "logger":
// 			service.WithLogHarbour(logharbour.NewLogger(logharbour.NewLoggerContext(logharbour.DefaultPriority),
// 				"testApp", io.Discard))
// 		case "sqlc":
// 			service.WithDatabase(queries)
// 		}
// 	}
// 	return service
// }

// // runMigrations takes migration path and pgx connection
// // run migrations on test db instance
// func RunMigrations(migrationPath string, db *pgx.Conn) error {
// 	m, err := migrate.NewMigrator(context.Background(), db, "db_version")
// 	if err != nil {
// 		return err
// 	}
// 	err = m.LoadMigrations(os.DirFS(migrationPath))
// 	if err != nil {
// 		return err
// 	}
// 	return m.Migrate(context.Background())
// }

// // testSetup setup test environment initialzation
// // and runs the handler and mock response is proccessed and
// // returned
// func TestSetupForIntegrationTest(method string, url string, payload *bytes.Buffer, testDb *pgx.Conn,
//
//		handler service.HandlerFunc,
//		dependencies []string, queries queriesFunc) (*httptest.ResponseRecorder, error) {
//		gin.SetMode(gin.TestMode)
//		r, err := router.SetupRouter(false, nil, nil)
//		if err != nil {
//			return nil, err
//		}
//		file, err := os.Open("./../errortypes.yaml")
//		if err != nil {
//			log.Fatalf("Failed to open error types file: %v", err)
//		}
//		defer file.Close()
//		wscutils.LoadErrorTypes(file)
//		w := httptest.NewRecorder()
//		s := setupDependencies(r, dependencies, testDb, nil)
//		err = queries(s.Database.(*sqlc.Queries))
//		if err != nil {
//			return nil, err
//		}
//		s.RegisterRoute(method, url, handler)
//		req, err := http.NewRequest(method, url, payload)
//		if err != nil {
//			return nil, err
//		}
//		r.ServeHTTP(w, req)
//		return w, nil
//	}
func MarshalJson(data any) []byte {
	jsonData, err := json.Marshal(&data)
	if err != nil {
		log.Fatal("error marshaling")
	}
	return jsonData
}

func ReadJsonFromFile(filepath string) ([]byte, error) {
	// var err error
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("testFile path is not exist")
	}
	defer file.Close()
	jsonData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

// func TestSetupForUnitTests(method string, url string,
// 	payload *bytes.Buffer, handler service.HandlerFunc, dependent []string, queries sqlc.DBQuerier) (*httptest.ResponseRecorder, error) {
// 	gin.SetMode(gin.TestMode)
// 	r, err := router.SetupRouter(false, nil, nil)
// 	s := setupDependencies(r, dependent, nil, queries)
// 	s.RegisterRoute(method, url, handler)
// 	if err != nil {
// 		return nil, err
// 	}
// 	file, err := os.Open("../../../errortypes.yaml")
// 	if err != nil {
// 		log.Fatalf("Failed to open error types file: %v", err)
// 	}
// 	defer file.Close()
// 	wscutils.LoadErrorTypes(file)
// 	w := httptest.NewRecorder()
// 	req, err := http.NewRequest(method, url, payload)
// 	if err != nil {
// 		return nil, err
// 	}
// 	r.ServeHTTP(w, req)
// 	return w, nil
// }
