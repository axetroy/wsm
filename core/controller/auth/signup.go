// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth

import (
	"errors"
	"github.com/axetroy/terminal/core/exception"
	"github.com/axetroy/terminal/core/helper"
	"github.com/axetroy/terminal/core/model"
	"github.com/axetroy/terminal/core/schema"
	"github.com/axetroy/terminal/core/service/database"
	"github.com/axetroy/terminal/core/util"
	"github.com/axetroy/terminal/core/validator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type SignUpWithUsernameParams struct {
	Username string `json:"username" valid:"required~请输入用户名"` // 用户名
	Password string `json:"password" valid:"required~请输入密码"`  // 密码
}

// 创建用户帐号，包括创建的邀请码，钱包数据等，继承到一起
func CreateUserTx(tx *gorm.DB, userInfo *model.User) (err error) {
	var (
		newTx bool
	)
	if tx == nil {
		tx = database.Db.Begin()
		newTx = true
	}

	defer func() {
		if newTx {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}
	}()

	if err = tx.Create(userInfo).Error; err != nil {
		return err
	}

	return nil
}

// 使用用户名注册
func SignUpWithUsername(input SignUpWithUsernameParams) (res schema.Response) {
	var (
		err  error
		data schema.Profile
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

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	if err = validator.ValidateUsername(input.Username); err != nil {
		return
	}

	tx = database.Db.Begin()

	u := model.User{Username: input.Username}

	if err = tx.Where("username = ?", input.Username).Find(&u).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}
	}

	if u.Id != "" {
		err = exception.UserExist
		return
	}

	userInfo := model.User{
		Username: input.Username,
		Nickname: &input.Username,
		Password: util.GeneratePassword(input.Password),
		Status:   model.UserStatusInit,
		Role:     pq.StringArray{model.DefaultUser.Name},
		Phone:    nil,
		Email:    nil,
		Gender:   model.GenderUnknown,
	}

	if err = CreateUserTx(tx, &userInfo); err != nil {
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func SignUpWithUsernameRouter(c *gin.Context) {
	var (
		input SignUpWithUsernameParams
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

	res = SignUpWithUsername(input)
}
