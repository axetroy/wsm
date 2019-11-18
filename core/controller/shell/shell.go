package shell

import (
	"net/http"

	"github.com/axetroy/terminal/session"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WebsocketStream struct {
	conn        *websocket.Conn
	messageType int
}

func (r WebsocketStream) Read(p []byte) (n int, err error) {
	t, message, err := r.conn.ReadMessage()

	r.messageType = t

	copy(p, message)

	n = len(message)

	return
}

func (r WebsocketStream) Write(p []byte) (n int, err error) {
	err = r.conn.WriteMessage(r.messageType, p)

	n = len(p)

	return
}

var upgrader = websocket.Upgrader{} // use default options

func DemoRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func StartRouter(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	terminal, err := session.New()

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	conn.SetCloseHandler(func(code int, text string) error {
		return terminal.Close()
	})

	stream := WebsocketStream{conn: conn, messageType: 1}

	err = terminal.Connect(stream, stream, stream)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}
