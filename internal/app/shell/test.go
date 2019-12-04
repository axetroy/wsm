package shell

import (
	"errors"
	"github.com/axetroy/terminal/internal/app/exception"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/library/controller"
	"github.com/axetroy/terminal/internal/library/helper"
	"github.com/axetroy/terminal/internal/library/session"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TestPublicServerParams struct {
	Username string `json:"username" valid:"required~请输入用户名"`
	Host     string `json:"host" valid:"required~请输入服务器,host~请输入正确的服务器地址"`
	Port     uint   `json:"port" valid:"required~请输入端口,port~请输入正确的端口,range(1|65535)"`
	Password string `json:"password" valid:"required~请输入密码"`
}

func (s *Service) TestPublicServerRouter(c *gin.Context) {
	var (
		input TestPublicServerParams
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

	res = s.TestPublicServer(controller.NewContextFromGinContext(c), input)
}

func (s *Service) TestPublicServer(c controller.Context, input TestPublicServerParams) (res schema.Response) {
	var (
		err  error
		data bool
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

		helper.Response(&res, data, nil, err)
	}()

	if err = c.Validator(input); err != nil {
		return
	}

	config := session.Config{
		Username: input.Username,
		Host:     input.Host,
		Port:     input.Port,
		Password: input.Password,
	}

	data = session.Test(config)

	return
}
