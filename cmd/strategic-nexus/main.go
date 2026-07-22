package main

// =============================================================================
// ESSENTIAL PROCESS:
// Boots and initializes the 01-Strategic-Nexus Go microservice.
// Automatically sets up Postgres connections, loads configurations, and
// routes either background server routines or direct CLI invocations.
//
// DATA FLOW:
// 1. Input:   Loads standalone.yaml capability config parameters.
// 2. Logic:   Initializes shared logger, establishes pgx pools, and maps
//             either REST HTTP/Telegram loops or runs active command iterations.
// 3. Output:  Runs CLI/REST service execution.
// =============================================================================

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Bastien-Antigravity/strategic-nexus/src/core"
	"github.com/Bastien-Antigravity/strategic-nexus/src/evolution"
	"github.com/Bastien-Antigravity/strategic-nexus/src/rest"
	"github.com/Bastien-Antigravity/strategic-nexus/src/store"
	"github.com/Bastien-Antigravity/strategic-nexus/src/telegram"

	toolbox_config "github.com/Bastien-Antigravity/microservice-toolbox/go/pkg/config"
	toolbox_lifecycle "github.com/Bastien-Antigravity/microservice-toolbox/go/pkg/lifecycle"
	unilog "github.com/Bastien-Antigravity/universal-logger/src/bootstrap"
	unilog_config "github.com/Bastien-Antigravity/universal-logger/src/config"
)

func main() {
	// 1. Initialize config loader matching standalone profile
	appConfig, err := toolbox_config.LoadConfig("standalone", []string{"prompt-name", "prompt-path", "dry-run"})
	if err != nil {
		fmt.Printf("Critical config loading error: %v\n", err)
		os.Exit(1)
	}

	// 2. Setup universal logging engine
	_, appLogger := unilog.Init("strategic-nexus", "standalone", "no_lock", "INFO", false, &unilog_config.DistConfig{Config: appConfig.Config})
	defer appLogger.Close()

	appConfig.Logger = appLogger
	appLogger.Info("Initializing Strategic-Nexus Go Microservice...")

	// 3. Resolve Timescale Database Configurations
	timescaleDBVal, exists := appConfig.Capabilities["timescale_db"]
	if !exists {
		appLogger.Error("Database config 'capabilities.timescale_db' is missing. Exiting.")
		os.Exit(1)
	}
	timescaleDB, ok := timescaleDBVal.(map[string]interface{})
	if !ok {
		appLogger.Error("Database config 'capabilities.timescale_db' is invalid. Exiting.")
		os.Exit(1)
	}

	dbHost := "127.0.0.1"
	if host, exists := timescaleDB["ip"]; exists {
		dbHost = fmt.Sprintf("%v", host)
	}
	dbPort := 5432
	if portVal, exists := timescaleDB["port"]; exists {
		fmt.Sscanf(fmt.Sprintf("%v", portVal), "%d", &dbPort)
	}
	dbUser := "dbuser"
	if user, exists := timescaleDB["user"]; exists {
		dbUser = fmt.Sprintf("%v", user)
	}
	dbPassword := "dbuser"
	if pass, exists := timescaleDB["password"]; exists {
		dbPassword = fmt.Sprintf("%v", pass)
		// Perform decryptions if configured in the handler
		if decPass, err := appConfig.DecryptSecret(dbPassword); err == nil {
			dbPassword = decPass
		}
	}
	dbName := "obsidiandb"

	// 4. Initialize Database Manager (pgx Pool & Migrations)
	dbManager, err := store.NewDBManager(dbHost, dbPort, dbUser, dbPassword, dbName, appLogger)
	if err != nil {
		appLogger.Warning("Database connection failed: %v. Running in sandbox volatile fallback mode.", err)
	} else {
		defer dbManager.Close()
	}

	// 5. Resolve Gemini API Key
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		geminiConfVal, exists := appConfig.Capabilities["gemini"]
		if exists {
			if geminiConf, ok := geminiConfVal.(map[string]interface{}); ok {
				if key, exists := geminiConf["api_key"]; exists {
					apiKey = fmt.Sprintf("%v", key)
				}
			}
		}
	}

	// 6. Setup Evolutionary Pipeline
	evolPipeline := evolution.NewEvolutionPipeline(dbManager, appLogger, apiKey)

	// 7. Setup Core Controller
	controller := core.NewStrategicController(dbManager, evolPipeline, appLogger)

	// 8. Handle Direct CLI Subcommand Routing
	// If subcommand arguments are passed (e.g. run-evolution or log-milestone)
	if len(os.Args) >= 2 && !strings.HasPrefix(os.Args[1], "-") {
		subcommand := os.Args[1]
		cmdArgs := os.Args[2:]
		
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		res, err := controller.ProcessCommand(ctx, subcommand, cmdArgs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Subcommand execution failed: %v\n", err)
			appLogger.Error("Subcommand execution failed: %v", err)
			os.Exit(1)
		}

		
		// Print successful JSON response
		fmt.Printf("Command Status: SUCCESS\nPayload:\n")
		for k, v := range res {
			fmt.Printf("  %s: %v\n", k, v)
		}
		return
	}

	// 9. Default Mode: Run background services
	restPort := 8092
	restConfVal, exists := appConfig.Capabilities["strategic_nexus"]
	if exists {
		if restConf, ok := restConfVal.(map[string]interface{}); ok {
			if portVal, exists := restConf["mfe_port"]; exists {
				fmt.Sscanf(fmt.Sprintf("%v", portVal), "%d", &restPort)
			}
		}
	}


	// Start Open MFE REST server
	mfeClient := rest.NewMFEClient(controller, appLogger, restPort)
	if err := mfeClient.Start(context.Background()); err != nil {
		appLogger.Error("Failed to start REST server: %v", err)
	}

	// Start Telegram remote subscriber
	teleClient := telegram.NewTelegramClient(controller, appLogger)
	if err := teleClient.Start(context.Background()); err != nil {
		appLogger.Error("Failed to start Telegram subscriber: %v", err)
	}

	// Initialize Graceful Shutdown Lifecycle manager
	lm := toolbox_lifecycle.NewManagerWithLogger(appLogger)
	lm.Register("StopRESTServer", func() error {
		return mfeClient.Stop(context.Background())
	})

	lm.Wait(context.Background())
}
