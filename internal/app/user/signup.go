package user

import (
	"errors"
	"github.com/lib/pq"
	"net/http"
	"time"

	"github.com/axetroy/terminal/internal/app/db"
	"github.com/axetroy/terminal/internal/app/exception"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/library/helper"
	"github.com/axetroy/terminal/internal/library/util"
	"github.com/axetroy/terminal/internal/library/validator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type SignUpWithUsernameParams struct {
	Username string `json:"username" valid:"required~请输入用户名"` // 用户名
	Password string `json:"password" valid:"required~请输入密码"`  // 密码
}

func (u *Service) CreateUserTx(tx *gorm.DB, userInfo *db.User) (err error) {
	var (
		newTx bool
	)
	if tx == nil {
		tx = db.Db.Begin()
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

func (u *Service) SignUpWithUsername(input SignUpWithUsernameParams) (res schema.Response) {
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

	tx = db.Db.Begin()

	u1 := db.User{Username: input.Username}

	if err = tx.Where("username = ?", input.Username).Find(&u1).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}
	}

	if u1.Id != "" {
		err = exception.UserExist
		return
	}

	userInfo := db.User{
		Username: input.Username,
		Nickname: &input.Username,
		Password: util.GeneratePassword(input.Password),
		Status:   db.UserStatusInit,
		Phone:    nil,
		Email:    nil,
		Gender:   db.GenderUnknown,
		Role:     pq.StringArray{},
	}

	if err = u.CreateUserTx(tx, &userInfo); err != nil {
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func (u *Service) SignUpWithUsernameRouter(c *gin.Context) {
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

	res = u.SignUpWithUsername(input)
}
