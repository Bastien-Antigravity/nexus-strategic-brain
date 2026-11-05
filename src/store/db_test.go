package store

import (
	"context"
	"testing"
)

func TestDBManager_NilReceiverSafe(t *testing.T) {
	// [NEXUS-1] Ensure that calling store.DBManager methods on nil receiver or nil Pool does not panic
	var db *DBManager = nil

	ctx := context.Background()

	if err := db.LogAction(ctx, "s1", "m1", "a1", "o1"); err != nil {
		t.Errorf("unexpected error on nil LogAction: %v", err)
	}

	if err := db.ProposeSkill(ctx, "skill1", "rules", "trace1"); err != nil {
		t.Errorf("unexpected error on nil ProposeSkill: %v", err)
	}

	if err := db.UpsertDecision(ctx, "D1", "T1", "C1", "content1"); err != nil {
		t.Errorf("unexpected error on nil UpsertDecision: %v", err)
	}

	decisions, err := db.GetDecisions(ctx)
	if err != nil {
		t.Errorf("unexpected error on nil GetDecisions: %v", err)
	}
	if len(decisions) != 0 {
		t.Errorf("expected empty decisions list, got %d", len(decisions))
	}

	// db.Close() must not panic
	db.Close()

	// Also test when db is non-nil pointer but db.Pool is nil
	dbWithNilPool := &DBManager{Pool: nil, Schema: "test"}
	if err := dbWithNilPool.LogAction(ctx, "s1", "m1", "a1", "o1"); err != nil {
		t.Errorf("unexpected error on nil pool LogAction: %v", err)
	}
	dbWithNilPool.Close()
}

