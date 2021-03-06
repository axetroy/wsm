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
	"github.com/mitchellh/mapstructure"
)

type queryMemberList struct {
	schema.Query
	Role *db.TeamRole `json:"role"` // 按角色来筛选
}

func GetTeamMembers(c *controller.Context) (res schema.Response) {
	var (
		err    error
		teamID = c.GetParam("team_id")
		input  queryMemberList
		data   = make([]schema.TeamMember, 0) // 输出到外部的结果
		list   = make([]db.TeamMember, 0)     // 数据库查询出来的原始结果
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

	filter := db.TeamMember{
		TeamID: teamID,
	}

	if input.Role != nil {
		filter.Role = *input.Role
	}

	if err = db.Db.Model(&filter).Where(&filter).Count(&total).Error; err != nil {
		return
	}

	if err = db.Db.Limit(query.Limit).Offset(query.Offset()).Order(query.Order()).Where(&filter).Preload("User").Find(&list).Error; err != nil {
		return
	}

	for _, v := range list {
		teamMemberInfo := schema.TeamMember{}
		if err = mapstructure.Decode(v.User, &teamMemberInfo.ProfilePublic); err != nil {
			return
		}
		teamMemberInfo.Role = v.Role
		teamMemberInfo.CreatedAt = v.CreatedAt.Format(time.RFC3339Nano)
		data = append(data, teamMemberInfo)
	}

	meta.Total = total
	meta.Num = len(data)
	meta.Page = query.Page
	meta.Limit = query.Limit
	meta.Sort = query.Sort

	return
}
