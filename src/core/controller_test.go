package core

import (
	"context"
	"testing"
)

type dummyLogger struct{}

func (d *dummyLogger) Info(format string, args ...interface{})    {}
func (d *dummyLogger) Warning(format string, args ...interface{}) {}
func (d *dummyLogger) Error(format string, args ...interface{})   {}

type mockEvolutionPipeline struct {
	lastPromptName string
	lastPromptPath string
	lastDryRun     bool
}

func (m *mockEvolutionPipeline) RunEvolution(ctx context.Context, promptName, currentPromptPath string, dryRun bool) error {
	m.lastPromptName = promptName
	m.lastPromptPath = currentPromptPath
	m.lastDryRun = dryRun
	return nil
}

func TestController_NilDB_LogMilestone(t *testing.T) {
	// [NEXUS-1] Test that calling LogMilestone with nil DB does not panic
	ctrl := NewStrategicController(nil, &mockEvolutionPipeline{}, &dummyLogger{})
	res, err := ctrl.LogMilestone(context.Background(), "session1", "m1", "action1", "ok")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res["status"] != "success" {
		t.Errorf("expected status 'success', got %v", res["status"])
	}
}

func TestController_ProcessCommand_RunEvolution_DryRunParsing(t *testing.T) {
	// [NEXUS-2] Test that --dry-run is parsed as a flag and not mistaken as promptName
	mockEvol := &mockEvolutionPipeline{}
	ctrl := NewStrategicController(nil, mockEvol, &dummyLogger{})

	// Scenario A: --dry-run only
	_, err := ctrl.ProcessCommand(context.Background(), "run-evolution", []string{"--dry-run"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mockEvol.lastDryRun {
		t.Errorf("expected dryRun to be true")
	}
	if mockEvol.lastPromptName != "Python-Integration-Specialist" {
		t.Errorf("expected default promptName 'Python-Integration-Specialist', got %q", mockEvol.lastPromptName)
	}
	if mockEvol.lastPromptPath != "" {
		t.Errorf("expected empty promptPath, got %q", mockEvol.lastPromptPath)
	}

	// Scenario B: --dry-run with prompt name and path
	mockEvol = &mockEvolutionPipeline{}
	ctrl = NewStrategicController(nil, mockEvol, &dummyLogger{})
	_, err = ctrl.ProcessCommand(context.Background(), "run-evolution", []string{"--dry-run", "CustomPrompt", "/path/to/prompt.md"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mockEvol.lastDryRun {
		t.Errorf("expected dryRun to be true")
	}
	if mockEvol.lastPromptName != "CustomPrompt" {
		t.Errorf("expected 'CustomPrompt', got %q", mockEvol.lastPromptName)
	}
	if mockEvol.lastPromptPath != "/path/to/prompt.md" {
		t.Errorf("expected '/path/to/prompt.md', got %q", mockEvol.lastPromptPath)
	}

	// Scenario C: flags after positional arguments
	mockEvol = &mockEvolutionPipeline{}
	ctrl = NewStrategicController(nil, mockEvol, &dummyLogger{})
	_, err = ctrl.ProcessCommand(context.Background(), "run-evolution", []string{"CustomPrompt", "--dry-run"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mockEvol.lastDryRun {
		t.Errorf("expected dryRun to be true")
	}
	if mockEvol.lastPromptName != "CustomPrompt" {
		t.Errorf("expected 'CustomPrompt', got %q", mockEvol.lastPromptName)
	}
}

