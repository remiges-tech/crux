package test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	pg "github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/services/schema"
	"github.com/remiges-tech/logharbour/logharbour"
	"github.com/remiges-tech/rigel"
	"github.com/remiges-tech/rigel/etcd"
)

// Declare a global variable to hold the Docker pool and resource.
var (
	network *dockertest.Network
	r       *gin.Engine
)

type testConfig struct {
	dbHost        string
	dbPort        int
	dbUser        string
	dbPassword    string
	dbName        string
	appServerPort int
}

func TestMain(m *testing.M) {
	// Create a new EtcdStorage instance
	etcdStorage, err := etcd.NewEtcdStorage([]string{"localhost:2379"})
	if err != nil {
		log.Fatalf("Failed to create EtcdStorage: %v", err)
	}
	// Create a new Rigel instance
	rigel := rigel.New(etcdStorage, "crux", "WFE", 1, "testConfig")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var testCon testConfig

	testCon.dbHost, err = rigel.GetString(ctx, "db_host")
	if err != nil {
		log.Fatalf("Failed to get data from rigel: %v", err)
	}
	testCon.dbPort, err = rigel.GetInt(ctx, "db_port")
	if err != nil {
		log.Fatalf("Failed to get data from rigel: %v", err)
	}
	testCon.dbUser, err = rigel.GetString(ctx, "db_user")
	if err != nil {
		log.Fatalf("Failed to get data from rigel: %v", err)
	}
	testCon.dbPassword, err = rigel.GetString(ctx, "db_password")
	if err != nil {
		log.Fatalf("Failed to get data from rigel: %v", err)
	}
	testCon.dbName, err = rigel.GetString(ctx, "db_name")
	if err != nil {
		log.Fatalf("Failed to get data from rigel: %v", err)
	}
	testCon.appServerPort, err = rigel.GetInt(ctx, "app_server_port")
	if err != nil {
		log.Fatalf("Failed to get data from rigel: %v", err)
	}
	fmt.Println("Retrieves the configuration data from rigel")

	// Initialize Docker pool to insure it close at the end
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to Docker: %v", err)
	}
	fmt.Println("Initialize Docker pool")

	// Create a Docker network for the tests.
	network, err = pool.CreateNetwork("test-network")
	if err != nil {
		log.Fatalf("Could not create network: %v", err)
	}
	fmt.Println("Create a Docker network")

	// Deploy the Postgres container.
	PostgresResource, err := deployPostgres(pool, testCon)
	if err != nil {
		log.Fatalf("Could not start resource: %v", err)
	}
	fmt.Println("Deploy the Postgres container")

	// Deploy the API container.
	r, err = deployAPIContainer(pool, testCon)
	if err != nil {
		log.Fatalf("Could not start resource: %v", err)
	}

	resources := []*dockertest.Resource{
		PostgresResource,
	}

	// Run the tests.
	exitCode := m.Run()

	// Exit with the appropriate code.
	err = TearDown(pool, resources)
	if err != nil {
		log.Fatalf("Could not purge resource: %v", err)
	}
	fmt.Println("Exit with the appropriate code")

	os.Exit(exitCode)

}

// deployPostgres builds and runs the Postgres container.
func deployPostgres(pool *dockertest.Pool, testCon testConfig) (*dockertest.Resource, error) {
	userName := fmt.Sprintf("POSTGRES_USER=%s", testCon.dbUser)
	dbName := fmt.Sprintf("POSTGRES_DB=%s", testCon.dbName)
	password := fmt.Sprintf("POSTGRES_PASSWORD=%s", testCon.dbPassword)
	dbPort := fmt.Sprintf("%d/tcp", testCon.dbPort)
	// pulls an image, creates a container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			password,
			userName,
			dbName,
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort(dbPort)
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", testCon.dbUser, testCon.dbPassword, hostAndPort, testCon.dbName)

	log.Println("Connecting to database on url: ", databaseUrl)

	// resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	// pool.MaxWait = 120 * time.Second

	// Ensure the Postgres container is ready to accept connections.
	if err = pool.Retry(func() error {
		db, err := sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Create a new Tern migration instance
	m, err := migrate.New("/db/migrations/001_crux.sql", databaseUrl)
	if err != nil {
		log.Fatal("Error creating migration instance:", err)
	}

	// Run migration
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error applying migration:", err)
	}
	fmt.Println("Migration successful!")

	// // Load data from SQL file
	// if err := loadDataFromFile(db, "path/to/data.sql"); err != nil {
	// 	log.Fatal("Error loading data:", err)
	// }
	// fmt.Println("Data loaded successfully!")
	return resource, nil

}

// deployAPIContainer builds and runs the API container.
func deployAPIContainer(pool *dockertest.Pool, testCon testConfig) (*gin.Engine, error) {
	// router
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// logger setup
	fallbackWriter := logharbour.NewFallbackWriter(os.Stdout, os.Stdout)
	lctx := logharbour.NewLoggerContext(logharbour.Info)
	l := logharbour.NewLogger(lctx, "crux", fallbackWriter)

	file, err := os.Open("./errortypes.yaml")
	if err != nil {
		log.Fatalf("Failed to open error types file: %v", err)
	}
	defer file.Close()
	wscutils.LoadErrorTypes(file)

	// Database connection
	connURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", testCon.dbHost, testCon.dbPort, testCon.dbUser, testCon.dbPassword, testCon.dbName)
	query, err := pg.NewProvider(connURL)
	if err != nil {
		l.LogActivity("Error while establishes a connection with database", err)
		log.Fatalln("Failed to establishes a connection with database", err)
	}

	// schema services
	schemaSvc := service.NewService(r).
		WithLogHarbour(l).
		WithDatabase(query)

	schemaSvc.RegisterRoute(http.MethodGet, "/WFschemaList", schema.SchemaList)
	schemaSvc.RegisterRoute(http.MethodPost, "/WFschemaNew", schema.SchemaNew)
	schemaSvc.RegisterRoute(http.MethodPut, "/WFschemaUpdate", schema.SchemaUpdate)

	appServerPortStr := strconv.Itoa(testCon.appServerPort)
	r.Run(":" + appServerPortStr)
	if err != nil {
		l.LogActivity("Failed to start server", err)
		log.Fatalf("Failed to start server: %v", err)
	}

	return r, nil

}

// TearDown purges the resources and removes the network.
func TearDown(pool *dockertest.Pool, resources []*dockertest.Resource) error {
	for _, resource := range resources {
		if err := pool.Purge(resource); err != nil {
			return fmt.Errorf("could not purge resource: %v", err)
		}
	}

	if err := pool.RemoveNetwork(network); err != nil {
		return fmt.Errorf("could not remove network: %v", err)
	}

	return nil
}
