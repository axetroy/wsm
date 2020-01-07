package shell

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/axetroy/wsm/internal/library/session"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 连接 WebSocket
func (s *Service) StartTerminalRouter(c *gin.Context) {
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

	if hostInfo.ConnectType == db.HostConnectTypePassword {
		terminalConfig.Password = hostInfo.Passport
	} else {
		terminalConfig.PrivateKey = hostInfo.Passport
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

	connectedAt := time.Now()

	stream := session.NewWebSocketSteam(conn)

	addingRecord := false

	terminal.SetCloseHandler(func() error {
		if addingRecord == true {
			return nil
		}

		addingRecord = true
		// 记录用户的操作
		recorders := stream.GetRecorder()

		if len(recorders) != 0 {
			recorderStr := make([]string, 0)

			for _, r := range recorders {
				t := r.Time.Format("2006-01-02 15:04:05.000")
				str := base64.StdEncoding.EncodeToString(r.Data)
				recorderStr = append(recorderStr, fmt.Sprintf("(%s) %s", t, str))
			}

			record := db.HostConnectionRecord{
				UserID:    ctx.Uid,
				HostID:    hostID,
				Records:   recorderStr,
				CreatedAt: connectedAt,
			}

			// TODO: 如果服务器进程退出，则来不及写入数据
			if err := db.Db.Create(&record).Error; err != nil {
				fmt.Println(err)
				fmt.Println("写入操作日志失败")
			} else {
				fmt.Println("写入记录成功...")
			}
		}

		return conn.Close()
	})

	conn.SetCloseHandler(func(code int, text string) error {
		return terminal.Close()
	})

	err = terminal.Connect(stream, stream, stream)

	if err != nil {
		_ = conn.WriteMessage(websocket.BinaryMessage, []byte(err.Error()))
		_ = conn.Close()
		return
	}

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

// 测试一个服务器是否可连接
func (s *Service) TestHostConnect(c controller.Context, hostID string) (res schema.Response) {
	var (
		err      error
		data     bool
		terminal *session.Terminal
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		if terminal != nil {
			if terminalCloseErr := terminal.Close(); terminalCloseErr != nil {
				if err == nil {
					err = terminalCloseErr
				}
			}
		}

		if err == nil {
			data = true
		}

		helper.Response(&res, data, nil, err)
	}()

	hostInfo := db.Host{
		Id: hostID,
	}

	if err = db.Db.Where(&hostInfo).First(&hostInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if hostInfo.OwnerType == db.HostOwnerTypeUser {
		// 如果是用户个人持有
		hostRecordInfo := db.HostRecord{
			HostID: hostID,
			UserID: c.Uid,
		}
		if err = db.Db.Where(&hostRecordInfo).First(&hostRecordInfo).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.NoPermission
			}
			return
		}

		if hostRecordInfo.Type != db.HostRecordTypeOwner && hostRecordInfo.Type != db.HostRecordTypeCollaborator {
			err = exception.NoPermission
			return
		}

	} else if hostInfo.OwnerType == db.HostOwnerTypeTeam {
		// 如果是团队持有
		memberInfo := db.TeamMember{
			TeamID: hostInfo.OwnerID,
			UserID: c.Uid,
		}

		if err = db.Db.Where(&memberInfo).First(&memberInfo).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.NoPermission
			}
			return
		}

		if memberInfo.Role == db.TeamRoleVisitor {
			err = exception.NoPermission
			return
		}
	}

	terminal, err = session.NewTerminal(session.Config{
		Host:     hostInfo.Host,
		Port:     hostInfo.Port,
		Username: hostInfo.Username,
		Password: hostInfo.Passport,
		Width:    80,
		Height:   25,
	})

	if err != nil {
		return
	}

	return
}

func (s *Service) TestHostConnectRouter(c *gin.Context) {
	c.JSON(http.StatusOK, s.TestHostConnect(controller.NewContextFromGinContext(c), c.Param("host_id")))
}
