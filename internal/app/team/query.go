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

type queryList struct {
	schema.Query
}

func GetTeamDetail(c *controller.Context) (res schema.Response) {
	var (
		err    error
		teamID = c.GetParam("team_id")
		data   = schema.Team{}
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

	if err = db.Db.Model(&teamMemberInfo).Where(&teamMemberInfo).Preload("Team").Preload("Team.Owner").First(&teamMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = mapstructure.Decode(teamMemberInfo.Team, &data.TeamPure); err != nil {
		return
	}

	data.CreatedAt = teamMemberInfo.Team.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = teamMemberInfo.Team.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func GetTeamList(c *controller.Context) (res schema.Response) {
	var (
		err   error
		input queryList
		data  = make([]schema.TeamWithMember, 0) // 输出到外部的结果
		list  = make([]db.TeamMember, 0)         // 数据库查询出来的原始结果
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

	filter := db.TeamMember{
		UserID: c.Uid,
	}

	if err = db.Db.Model(&filter).Where(&filter).Count(&total).Error; err != nil {
		return
	}

	if err = db.Db.Model(&filter).Limit(query.Limit).Offset(query.Offset()).Order(query.Order()).Where(&filter).Preload("Team").Preload("Team.Owner").Find(&list).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.TeamWithMember{}
		if err = mapstructure.Decode(v.Team, &d.TeamPure); err != nil {
			return
		}
		if err = mapstructure.Decode(v.Team.Owner, &d.Owner); err != nil {
			return
		}
		d.UserID = v.UserID
		d.Role = v.Role
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

func GetAllTeamList(c *controller.Context) (res schema.Response) {
	var (
		err   error
		input queryList
		data  = make([]schema.TeamWithMember, 0) // 输出到外部的结果
		list  = make([]db.TeamMember, 0)         // 数据库查询出来的原始结果
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

	filter := db.TeamMember{
		UserID: c.Uid,
	}

	if err = db.Db.Model(&filter).Where(&filter).Count(&total).Error; err != nil {
		return
	}

	if err = db.Db.Model(&filter).Offset(query.Offset()).Order(query.Order()).Where(&filter).Preload("Team").Preload("Team.Owner").Find(&list).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.TeamWithMember{}
		if err = mapstructure.Decode(v.Team, &d.TeamPure); err != nil {
			return
		}
		if err = mapstructure.Decode(v.Team.Owner, &d.Owner); err != nil {
			return
		}
		d.UserID = v.UserID
		d.Role = v.Role
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
