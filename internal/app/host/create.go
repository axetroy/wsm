// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package host

import (
	"errors"
	"time"

	"github.com/axetroy/wsm/internal/app/config"
	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/crypto"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type createHostParams struct {
	Name        string             `json:"name" valid:"required~请输入名称"`
	Host        string             `json:"host" valid:"required~请输入地址,host~请输入正确的服务器地址"`
	Port        uint               `json:"port" valid:"required~请输入端口,port~请输入正确的端口,range(1|65535)"`
	Username    string             `json:"username" valid:"required~请输入用户名"`
	ConnectType db.HostConnectType `json:"connect_type" valid:"required~请选择服务器连接方式"` // 服务器的连接方式
	Passport    string             `json:"passport" valid:"required~请输入连接口令"`        // 对应连接方式的口令，密码/私钥
	Remark      *string            `json:"remark"`
}

func CreateHostByUser(c *controller.Context) (res schema.Response) {
	var (
		err     error
		input   createHostParams
		ownerID = c.Uid
		data    schema.Host
		tx      *gorm.DB
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

	if err = c.ShouldBindJSON(&input); err != nil {
		return
	}

	switch input.ConnectType {
	case db.HostConnectTypePassword:
		fallthrough
	case db.HostConnectTypePrivateKey:
		break
	default:
		err = exception.InvalidParams
		return
	}

	tx = db.Db.Begin()

	hostPassport := crypto.EncryptHostPassport(input.Passport, config.Common.Secret)

	hostInfo := db.Host{
		OwnerID:     ownerID,
		OwnerType:   db.HostOwnerTypeTeam,
		Name:        input.Name,
		Host:        input.Host,
		Port:        input.Port,
		Username:    input.Username,
		ConnectType: input.ConnectType,
		Passport:    hostPassport,
		Remark:      input.Remark,
	}

	if err = tx.Create(&hostInfo).Error; err != nil {
		return
	}

	if err = tx.Create(&db.HostRecord{UserID: ownerID, HostID: hostInfo.Id, Type: db.HostRecordTypeOwner}).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(hostInfo, &data.HostPure); err != nil {
		return
	}

	data.CreatedAt = hostInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = hostInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func CreateHostByTeam(c *controller.Context) (res schema.Response) {
	var (
		err     error
		input   createHostParams
		ownerID = c.GetParam("team_id")
		data    schema.Host
		tx      *gorm.DB
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

	if err = c.ShouldBindJSON(&input); err != nil {
		return
	}

	switch input.ConnectType {
	case db.HostConnectTypePassword:
		fallthrough
	case db.HostConnectTypePrivateKey:
		break
	default:
		err = exception.InvalidParams
		return
	}

	tx = db.Db.Begin()

	hostPassport := crypto.EncryptHostPassport(input.Passport, config.Common.Secret)

	hostInfo := db.Host{
		OwnerID:     ownerID,
		OwnerType:   db.HostOwnerTypeTeam,
		Name:        input.Name,
		Host:        input.Host,
		Port:        input.Port,
		Username:    input.Username,
		ConnectType: input.ConnectType,
		Passport:    hostPassport,
		Remark:      input.Remark,
	}

	if err = tx.Create(&hostInfo).Error; err != nil {
		return
	}

	// 如果是组织，则看看是否有权限
	teamMemberInfo := db.TeamMember{
		TeamID: ownerID,
		UserID: c.Uid,
	}
	if err = tx.Where(&teamMemberInfo).First(&teamMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	// 校验是否拥有权限
	if teamMemberInfo.Role != db.TeamRoleAdmin && teamMemberInfo.Role != db.TeamRoleOwner {
		err = exception.NoPermission
		return
	}

	if err = mapstructure.Decode(hostInfo, &data.HostPure); err != nil {
		return
	}

	data.CreatedAt = hostInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = hostInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}
