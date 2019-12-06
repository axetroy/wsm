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

func (s *Service) QuitTeam(c controller.Context, teamID string) (res schema.Response) {
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

	teamMemberInfo := db.TeamMember{}

	// 查询用户是否在团队中
	if err = tx.Where(db.TeamMember{TeamID: teamID, UserID: c.Uid}).First(&teamMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	// 拥有者无法退出团队
	if teamMemberInfo.Role == db.TeamRoleOwner {
		err = exception.NoPermission
		return
	}

	// 删除组成员
	if err = tx.Where(db.TeamMember{TeamID: teamID, UserID: c.Uid}).Delete(&db.TeamMember{}).Error; err != nil {
		return
	}

	return
}

func (s *Service) QuitTeamRouter(c *gin.Context) {
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

	teamID := c.Param("team_id")

	res = s.QuitTeam(controller.NewContextFromGinContext(c), teamID)
}
