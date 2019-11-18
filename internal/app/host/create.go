package host

import (
	"errors"
	"time"

	"github.com/axetroy/terminal/internal/app/exception"
	"github.com/axetroy/terminal/internal/app/model"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/library/controller"
	"github.com/axetroy/terminal/internal/library/database"
	"github.com/axetroy/terminal/internal/library/helper"
	"github.com/axetroy/terminal/internal/library/validator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type CreateHostParams struct {
	Host     string  `json:"host"`
	Port     uint    `json:"port"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Remark   *string `json:"remark"`
}

func (s *Service) CreateHostRouter(c *gin.Context) {
	var (
		input CreateHostParams
		err   error
		res   = schema.Response{}
	)

	defer helper.Response(&res, nil, err)

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = s.CreateHost(controller.NewContextFromGinContext(c), input)
}

func (s *Service) CreateHost(c controller.Context, input CreateHostParams) (res schema.Response) {
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

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	tx = database.Db.Begin()

	hostInfo := model.Host{
		OwnerID:  c.Uid,
		Host:     input.Host,
		Port:     input.Port,
		Username: input.Username,
		Password: input.Password,
		Remark:   input.Remark,
	}

	if err = tx.Create(&hostInfo).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(hostInfo, &data.HostPure); err != nil {
		return
	}

	data.CreatedAt = hostInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = hostInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}
