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

type UpdateHostParams struct {
	Name     *string `json:"name"`
	Host     *string `json:"host" valid:"host~请输入正确的服务器地址"`
	Port     *uint   `json:"port" valid:"port~请输入正确的端口,range(1|65535)"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	Remark   *string `json:"remark"`
}

func (s *Service) UpdateHostRouter(c *gin.Context) {
	var (
		input UpdateHostParams
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

	hostID := c.Param("host_id")

	res = s.UpdateHost(controller.NewContextFromGinContext(c), hostID, input)
}

func (s *Service) UpdateHost(c controller.Context, hostID string, input UpdateHostParams) (res schema.Response) {
	var (
		err          error
		data         schema.Host
		tx           *gorm.DB
		shouldUpdate bool
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
			if err != nil || !shouldUpdate {
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

	hostInfo := db.Host{Id: hostID, OwnerID: c.Uid}

	updateModel := db.Host{}

	if input.Name != nil {
		updateModel.Name = *input.Name
		shouldUpdate = true
	}

	if input.Host != nil {
		updateModel.Host = *input.Host
		shouldUpdate = true
	}

	if input.Port != nil {
		updateModel.Port = *input.Port
		shouldUpdate = true
	}

	if input.Username != nil {
		updateModel.Username = *input.Username
		shouldUpdate = true
	}

	if input.Password != nil {
		updateModel.Password = *input.Password
		shouldUpdate = true
	}

	if input.Remark != nil {
		updateModel.Remark = input.Remark
		shouldUpdate = true
	}

	if shouldUpdate {
		if err = tx.Model(&hostInfo).Updates(&updateModel).First(&hostInfo).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.NoData
			}
			return
		}
	}

	if err = mapstructure.Decode(hostInfo, &data.HostPure); err != nil {
		return
	}

	data.CreatedAt = hostInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = hostInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}
