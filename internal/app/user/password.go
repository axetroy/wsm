package user

import (
	"errors"
	"net/http"

	"github.com/axetroy/terminal/internal/app/exception"
	"github.com/axetroy/terminal/internal/app/middleware"
	"github.com/axetroy/terminal/internal/app/model"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/library/controller"
	"github.com/axetroy/terminal/internal/library/database"
	"github.com/axetroy/terminal/internal/library/helper"
	"github.com/axetroy/terminal/internal/library/util"
	"github.com/axetroy/terminal/internal/library/validator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type UpdatePasswordParams struct {
	OldPassword string `json:"old_password" valid:"required~请输入旧密码"`
	NewPassword string `json:"new_password" valid:"required~请输入新密码"`
}

func (u *Service) UpdatePasswordRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdatePasswordParams
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

	res = u.UpdatePassword(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, input)
}

func (u *Service) UpdatePassword(c controller.Context, input UpdatePasswordParams) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
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

		helper.Response(&res, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	if input.OldPassword == input.NewPassword {
		err = exception.PasswordDuplicate
		return
	}

	tx = database.Db.Begin()

	userInfo := model.User{Id: c.Uid}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	// 验证密码是否正确
	if userInfo.Password != util.GeneratePassword(input.OldPassword) {
		err = exception.InvalidPassword
		return
	}

	newPassword := util.GeneratePassword(input.NewPassword)

	if err = tx.Model(&userInfo).Update(model.User{Password: newPassword}).Error; err != nil {
		return
	}

	return
}
