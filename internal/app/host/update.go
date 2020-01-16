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

type UpdateHostParams struct {
	Name        *string             `json:"name"`
	Host        *string             `json:"host" valid:"host~请输入正确的服务器地址"`
	Port        *uint               `json:"port" valid:"port~请输入正确的端口,range(1|65535)"`
	Username    *string             `json:"username"`
	ConnectType *db.HostConnectType `json:"connect_type"`
	Passport    *string             `json:"passport"`
	Remark      *string             `json:"remark"`
}

func UpdateHostForUser(c *controller.Context) (res schema.Response) {
	var (
		err          error
		hostID       = c.GetParam("host_id")
		input        UpdateHostParams
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

	if err = c.ShouldBindJSON(&input); err != nil {
		return
	}

	tx = db.Db.Begin()

	hostInfo := db.Host{Id: hostID, OwnerID: c.Uid, OwnerType: db.HostOwnerTypeUser}

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

	if input.ConnectType != nil {
		switch *input.ConnectType {
		case db.HostConnectTypePassword:
			fallthrough
		case db.HostConnectTypePrivateKey:
			break
		default:
			err = exception.InvalidParams
			return
		}
		updateModel.ConnectType = *input.ConnectType
		shouldUpdate = true
	}

	if input.Passport != nil {
		updateModel.Passport = *input.Passport
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

func UpdateHostForTeam(c *controller.Context) (res schema.Response) {
	var (
		err          error
		hostID       = c.GetParam("host_id")
		teamID       = c.GetParam("team_id")
		input        UpdateHostParams
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

	if err = c.ShouldBindJSON(&input); err != nil {
		return
	}

	tx = db.Db.Begin()

	teamMemberInfo := db.TeamMember{}

	// 查询用户是否在团队中
	if err = db.Db.Where(db.TeamMember{TeamID: teamID, UserID: c.Uid}).First(&teamMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	// 只有拥有者和管理员才能删除
	if teamMemberInfo.Role != db.TeamRoleOwner && teamMemberInfo.Role != db.TeamRoleAdmin {
		err = exception.NoPermission
		return
	}

	hostInfo := db.Host{Id: hostID, OwnerID: teamID, OwnerType: db.HostOwnerTypeTeam}

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

	if input.Passport != nil {
		updateModel.Passport = *input.Passport
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
