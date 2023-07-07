package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sedwards2009/llm-workbench/internal/data"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SessionStorage struct {
	storagePath      string
	sessions         map[string]*data.Session
	sessionSummaries map[string]*data.SessionSummary
}

func NewSessionStorage(storagePath string) *SessionStorage {
	return &SessionStorage{
		storagePath:      storagePath,
		sessions:         make(map[string]*data.Session, 16),
		sessionSummaries: make(map[string]*data.SessionSummary, 16),
	}
}

func (this *SessionStorage) Scan() {
	entries, err := os.ReadDir(this.storagePath)
	if err != nil {
		log.Panicf("Error occurred while scanning storage: %v", err)
	}

	for _, entry := range entries {
		if entry.Type().IsRegular() {
			if strings.HasSuffix(entry.Name(), ".json") {
				jsonPath := filepath.Join(this.storagePath, entry.Name())
				newSession := this.readSession(jsonPath)
				if newSession != nil {
					this.cacheSession(newSession)
				}
			}
		}
	}
}

func (this *SessionStorage) NewSession() *data.Session {
	now := time.Now().UTC()
	session := &data.Session{
		ID:                uuid.NewString(),
		Title:             "(new session)",
		CreationTimestamp: now.Format(time.RFC3339),
		Responses:         []*data.Response{},
	}
	this.WriteSession(session)
	return session
}

func (this *SessionStorage) readSession(filePath string) *data.Session {
	content, err := ioutil.ReadFile(filePath)
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
	return &session
}

func (this *SessionStorage) ReadSession(id string) *data.Session {
	session := this.sessions[id]
	return session
}

func (this *SessionStorage) WriteSession(session *data.Session) {
	this.cacheSession(session)
	this.writeSession(session)
}

func (this *SessionStorage) cacheSession(session *data.Session) {
	this.sessions[session.ID] = session
	delete(this.sessionSummaries, session.ID)
}

func (this *SessionStorage) writeSession(session *data.Session) {
	fileName := session.ID + ".json"

	jsonData, err := json.Marshal(session)
	if err != nil {
		log.Fatalf("Couldn't marshal Session object: %v", err)
		panic(err)
	}

	filePath := filepath.Join(this.storagePath, fileName)
	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		log.Fatalf("Couldn't write Session object to '%s': %v", filePath, err)
		panic(err)
	}
}

func (this *SessionStorage) SessionOverview() *data.SessionOverview {
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
		Responses:         srcSession.Responses,
	}
	return copy
}

func (this *SessionStorage) sessionSummary(session *data.Session) *data.SessionSummary {
	summary, present := this.sessionSummaries[session.ID]
	if present {
		return summary
	}

	newSummary := &data.SessionSummary{
		ID:    session.ID,
		Title: session.Title,
	}
	this.sessionSummaries[session.ID] = newSummary
	return newSummary
}

func (this *SessionStorage) NewResponse(sessionId string) (*data.Response, error) {
	session := this.ReadSession(sessionId)
	if session == nil {
		return nil, errors.New("Session not found for sessionId.")
	}

	newResponsesSlice := session.Responses[:]
	now := time.Now().UTC()
	newResponse := &data.Response{
		ID:                uuid.NewString(),
		CreationTimestamp: now.Format(time.RFC3339),
		Status:            data.ResponseStatus_Pending,
		Prompt:            session.Prompt,
		Text:              "",
	}
	newResponsesSlice = append(newResponsesSlice, newResponse)

	newSession := copySession(session)
	newSession.Responses = newResponsesSlice

	this.WriteSession(newSession)
	return newResponse, nil
}

func (this *SessionStorage) DeleteResponse(sessionId string, responseId string) error {
	session := this.ReadSession(sessionId)
	if session == nil {
		return errors.New("Session not found for sessionId.")
	}

	newResponses := make([]*data.Response, 0)
	for _, response := range session.Responses {
		if response.ID != responseId {
			newResponses = append(newResponses, response)
		}
	}

	if len(newResponses) == len(session.Responses) {
		return errors.New("Response not found for responseId.")
	}

	newSession := copySession(session)
	newSession.Responses = newResponses
	this.WriteSession(newSession)
	return nil
}
