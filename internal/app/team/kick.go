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

func KickUserOutOfTeam(c *controller.Context) (res schema.Response) {
	var (
		err    error
		teamID = c.GetParam("team_id")
		UserID = c.GetParam("user_id")
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

	ownerInfo := db.TeamMember{
		TeamID: teamID,
		UserID: c.Uid,
	}

	if err := tx.Where(&ownerInfo).Find(&ownerInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	// 仅有管理员和拥有者有权限踢其他团队成员
	if ownerInfo.Role != db.TeamRoleOwner && ownerInfo.Role != db.TeamRoleAdmin {
		err = exception.NoPermission
		return
	}

	targetMemberInfo := db.TeamMember{TeamID: teamID, UserID: UserID}
	if err = tx.Where(&targetMemberInfo).First(&targetMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	// 不能将自己踢出团队
	if targetMemberInfo.UserID == ownerInfo.UserID {
		err = exception.InvalidParams
		return
	}

	// 校验是否有对应的权限
	switch targetMemberInfo.Role {
	case db.TeamRoleOwner:
		err = exception.NoPermission
		return
	case db.TeamRoleAdmin:
		if ownerInfo.Role != db.TeamRoleOwner {
			err = exception.NoPermission
			return
		}
		break
	case db.TeamRoleMember:
	case db.TeamRoleVisitor:
		if ownerInfo.Role != db.TeamRoleOwner && ownerInfo.Role != db.TeamRoleAdmin {
			err = exception.NoPermission
			return
		}
		break
	}

	// 删除成员
	if err := tx.Where(&db.TeamMember{TeamID: teamID, UserID: UserID}).Not("role", db.TeamRoleOwner).Delete(&db.TeamMember{}).Error; err != nil {
		return
	}

	return
}
