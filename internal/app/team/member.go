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
	"github.com/mitchellh/mapstructure"
)

type QueryMemberList struct {
	schema.Query
	TeamID *string      `json:"team_id"` // 根据团队ID获取成员列表
	Role   *db.TeamRole `json:"role"`    // 按角色来筛选
}

func (s *Service) QueryTeamMembers(c controller.Context, input QueryMemberList) (res schema.Response) {
	var (
		err   error
		data  = make([]schema.ProfilePublic, 0) // 输出到外部的结果
		list  = make([]db.TeamMember, 0)        // 数据库查询出来的原始结果
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

	filter := db.TeamMember{}

	if input.TeamID != nil {
		filter.TeamID = *input.TeamID
	}

	if input.Role != nil {
		filter.Role = *input.Role
	}

	if err = db.Db.Where(&filter).Find(&list).Count(&total).Error; err != nil {
		return
	}

	if err = db.Db.Limit(query.Limit).Offset(query.Offset()).Order(query.Order()).Where(&filter).Preload("User").Find(&list).Error; err != nil {
		return
	}

	for _, v := range list {
		memberInfoPublic := schema.ProfilePublic{}
		if err = mapstructure.Decode(v.Team, &memberInfoPublic); err != nil {
			return
		}
		data = append(data, memberInfoPublic)
	}

	meta.Total = total
	meta.Num = len(data)
	meta.Page = query.Page
	meta.Limit = query.Limit
	meta.Sort = query.Sort

	return
}

func (s *Service) QueryTeamMembersRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input QueryMemberList
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

	teamID := c.Param("team_id")

	input.TeamID = &teamID

	res = s.QueryTeamMembers(controller.NewContextFromGinContext(c), input)
}
