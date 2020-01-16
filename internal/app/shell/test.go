// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package shell

import (
	"errors"

	"github.com/axetroy/wsm/internal/app/config"
	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/crypto"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/axetroy/wsm/internal/library/session"
	"github.com/jinzhu/gorm"
)

type TestPublicServerParams struct {
	Username string `json:"username" valid:"required~请输入用户名"`
	Host     string `json:"host" valid:"required~请输入服务器,host~请输入正确的服务器地址"`
	Port     uint   `json:"port" valid:"required~请输入端口,port~请输入正确的端口,range(1|65535)"`
	Password string `json:"password" valid:"required~请输入密码"`
}

// 测试一个服务器是否可连接
func TestHostConnect(c *controller.Context) (res schema.Response) {
	var (
		err      error
		hostID   = c.GetParam("host_id")
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

	passport := crypto.DecryptHostPassport(hostInfo.Passport, config.Common.Secret)

	terminal, err = session.NewTerminal(session.Config{
		Host:     hostInfo.Host,
		Port:     hostInfo.Port,
		Username: hostInfo.Username,
		Password: passport,
		Width:    80,
		Height:   25,
	})

	if err != nil {
		return
	}

	return
}

// 测试一个公开的服务器是否可用
func TestPublicServer(c *controller.Context) (res schema.Response) {
	var (
		err   error
		input TestPublicServerParams
		data  bool
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

	if err = c.ShouldBindJSON(&input); err != nil {
		return
	}

	sshConfig := session.Config{
		Username: input.Username,
		Host:     input.Host,
		Port:     input.Port,
		Password: input.Password,
	}

	data = session.Test(sshConfig)

	return
}
