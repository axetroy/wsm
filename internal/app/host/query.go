// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package host

import (
	"errors"
	"time"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type queryList struct {
	schema.Query
}

// 获取我的服务器详情
func GetHostDetailByIdForUser(c *controller.Context) (res schema.Response) {
	var (
		err    error
		hostID = c.GetParam("host_id")
		data   = schema.Host{}
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

// 获取我可操作的服务器列表
func GetMyOperationalHostForUser(c *controller.Context) (res schema.Response) {
	var (
		err   error
		input queryList
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

	if err = c.ShouldBindQuery(&input); err != nil {
		return
	}

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

// 获取团队服务器详情
func GetHostDetailByIdForTeam(c *controller.Context) (res schema.Response) {
	var (
		err    error
		teamID = c.GetParam("team_id")
		hostID = c.GetParam("host_id")
		data   = schema.Host{}
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

	// 查询用户是否在团队中
	if err = db.Db.Where(db.TeamMember{TeamID: teamID, UserID: c.Uid}).First(&db.TeamMember{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	hostInfo := db.Host{
		Id:        hostID,
		OwnerID:   teamID,
		OwnerType: db.HostOwnerTypeTeam,
	}

	if err = db.Db.Where(&hostInfo).First(&hostInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = mapstructure.Decode(hostInfo, &data.HostPure); err != nil {
		return
	}

	data.CreatedAt = hostInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = hostInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

// 获取团队可操作的服务器列表
func GetMyOperationalHostForTeam(c *controller.Context) (res schema.Response) {
	var (
		err    error
		teamID = c.GetParam("team_id")
		input  queryList
		data   = make([]schema.Host, 0) // 输出到外部的结果
		list   = make([]db.Host, 0)     // 数据库查询出来的原始结果
		total  int64
		meta   = &schema.Meta{}
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

	if err = c.ShouldBindQuery(&input); err != nil {
		return
	}

	query := input.Query

	query.Normalize()

	filter := db.Host{
		OwnerID:   teamID,
		OwnerType: db.HostOwnerTypeTeam,
	}

	teamMemberInfo := db.TeamMember{}

	// 查询用户是否在团队中
	if err = db.Db.Where(db.TeamMember{TeamID: teamID, UserID: c.Uid}).First(&teamMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = db.Db.Where(&filter).Find(&list).Count(&total).Error; err != nil {
		return
	}

	if err = db.Db.Limit(query.Limit).Offset(query.Offset()).Order(query.Order()).Where(&filter).Preload("Team").Find(&list).Error; err != nil {
		err = exception.DataBase.New(err.Error())
		return
	}

	for _, v := range list {
		d := schema.Host{}
		if err = mapstructure.Decode(v, &d.HostPure); err != nil {
			err = exception.DataBinding.New(err.Error())
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
