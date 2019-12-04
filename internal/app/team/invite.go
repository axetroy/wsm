package team

import (
	"errors"
	"net/http"

	"github.com/axetroy/terminal/internal/app/db"
	"github.com/axetroy/terminal/internal/app/exception"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/library/controller"
	"github.com/axetroy/terminal/internal/library/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type InviteTeamParams struct {
	Members []Member `json:"members" valid:"required~请添加成员列表"` // 组成员的 ID 列表
}

type InviteResoleParams struct {
	Confirm bool `json:"confirm" valid:"required~请确定是否加入团队"` // 是否确定加入团队
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
		UserID: teamID,
	}

	if err := tx.Where(&ownerMemberInfo).Find(&ownerMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
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

		if err = tx.Where(&userInfo).Find(&userInfo).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.UserNotExist
			}
			return
		}

		memberInviteInfo := db.TeamMemberInvite{
			TeamID:    teamID,
			UserID:    member.ID,
			Role:      member.Role,
			Available: true,
		}

		if err = tx.Where(&db.TeamMemberInvite{
			TeamID:    teamID,
			UserID:    member.ID,
			Available: true,
		}).Error; err != nil {
			// 如果不存在以前的邀请记录，那么一切正常，不用干什么
			if err == gorm.ErrRecordNotFound {
				err = nil
			} else {
				return
			}
		} else {
			// 如果前面已有这个团队的邀请，那么应该使其失效, 然后再创建一条新的记录
			if err = tx.Where(&db.TeamMemberInvite{
				TeamID:    teamID,
				UserID:    member.ID,
				Available: true,
			}).Update(&db.TeamMemberInvite{Available: false}).Error; err != nil {
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

	inviteInfo := db.TeamMemberInvite{
		Id:        inviteID,
		UserID:    teamID,
		Available: true,
	}

	// 查找邀请记录
	if err = tx.Where(&inviteInfo).Find(&inviteInfo).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	// 更新邀请记录, 让其不可用
	if err = tx.Where(db.TeamMemberInvite{Id: inviteID}).Update(db.TeamMemberInvite{Available: input.Confirm}).Error; err != nil {
		return
	}

	if input.Confirm {
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
