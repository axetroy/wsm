// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package team

import (
	"errors"
	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/jinzhu/gorm"
)

func QuitTeam(c *controller.Context) (res schema.Response) {
	var (
		err    error
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
	if err = tx.Where(db.TeamMember{TeamID: teamID, UserID: c.Uid}).First(&teamMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	// 拥有者无法退出团队
	if teamMemberInfo.Role == db.TeamRoleOwner {
		err = exception.NoPermission
		return
	}

	// 删除组成员
	if err = tx.Where(db.TeamMember{TeamID: teamID, UserID: c.Uid}).Delete(&db.TeamMember{}).Error; err != nil {
		return
	}

	return
}
