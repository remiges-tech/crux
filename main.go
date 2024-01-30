package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/remiges-tech/alya/service"
	"github.com/remiges-tech/alya/wscutils"
	pg "github.com/remiges-tech/crux/db"
	"github.com/remiges-tech/crux/services/schema"
	"github.com/remiges-tech/logharbour/logharbour"
	"github.com/remiges-tech/rigel"
	"github.com/remiges-tech/rigel/etcd"
)

func main() {
	// logger setup
	fallbackWriter := logharbour.NewFallbackWriter(os.Stdout, os.Stdout)
	lctx := logharbour.NewLoggerContext(logharbour.Info)
	l := logharbour.NewLogger(lctx, "crux", fallbackWriter)

	rigelAppName := flag.String("appName", "crux", "The name of the application")
	rigelModuleName := flag.String("moduleName", "WFE", "The name of the module")
	rigelVersionNumber := flag.Int("versionNumber", 1, "The number of the version")
	rigelConfigName := flag.String("configName", "devConfig", "The name of the configuration")
	etcdEndpoints := flag.String("etcdEndpoints", "localhost:2379", "Comma-separated list of etcd endpoints")

	flag.Parse()
	// Create a new EtcdStorage instance
	etcdStorage, err := etcd.NewEtcdStorage([]string{*etcdEndpoints})
	if err != nil {
		l.LogActivity("Error while Creating new instance of EtcdStorage", err)
		log.Fatalf("Failed to create EtcdStorage: %v", err)
	}
	l.LogActivity("Creates a new instance of EtcdStorage with endpoints", "localhost:2379")

	// Create a new Rigel instance
	rigel := rigel.New(etcdStorage, *rigelAppName, *rigelModuleName, *rigelVersionNumber, *rigelConfigName)
	l.LogActivity("Creates a new instance of rigel", rigel)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbHost, err := rigel.GetString(ctx, "db_host")
	if err != nil {
		l.LogActivity("Error while getting data from rigel", err)
		log.Fatalf("Failed to get data from rigel: %v", err)
	}
	dbPort, err := rigel.GetInt(ctx, "db_port")
	if err != nil {
		l.LogActivity("Error while getting data from rigel", err)
		log.Fatalf("Failed to get data from rigel: %v", err)
	}
	dbUser, err := rigel.GetString(ctx, "db_user")
	if err != nil {
		l.LogActivity("Error while getting data from rigel", err)
		log.Fatalf("Failed to get data from rigel: %v", err)
	}
	dbPassword, err := rigel.GetString(ctx, "db_password")
	if err != nil {
		l.LogActivity("Error while getting data from rigel", err)
		log.Fatalf("Failed to get data from rigel: %v", err)
	}
	dbName, err := rigel.GetString(ctx, "db_name")
	if err != nil {
		l.LogActivity("Error while getting data from rigel", err)
		log.Fatalf("Failed to get data from rigel: %v", err)
	}
	appServerPort, err := rigel.GetInt(ctx, "app_server_port")
	if err != nil {
		l.LogActivity("Error while getting data from rigel", err)
		log.Fatalf("Failed to get data from rigel: %v", err)
	}
	errorTypeFile, err := rigel.GetString(ctx, "error_type_file")
	if err != nil {
		l.LogActivity("Error while getting data from rigel", err)
		log.Fatalf("Failed to get data from rigel: %v", err)
	}
	l.Log("Retrieves the configuration data from rigel")

	// Open the error types file
	file, err := os.Open(errorTypeFile)
	if err != nil {
		l.LogDebug("Failed to open error types file:", err)
		log.Fatalf("Failed to open error types file: %v", err)
	}
	defer file.Close()
	// Load the error types
	wscutils.LoadErrorTypes(file)

	// Database connection
	connURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	query, err := pg.NewProvider(connURL)
	if err != nil {
		l.LogActivity("Error while establishes a connection with database", err)
		log.Fatalln("Failed to establishes a connection with database", err)
	}

	// router
	r := gin.Default()

	// schema services
	schemaSvc := service.NewService(r).
		WithLogHarbour(l).
		WithDatabase(query)

	schemaSvc.RegisterRoute(http.MethodGet, "/WFschemaList", schema.SchemaList)
	schemaSvc.RegisterRoute(http.MethodGet, "/WFschemaNew", schema.SchemaNew)
	schemaSvc.RegisterRoute(http.MethodGet, "/WFschemaUpdate", schema.SchemaUpdate)

	appServerPortStr := strconv.Itoa(appServerPort)
	r.Run(":" + appServerPortStr)
	if err != nil {
		l.LogActivity("Failed to start server", err)
		log.Fatalf("Failed to start server: %v", err)
	}

}
