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

type Member struct {
	ID   string      `json:"id"`
	Role db.TeamRole `json:"role"`
}

type CreateTeamParams struct {
	Name    string   `json:"name" valid:"required~请输入名称"`      // team 名称
	Members []string `json:"members" valid:"required~请添加成员列表"` // 组成员的 ID 列表
	Remark  *string  `json:"remark"`                           // 备注
}

func (s *Service) CreateTeamRouter(c *gin.Context) {
	var (
		input CreateTeamParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = s.CreateTeam(controller.NewContextFromGinContext(c), input)
}

func (s *Service) CreateTeam(c controller.Context, input CreateTeamParams) (res schema.Response) {
	var (
		err  error
		data schema.Team
		tx   *gorm.DB
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

		helper.Response(&res, data, nil, err)
	}()

	if err = c.Validator(input); err != nil {
		return
	}

	tx = db.Db.Begin()

	teamInfo := db.Team{
		OwnerID: c.Uid,
		Name:    input.Name,
		Remark:  input.Remark,
	}

	if err = tx.Create(&teamInfo).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(teamInfo, &data.TeamPure); err != nil {
		return
	}

	data.CreatedAt = teamInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = teamInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}
