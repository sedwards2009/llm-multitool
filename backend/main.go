package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"sedwards2009/llm-workbench/internal/argsparser"
	"sedwards2009/llm-workbench/internal/broadcaster"
	"sedwards2009/llm-workbench/internal/data"
	"sedwards2009/llm-workbench/internal/data/responsestatus"
	"sedwards2009/llm-workbench/internal/data/role"
	"sedwards2009/llm-workbench/internal/engine"
	"sedwards2009/llm-workbench/internal/mem_storage"
	"sedwards2009/llm-workbench/internal/presets"
	"sedwards2009/llm-workbench/internal/template"

	"github.com/bobg/go-generics/v2/slices"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

//go:embed resources/*
var staticFS embed.FS

var logger gin.HandlerFunc = nil
var sessionStorage *mem_storage.SimpleStorage = nil
var llmEngine *engine.Engine = nil
var presetDatabase *presets.PresetDatabase = nil
var sessionBroadcaster *broadcaster.Broadcaster = nil
var templates *template.TemplateDatabase = nil

func setupStorage(storagePath string) *mem_storage.SimpleStorage {
	return mem_storage.New(storagePath)
}

func setupEngine(configPath string, presetDatabase *presets.PresetDatabase) *engine.Engine {
	return engine.NewEngine(configPath, presetDatabase)
}

func setupTemplates(templatesPath string) *template.TemplateDatabase {
	return template.NewTemplateDatabase(templatesPath)
}

func setupPresets(presetsPath string) *presets.PresetDatabase {
	return presets.MakePresetDatabase(presetsPath)
}

func setupBroadcaster() *broadcaster.Broadcaster {
	return broadcaster.NewBroadcaster()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	logger := gin.Logger()
	r.Use(logger)

	r.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"*"},
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           10 * time.Second,
	}))

	r.GET("/", handleIndex)
	r.GET("/assets/*filepath", handleAssets)
	r.GET("/session/:sessionId", handleIndex)
	r.GET("/api/ping", handlePing)
	r.GET("/api/session", handleSessionOverview)
	r.POST("/api/session", handleNewSession)
	r.GET("/api/session/:sessionId", handleSessionGet)
	r.PUT("/api/session/:sessionId/prompt", handleSessionPromptPut)
	r.DELETE("/api/session/:sessionId", handleSessionDelete)
	r.POST("/api/session/:sessionId/response", handleResponsePost)
	r.GET("/api/session/:sessionId/changes", handleSessionChangesGet)
	r.DELETE("/api/session/:sessionId/response/:responseId", handleResponseDelete)
	r.GET("/api/model", handleModelOverviewGet)
	r.POST("/api/model/scan", handleModelScanPost)
	r.PUT("/api/session/:sessionId/modelSettings", handleSessionModelSettingsPut)
	r.POST("/api/session/:sessionId/response/:responseId/message", handleNewMessagePost)
	r.POST("/api/session/:sessionId/response/:responseId/continue", handleMessageContinuePost)
	r.GET("/api/template", handleTemplateOverviewGet)
	r.GET("/api/preset", handlePresetOverviewGet)

	return r
}

func handleIndex(c *gin.Context) {
	// Work-around for one of the dumbest problems regarding index.html
	// See: https://github.com/gin-gonic/gin/issues/2654
	contents, _ := staticFS.ReadFile("resources/index.html")
	c.Header("Content-Type", "text/html")
	c.Data(http.StatusOK, "text/html", contents)
}

func handleAssets(c *gin.Context) {
	c.FileFromFS(path.Join("resources/", c.Request.URL.Path), http.FS(staticFS))
}

