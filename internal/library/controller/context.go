// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package controller

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/middleware"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/gin-gonic/gin"
)

// 控制器的上下文
type Context struct {
	ctx       *gin.Context
	Uid       string `json:"uid"`        // 操作人的用户 ID
	UserAgent string `json:"user_agent"` // 用户代理
	Ip        string `json:"ip"`         // IP地址
}

type controllerFunc func(c *Context) schema.Response

// 校验 body 中的 JSON 字段
func (c *Context) ShouldBindJSON(inputPointer interface{}) error {
	if err := c.ctx.ShouldBindJSON(inputPointer); err != nil {
		return exception.InvalidParams.New(err.Error())
	}

	if isValid, err := govalidator.ValidateStruct(inputPointer); err != nil {
		return exception.InvalidParams.New(err.Error())
	} else if !isValid {
		return exception.InvalidParams
	}

	return nil
}

// 校验 url 中的 query
func (c *Context) ShouldBindQuery(inputPointer interface{}) error {
	if err := c.ctx.ShouldBindQuery(inputPointer); err != nil {
		return exception.InvalidParams.New(err.Error())
	}

	if isValid, err := govalidator.ValidateStruct(inputPointer); err != nil {
		return exception.InvalidParams.New(err.Error())
	} else if !isValid {
		return exception.InvalidParams
	}

	return nil
}

func (c *Context) response(data interface{}) {
	c.ctx.JSON(http.StatusOK, data)
}

func (c *Context) SetHeader(key, value string) {
	c.ctx.Header(key, value)
}

func (c *Context) GetHeader(key string) string {
	return c.ctx.GetHeader(key)
}

func (c *Context) GetParam(key string) string {
	return c.ctx.Param(key)
}

func (c *Context) GetQuery(key string) string {
	return c.ctx.Query(key)
}

func NewContext(c *gin.Context) Context {
	return Context{
		ctx:       c,
		Uid:       c.GetString(middleware.ContextUidField),
		UserAgent: c.GetHeader("user-agent"),
		Ip:        c.ClientIP(),
	}
}

func Router(ctrl controllerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := NewContext(c)

		ctx.response(ctrl(&ctx))
	}
}
