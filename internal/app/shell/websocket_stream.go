package shell

import (
	"time"

	"github.com/gorilla/websocket"
)

// 定义书记记录器
type Recorder struct {
	Time time.Time `json:"time"` // 时间戳
	Data []byte    `json:"data"` // 数据
}

func NewWebSocketSteam(connection *websocket.Conn) *WebsocketStream {
	return &WebsocketStream{
		conn:        connection,
		messageType: websocket.BinaryMessage,
		UpdatedAt:   time.Now(),
	}
}

type WebsocketStream struct {
	conn        *websocket.Conn
	messageType int
	UpdatedAt   time.Time // 最新的更新时间
}

func (r *WebsocketStream) Read(p []byte) (n int, err error) {
	t, message, err := r.conn.ReadMessage()

	r.UpdatedAt = time.Now() // 更新时间

	r.messageType = t

	copy(p, message)

	n = len(message)

	return
}

func (r *WebsocketStream) Write(p []byte) (n int, err error) {
	err = r.conn.WriteMessage(r.messageType, p)

	r.UpdatedAt = time.Now() // 更新时间

	n = len(p)

	return
}
