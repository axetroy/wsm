package user

import (
	"errors"
	"net/http"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/axetroy/wsm/internal/library/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type SearchUserParams struct {
	Account string `json:"account" form:"account" valid:"required~请输入搜索字段"` // 按用户名来搜索
}

func (u *Service) SearchUserRouter(c *gin.Context) {
	var (
		input SearchUserParams
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

	if err = c.ShouldBindQuery(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = u.SearchUser(controller.NewContextFromGinContext(c), input)
}

func (u *Service) SearchUser(c controller.Context, input SearchUserParams) (res schema.Response) {
	var (
		err  error
		data = make([]schema.ProfilePublic, 0) // 输出到外部的结果
		list = make([]db.User, 0)              // 数据库查询出来的原始结果
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

	userInfo := db.User{}

	if err = db.Db.Model(&userInfo).Limit(100).Where("username LIKE ?", "%"+input.Account+"%").Find(&list).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.ProfilePublic{}
		if err = mapstructure.Decode(v, &d); err != nil {
			err = exception.DataBinding.New(err.Error())
			return
		}
		data = append(data, d)
	}

	return
}
