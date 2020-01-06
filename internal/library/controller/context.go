// Copyright 2020 Axetroy. All rights reserved. Apache License 2.0.
package controller

import (
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

// 控制器的上下文
type Context struct {
	Uid       string `json:"uid"`        // 操作人的用户 ID
	UserAgent string `json:"user_agent"` // 用户代理
	Ip        string `json:"ip"`         // IP地址
}

func (c *Context) Validator(input interface{}) error {
	if isValid, err := govalidator.ValidateStruct(input); err != nil {
		return exception.New(err.Error(), exception.InvalidParams.Code())
	} else if !isValid {
		return exception.InvalidParams
	}
	return nil
}

func NewContextFromGinContext(c *gin.Context) Context {
	return Context{
		Uid:       c.GetString(middleware.ContextUidField),
		UserAgent: c.GetHeader("user-agent"),
		Ip:        c.ClientIP(),
	}
}

func NewContext(uid string, userAgent string, ip string) Context {
	return Context{
		Uid:       uid,
		UserAgent: userAgent,
		Ip:        ip,
	}
}
