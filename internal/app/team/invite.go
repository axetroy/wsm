package team

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/axetroy/terminal/internal/app/db"
	"github.com/axetroy/terminal/internal/app/exception"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/library/controller"
	"github.com/axetroy/terminal/internal/library/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type InviteTeamParams struct {
	Members []Member `json:"members" valid:"required~请添加成员列表"` // 组成员的 ID 列表
}

type InviteResoleParams struct {
	State db.InviteState `json:"state" valid:"required~请输入要更改的状态"`
}

func (s *Service) InviteTeamRouter(c *gin.Context) {
	var (
		input InviteTeamParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = s.InviteTeam(controller.NewContextFromGinContext(c), c.Param("team_id"), input)
}

func (s *Service) InviteTeam(c controller.Context, teamID string, input InviteTeamParams) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
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

	if err = c.Validator(input); err != nil {
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

func (s *Service) ResolveInviteTeamRouter(c *gin.Context) {
	var (
		input InviteResoleParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = s.ResolveInviteTeam(controller.NewContextFromGinContext(c), c.Param("team_id"), c.Param("invite_id"), input)
}

// 受邀者 接受/拒绝 团队邀请
func (s *Service) ResolveInviteTeam(c controller.Context, teamID string, inviteID string, input InviteResoleParams) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
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

	if err = c.Validator(input); err != nil {
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
func (s *Service) CancelInviteTeam(c controller.Context, teamID string, inviteID string) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
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

func (s *Service) CancelInviteTeamRouter(c *gin.Context) {
	c.JSON(http.StatusOK, s.CancelInviteTeam(controller.NewContextFromGinContext(c), c.Param("team_id"), c.Param("invite_id")))
}

// 获取当前团队的邀请记录
func (s *Service) GetTeamInviteRecord(c controller.Context, teamID string, input QueryList) (res schema.Response) {
	var (
		err   error
		data  = make([]schema.TeamMemberInvite, 0) // 输出到外部的结果
		list  = make([]db.TeamMemberInvite, 0)     // 数据库查询出来的原始结果
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

	if err = db.Db.Limit(query.Limit).Offset(query.Offset()).Order(query.Order()).Where(&filter).Preload("User").Preload("Team").Find(&list).Error; err != nil {
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

func (s *Service) GetTeamInviteRecordRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input QueryList
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindQuery(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = s.GetTeamInviteRecord(controller.NewContextFromGinContext(c), c.Param("team_id"), input)
}

// 获取我的受邀列表
func (s *Service) GetMyInvitedRecord(c controller.Context, teamID string, input QueryList) (res schema.Response) {
	var (
		err   error
		data  = make([]schema.TeamMemberInvite, 0) // 输出到外部的结果
		list  = make([]db.TeamMemberInvite, 0)     // 数据库查询出来的原始结果
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

	query := input.Query

	query.Normalize()

	filter := db.TeamMemberInvite{
		TeamID: teamID,
		UserID: c.Uid,
	}

	if err = db.Db.Where(&filter).Find(&list).Count(&total).Error; err != nil {
		return
	}

	if err = db.Db.Limit(query.Limit).Offset(query.Offset()).Order(query.Order()).Where(&filter).Preload("User").Preload("Team").Preload("Team.Owner").Find(&list).Error; err != nil {
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

func (s *Service) GetMyInvitedRecordRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input QueryList
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindQuery(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = s.GetMyInvitedRecord(controller.NewContextFromGinContext(c), c.Param("team_id"), input)
}
