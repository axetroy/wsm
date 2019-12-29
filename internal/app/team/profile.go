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
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func (s *Service) GetMyProfile(c controller.Context, teamID string) (res schema.Response) {
	var (
		data schema.TeamMember
		err  error
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

		helper.Response(&res, &data, nil, err)
	}()

	teamMemberInfo := db.TeamMember{}

	if err = db.Db.Where(db.TeamMember{TeamID: teamID, UserID: c.Uid}).Preload("User").First(&teamMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	if err = mapstructure.Decode(teamMemberInfo.User, &data.ProfilePublic); err != nil {
		return
	}

	data.Role = teamMemberInfo.Role
	data.CreatedAt = teamMemberInfo.CreatedAt.Format(time.RFC3339Nano)

	return
}

func (s *Service) GetMyProfileRouter(c *gin.Context) {
	c.JSON(http.StatusOK, s.GetMyProfile(controller.NewContextFromGinContext(c), c.Param("team_id")))
}
