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

type QueryList struct {
	schema.Query
}

func (s *Service) QueryMyTeam(c controller.Context, teamID string) (res schema.Response) {
	var (
		err  error
		data = schema.Team{}
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

	if err = mapstructure.Decode(teamMemberInfo.Team, &data.TeamPure); err != nil {
		return
	}

	data.CreatedAt = teamMemberInfo.Team.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = teamMemberInfo.Team.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func (s *Service) QueryMyTeamRouter(c *gin.Context) {
	c.JSON(http.StatusOK, s.QueryMyTeam(controller.NewContextFromGinContext(c), c.Param("team_id")))
}

func (s *Service) QueryMyTeams(c controller.Context, input QueryList) (res schema.Response) {
	var (
		err   error
		data  = make([]schema.Team, 0)   // 输出到外部的结果
		list  = make([]db.TeamMember, 0) // 数据库查询出来的原始结果
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

	filter := db.TeamMember{
		UserID: c.Uid,
	}

	if err = db.Db.Where(&filter).Find(&list).Count(&total).Error; err != nil {
		return
	}

	if err = db.Db.Limit(query.Limit).Offset(query.Offset()).Order(query.Order()).Where(&filter).Preload("Team").Find(&list).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.Team{}
		if err = mapstructure.Decode(v.Team, &d.TeamPure); err != nil {
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

func (s *Service) QueryMyTeamsRouter(c *gin.Context) {
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

	res = s.QueryMyTeams(controller.NewContextFromGinContext(c), input)
}
