// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package host

import (
	"errors"
	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/jinzhu/gorm"
)

// 用户个人删除服务器
func DeleteHostByIdForUser(c *controller.Context) (res schema.Response) {
	var (
		err    error
		tx     *gorm.DB
		hostID = c.GetParam("host_id")
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

		helper.Response(&res, nil, nil, err)
	}()

	tx = db.Db.Begin()

	hostInfo := db.Host{
		Id:        hostID,
		OwnerID:   c.Uid,
		OwnerType: db.HostOwnerTypeUser,
	}

	hostRecordInfo := db.HostRecord{
		HostID: hostID,
	}

	hostConnectionInfo := db.HostConnectionRecord{HostID: hostID}

	// 删除服务器
	if err := tx.Where(&hostInfo).Delete(db.Host{}).Error; err != nil {
		return
	}

	// 删除服务器信息
	if err := tx.Where(&hostRecordInfo).Delete(db.HostRecord{}).Error; err != nil {
		return
	}

	// 删除服务器的连接记录
	if err := tx.Where(&hostConnectionInfo).Delete(db.HostConnectionRecord{}).Error; err != nil {
		return
	}

	return
}

func DeleteHostByIdForTeam(c *controller.Context) (res schema.Response) {
	var (
		err    error
		hostID = c.GetParam("host_id")
		teamID = c.GetParam("team_id")
		tx     *gorm.DB
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

		helper.Response(&res, nil, nil, err)
	}()

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

	hostInfo := db.Host{
		Id:        hostID,
		OwnerID:   teamID,
		OwnerType: db.HostOwnerTypeTeam,
	}

	// 删除服务器
	if err := tx.Where(&hostInfo).Delete(db.Host{}).Error; err != nil {
		return
	}

	// 删除服务器的连接记录
	if err := tx.Where(&db.HostConnectionRecord{HostID: hostID}).Delete(db.HostConnectionRecord{}).Error; err != nil {
		return
	}

	return
}
