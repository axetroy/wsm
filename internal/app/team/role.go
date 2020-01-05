package team

import (
	"errors"
	"net/http"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UpdateMemberRoleParams struct {
	Role db.TeamRole `json:"role" valid:"required~请输入新的团队成员身份"` // 新的角色信息
}

func (s *Service) UpdateMemberRole(c controller.Context, teamID string, userID string, input UpdateMemberRoleParams) (res schema.Response) {
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

	teamMemberInfo := db.TeamMember{}

	// 查询用户是否在团队中
	if err = tx.Where(db.TeamMember{TeamID: teamID, UserID: userID}).First(&teamMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	operatorMemberInfo := db.TeamMember{}

	// 查询操作者的成员信息
	if err = tx.Where(db.TeamMember{TeamID: teamID, UserID: c.Uid}).First(&operatorMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	// 验证操作者是否有权限
	if operatorMemberInfo.Role != db.TeamRoleOwner && operatorMemberInfo.Role != db.TeamRoleAdmin {
		err = exception.NoPermission
		return
	}

	// 保证权限不跨级
	switch input.Role {
	case db.TeamRoleOwner:
		err = exception.NoPermission
	case db.TeamRoleAdmin:
		if operatorMemberInfo.Role != db.TeamRoleOwner {
			err = exception.NoPermission
			return
		}
		break
	}

	// 更新组成员的角色
	if err = tx.Model(db.TeamMember{}).Where(db.TeamMember{TeamID: teamID, UserID: userID}).Update(&db.TeamMember{Role: input.Role}).Error; err != nil {
		return
	}

	return
}

func (s *Service) UpdateMemberRoleRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdateMemberRoleParams
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	teamID := c.Param("team_id")
	userID := c.Param("user_id")

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = s.UpdateMemberRole(controller.NewContextFromGinContext(c), teamID, userID, input)
}
