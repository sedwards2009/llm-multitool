package mem_storage

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sedwards2009/llm-workbench/internal/data"
	"sedwards2009/llm-workbench/internal/data/responsestatus"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type SimpleStorage struct {
	storagePath string
	sessions    map[string]*data.Session
	lock        sync.Mutex
}

func New(storagePath string) *SimpleStorage {
	instance := &SimpleStorage{
		storagePath: storagePath,
		sessions:    make(map[string]*data.Session, 16),
	}
	instance.scan()
	return instance
}

func (this *SimpleStorage) scan() {
	entries, err := os.ReadDir(this.storagePath)
	if err != nil {
		log.Panicf("Error occurred while scanning storage: %v", err)
	}

	for _, entry := range entries {
		if entry.Type().IsRegular() {
			if strings.HasSuffix(entry.Name(), ".json") {
				jsonPath := filepath.Join(this.storagePath, entry.Name())
				newSession := this.readSessionFromFile(jsonPath)
				if newSession != nil {
					this.cacheSession(newSession)
				}
			}
		}
	}
}

func (this *SimpleStorage) readSessionFromFile(filePath string) *data.Session {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Println("Error reading file:", err)
		return nil
	}

	var session data.Session
	err = json.Unmarshal([]byte(content), &session)
	if err != nil {
		log.Println("Error unmarshalling JSON", err)
		return nil
	}

	if session.ModelSettings == nil {
		session.ModelSettings = &data.ModelSettings{}
	}
	for _, r := range session.Responses {
		if r.Messages == nil {
			r.Messages = []data.Message{}
		}
	}
	return &session
}

func (this *SimpleStorage) cacheSession(session *data.Session) {
	this.sessions[session.ID] = session
}

func (this *SimpleStorage) DeleteSession(id string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	delete(this.sessions, id)
	os.Remove(this.sessionFilepath(id))
}

func (this *SimpleStorage) sessionFilepath(sessionId string) string {
	return filepath.Join(this.storagePath, sessionId+".json")
}

func (this *SimpleStorage) sessionSummary(session *data.Session) *data.SessionSummary {
	newSummary := &data.SessionSummary{
		ID:                session.ID,
		Title:             session.Title,
		CreationTimestamp: session.CreationTimestamp,
	}
	return newSummary
}

func (this *SimpleStorage) SessionOverview() *data.SessionOverview {
	this.lock.Lock()
	defer this.lock.Unlock()

	sessionOverview := new(data.SessionOverview)
	sessionOverview.SessionSummaries = make([]*data.SessionSummary, 0)

	for s := range this.sessions {
		sessionSummary := this.sessionSummary(this.sessions[s])
		sessionOverview.SessionSummaries = append(sessionOverview.SessionSummaries, sessionSummary)
	}

	sort.SliceStable(sessionOverview.SessionSummaries,
		func(i, j int) bool {
			return sessionOverview.SessionSummaries[i].CreationTimestamp < sessionOverview.SessionSummaries[j].CreationTimestamp
		})

	return sessionOverview
}

func copySession(srcSession *data.Session) *data.Session {
	copy := &data.Session{
		ID:                srcSession.ID,
		CreationTimestamp: srcSession.CreationTimestamp,
		Title:             srcSession.Title,
		Prompt:            srcSession.Prompt,
		Responses:         copyResponses(srcSession.Responses),
		ModelSettings:     copyModelSettings(srcSession.ModelSettings),
	}
	return copy
}

func copyResponses(srcResponses []*data.Response) []*data.Response {
	result := []*data.Response{}
	for _, r := range srcResponses {
		result = append(result, copyResponse(r))
	}
	return result
}

func copyResponse(srcResponse *data.Response) *data.Response {
	return &data.Response{
		ID:                srcResponse.ID,
		CreationTimestamp: srcResponse.CreationTimestamp,
		Status:            srcResponse.Status,
		Messages:          copyMessages(srcResponse.Messages),
	}
}

func copyMessages(srcMessages []data.Message) []data.Message {
	result := []data.Message{}
	for _, m := range srcMessages {
		result = append(result, m)
	}
	return result
}

func copyModelSettings(settings *data.ModelSettings) *data.ModelSettings {
	return &data.ModelSettings{
		ModelID:    settings.ModelID,
		PresetID:   settings.PresetID,
		TemplateID: settings.TemplateID,
	}
}

func (this *SimpleStorage) NewSession() *data.Session {
	now := time.Now().UTC()
	session := &data.Session{
		ID:                uuid.NewString(),
		Title:             "(new session)",
		CreationTimestamp: now.Format(time.RFC3339),
		Responses:         []*data.Response{},
		ModelSettings:     &data.ModelSettings{},
	}
	this.WriteSession(session)
	return session
}

func (this *SimpleStorage) ReadSession(id string) *data.Session {
	session := this.sessions[id]
	return copySession(session)
}

func (this *SimpleStorage) WriteSession(session *data.Session) {
	this.lock.Lock()
	defer this.lock.Unlock()

	sessionCopy := copySession(session)
	this.cacheSession(sessionCopy)
	this.writeToDisk(sessionCopy)
}

func (this *SimpleStorage) writeToDisk(session *data.Session) {
	jsonData, err := json.Marshal(session)
	if err != nil {
		log.Fatalf("Couldn't marshal Session object: %v", err)
		panic(err)
	}

	filePath := this.sessionFilepath(session.ID)
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		log.Fatalf("Couldn't write Session object to '%s': %v", filePath, err)
		panic(err)
	}
}

func (this *SimpleStorage) NewResponse(session *data.Session) *data.Response {
	now := time.Now().UTC()
	newResponse := &data.Response{
		ID:                uuid.NewString(),
		CreationTimestamp: now.Format(time.RFC3339),
		Status:            responsestatus.Pending,
		Messages:          []data.Message{},
	}
	session.Responses = append(session.Responses, newResponse)
	return newResponse
}
