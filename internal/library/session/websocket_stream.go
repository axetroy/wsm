// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package session

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// 数据的时间线
type timeline struct {
	Time time.Time `json:"time"` // 时间戳
	Data []byte    `json:"data"` // 数据
}

func NewWebSocketSteam(connection *websocket.Conn) *WebsocketStream {
	return &WebsocketStream{
		conn:        connection,
		messageType: websocket.BinaryMessage,
		UpdatedAt:   time.Now(),
		recorder:    make([]*timeline, 0),
	}
}

type WebsocketStream struct {
	lock        sync.RWMutex
	conn        *websocket.Conn
	messageType int
	UpdatedAt   time.Time // 最新的更新时间
	recorder    []*timeline
}

func (r *WebsocketStream) Read(p []byte) (n int, err error) {
	t, message, err := r.conn.ReadMessage()

	copy(p, message)

	r.lock.Lock()

	defer r.lock.Unlock()

	r.UpdatedAt = time.Now() // 更新时间
	r.messageType = t

	n = len(message)

	return
}

func (r *WebsocketStream) Write(p []byte) (n int, err error) {
	r.lock.Lock()

	var data = make([]byte, len(p))

	copy(data, p)

	defer r.lock.Unlock()
	err = r.conn.WriteMessage(r.messageType, p)

	r.UpdatedAt = time.Now() // 更新时间
	r.recorder = append(r.recorder, &timeline{
		Time: r.UpdatedAt,
		Data: data,
	})

	n = len(p)

	return
}

func (r *WebsocketStream) GetRecorder() []*timeline {
	return r.recorder
}
