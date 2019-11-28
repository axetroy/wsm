package shell

import (
	"net/http"

	"github.com/axetroy/terminal/internal/app/db"
	"github.com/axetroy/terminal/internal/app/exception"
	"github.com/axetroy/terminal/internal/library/controller"
	"github.com/axetroy/terminal/internal/library/session"
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

func (s *Service) StartTerminalRouter(c *gin.Context) {
	// search for host
	hostID := c.Query("host_id")

	ctx := controller.NewContextFromGinContext(c)

	recordInfo := db.HostRecord{Id: hostID, UserID: ctx.Uid}

	if err := db.Db.Where(&recordInfo).Preload("Host").First(&recordInfo).Error; err != nil {
		c.String(http.StatusBadRequest, exception.NoPermission.Error())
		return
	}

	hostInfo := recordInfo.Host

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	terminal, err := session.NewTerminal(session.Config{
		Host:     hostInfo.Host,
		Port:     hostInfo.Port,
		Username: hostInfo.Username,
		Password: hostInfo.Password,
		Width:    200,
		Height:   50,
	})

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
