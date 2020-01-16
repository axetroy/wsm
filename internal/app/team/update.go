package team

import (
	"errors"
	"time"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type updateTeamParams struct {
	Name   *string `json:"name"`
	Remark *string `json:"remark"`
}

func UpdateTeam(c *controller.Context) (res schema.Response) {
	var (
		err          error
		teamID       = c.GetParam("team_id")
		input        updateTeamParams
		data         schema.Team
		tx           *gorm.DB
		shouldUpdate bool
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
			if err != nil || !shouldUpdate {
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

	ownerInfo := db.TeamMember{TeamID: teamID, UserID: c.Uid}

	if err = tx.Where(&ownerInfo).First(&ownerInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoPermission
		}
		return
	}

	// 仅有管理员和拥有者有权限修改团队信息
	if ownerInfo.Role != db.TeamRoleOwner && ownerInfo.Role != db.TeamRoleAdmin {
		err = exception.NoPermission
		return
	}

	updateModel := db.Team{
		Id: teamID,
	}

	if input.Name != nil {
		updateModel.Name = *input.Name
		shouldUpdate = true
	}

	if input.Remark != nil {
		updateModel.Remark = input.Remark
		shouldUpdate = true
	}

	if shouldUpdate {
		if err = tx.Model(&updateModel).Where(&db.Team{Id: teamID}).Updates(&updateModel).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.NoData
			}
			return
		}
	}

	if err = mapstructure.Decode(updateModel, &data.TeamPure); err != nil {
		return
	}

	data.CreatedAt = updateModel.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = updateModel.UpdatedAt.Format(time.RFC3339Nano)

	return
}
