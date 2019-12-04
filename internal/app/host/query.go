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

type QueryList struct {
	schema.Query
}

func (s *Service) QueryHostByID(c controller.Context, hostID string) (res schema.Response) {
	var (
		err  error
		data = schema.Host{}
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

	hostRecordInfo := db.HostRecord{
		HostID: hostID,
		UserID: c.Uid,
	}

	if err = db.Db.Where(&hostRecordInfo).Preload("Host", "host.owner_type = ?", string(db.HostOwnerTypeUser)).First(&hostRecordInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = mapstructure.Decode(hostRecordInfo.Host, &data.HostPure); err != nil {
		return
	}

	data.CreatedAt = hostRecordInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = hostRecordInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func (s *Service) QueryHostByIDRouter(c *gin.Context) {
	c.JSON(http.StatusOK, s.QueryHostByID(controller.NewContextFromGinContext(c), c.Param("host_id")))
}

func (s *Service) QueryOperationalServer(c controller.Context, input QueryList) (res schema.Response) {
	var (
		err   error
		data  = make([]schema.Host, 0)   // 输出到外部的结果
		list  = make([]db.HostRecord, 0) // 数据库查询出来的原始结果
		total int64
		meta  = &schema.Meta{}
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

		helper.Response(&res, data, meta, err)
	}()

	query := input.Query

	query.Normalize()

	filter := db.HostRecord{
		UserID: c.Uid,
	}

	if err = db.Db.Where(&filter).Find(&list).Count(&total).Error; err != nil {
		return
	}

	if err = db.Db.Limit(query.Limit).Offset(query.Offset()).Order(query.Order()).Where(&filter).Preload("Host").Find(&list).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.Host{}
		if err = mapstructure.Decode(v.Host, &d.HostPure); err != nil {
			return
		}
		d.CreatedAt = v.CreatedAt.Format(time.RFC3339Nano)
		d.UpdatedAt = v.UpdatedAt.Format(time.RFC3339Nano)
		data = append(data, d)
	}

	meta.Total = total
	meta.Num = len(data)
	meta.Page = query.Page
	meta.Limit = query.Limit
	meta.Sort = query.Sort

	return
}

func (s *Service) QueryOperationalServerRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input QueryList
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindQuery(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = s.QueryOperationalServer(controller.NewContextFromGinContext(c), input)
}