func handlePing(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

var upgrader = websocket.Upgrader{
	//Solve "request origin not allowed by Upgrader.CheckOrigin"
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	// Time allowed to write a message to the websocket.
	websocketWriteWait  = 10 * time.Second
	websocketPongWait   = 5 * time.Second
	changeThrottleDelay = 250 * time.Millisecond
	websocketPingPeriod = (websocketPongWait * 9) / 10
)

func handleSessionChangesGet(c *gin.Context) {
	sessionId := c.Params.ByName("sessionId")
	session := sessionStorage.ReadSession(sessionId)
	if session == nil {
		c.String(http.StatusNotFound, "Session not found")
	}

	wsSession, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer wsSession.Close()

	changeChan := make(chan string, 16)
	sessionBroadcaster.Register(sessionId, changeChan)
	pingTicker := time.NewTicker(websocketPingPeriod)
	defer func() {
		sessionBroadcaster.Unregister(changeChan)
		pingTicker.Stop()
		close(changeChan)
	}()

	go websocketReader(wsSession)

	throttleTimer := time.NewTimer(changeThrottleDelay)
	isChangeWaiting := false
	for {
		select {
		case <-changeChan:
			isChangeWaiting = true

		case <-throttleTimer.C:
			if isChangeWaiting {
				isChangeWaiting = false
				wsSession.SetWriteDeadline(time.Now().Add(websocketWriteWait))

				if err := wsSession.WriteMessage(websocket.TextMessage, []byte("changed")); err != nil {
					wsSession.Close()
					if websocket.IsCloseError(err, websocket.CloseGoingAway) {
						log.Printf("Client disconnected for session ID %s.", sessionId)
					} else {
						log.Printf("Writing error for session ID %s: %v.", sessionId, err)
					}
					return
				}
			}
			throttleTimer.Reset(changeThrottleDelay)

		case <-pingTicker.C:
			wsSession.SetWriteDeadline(time.Now().Add(websocketWriteWait))
			if err := wsSession.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("Client disconnected for session ID %s.", sessionId)
				return
			}
		}
	}
}

func websocketReader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(websocketPongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(websocketPongWait))
		return nil
	})
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func handleSessionOverview(c *gin.Context) {
	sessionOverview := sessionStorage.SessionOverview()
	c.JSON(http.StatusOK, sessionOverview)
}

// Create a new session.
func handleNewSession(c *gin.Context) {
	session := sessionStorage.NewSession()
	if session == nil {
		c.String(http.StatusNotFound, "Session couldn't be created")
		return
	}

	session.ModelSettings.ModelID = llmEngine.DefaultID()
	session.ModelSettings.PresetID = presetDatabase.DefaultID()
	session.ModelSettings.TemplateID = templates.DefaultID()
	log.Printf("presetDatabase.DefaultID(): %s", presetDatabase.DefaultID())
	sessionStorage.WriteSession(session)
	c.JSON(http.StatusOK, session)
}

// Get a full session and its data.
func handleSessionGet(c *gin.Context) {
	sessionId := c.Params.ByName("sessionId")

	session := sessionStorage.ReadSession(sessionId)
	if session != nil {
		c.JSON(http.StatusOK, session)
		return
	}
	c.String(http.StatusNotFound, "Session not found")
}

func handleSessionDelete(c *gin.Context) {
	sessionId := c.Params.ByName("sessionId")
	session := sessionStorage.ReadSession(sessionId)
	if session == nil {
		c.String(http.StatusNotFound, "Session not found")
		return
	}
	sessionStorage.DeleteSession(sessionId)
	c.Status(http.StatusNoContent)
}

