// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package shell

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/axetroy/wsm/internal/app/config"
	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/crypto"
	"github.com/axetroy/wsm/internal/library/session"
	"github.com/axetroy/wsm/internal/library/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	SteamMap = NewStreamMap() // websocket 连接流的 map
)

type streamMap struct {
	sync.RWMutex
	innerMap map[string]*session.WebsocketStream
}

func NewStreamMap() *streamMap {
	sm := new(streamMap)
	sm.innerMap = make(map[string]*session.WebsocketStream)
	return sm
}

func (sm *streamMap) Set(key string, value *session.WebsocketStream) {
	sm.Lock()
	sm.innerMap[key] = value
	sm.Unlock()
}

func (sm *streamMap) Get(key string) *session.WebsocketStream {
	sm.RLock()
	value := sm.innerMap[key]
	sm.RUnlock()
	return value
}

func (sm *streamMap) Remove(key string) {
	sm.RLock()
	delete(sm.innerMap, key)
	sm.RUnlock()
	return
}

func (sm *streamMap) ForEach(cb func(key string, value *session.WebsocketStream)) {
	sm.RLock()

	defer sm.RUnlock()

	for k, v := range sm.innerMap {
		cb(k, v)
	}
}

// 连接 WebSocket
func Connect(c *gin.Context) {
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

	ctx := controller.NewContext(c)

	hostInfo := db.Host{Id: hostID}

	if err := db.Db.Where(&hostInfo).First(&hostInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.String(http.StatusNotFound, exception.NoData.Error())
			return
		}
		c.String(http.StatusInternalServerError, exception.NoPermission.Error())
		return
	}

	if hostInfo.OwnerType == db.HostOwnerTypeUser {
		// 如果是用户持有的的服务器
		recordInfo := db.HostRecord{HostID: hostID, UserID: ctx.Uid}
		if err := db.Db.Where(&recordInfo).First(&recordInfo).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.String(http.StatusNotFound, exception.NoPermission.Error())
				return
			}
			c.String(http.StatusInternalServerError, exception.Unknown.Error())
			return
		}
	} else if hostInfo.OwnerType == db.HostOwnerTypeTeam {
		// 查询操作者是否属于该组织
		memberInfo := db.TeamMember{TeamID: hostInfo.OwnerID, UserID: ctx.Uid}
		if err := db.Db.Where(&memberInfo).First(&memberInfo).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.String(http.StatusNotFound, exception.NoPermission.Error())
				return
			}
			c.String(http.StatusInternalServerError, exception.Unknown.Error())
			return
		}

		// 校验权限是否足够
		if memberInfo.Role == db.TeamRoleVisitor {
			c.String(http.StatusBadRequest, exception.NoPermission.Error())
			return
		}
	} else {
		c.String(http.StatusBadRequest, exception.NoData.Error())
		return
	}

	terminalConfig := session.Config{
		Host:     hostInfo.Host,
		Port:     hostInfo.Port,
		Username: hostInfo.Username,
		Width:    cols,
		Height:   rows,
	}

	passport := crypto.DecryptHostPassport(hostInfo.Passport, config.Common.Secret)

	if hostInfo.ConnectType == db.HostConnectTypePassword {
		terminalConfig.Password = passport
	} else {
		terminalConfig.PrivateKey = passport
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	terminal, err := session.NewTerminal(terminalConfig)

	if err != nil {
		_ = conn.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
		_ = conn.Close()
		return
	}

	// 随机生成一个 ID
	connectionID := util.RandomString(64)

	stream := session.NewWebSocketSteam(conn, session.Meta{
		Uid:    ctx.Uid,
		Ip:     ctx.Ip,
		HostID: hostID,
	})

	addingRecord := false

	terminal.SetCloseHandler(func() error {
		fmt.Println("SSH 关闭，开始记录")
		if addingRecord == true {
			return nil
		}

		addingRecord = true

		// 记录用户的操作
		if err := stream.Write2Log(); err != nil {
			return err
		}

		return conn.Close()
	})

	conn.SetCloseHandler(func(code int, text string) error {
		SteamMap.Remove(connectionID)
		return terminal.Close()
	})

	err = terminal.Connect(stream, stream, stream)

	if err != nil {
		_ = conn.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
		_ = conn.Close()
		return
	}

	SteamMap.Set(connectionID, stream)

	go func() {
		for {
			timer := time.NewTimer(5 * time.Second)
			<-timer.C

			if terminal.IsClosed() {
				_ = timer.Stop()
				break
			}

			// 如果有 10 分钟没有数据流动，则断开连接
			if time.Now().Unix()-stream.UpdatedAt.Unix() > 60*10 {
				_ = conn.WriteMessage(websocket.BinaryMessage, []byte("检测到终端闲置，已断开连接..."))
				_ = conn.Close()
				_ = timer.Stop()
				break
			}
		}
	}()
}
