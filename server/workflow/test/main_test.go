package workflow_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	pg "github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server/workflow"
	"github.com/remiges-tech/logharbour/logharbour"
)

var r *gin.Engine
var versionTable string = "schema_version_non_default"
var timeout = 100 * time.Second

func TestMain(m *testing.M) {

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	// Deploy the Postgres container.
	databaseUrl, err := deployPostgres(ctx)
	if err != nil {
		log.Fatalf("Could not start resource: %v", err)
	}
	fmt.Println("Deploy the Postgres container")

	ternMigrate(databaseUrl)
	fmt.Println("tern migrate")

	// Register routes.
	r, err = registerRoutes(databaseUrl)
	if err != nil {
		log.Fatalf("Could not start resource: %v", err)
	}
	fmt.Println("Register routes")
	// Run the tests.
	exitCode := m.Run()

	os.Exit(exitCode)

}

// deployPostgres builds and runs the Postgres container.
func deployPostgres(ctx context.Context) (string, error) {
	// pulls an image, creates a container
	postgresContainer, err := postgres.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		postgres.WithDatabase("crux"),
		postgres.WithUsername("root"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)))

	if err != nil {
		return "", fmt.Errorf("Could not start resource: %s", err)
	}

	dbURL, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return "", fmt.Errorf("Could not get dbURL: %s", err)
	}

	log.Println("Connecting to database on url: ", dbURL)

	// Ensure the Postgres container is ready to accept connections.

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return "", fmt.Errorf("error while connecting db: %s", err)
	}
	if err := db.Ping(); err != nil {
		return "", fmt.Errorf("error while pinging db: %s", err)
	}

	return dbURL, nil

}

// registerRoutes register and runs.
func registerRoutes(databaseUrl string) (*gin.Engine, error) {
	// router
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// logger setup
	fallbackWriter := logharbour.NewFallbackWriter(os.Stdout, os.Stdout)
	lctx := logharbour.NewLoggerContext(logharbour.Info)
	l := logharbour.NewLogger(lctx, "crux", fallbackWriter)
	// Define a custom validation tag-to-message ID map
	customValidationMap := map[string]int{
		"required":  101,
		"gt":        102,
		"alpha":     103,
		"lowercase": 104,
		"max":       105,
		"lt":        106,
	}
	// Custom validation tag-to-error code map
	customErrCodeMap := map[string]string{
		"required":  "required",
		"gt":        "greater",
		"alpha":     "alphabet",
		"lowercase": "lowercase",
		"max":       "exceed the maximum value allowed",
		"lt":        "exceed the limit value allowed",
	}
	// Register the custom map with wscutils
	wscutils.SetValidationTagToMsgIDMap(customValidationMap)
	wscutils.SetValidationTagToErrCodeMap(customErrCodeMap)

	// Set default message ID and error code if needed
	wscutils.SetDefaultMsgID(100)
	wscutils.SetDefaultErrCode("validation_error")

	connPool, err := pg.NewProvider(databaseUrl)
	if err != nil {
		l.LogActivity("Error while establishes a connection with database", err)
		log.Fatalln("Failed to establishes a connection with database", err)
	}
	queries := sqlc.New(connPool)

	// schema services
	s := service.NewService(r).
		WithLogHarbour(l).
		WithDatabase(connPool).
		WithDependency("queries", queries)

	s.RegisterRoute(http.MethodPut, "/workflowUpdate", workflow.WorkFlowUpdate)
	s.RegisterRoute(http.MethodDelete, "/workflowdelete", workflow.WorkflowDelete)
	s.RegisterRoute(http.MethodPost, "/workflowlist", workflow.WorkflowList)
	s.RegisterRoute(http.MethodPost, "/workflowget", workflow.WorkflowGet)
	s.RegisterRoute(http.MethodPost, "/workflownew", workflow.WorkFlowNew)

	return r, nil

}

func ternMigrate(databaseUrl string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	conn, err := pgx.Connect(ctx, databaseUrl)
	if err != nil {
		log.Fatalln("unable to connect", err)
	}

	// Create a new Tern migration instance
	m, err := migrate.NewMigrator(ctx, conn, versionTable)
	if err != nil {
		log.Fatal("Error creating migration instance:", err)
	}
	if err := m.LoadMigrations("../../../db/migrations/"); err != nil {
		log.Fatal("Error loading data:", err)
	}
	if err = m.Migrate(ctx); err != nil {
		log.Fatal("Error loading data:", err)
	}
}


