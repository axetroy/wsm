package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		_ = conn.WriteMessage(websocket.BinaryMessage, []byte("Hello, world!"))
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
