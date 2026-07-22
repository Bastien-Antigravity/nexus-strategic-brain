package store

// =============================================================================
// ESSENTIAL PROCESS:
// Manages the PostgreSQL database connection pool and schema lifecycle
// for the 01-Strategic-Nexus.
//
// DATA FLOW:
// 1. Input: Database connection credentials (host, port, user, password).
// 2. Logic: Checks/creates obsidiandb, runs migrations for strategic tables.
// 3. Output: Exposed connection pool and thread-safe CRUD methods.
// =============================================================================

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBManager struct {
	Pool   *pgxpool.Pool
	Schema string
	logger Logger
	mu     sync.Mutex
}

type Logger interface {
	Info(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
}

func NewDBManager(host string, port int, user, password, dbName string, logger Logger) (*DBManager, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Ensure target database exists by connecting to postgres system DB first
	pgConnString := fmt.Sprintf("postgres://%s:%s@%s:%d/postgres?sslmode=disable", user, password, host, port)
	pgPool, err := pgxpool.New(ctx, pgConnString)
	if err == nil {
		var exists bool
		query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
		_ = pgPool.QueryRow(ctx, query, dbName).Scan(&exists)
		if !exists {
			logger.Info("[Strategic-Nexus] Database '%s' not found. Creating dynamically...", dbName)
			_, _ = pgPool.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbName))
		}
		pgPool.Close()
	}

	// 2. Connect to the target database pool
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=10", user, password, host, port, dbName)
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	mgr := &DBManager{
		Pool:   pool,
		Schema: "01-Strategic-Nexus",
		logger: logger,
	}

	if err := mgr.InitSchema(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return mgr, nil
}

func (db *DBManager) InitSchema(ctx context.Context) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Declare schema context
	_, err = tx.Exec(ctx, fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, db.Schema))
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	// Create strategic_logs table
	logsTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS "%s".strategic_logs (
			id SERIAL PRIMARY KEY,
			session_id VARCHAR(255) NOT NULL,
			milestone VARCHAR(255) NOT NULL,
			action_taken TEXT NOT NULL,
			outcome TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`, db.Schema)
	if _, err = tx.Exec(ctx, logsTable); err != nil {
		return fmt.Errorf("failed to create strategic_logs table: %w", err)
	}

	// Create evolutive_skills table
	skillsTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS "%s".evolutive_skills (
			id SERIAL PRIMARY KEY,
			skill_name VARCHAR(255) UNIQUE NOT NULL,
			proposed_rules TEXT NOT NULL,
			source_trace_id VARCHAR(255),
			status VARCHAR(50) DEFAULT 'pending',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`, db.Schema)
	if _, err = tx.Exec(ctx, skillsTable); err != nil {
		return fmt.Errorf("failed to create evolutive_skills table: %w", err)
	}

	// Create strategic_decisions table
	decisionsTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS "%s".strategic_decisions (
			id SERIAL PRIMARY KEY,
			decision_id VARCHAR(50) UNIQUE NOT NULL,
			title VARCHAR(255) NOT NULL,
			category VARCHAR(50) NOT NULL,
			content TEXT NOT NULL,
			status VARCHAR(50) DEFAULT 'active',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`, db.Schema)
	if _, err = tx.Exec(ctx, decisionsTable); err != nil {
		return fmt.Errorf("failed to create strategic_decisions table: %w", err)
	}

	return tx.Commit(ctx)
}

func (db *DBManager) LogAction(ctx context.Context, sessionID, milestone, action, outcome string) error {
	query := fmt.Sprintf(`
		INSERT INTO "%s".strategic_logs (session_id, milestone, action_taken, outcome)
		VALUES ($1, $2, $3, $4)
	`, db.Schema)
	_, err := db.Pool.Exec(ctx, query, sessionID, milestone, action, outcome)
	return err
}

func (db *DBManager) ProposeSkill(ctx context.Context, skillName, proposedRules, sourceTraceID string) error {
	query := fmt.Sprintf(`
		INSERT INTO "%s".evolutive_skills (skill_name, proposed_rules, source_trace_id, status)
		VALUES ($1, $2, $3, 'pending')
		ON CONFLICT (skill_name) 
		DO UPDATE SET proposed_rules = EXCLUDED.proposed_rules, status = 'pending', created_at = CURRENT_TIMESTAMP
	`, db.Schema)
	_, err := db.Pool.Exec(ctx, query, skillName, proposedRules, strings.TrimSpace(sourceTraceID))
	return err
}

func (db *DBManager) UpsertDecision(ctx context.Context, decisionID, title, category, content string) error {
	query := fmt.Sprintf(`
		INSERT INTO "%s".strategic_decisions (decision_id, title, category, content)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (decision_id) DO UPDATE SET
			title = EXCLUDED.title,
			category = EXCLUDED.category,
			content = EXCLUDED.content,
			updated_at = CURRENT_TIMESTAMP
	`, db.Schema)
	_, err := db.Pool.Exec(ctx, query, decisionID, title, category, content)
	return err
}

func (db *DBManager) GetDecisions(ctx context.Context) ([]map[string]interface{}, error) {
	query := fmt.Sprintf(`
		SELECT decision_id, title, category, content, status, created_at, updated_at
		FROM "%s".strategic_decisions
	`, db.Schema)
	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var dID, title, category, content, status string
		var cAt, uAt time.Time
		if err := rows.Scan(&dID, &title, &category, &content, &status, &cAt, &uAt); err != nil {
			return nil, err
		}
		item := map[string]interface{}{
			"decision_id": dID,
			"title":       title,
			"category":    category,
			"content":     content,
			"status":      status,
			"created_at":  cAt,
			"updated_at":  uAt,
		}
		results = append(results, item)
	}
	return results, nil
}

func (db *DBManager) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
