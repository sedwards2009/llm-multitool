package mem_storage

import (
	"testing"
)

func TestSave(t *testing.T) {
	tempDir := t.TempDir()
	storage := New(tempDir)

	if len(storage.SessionOverview().SessionSummaries) != 0 {
		t.Errorf("Empty SessionOverview didn't have empty array.")
	}
	session := storage.NewSession()

	sessionOverview := storage.SessionOverview()
	if sessionOverview.SessionSummaries[0].ID != session.ID {
		t.Errorf("Couldn't find new session in the SessionOverview")
	}
}

func TestRoundTrip(t *testing.T) {
	tempDir := t.TempDir()
	storage := New(tempDir)
	session := storage.NewSession()
	session2 := storage.NewSession()

	storage2 := New(tempDir)
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

func TestNewWriteRead(t *testing.T) {
	tempDir := t.TempDir()
	storage := New(tempDir)
	session := storage.NewSession()
	session.Title = "A test"

	storage.WriteSession(session)

	id := session.ID
	session2 := storage.ReadSession(id)
	if session2.Title != "A test" {
		t.Errorf("TestNewWriteRead failed. Expected '%s', got '%s'", "A test", session2.Title)
	}
}

func TestNewResponseWriteRead(t *testing.T) {
	tempDir := t.TempDir()
	storage := New(tempDir)
	session := storage.NewSession()
	session.Title = "A test"

	id := session.ID
	response := storage.NewResponse(session)
	storage.WriteSession(session)

	session2 := storage.ReadSession(id)
	if len(session2.Responses) != 1 {
		t.Errorf("TestNewResponseWriteRead failed. Expected len(session2.Responses), got %d", len(session2.Responses))
	}

	if session2.Responses[0].ID != response.ID {
		t.Errorf("TestNewResponseWriteRead failed. Expected to find our response struct.")
	}
}
