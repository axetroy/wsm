// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/axetroy/terminal/core/controller"
	"github.com/axetroy/terminal/core/exception"
	"github.com/axetroy/terminal/core/helper"
	"github.com/axetroy/terminal/core/model"
	"github.com/axetroy/terminal/core/schema"
	"github.com/axetroy/terminal/core/service/database"
	"github.com/axetroy/terminal/core/service/redis"
	"github.com/axetroy/terminal/core/service/token"
	"github.com/axetroy/terminal/core/util"
	"github.com/axetroy/terminal/core/validator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type SignInParams struct {
	Account  string `json:"account" valid:"required~请输入登陆账号"`
	Password string `json:"password" valid:"required~请输入密码"`
}

type SignInWithOAuthParams struct {
	Code string `json:"code" valid:"required~请输入授权代码"` // oAuth 授权之后回调返回的 code
}

// 普通帐号登陆
func SignIn(c controller.Context, input SignInParams) (res schema.Response) {
	var (
		err  error
		data = &schema.ProfileWithToken{}
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

	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	userInfo := model.User{
		Password: util.GeneratePassword(input.Password),
	}

	if validator.IsPhone(input.Account) {
		// 用手机号登陆
		userInfo.Phone = &input.Account
	} else if validator.IsEmail(input.Account) {
		// 用邮箱登陆
		userInfo.Email = &input.Account
	} else {
		// 用用户名
		userInfo.Username = input.Account
	}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.InvalidAccountOrPassword
		}
		return
	}

	if err = userInfo.CheckStatusValid(); err != nil {
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	// generate token
	if t, er := token.Generate(userInfo.Id, false); er != nil {
		err = er
		return
	} else {
		data.Token = t
	}

	// 写入登陆记录
	log := model.LoginLog{
		Uid:     userInfo.Id,                       // 用户ID
		Type:    model.LoginLogTypeUserName,        // 默认用户名登陆
		Command: model.LoginLogCommandLoginSuccess, // 登陆成功
		Client:  c.UserAgent,                       // 用户的 userAgent
		LastIp:  c.Ip,                              // 用户的IP
	}

	if err = tx.Create(&log).Error; err != nil {
		return
	}

	return
}

// 使用 oAuth 认证方式登陆
func SignInWithOAuth(c controller.Context, input SignInWithOAuthParams) (res schema.Response) {
	var (
		err  error
		data = &schema.ProfileWithToken{}
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

	uid, err := redis.ClientOAuthCode.Get(input.Code).Result()

	if err != nil {
		return
	}

	var userInfo = model.User{
		Id: uid,
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	// generate token
	if t, er := token.Generate(userInfo.Id, false); er != nil {
		err = er
		return
	} else {
		data.Token = t
	}

	// 写入登陆记录
	log := model.LoginLog{
		Uid:     userInfo.Id,                       // 用户ID
		Type:    model.LoginLogTypeUserName,        // 默认用户名登陆
		Command: model.LoginLogCommandLoginSuccess, // 登陆成功
		Client:  c.UserAgent,                       // 用户的 userAgent
		LastIp:  c.Ip,                              // 用户的IP
	}

	if err = tx.Create(&log).Error; err != nil {
		return
	}

	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func SignInRouter(c *gin.Context) {
	var (
		input SignInParams
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

	res = SignIn(controller.NewContext(c), input)
}

func SignInWithOAuthRouter(c *gin.Context) {
	var (
		input SignInWithOAuthParams
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

	res = SignInWithOAuth(controller.NewContext(c), input)
}
