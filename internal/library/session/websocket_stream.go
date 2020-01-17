// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package session

import (
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/gorilla/websocket"
)

// 数据的时间线
type timeline struct {
	Time time.Time `json:"time"` // 时间戳
	Data []byte    `json:"data"` // 数据
}

type Meta struct {
	Uid    string
	Ip     string
	HostID string
}

type WebsocketStream struct {
	sync.RWMutex
	conn        *websocket.Conn // socket 连接
	messageType int             // 发送的数据类型
	recorder    []*timeline     // 记录的时间轴
	CreatedAt   time.Time       // 创建时间
	UpdatedAt   time.Time       // 最新的更新时间
	Meta        Meta            // 元信息
	written     bool            // 是否已写入记录, 一个流只允许写入一次
}

func NewWebSocketSteam(connection *websocket.Conn, meta Meta) *WebsocketStream {
	now := time.Now()

	return &WebsocketStream{
		conn:        connection,
		messageType: websocket.BinaryMessage,
		CreatedAt:   now,
		UpdatedAt:   now,
		recorder:    make([]*timeline, 0),
		Meta:        meta,
	}
}

func (r *WebsocketStream) Read(p []byte) (n int, err error) {
	t, message, err := r.conn.ReadMessage()

	copy(p, message)

	r.Lock()

	defer r.Unlock()

	r.UpdatedAt = time.Now() // 更新时间
	r.messageType = t

	n = len(message)

	return
}

func (r *WebsocketStream) Write(p []byte) (n int, err error) {
	r.Lock()

	var data = make([]byte, len(p))

	copy(data, p)

	defer r.Unlock()
	err = r.conn.WriteMessage(r.messageType, p)

	r.UpdatedAt = time.Now() // 更新时间
	r.recorder = append(r.recorder, &timeline{
		Time: r.UpdatedAt,
		Data: data,
	})

	n = len(p)

	return
}

func (r *WebsocketStream) Write2Log() error {
	// 记录用户的操作
	r.Lock()

	defer r.Unlock()

	if r.written {
		return nil
	}

	recorders := r.recorder

	if len(recorders) != 0 {
		recorderStr := make([]string, 0)

		for _, r := range recorders {
			t := r.Time.Format("2006-01-02 15:04:05.000")
			str := base64.StdEncoding.EncodeToString(r.Data)
			recorderStr = append(recorderStr, fmt.Sprintf("(%s) %s", t, str))
		}

		record := db.HostConnectionRecord{
			UserID:    r.Meta.Uid,
			Ip:        r.Meta.Ip,
			HostID:    r.Meta.HostID,
			Records:   recorderStr,
			CreatedAt: r.CreatedAt,
		}

		if err := db.Db.Create(&record).Error; err != nil {
			return err
		}

		r.written = true
	}

	return nil
}
