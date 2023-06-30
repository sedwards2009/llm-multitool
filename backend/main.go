package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var logger gin.HandlerFunc = nil

func setupRouter() *gin.Engine {
	r := gin.Default()
	logger := gin.Logger()
	r.Use(logger)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/ping", handlePing)
	r.Any("/websocket", wsHandler)
	r.GET("/session", handleSessionOverview)
	r.GET("/session/:sessionId", handleSessionGet)

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

var sessionOverview SessionOverview = SessionOverview{
	[]SessionSummary{
		{ID: "1111", Title: "Qt event questions"},
		{ID: "2222", Title: "Simple React Component"},
	},
}

func handleSessionOverview(c *gin.Context) {
	c.JSON(http.StatusOK, sessionOverview)
}

func handleSessionGet(c *gin.Context) {
	sessionId := c.Params.ByName("sessionId")

	if sessionId == "1111" {
		c.JSON(http.StatusOK, Session{
			ID:        "1111",
			Title:     "Qt event questions",
			Prompt:    "Which Qt event is for a window gaining focus?",
			Responses: []Response{},
		})
		return
	}
	if sessionId == "2222" {
		c.JSON(http.StatusOK, Session{
			ID:        "2222",
			Title:     "Simple React component",
			Prompt:    "Write out a simple React component.",
			Responses: []Response{},
		})
		return
	}
	c.String(http.StatusNotFound, "Session not found")
}

func main() {
	// parsedArgs, errorString := argsparser.Parse(&os.Args)

	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
