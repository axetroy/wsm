// Copyright 2019 Axetroy. All rights reserved. MIT license.
package rbac

import (
	"github.com/axetroy/terminal/core/exception"
	"github.com/axetroy/terminal/core/model"
	"github.com/axetroy/terminal/core/rbac/accession"
	"github.com/axetroy/terminal/core/rbac/role"
	"github.com/axetroy/terminal/core/schema"
	"github.com/axetroy/terminal/core/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Controller struct {
	Roles []*role.Role
}

func New(uid string) (c *Controller, err error) {
	c = &Controller{}

	userInfo := model.User{
		Id: uid,
	}

	if err = database.Db.First(&userInfo).Error; err != nil {
		return
	}

	if len(userInfo.Role) == 0 {
		err = exception.NoPermission
		return
	}

	for _, roleName := range userInfo.Role {
		roleInfo := model.Role{
			Name: roleName,
		}

		if err = database.Db.First(&roleInfo).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				continue
			}
			return
		}

		r := role.New(roleInfo.Name, roleInfo.Description, accession.Normalize(roleInfo.Accession))

		c.Roles = append(c.Roles, r)
	}

	return c, nil
}

// 验证是否有这些权限
func (c *Controller) Require(a []accession.Accession) bool {
	for _, v := range a {
		if c.Has(v) {
			return true
		}
	}
	return false
}

// 检验是否拥有单独的权限
func (c *Controller) Has(a accession.Accession) bool {
	for _, r := range c.Roles {
		for _, v := range r.Accession {
			if v.Name == a.Name {
				return true
			}
		}
	}
	return false
}

// 根据 RBAC 鉴权的中间件
func Require(accesions ...accession.Accession) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err error
			uid = c.GetString("uid") // 这个中间件必须安排在JWT的中间件后面, 所以这里是拿的到 UID 的
			cc  *Controller
		)

		defer func() {
			if err != nil {
				c.JSON(http.StatusOK, schema.Response{
					Message: err.Error(),
					Data:    nil,
				})
				c.Abort()
			}
		}()

		if uid == "" {
			err = exception.NoPermission
		}

		if cc, err = New(uid); err != nil {
			return
		}

		if cc.Require(accesions) == false {
			err = exception.NoPermission
		}
	}
}
