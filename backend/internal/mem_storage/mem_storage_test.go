package mem_storage

import (
	"os"
	"sedwards2009/llm-multitool/internal/data"
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
	storage.Stop()

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

func TestFileUpload(t *testing.T) {
	tempDir := t.TempDir()
	t.Logf("tempDir: %s\n", tempDir)
	storage := New(tempDir)
	session := storage.NewSession()
	sessionID := session.ID
	session.Title = "A test"

	storage.WriteSession(session)
	storage.Stop() // Force everything to disk.

	if countFiles(t, tempDir) != 1 {
		t.Errorf("Wrong number of files found in %s.", tempDir)
	}

	storage = New(tempDir)
	session2 := storage.ReadSession(sessionID)

	filename, filepath := storage.SessionMakeAttachedFileFilepath(sessionID, "empty.txt")
	os.WriteFile(filepath, []byte("Some text"), 0666)
	session2.AttachedFiles = append(session2.AttachedFiles,
		&data.AttachedFile{Filename: filename, MimeType: "text/plain", OriginalFilename: "empty.txt"})

	storage.WriteSession(session2)
	storage.Stop()

	expectCountFiles(t, tempDir, 2)

	storage = New(tempDir)
	session3 := storage.ReadSession(sessionID)
	session3.AttachedFiles = []*data.AttachedFile{}
	storage.WriteSession(session3)
	storage.Stop()

	expectCountFiles(t, tempDir, 1)
}

func expectCountFiles(t *testing.T, dirpath string, expected int) {
	found := countFiles(t, dirpath)
	if found != expected {
		t.Errorf("Expected %d entries in %s, but found %d", expected, dirpath, found)
	}
}

func countFiles(t *testing.T, dirpath string) int {
	entries, err := os.ReadDir(dirpath)
	if err != nil {
		t.Errorf("Unable to count files in %s, %v", dirpath, err)
		return -1
	}
	return len(entries)
}
