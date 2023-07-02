package storage

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestScan(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("Failed to retrieve caller information")
	}

	absFilepath, err := filepath.Abs(filename)
	if err != nil {
		panic("Error getting absolutepath:" + err.Error())
	}
	baseDir := filepath.Join(absFilepath, "../", "storage_testdata/")

	storage := NewSessionStorage(baseDir)
	storage.Scan()
}

func TestSave(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewSessionStorage(tempDir)
	session := storage.NewSession()

	sessionOverview := storage.SessionOverview()
	if sessionOverview.SessionSummaries[0].ID != session.ID {
		t.Errorf("Couldn't find new session in the SessionOverview")
	}
}

func TestRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	storage := NewSessionStorage(tempDir)
	session := storage.NewSession()
	session2 := storage.NewSession()

	storage2 := NewSessionStorage(tempDir)
	storage2.Scan()
	overview := storage2.SessionOverview()

	if len(overview.SessionSummaries) != 2 {
		t.Errorf("Round-trip failed: Overview length is wrong. Expected %d, got %d", 2, len(overview.SessionSummaries))
	}
	firstID := overview.SessionSummaries[0].ID
	secondID := overview.SessionSummaries[1].ID

	if firstID != session.ID && secondID != session.ID {
		t.Errorf("Round-trip failed: Expected %s, got %s", session.ID, overview.SessionSummaries[0].ID)
	}
	if firstID != session2.ID && secondID != session2.ID {
		t.Errorf("Round-trip failed: Expected %s, got %s", session2.ID, overview.SessionSummaries[1].ID)
	}
}
