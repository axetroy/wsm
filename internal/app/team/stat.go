package team

import (
	"errors"
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

func (s *Service) StatTeam(c controller.Context, teamID string) (res schema.Response) {
	var (
		err  error
		data = schema.TeamStat{}
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

		helper.Response(&res, data, nil, err)
	}()

	teamMemberInfo := db.TeamMember{
		TeamID: teamID,
		UserID: c.Uid,
	}

	if err = db.Db.Model(&teamMemberInfo).Where(&teamMemberInfo).Preload("Team").First(&teamMemberInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = mapstructure.Decode(teamMemberInfo.Team, &data.Team.TeamPure); err != nil {
		return
	}

	data.Team.CreatedAt = teamMemberInfo.Team.CreatedAt.Format(time.RFC3339Nano)
	data.Team.UpdatedAt = teamMemberInfo.Team.UpdatedAt.Format(time.RFC3339Nano)

	// 统计拥有的成员数量
	teamMember := db.TeamMember{TeamID: teamID}
	if err = db.Db.Model(teamMember).Where(&teamMember).Count(&data.MemberNum).Error; err != nil {
		return
	}

	// 统计拥有的服务器数量
	teamHost := db.Host{OwnerID: teamID, OwnerType: db.HostOwnerTypeTeam}
	if err = db.Db.Model(teamHost).Where(&teamHost).Count(&data.HostNum).Error; err != nil {
		return
	}

	return
}

func (s *Service) StatTeamRouter(c *gin.Context) {
	c.JSON(http.StatusOK, s.StatTeam(controller.NewContextFromGinContext(c), c.Param("team_id")))
}