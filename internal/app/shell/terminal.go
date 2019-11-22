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
	hostID := c.Query("id")

	controllerContext := controller.NewContextFromGinContext(c)

	if hostID == "" {
		c.String(http.StatusNotFound, "Host not provide")
		return
	}

	hostInfo := db.Host{
		Id: hostID,
	}

	if err := db.Db.Model(&hostInfo).Where("id = ?", hostID).First(&hostInfo).Error; err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	if hostInfo.OwnerID != controllerContext.Uid {
		isCollaboration := false
		for _, v := range hostInfo.Collaborations {
			if v == controllerContext.Uid {
				isCollaboration = true
			}
		}

		if isCollaboration == false {
			c.String(http.StatusNotFound, exception.NoPermission.Error())
			return
		}
	}

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
