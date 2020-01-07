package host

import (
	"errors"
	"net/http"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

// TODO: 校验身份
func (s *Service) QueryHostConnectionRecord(c controller.Context, recordId string) (res schema.Response) {
	var (
		err  error
		data = schema.HostConnectionRecord{}
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

		helper.Response(&res, data, nil, err)
	}()

	connectionInfo := db.HostConnectionRecord{
		Id: recordId,
	}

	if err = db.Db.Model(&connectionInfo).Where(&connectionInfo).First(&connectionInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = mapstructure.Decode(connectionInfo, &data); err != nil {
		return
	}
	return
}

func (s *Service) QueryHostConnectionRecordRouter(c *gin.Context) {
	c.JSON(http.StatusOK, s.QueryHostConnectionRecord(controller.NewContextFromGinContext(c), c.Param("record_id")))
}
