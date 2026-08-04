package rest

// =============================================================================
// ESSENTIAL PROCESS:
// HTTP REST server listener for frontend Microfrontends to query status
// and trigger evolution/planning events.
//
// DATA FLOW:
// 1. Input:   HTTP POST/GET payloads on /status or /command.
// 2. Logic:   Invokes ProcessCommand on StrategicController.
// 3. Output:  JSON-serialized response returned to the HTTP socket.
// =============================================================================

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type StrategicController interface {
	ProcessCommand(ctx context.Context, command string, args []string) (map[string]interface{}, error)
}

type Logger interface {
	Info(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
}

type MFEClient struct {
	controller StrategicController
	logger     Logger
	port       int
	server     *http.Server
}

func NewMFEClient(controller StrategicController, logger Logger, port int) *MFEClient {
	return &MFEClient{
		controller: controller,
		logger:     logger,
		port:       port,
	}
}

func (mfe *MFEClient) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// GET /status endpoint
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		res := map[string]interface{}{
			"status":     "active",
			"repository": "01-Strategic-Nexus",
			"healthy":    true,
		}
		_ = json.NewEncoder(w).Encode(res)
	})

	// POST /command endpoint
	mux.HandleFunc("/command", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var payload struct {
			Command string   `json:"command"`
			Args    []string `json:"args"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Run via controller
		res, err := mfe.controller.ProcessCommand(r.Context(), payload.Command, payload.Args)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "error", "message": err.Error()})
			return
		}

		_ = json.NewEncoder(w).Encode(res)
	})

	mfe.server = &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", mfe.port),
		Handler: mux,
	}

	mfe.logger.Info("MFEClient: Starting HTTP REST API portal on port %d...", mfe.port)

	// Start server in background goroutine
	go func() {
		if err := mfe.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			mfe.logger.Error("MFEClient: Server failed: %v", err)
		}
	}()

	return nil
}

func (mfe *MFEClient) Stop(ctx context.Context) error {
	if mfe.server != nil {
		shutdownCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		return mfe.server.Shutdown(shutdownCtx)
	}
	return nil
}
