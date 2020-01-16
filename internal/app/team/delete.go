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

func DeleteTeamByID(c *controller.Context) (res schema.Response) {
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

	teamInfo := db.Team{
		Id:      teamID,
		OwnerID: c.Uid,
	}

	teamHostInfo := db.Host{
		OwnerID:   teamID,
		OwnerType: db.HostOwnerTypeTeam,
	}

	teamMemberRecordInfo := db.TeamMember{
		TeamID: teamID,
	}

	teamMemberInviteRecordInfo := db.TeamMemberInvite{
		TeamID: teamID,
	}

	// 删除团队, 仅是拥有者才有权限删除团队
	if err = tx.Where(&teamInfo).Delete(&teamInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	// 删除团队记录
	if err = tx.Where(&teamMemberRecordInfo).Delete(&teamMemberRecordInfo).Error; err != nil {
		return
	}

	// 删除团队的邀请记录
	if err = tx.Where(&teamMemberInviteRecordInfo).Delete(&teamMemberInviteRecordInfo).Error; err != nil {
		return
	}

	// 删除团队的服务器
	if err := tx.Where(&teamHostInfo).Delete(&teamHostInfo).Error; err != nil {
		return
	}
	return
}
