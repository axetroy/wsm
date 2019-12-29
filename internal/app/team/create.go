package team

import (
	"errors"
	"net/http"
	"time"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type Member struct {
	ID   string      `json:"id"`
	Role db.TeamRole `json:"role"`
}

type CreateTeamParams struct {
	Name    string   `json:"name" valid:"required~请输入名称"` // 团队名称
	Members []string `json:"members"`                     // 成员的 ID 列表
	Remark  *string  `json:"remark"`                      // 备注
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

	// 创建团队
	if err = tx.Create(&teamInfo).Error; err != nil {
		return
	}

	memberInfo := db.TeamMember{
		TeamID: teamInfo.Id,
		UserID: c.Uid,
		Role:   db.TeamRoleOwner,
	}

	// 创建 Member
	if err = tx.Create(&memberInfo).Error; err != nil {
		return
	}

	if input.Members != nil && len(input.Members) > 0 {
		memberMap := map[string]string{}
		for _, memberID := range input.Members {
			// 防止重复
			if _, ok := memberMap[memberID]; ok {
				continue
			}

			memberMap[memberID] = "1"
			userInfo := db.User{}

			if err = tx.Where(&db.User{Id: memberID}).First(&userInfo).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					err = exception.UserNotExist
				}
				return
			}

			// 创建 Member
			if err = tx.Create(&db.TeamMember{
				TeamID: teamInfo.Id,
				UserID: memberID,
				Role:   db.TeamRoleMember,
			}).Error; err != nil {
				return
			}
		}
	}

	if err = mapstructure.Decode(teamInfo, &data.TeamPure); err != nil {
		return
	}

	data.CreatedAt = teamInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = teamInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}
