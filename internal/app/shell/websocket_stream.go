package shell

import (
	"time"

	"github.com/gorilla/websocket"
)

type DataType int

const (
	DataTypeInput  DataType = 0 // 数据的输入
	DataTypeOutput DataType = 1 // 数据的输出
)

// 数据的时间线
type Timeline struct {
	Type DataType  `json:"type"` // 输入/输出
	Time time.Time `json:"time"` // 时间戳
	Data []byte    `json:"data"` // 数据
}

func NewWebSocketSteam(connection *websocket.Conn) *WebsocketStream {
	return &WebsocketStream{
		conn:        connection,
		messageType: websocket.BinaryMessage,
		UpdatedAt:   time.Now(),
		recorder:    make([]*Timeline, 0),
	}
}

type WebsocketStream struct {
	conn        *websocket.Conn
	messageType int
	UpdatedAt   time.Time // 最新的更新时间
	recorder    []*Timeline
}

func (r *WebsocketStream) Read(p []byte) (n int, err error) {
	t, message, err := r.conn.ReadMessage()

	r.UpdatedAt = time.Now() // 更新时间
	r.messageType = t

	r.recorder = append(r.recorder, &Timeline{
		Type: DataTypeInput,
		Time: r.UpdatedAt,
		Data: message,
	})

	copy(p, message)

	n = len(message)

	return
}

func (r *WebsocketStream) Write(p []byte) (n int, err error) {
	err = r.conn.WriteMessage(r.messageType, p)

	r.UpdatedAt = time.Now() // 更新时间
	r.recorder = append(r.recorder, &Timeline{
		Type: DataTypeOutput,
		Time: r.UpdatedAt,
		Data: p,
	})

	n = len(p)

	return
}

func (r *WebsocketStream) GetRecorder() []*Timeline {
	return r.recorder
}
