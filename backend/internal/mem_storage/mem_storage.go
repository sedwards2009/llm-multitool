package mem_storage

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sedwards2009/llm-multitool/internal/data"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type writeMessage struct {
	session  *data.Session
	stopChan chan bool
}

type SimpleStorage struct {
	storagePath string
	sessions    map[string]*data.Session
	lock        sync.Mutex
	writeChan   chan writeMessage
}

const WRITE_BACK_QUEUE_LENGTH = 128

func New(storagePath string) *SimpleStorage {
	instance := &SimpleStorage{
		storagePath: storagePath,
		sessions:    make(map[string]*data.Session, 16),
		writeChan:   make(chan writeMessage, WRITE_BACK_QUEUE_LENGTH),
	}
	instance.scan()
	go instance.writer(instance.writeChan)
	return instance
}

func (this *SimpleStorage) GetStoragePath() string {
	return this.storagePath
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
		log.Printf("Error unmarshalling JSON from file '%s': %v\n", filePath, err)
		return nil
	}

	if session.ModelSettings == nil {
		session.ModelSettings = &data.ModelSettings{}
	}
	if session.AttachedFiles == nil {
		session.AttachedFiles = make([]*data.AttachedFile, 0)
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

	session := this.sessions[id]
	if session == nil {
		return
	}

	delete(this.sessions, id)
	os.Remove(this.sessionFilepath(id))

	for _, attachedFile := range session.AttachedFiles {
		os.Remove(filepath.Join(this.storagePath, attachedFile.Filename))
	}
}

func (this *SimpleStorage) sessionFilepath(sessionId string) string {
	return filepath.Join(this.storagePath, sessionId+".json")
}

func (this *SimpleStorage) SessionMakeAttachedFileFilepath(sessionId string, originalFilename string) (string, string) {
	extension := ""
	if originalFilename != "" {
		extension = filepath.Ext(originalFilename)
	}

	filename := sessionId + "_" + uuid.NewString() + extension
	return filename, filepath.Join(this.storagePath, filename)
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
	newAttachedFiles := make([]*data.AttachedFile, len(srcSession.AttachedFiles))
	copy(newAttachedFiles, srcSession.AttachedFiles)

	copy := &data.Session{
		ID:                srcSession.ID,
		CreationTimestamp: srcSession.CreationTimestamp,
		Title:             srcSession.Title,
		Prompt:            srcSession.Prompt,
		Responses:         copyResponses(srcSession.Responses),
		ModelSettings:     copyModelSettings(srcSession.ModelSettings),
		AttachedFiles:     newAttachedFiles,
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
		ID:                    srcResponse.ID,
		CreationTimestamp:     srcResponse.CreationTimestamp,
		Status:                srcResponse.Status,
		Messages:              copyMessages(srcResponse.Messages),
		ModelSettingsSnapshot: copyModelSettingsSnapshot(srcResponse.ModelSettingsSnapshot),
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

func copyModelSettingsSnapshot(snapshot *data.ModelSettingsSnapshot) *data.ModelSettingsSnapshot {
	if snapshot == nil {
		return nil
	}
	return &data.ModelSettingsSnapshot{
		ModelSettings: data.ModelSettings{
			ModelID:    snapshot.ModelSettings.ModelID,
			PresetID:   snapshot.PresetID,
			TemplateID: snapshot.TemplateID,
		},
		ModelName:    snapshot.ModelName,
		PresetName:   snapshot.PresetName,
		TemplateName: snapshot.TemplateName,
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
		AttachedFiles:     []*data.AttachedFile{},
	}
	this.WriteSession(session)
	return session
}

func (this *SimpleStorage) ReadSession(id string) *data.Session {
	this.lock.Lock()
	defer this.lock.Unlock()

	session := this.sessions[id]
	if session == nil {
		return nil
	}
	return copySession(session)
}

func (this *SimpleStorage) WriteSession(session *data.Session) {
	this.lock.Lock()
	defer this.lock.Unlock()

	storedSession := this.sessions[session.ID]
	if storedSession != nil {
		this.garbageCollectFiles(storedSession, session)
	}

	sessionCopy := copySession(session)
	this.cacheSession(sessionCopy)

	this.writeChan <- writeMessage{session: sessionCopy}
}

func (this *SimpleStorage) garbageCollectFiles(previousSession *data.Session, newSesion *data.Session) {
	previousFilenames := previousSession.GetAttachedFileNames()
	currentFilenames := newSesion.GetAttachedFileNames()
	sort.Strings(currentFilenames)
	for _, filename := range previousFilenames {
		isFound := sort.SearchStrings(currentFilenames, filename) < len(currentFilenames)
		if !isFound {
			os.Remove(filepath.Join(this.GetStoragePath(), filename))
		}
	}
}

func (this *SimpleStorage) Stop() {
	this.lock.Lock()
	defer this.lock.Unlock()

	stopChan := make(chan bool, 0)
	this.writeChan <- writeMessage{session: nil, stopChan: stopChan}
	<-stopChan
}

type SessionDeadline struct {
	deadline time.Time
	session  *data.Session
}

const WRITE_BACK_DELAY = 3 * time.Second

/*
 * Goroutine which receives Sessions to save and writes them back to disk.
 */
func (this *SimpleStorage) writer(workChan chan writeMessage) {
	workPool := make(map[string]SessionDeadline)

	setSession := func(session *data.Session) {
		deadlineSession, ok := workPool[session.ID]
		if ok {
			deadlineSession.session = session
			workPool[session.ID] = deadlineSession
		} else {
			workPool[session.ID] = SessionDeadline{deadline: time.Now().Add(WRITE_BACK_DELAY), session: session}
		}
	}

	running := true
	var stopChan chan bool

	for running {
		timeout := time.After(1 * time.Second)
		select {
		case msg := <-workChan:
			if msg.session != nil {
				setSession(msg.session)
			} else {
				running = false
				stopChan = msg.stopChan
			}
		case <-timeout:
			break
		}

		for running {
			select {
			case msg := <-workChan:
				if msg.session != nil {
					setSession(msg.session)
				} else {
					running = false
					stopChan = msg.stopChan
				}
			default:
				break
			}

			if len(workPool) == 0 {
				break
			}

			now := time.Now()
			for _, deadlineSession := range workPool {
				if now.After(deadlineSession.deadline) {
					log.Printf("Writing %s back to disk.\n", deadlineSession.session.ID)
					this.writeToDisk(deadlineSession.session)
					delete(workPool, deadlineSession.session.ID)
					break
				}
			}
		}
	}

	for _, deadlineSession := range workPool {
		log.Printf("Writing %s back to disk.\n", deadlineSession.session.ID)
		this.writeToDisk(deadlineSession.session)
	}
	stopChan <- true
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
