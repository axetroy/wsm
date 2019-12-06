package host

import (
	"errors"
	"net/http"
	"time"

	"github.com/axetroy/terminal/internal/app/db"
	"github.com/axetroy/terminal/internal/app/exception"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/library/controller"
	"github.com/axetroy/terminal/internal/library/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type CreateHostCommonParams struct {
	OwnerId   string           `json:"owner_id" valid:"请输入拥有者的ID，可以是用户 ID，可以是组织ID"`
	OwnerType db.HostOwnerType `json:"owner_type" valid:"required~请输入拥有者的类型"`
	Name      string           `json:"name" valid:"required~请输入名称"`
	Host      string           `json:"host" valid:"required~请输入地址,host~请输入正确的服务器地址"`
	Port      uint             `json:"port" valid:"required~请输入端口,port~请输入正确的端口,range(1|65535)"`
	Username  string           `json:"username" valid:"required~请输入用户名"`
	Password  string           `json:"password" valid:"required~请输入密码"`
	Remark    *string          `json:"remark"`
}

type CreateHostByUserParams struct {
	Name     string  `json:"name" valid:"required~请输入名称"`
	Host     string  `json:"host" valid:"required~请输入地址,host~请输入正确的服务器地址"`
	Port     uint    `json:"port" valid:"required~请输入端口,port~请输入正确的端口,range(1|65535)"`
	Username string  `json:"username" valid:"required~请输入用户名"`
	Password string  `json:"password" valid:"required~请输入密码"`
	Remark   *string `json:"remark"`
}

func (s *Service) CreateHostByUserRouter(c *gin.Context) {
	var (
		input CreateHostByUserParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = s.CreateHostByUser(controller.NewContextFromGinContext(c), input)
}

func (s *Service) CreateHostByUser(c controller.Context, input CreateHostByUserParams) (res schema.Response) {
	var (
		err  error
		data schema.Host
		tx   *gorm.DB
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

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, data, nil, err)
	}()

	if err = c.Validator(input); err != nil {
		return
	}

	tx = db.Db.Begin()

	hostInfo := db.Host{
		OwnerID:    c.Uid,
		OwnerType:  db.HostOwnerTypeUser,
		Name:       input.Name,
		Host:       input.Host,
		Port:       input.Port,
		Username:   input.Username,
		Password:   input.Password,
		PrivateKey: "", // TODO: finish this
		Remark:     input.Remark,
	}

	if err = tx.Create(&hostInfo).Error; err != nil {
		return
	}

	if err = tx.Create(&db.HostRecord{UserID: c.Uid, HostID: hostInfo.Id, Type: db.HostRecordTypeOwner}).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(hostInfo, &data.HostPure); err != nil {
		return
	}

	data.CreatedAt = hostInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = hostInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}
