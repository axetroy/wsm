package team

import (
	"errors"
	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func (s *Service) TransferTeam(c controller.Context, teamID string, userID string) (res schema.Response) {
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

	teamInfo := db.Team{}

	// 查询用户是否真的拥有该团队
	if err = tx.Where(db.Team{Id: teamID, OwnerID: c.Uid}).First(&teamInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	// 更新团队 owner 信息
	if err = tx.Model(db.Team{}).Where(db.Team{Id: teamID}).Update(db.Team{OwnerID: userID}).Error; err != nil {
		return
	}

	// 把旧的拥有者更改为普通成员
	if err = tx.Model(db.TeamMember{}).Where(db.TeamMember{UserID: c.Uid, Role: db.TeamRoleOwner}).Update(db.TeamMember{Role: db.TeamRoleMember}).Error; err != nil {
		return
	}

	// 查询新的拥有者是否在团队里面了
	if err = tx.Where(db.TeamMember{TeamID: teamID, UserID: userID}).First(&db.TeamMember{}).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果新的拥有者没有在团队里面，那么先加入组织
			if err = tx.Create(&db.TeamMember{
				TeamID: teamID,
				UserID: userID,
				Role:   db.TeamRoleOwner,
			}).Error; err != nil {
				return
			}
			return
		}

		return
	}

	// 如果 新的拥有者 之前已经是团队成员，则更新身份信息
	if err = tx.Model(db.TeamMember{}).Where(db.TeamMember{TeamID: teamID, UserID: userID}).Update(db.TeamMember{Role: db.TeamRoleOwner}).Error; err != nil {
		return
	}

	return
}

func (s *Service) TransferTeamRouter(c *gin.Context) {
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
	userID := c.Param("user_id")

	res = s.TransferTeam(controller.NewContextFromGinContext(c), teamID, userID)
}
