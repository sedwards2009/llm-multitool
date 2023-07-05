package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"sedwards2009/llm-workbench/internal/storage"
)

var logger gin.HandlerFunc = nil
var sessionStorage *storage.ConcurrentSessionStorage = nil

func setupStorage() {
	sessionStorage = storage.NewConcurrentSessionStorage("/home/sbe/devel/llm-workbench/data")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	logger := gin.Logger()
	r.Use(logger)

	r.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"*"},
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           10 * time.Second,
	}))

	r.GET("/ping", handlePing)
	r.Any("/websocket", wsHandler)
	r.GET("/session", handleSessionOverview)
	r.POST("/session", handleNewSession)
	r.GET("/session/:sessionId", handleSessionGet)
	r.PUT("/session/:sessionId/prompt", handleSessionPromptPut)
	r.POST("/session/:sessionId/response", handleResponsePost)

	return r
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

func wsHandler(c *gin.Context) { //Usually use c *gin.Context
	wsSession, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer wsSession.Close()
	echo(wsSession)
}

func echo(wsSession *websocket.Conn) {
	for { //An endlessloop
		messageType, messageContent, err := wsSession.ReadMessage()
		if err != nil {
			wsSession.Close()
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Client disconnected")
			} else {
				log.Printf("Reading Error in %s.", err)
			}
			break //To escape from the endless loop
		}
		if messageType == 1 {
			log.Printf("Recv:%s", messageContent)
		}
	}
}

func handleSessionOverview(c *gin.Context) {
	sessionOverview := sessionStorage.SessionOverview()
	c.JSON(http.StatusOK, sessionOverview)
}

func handleNewSession(c *gin.Context) {
	session := sessionStorage.NewSession()
	if session != nil {
		c.JSON(http.StatusOK, session)
		return
	}
	c.String(http.StatusNotFound, "Session couldn't be created")
}

func handleSessionGet(c *gin.Context) {
	sessionId := c.Params.ByName("sessionId")

	session := sessionStorage.ReadSession(sessionId)
	if session != nil {
		c.JSON(http.StatusOK, session)
		return
	}
	c.String(http.StatusNotFound, "Session not found")
}

func handleSessionPromptPut(c *gin.Context) {
	sessionId := c.Params.ByName("sessionId")
	session := sessionStorage.ReadSession(sessionId)
	if session == nil {
		c.String(http.StatusNotFound, "Session not found")
		return
	}

	var data struct {
		Prompt string `json:"prompt"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.String(http.StatusBadRequest, "Couldn't parse the JSON PUT body.")
		return
	}

	session.Prompt = data.Prompt
	sessionStorage.WriteSession(session)

	c.JSON(http.StatusOK, session)
}

func handleResponsePost(c *gin.Context) {
	sessionId := c.Params.ByName("sessionId")
	session := sessionStorage.ReadSession(sessionId)
	if session == nil {
		c.String(http.StatusNotFound, "Session not found")
		return
	}

	response, err := sessionStorage.NewResponse(sessionId)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error occured while creating new response: %v", err))
		return
	}

	c.JSON(http.StatusOK, response)
}

func main() {
	// parsedArgs, errorString := argsparser.Parse(&os.Args)
	setupStorage()
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
