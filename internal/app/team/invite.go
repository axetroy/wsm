// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package team

import (
	"errors"
	"fmt"
	"time"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type inviteMember struct {
	ID   string      `json:"id"`
	Role db.TeamRole `json:"role"`
}

type queryInviteList struct {
	schema.Query
}

type inviteTeamParams struct {
	Members []inviteMember `json:"members" valid:"required~请添加成员列表"` // 组成员的 ID 列表
}

type inviteResoleParams struct {
	State db.InviteState `json:"state" valid:"required~请输入要更改的状态"`
}

// 邀请加入团队
func InviteTeam(c *controller.Context) (res schema.Response) {
	var (
		err    error
		teamID = c.GetParam("team_id")
		input  inviteTeamParams
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

	if err = c.ShouldBindJSON(&input); err != nil {
		return
	}

	tx = db.Db.Begin()

	ownerMemberInfo := db.TeamMember{
		TeamID: teamID,
		UserID: c.Uid,
	}

	if err := tx.Where(&ownerMemberInfo).Find(&ownerMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	// 只有管理员和拥有者有权限邀请人
	if ownerMemberInfo.Role != db.TeamRoleOwner && ownerMemberInfo.Role != db.TeamRoleAdmin {
		err = exception.NoPermission
		return
	}

	// 生成邀请记录
	// TODO: 生成通知，当前暂未有通知相关的接口
	for _, member := range input.Members {
		userInfo := db.User{
			Id: member.ID,
		}

		// 确保这个用户存在
		if err = tx.Where(&userInfo).Find(&userInfo).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.UserNotExist
			}
			return
		}

		memberInviteInfo := db.TeamMemberInvite{
			TeamID:    teamID,
			InvitorID: c.Uid,
			UserID:    member.ID,
			Role:      member.Role,
			State:     db.InviteStateInit,
		}

		// 如果邀请的这个用户已存在团队中，那么报错
		memberInfo := db.TeamMember{TeamID: teamID, UserID: member.ID}

		if er := tx.Where(&memberInfo).First(&memberInfo).Error; er != nil {
			if er != gorm.ErrRecordNotFound {
				err = er
				return
			}
		} else {
			err = exception.Duplicate.New(fmt.Sprintf("用户 `%s` 已经存在团队中", memberInfo.UserID))
			return
		}

		if err = tx.Where(&db.TeamMemberInvite{
			TeamID: teamID,
			UserID: member.ID,
			State:  db.InviteStateInit,
		}).Error; err != nil {
			// 如果不存在以前的邀请记录，那么一切正常，不用干什么
			if err == gorm.ErrRecordNotFound {
				err = nil
			} else {
				return
			}
		} else {
			// 如果前面已有这个团队的邀请，那么应该使其失效, 然后再创建一条新的记录
			if err = tx.Model(&db.TeamMemberInvite{}).Where(&db.TeamMemberInvite{
				TeamID: teamID,
				UserID: member.ID,
				State:  db.InviteStateInit,
			}).Update(&db.TeamMemberInvite{State: db.InviteStateDeprecated}).Error; err != nil {
				return
			}
		}

		if err = tx.Create(&memberInviteInfo).Error; err != nil {
			return
		}
	}
	return
}

// 受邀者 接受/拒绝 团队邀请
func ResolveInviteTeam(c *controller.Context) (res schema.Response) {
	var (
		err      error
		teamID   = c.GetParam("team_id")
		inviteID = c.GetParam("invite_id")
		input    inviteResoleParams
		tx       *gorm.DB
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

	if err = c.ShouldBindJSON(&input); err != nil {
		return
	}

	tx = db.Db.Begin()

	inviteInfo := db.TeamMemberInvite{}

	// 查找邀请记录
	if err = tx.Where(&db.TeamMemberInvite{
		Id:     inviteID,
		TeamID: teamID,
		UserID: c.Uid,
		State:  db.InviteStateInit,
	}).First(&inviteInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	switch input.State {
	case db.InviteStateAccept:
		fallthrough
	case db.InviteStateRefuse:
		break
	default:
		err = exception.InvalidParams
		return
	}

	// 更新邀请记录
	if err = tx.Model(&db.TeamMemberInvite{}).Where(&db.TeamMemberInvite{
		Id:     inviteID,
		TeamID: teamID,
		UserID: c.Uid,
		State:  db.InviteStateInit,
	}).Update(&db.TeamMemberInvite{State: input.State}).Error; err != nil {
		return
	}

	// 如果是接受邀请
	if input.State == db.InviteStateAccept {
		// 加入团队
		memberInfo := db.TeamMember{
			TeamID: teamID,
			UserID: c.Uid,
			Role:   inviteInfo.Role,
		}

		if err = tx.Create(&memberInfo).Error; err != nil {
			return
		}
	}

	return
}

// 团队管理者取消邀请
func CancelInviteTeam(c *controller.Context) (res schema.Response) {
	var (
		err      error
		teamID   = c.GetParam("team_id")
		inviteID = c.GetParam("invite_id")
		tx       *gorm.DB
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

	teamMemberInfo := db.TeamMember{TeamID: teamID, UserID: c.Uid}

	if err = tx.Where(&teamMemberInfo).First(&teamMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	if teamMemberInfo.Role != db.TeamRoleOwner && teamMemberInfo.Role != db.TeamRoleAdmin {
		err = exception.NoPermission
		return
	}

	inviteInfo := db.TeamMemberInvite{
		Id:    inviteID,
		State: db.InviteStateInit,
	}

	// 查找邀请记录
	if err = tx.Where(&inviteInfo).Find(&inviteInfo).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	// 更新邀请记录
	if err = tx.Model(&db.TeamMemberInvite{}).Where(&db.TeamMemberInvite{Id: inviteID}).Update(&db.TeamMemberInvite{State: db.InviteStateCancel}).Error; err != nil {
		return
	}

	return
}

// 获取团队的邀请记录
func GetTeamInviteList(c *controller.Context) (res schema.Response) {
	var (
		err    error
		teamID = c.GetParam("team_id")
		input  queryInviteList
		data   = make([]schema.TeamMemberInvite, 0) // 输出到外部的结果
		list   = make([]db.TeamMemberInvite, 0)     // 数据库查询出来的原始结果
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

	memberInfo := db.TeamMember{
		TeamID: teamID,
		UserID: c.Uid,
	}

	if err = db.Db.Where(&memberInfo).First(&memberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	query := input.Query

	query.Normalize()

	filter := db.TeamMemberInvite{
		TeamID: teamID,
	}

	if err = db.Db.Where(&filter).Find(&list).Count(&total).Error; err != nil {
		return
	}

	if err = db.Db.Limit(query.Limit).Offset(query.Offset()).Order(query.Order()).Where(&filter).Preload("User").Preload("Team").Preload("Team.Owner").Preload("Invitor").Find(&list).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	for _, v := range list {
		d := schema.TeamMemberInvite{}
		if err = mapstructure.Decode(v, &d.TeamMemberInvitePure); err != nil {
			err = exception.DataBinding.New(err.Error())
			return
		}
		if err = mapstructure.Decode(v.User, &d.User); err != nil {
			err = exception.DataBinding.New(err.Error())
			return
		}
		if err = mapstructure.Decode(v.Invitor, &d.Invitor); err != nil {
			err = exception.DataBinding.New(err.Error())
			return
		}
		if err = mapstructure.Decode(v.Team, &d.Team); err != nil {
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

// 获取我的受邀列表
func GetMyInvitedList(c *controller.Context) (res schema.Response) {
	var (
		err    error
		teamID = c.GetParam("team_id")
		input  queryInviteList
		data   = make([]schema.TeamMemberInvite, 0) // 输出到外部的结果
		list   = make([]db.TeamMemberInvite, 0)     // 数据库查询出来的原始结果
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

	filter := db.TeamMemberInvite{
		TeamID: teamID,
		UserID: c.Uid,
	}

	if err = db.Db.Where(&filter).Find(&list).Count(&total).Error; err != nil {
		return
	}

	if err = db.Db.Limit(query.Limit).Offset(query.Offset()).Order(query.Order()).Where(&filter).Preload("User").Preload("Team").Preload("Team.Owner").Preload("Invitor").Find(&list).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	for _, v := range list {
		d := schema.TeamMemberInvite{}
		if err = mapstructure.Decode(v, &d.TeamMemberInvitePure); err != nil {
			err = exception.DataBinding.New(err.Error())
			return
		}

		if err = mapstructure.Decode(v.User, &d.User); err != nil {
			err = exception.DataBinding.New(err.Error())
			return
		}

		if err = mapstructure.Decode(v.Invitor, &d.Invitor); err != nil {
			err = exception.DataBinding.New(err.Error())
			return
		}

		if err = mapstructure.Decode(v.Team, &d.Team); err != nil {
			err = exception.DataBinding.New(err.Error())
			return
		}

		if err = mapstructure.Decode(v.Team.Owner, &d.Team.Owner); err != nil {
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