func handleSessionPromptPut(c *gin.Context) {
	sessionId := c.Params.ByName("sessionId")
	session := sessionStorage.ReadSession(sessionId)
	if session == nil {
		c.String(http.StatusNotFound, "Session not found")
		return
	}

	var data struct {
		Value string `json:"value"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the JSON PUT body.")
		return
	}

	session.Prompt = data.Value
	sessionStorage.WriteSession(session)

	c.JSON(http.StatusOK, session)
}

// Trigger the generation of a new response in a session used the current model and prompt.
func handleResponsePost(c *gin.Context) {
	sessionId := c.Params.ByName("sessionId")
	session := sessionStorage.ReadSession(sessionId)
	if session == nil {
		c.String(http.StatusNotFound, "Session not found")
		return
	}

	session.Title = templates.MakeTitle(session.ModelSettings.TemplateID, session.Prompt)
	response := CreateNewResponse(session)

	responseId := response.ID

	formattedPrompt := templates.ApplyTemplate(session.ModelSettings.TemplateID, session.Prompt)

	response.Messages = append(response.Messages, data.Message{
		ID:   uuid.NewString(),
		Role: role.User,
		Text: formattedPrompt,
	})
	response.Messages = append(response.Messages, data.Message{
		ID:   uuid.NewString(),
		Role: role.Assistant,
		Text: "",
	})
	sessionStorage.WriteSession(session)

	appendFunc := func(text string) bool {
		success := appendToLastMessage(sessionId, responseId, text)
		sessionBroadcaster.Send(sessionId, "changed")
		return success
	}

	completeFunc := func() {
		sessionBroadcaster.Send(sessionId, "changed")
	}

	setStatusFunc := func(status responsestatus.ResponseStatus) {
		editResponse(sessionId, responseId, func(session *data.Session, response *data.Response) bool {
			response.Status = status
			return true
		})
		sessionBroadcaster.Send(sessionId, "changed")
	}

	llmEngine.Enqueue(response.Messages, appendFunc, completeFunc, setStatusFunc, session.ModelSettings)
	c.JSON(http.StatusOK, response)
}

func editResponse(sessionId string, responseId string, callback func(*data.Session, *data.Response) bool) bool {
	session := sessionStorage.ReadSession(sessionId)
	if session == nil {
		return false
	}
	for _, r := range session.Responses {
		if r.ID == responseId {
			if callback(session, r) {
				sessionStorage.WriteSession(session)
			}
			return true
		}
	}
	return false
}

func appendToLastMessage(sessionId string, responseId string, text string) bool {
	return editResponse(sessionId, responseId, func(session *data.Session, response *data.Response) bool {
		response.Messages[len(response.Messages)-1].Text += text
		return true
	})
}

func handleResponseDelete(c *gin.Context) {
	sessionId := c.Params.ByName("sessionId")
	responseId := c.Params.ByName("responseId")

	session := sessionStorage.ReadSession(sessionId)
	if session == nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to find session with ID %s\n", sessionId))
		return
	}

	originalLength := len(session.Responses)
	session.Responses = slices.Filter(session.Responses, func(r *data.Response) bool {
		return r.ID != responseId
	})
	if originalLength == len(session.Responses) {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Unable to find response with ID %s\n", responseId))
		return
	}
	sessionStorage.WriteSession(session)

	c.Status(http.StatusNoContent)
}

func handleMessageContinuePost(c *gin.Context) {
	sessionId := c.Params.ByName("sessionId")
	session := sessionStorage.ReadSession(sessionId)
	if session == nil {
		c.String(http.StatusNotFound, "Session not found")
		return
	}

	responseId := c.Params.ByName("responseId")
	response := getResponseFromSessionByID(session, responseId)
	if response == nil {
		c.String(http.StatusNotFound, "Response not found")
		return
	}

	appendFunc := func(text string) bool {
		success := appendToLastMessage(sessionId, responseId, text)
		sessionBroadcaster.Send(sessionId, "changed")
		return success
	}

	completeFunc := func() {
		sessionBroadcaster.Send(sessionId, "changed")
	}

	setStatusFunc := func(status responsestatus.ResponseStatus) {
		editResponse(sessionId, responseId, func(session *data.Session, response *data.Response) bool {
			response.Status = status
			return true
		})
		sessionBroadcaster.Send(sessionId, "changed")
	}

	llmEngine.Enqueue(response.Messages, appendFunc, completeFunc, setStatusFunc,
		&response.ModelSettingsSnapshot.ModelSettings)
	c.JSON(http.StatusOK, response)
}

func handleModelOverviewGet(c *gin.Context) {
	modelOverview := llmEngine.ModelOverview()
	c.JSON(http.StatusOK, modelOverview)
}

func handleModelScanPost(c *gin.Context) {
	llmEngine.ScanModels()
	handleModelOverviewGet(c)
}

func handleSessionModelSettingsPut(c *gin.Context) {
	sessionId := c.Params.ByName("sessionId")
	session := sessionStorage.ReadSession(sessionId)
	if session == nil {
		c.String(http.StatusNotFound, "Session not found")
		return
	}

	data := &data.ModelSettings{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the JSON PUT body.")
		return
	}

	if !llmEngine.ValidateModelSettings(data) {
		c.String(http.StatusBadRequest, "An invalid ModelID value was given in the PUT body.")
		return
	}

	if !presetDatabase.Exists(data.PresetID) {
		c.String(http.StatusBadRequest, "An invalid PresetID was given in the PUT body.")
		return
	}

	session.ModelSettings = data
	sessionStorage.WriteSession(session)

	c.JSON(http.StatusOK, session)
}

func handleNewMessagePost(c *gin.Context) {
	sessionId := c.Params.ByName("sessionId")
	responseId := c.Params.ByName("responseId")
	var postData struct {
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&postData); err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the JSON POST body.")
		return
	}

	var foundSession *data.Session
	var foundResponse *data.Response

	if !editResponse(sessionId, responseId, func(session *data.Session, response *data.Response) bool {
		foundResponse = response
		foundSession = session
		response.Messages = append(response.Messages, data.Message{
			ID:   uuid.NewString(),
			Role: role.User,
			Text: postData.Value,
		})
		response.Messages = append(response.Messages, data.Message{
			ID:   uuid.NewString(),
			Role: role.Assistant,
			Text: "",
		})
		return true
	}) {
		c.String(http.StatusNotFound, "Response not found")
		return
	}

	appendFunc := func(text string) bool {
		success := appendToLastMessage(sessionId, responseId, text)
		sessionBroadcaster.Send(sessionId, "changed")
		return success
	}

	completeFunc := func() {
		sessionBroadcaster.Send(sessionId, "changed")
	}

	setStatusFunc := func(status responsestatus.ResponseStatus) {
		editResponse(sessionId, responseId, func(session *data.Session, response *data.Response) bool {
			response.Status = status
			return true
		})
		sessionBroadcaster.Send(sessionId, "changed")
	}

	llmEngine.Enqueue(foundResponse.Messages, appendFunc, completeFunc, setStatusFunc, foundSession.ModelSettings)
	c.JSON(http.StatusOK, foundResponse)
}

func handleTemplateOverviewGet(c *gin.Context) {
	templateOverview := templates.TemplateOverview()
	c.JSON(http.StatusOK, templateOverview)
}

func getResponseFromSessionByID(session *data.Session, responseID string) *data.Response {
	responseIndex := slices.IndexFunc(session.Responses, func(r *data.Response) bool {
		return responseID == r.ID
	})
	if responseIndex == -1 {
		return nil
	}

	return session.Responses[responseIndex]
}

func handlePresetOverviewGet(c *gin.Context) {
	presetOverview := presetDatabase.PresetOverview()
	c.JSON(http.StatusOK, presetOverview)
}

func CreateNewResponse(session *data.Session) *data.Response {
	now := time.Now().UTC()

	preset := presetDatabase.Get(session.ModelSettings.PresetID)
	template := templates.Get(session.ModelSettings.TemplateID)
	model := llmEngine.GetModel(session.ModelSettings.ModelID)

	newResponse := &data.Response{
		ID:                uuid.NewString(),
		CreationTimestamp: now.Format(time.RFC3339),
		Status:            responsestatus.Pending,
		Messages:          []data.Message{},
		ModelSettingsSnapshot: &data.ModelSettingsSnapshot{
			ModelSettings: data.ModelSettings{
				ModelID:    session.ModelSettings.ModelID,
				PresetID:   session.ModelSettings.PresetID,
				TemplateID: session.ModelSettings.TemplateID,
			},
			ModelName:    model.Name,
			PresetName:   preset.Name,
			TemplateName: template.Name,
		},
	}
	session.Responses = append(session.Responses, newResponse)
	return newResponse
}

func main() {
	config := argsparser.Parse()
	if config == nil {
		return
	}

	sessionStorage = setupStorage(config.StoragePath)
	presetDatabase = setupPresets(config.PresetsPath)
	llmEngine = setupEngine(config.ConfigFilePath, presetDatabase)
	sessionBroadcaster = setupBroadcaster()
	templates = setupTemplates(config.TemplatesPath)
	r := setupRouter()
	fmt.Printf("\n    Starting server on http://%s\n\n", config.Address)
	r.Run(config.Address)
}
