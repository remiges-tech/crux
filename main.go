package main

import (
	"context"
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
	"github.com/remiges-tech/crux/db/sqlc-gen"
	"github.com/remiges-tech/crux/server/schema"
	"github.com/remiges-tech/crux/types"
	"github.com/remiges-tech/logharbour/logharbour"
	"github.com/remiges-tech/rigel"
	"github.com/remiges-tech/rigel/etcd"
)

func main() {
	// logger setup
	fallbackWriter := logharbour.NewFallbackWriter(os.Stdout, os.Stdout)
	lctx := logharbour.NewLoggerContext(logharbour.Info)
	l := logharbour.NewLogger(lctx, "crux", fallbackWriter)

	// Create a new EtcdStorage instance
	etcdStorage, err := etcd.NewEtcdStorage([]string{"localhost:2379"})
	if err != nil {
		l.LogActivity("Error while Creating new instance of EtcdStorage", err)
		log.Fatalf("Failed to create EtcdStorage: %v", err)
	}
	l.LogActivity("Creates a new instance of EtcdStorage with endpoints", "localhost:2379")

	// Create a new Rigel instance
	rigelClient := rigel.New(etcdStorage, "crux", "WFE", 1, "devConfig")
	l.LogActivity("Creates a new instance of rigel", rigelClient)

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var config types.AppConfig

	// Load the config
	err = rigelClient.LoadConfig(ctx, &config)
	if err != nil {
		l.LogActivity("Error while loading config", err)
		log.Fatalf("Failed to load config: %v", err)
	}
	l.LogActivity("Retrieves the configuration data from rigel", config)

	// Open the error types file
	file, err := os.Open(config.ErrorTypeFile)
	if err != nil {
		l.LogDebug("Failed to open error types file:", err)
		log.Fatalf("Failed to open error types file: %v", err)
	}
	defer file.Close()
	// Load the error types
	wscutils.LoadErrorTypes(file)

	// Database connection
	connURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)
	db, err := pg.Connect(config.DriverName, connURL)
	if err != nil {
		l.LogActivity("Error while establishes a connection with database", err)
		log.Fatalln("Failed to establishes a connection with database", err)
	}
	query := sqlc.New(db)

	// router
	r := gin.Default()

	// schema services
	schemaSvc := service.NewService(r).
		WithLogHarbour(l).
		WithDatabase(query).
		WithRigelConfig(rigelClient)

	schemaSvc.RegisterRoute(http.MethodGet, "/WFschemaList", schema.SchemaList)

	appServerPort := strconv.Itoa(config.AppServerPort)
	r.Run(":" + appServerPort)
	if err != nil {
		l.LogActivity("Failed to start server", err)
		log.Fatalf("Failed to start server: %v", err)
	}

}
