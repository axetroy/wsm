package team

import (
	"errors"
	"github.com/axetroy/terminal/internal/app/db"
	"github.com/axetroy/terminal/internal/app/exception"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/library/controller"
	"github.com/axetroy/terminal/internal/library/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func (s *Service) KickOutByUID(c controller.Context, teamID string, UserID string) (res schema.Response) {
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

	// 删除成员
	if err := tx.Where(&db.TeamMember{TeamID: teamID, UserID: UserID}).Delete(&db.TeamMember{}).Error; err != nil {
		return
	}

	return
}

func (s *Service) KickOutByUIDRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	res = s.KickOutByUID(controller.NewContextFromGinContext(c), c.Param("team_id"), c.Param("user_id"))
}
