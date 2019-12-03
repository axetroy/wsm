package shell

import (
	"net/http"
	"strconv"

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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Service) StartTerminalRouter(c *gin.Context) {
	// search for host

	var (
		hostID = c.Param("host_id")
		rows   = 25
		cols   = 80
		// internal
		rowsStr = c.Param("rows")
		colsStr = c.Param("cols")
	)

	if hostID == "" {
		c.String(http.StatusNotFound, exception.NoData.Error())
		return
	}

	if rowsStr != "" {
		if i, err := strconv.ParseInt(rowsStr, 0, 10); err != nil {
			c.String(http.StatusNotFound, exception.InvalidParams.Error())
		} else {
			rows = int(i)
		}
	}

	if colsStr != "" {
		if i, err := strconv.ParseInt(colsStr, 0, 10); err != nil {
			c.String(http.StatusNotFound, exception.InvalidParams.Error())
		} else {
			cols = int(i)
		}
	}

	ctx := controller.NewContextFromGinContext(c)

	recordInfo := db.HostRecord{HostID: hostID, UserID: ctx.Uid}

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
		Width:    cols,
		Height:   rows,
	})

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	terminal.SetCloseHandler(func() error {
		return conn.Close()
	})

	conn.SetCloseHandler(func(code int, text string) error {
		return terminal.Close()
	})

	stream := WebsocketStream{conn: conn, messageType: websocket.BinaryMessage}

	err = terminal.Connect(stream, stream, stream)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
}
