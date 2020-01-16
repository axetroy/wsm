// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package team

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

func StatTeam(c *controller.Context) (res schema.Response) {
	var (
		err    error
		teamID = c.GetParam("team_id")
		data   = schema.TeamStat{}
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

	teamMemberInfo := db.TeamMember{
		TeamID: teamID,
		UserID: c.Uid,
	}

	if err = db.Db.Model(&teamMemberInfo).Where(&teamMemberInfo).Preload("Team").First(&teamMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = mapstructure.Decode(teamMemberInfo.Team, &data.Team.TeamPure); err != nil {
		return
	}

	data.Team.CreatedAt = teamMemberInfo.Team.CreatedAt.Format(time.RFC3339Nano)
	data.Team.UpdatedAt = teamMemberInfo.Team.UpdatedAt.Format(time.RFC3339Nano)

	// 统计拥有的成员数量
	teamMember := db.TeamMember{TeamID: teamID}
	if err = db.Db.Model(teamMember).Where(&teamMember).Count(&data.MemberNum).Error; err != nil {
		return
	}

	// 统计拥有的服务器数量
	teamHost := db.Host{OwnerID: teamID, OwnerType: db.HostOwnerTypeTeam}
	if err = db.Db.Model(teamHost).Where(&teamHost).Count(&data.HostNum).Error; err != nil {
		return
	}

	return
}
