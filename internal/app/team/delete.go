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

func (s *Service) DeleteTeamByID(c controller.Context, teamID string) (res schema.Response) {
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

	teamInfo := db.Team{
		Id:      teamID,
		OwnerID: c.Uid,
	}

	teamMemberRecordInfo := db.TeamMember{
		TeamID: teamID,
	}

	// 删除团队, 仅是拥有者才有权限删除团队
	if err := tx.Where(&teamInfo).Delete(&teamInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	// 删除团队记录
	if err := tx.Where(&teamMemberRecordInfo).Delete(&teamMemberRecordInfo).Error; err != nil {
		return
	}

	return
}

func (s *Service) DeleteTeamByIDRouter(c *gin.Context) {
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

	res = s.DeleteTeamByID(controller.NewContextFromGinContext(c), c.Param("team_id"))
}
