// Copyright 2020 Axetroy. All rights reserved. Apache License 2.0.
package host

import (
	"errors"
	"net/http"
	"time"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type queryListRecord struct {
	schema.Query
	HostID *string `json:"host_id" form:"host_id"` // 指定获取某个服务器的连接记录
}

// 获取连接记录详情
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

	if err = db.Db.Model(&connectionInfo).Where(&connectionInfo).Preload("Host").First(&connectionInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	// 如果是用户持有的记录
	if connectionInfo.Host.OwnerType == db.HostOwnerTypeUser {
		if connectionInfo.UserID != c.Uid && connectionInfo.Host.OwnerID != c.Uid {
			err = exception.NoPermission
			return
		}
	} else if connectionInfo.Host.OwnerType == db.HostOwnerTypeTeam {
		// 如果是团队持有的记录
		memberInfo := db.TeamMember{
			UserID: c.Uid,
		}

		if err = db.Db.Where(&memberInfo).First(&memberInfo).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.NoPermission
			}
			return
		}

		// 如果不是自己的连接
		if connectionInfo.UserID != c.Uid {
			if memberInfo.IsAdmin() == false {
				err = exception.NoPermission
				return
			}
		}
	}

	if err = mapstructure.Decode(connectionInfo, &data.HostConnectionRecordPure); err != nil {
		return
	}

	data.CreatedAt = connectionInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = connectionInfo.UpdatedAt.Format(time.RFC3339Nano)
	return
}

func (s *Service) QueryHostConnectionRecordRouter(c *gin.Context) {
	c.JSON(http.StatusOK, s.QueryHostConnectionRecord(controller.NewContextFromGinContext(c), c.Param("record_id")))
}

// 获取连接记录列表
func (s *Service) QueryHostConnectionRecordList(c controller.Context, input queryListRecord) (res schema.Response) {
	var (
		err   error
		data  = make([]schema.HostConnectionRecord, 0) // 输出到外部的结果
		list  = make([]db.HostConnectionRecord, 0)     // 数据库查询出来的原始结果
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

	filter := db.HostConnectionRecord{
		UserID: c.Uid,
	}

	if input.HostID != nil {
		filter.HostID = *input.HostID
	}

	if err = db.Db.Where(&filter).Find(&list).Count(&total).Error; err != nil {
		return
	}

	if err = db.Db.Limit(query.Limit).Offset(query.Offset()).Order(query.Order()).Where(&filter).Preload("User").Preload("Host", "Host.owner_type = ?", db.HostOwnerTypeUser).Find(&list).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.HostConnectionRecord{}
		if err = mapstructure.Decode(v, &d.HostConnectionRecordPure); err != nil {
			return
		}
		if err = mapstructure.Decode(v.User, &d.User); err != nil {
			return
		}
		if err = mapstructure.Decode(v.Host, &d.Host); err != nil {
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

func (s *Service) QueryHostConnectionRecordListRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input queryListRecord
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

	hostID := c.Param("host_id")

	if hostID != "" {
		input.HostID = &hostID
	}

	res = s.QueryHostConnectionRecordList(controller.NewContextFromGinContext(c), input)
}

// 获取团队的服务器连接记录列表
func (s *Service) QueryTeamHostConnectionRecordList(c controller.Context, teamID string, input queryListRecord) (res schema.Response) {
	var (
		err   error
		data  = make([]schema.HostConnectionRecord, 0) // 输出到外部的结果
		list  = make([]db.HostConnectionRecord, 0)     // 数据库查询出来的原始结果
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

	memberInfo := db.TeamMember{
		UserID: c.Uid,
	}

	if err = db.Db.Where(&memberInfo).First(&memberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
	}

	teamOwnHosts := make([]db.Host, 0)

	if err := db.Db.Where(&db.Host{
		OwnerID:   teamID,
		OwnerType: db.HostOwnerTypeTeam,
	}).Find(&teamOwnHosts).Error; err != nil {
		return
	}

	hostIDs := make([]string, 0)

	for _, host := range teamOwnHosts {
		hostIDs = append(hostIDs, host.Id)
	}

	filter := db.HostConnectionRecord{}

	if input.HostID != nil {
		filter.HostID = *input.HostID
	}

	if err = db.Db.Where("host_id IN (?)", hostIDs).Where(&filter).Find(&list).Count(&total).Error; err != nil {
		return
	}

	if err = db.Db.Limit(query.Limit).Offset(query.Offset()).Order(query.Order()).Where("host_id IN (?)", hostIDs).Where(&filter).Preload("Host").Preload("User").Find(&list).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.HostConnectionRecord{}
		if err = mapstructure.Decode(v, &d.HostConnectionRecordPure); err != nil {
			return
		}
		if err = mapstructure.Decode(v.User, &d.User); err != nil {
			return
		}
		if err = mapstructure.Decode(v.Host, &d.Host); err != nil {
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

func (s *Service) QueryTeamHostConnectionRecordListRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input queryListRecord
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

	hostID := c.Param("host_id")

	if hostID != "" {
		input.HostID = &hostID
	}

	res = s.QueryTeamHostConnectionRecordList(controller.NewContextFromGinContext(c), c.Param("team_id"), input)
}
