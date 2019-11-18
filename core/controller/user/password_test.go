// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user_test

import (
	"encoding/json"
	"github.com/axetroy/mocker"
	"github.com/axetroy/terminal/core/controller"
	"github.com/axetroy/terminal/core/controller/auth"
	"github.com/axetroy/terminal/core/controller/user"
	"github.com/axetroy/terminal/core/exception"
	"github.com/axetroy/terminal/core/model"
	"github.com/axetroy/terminal/core/schema"
	"github.com/axetroy/terminal/core/service/database"
	"github.com/axetroy/terminal/core/service/token"
	"github.com/axetroy/terminal/core/util"
	"github.com/axetroy/terminal/tester"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUpdatePassword(t *testing.T) {
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	{
		// 2. 更新测试账号的密码, 旧密码错误
		r := user.UpdatePassword(controller.Context{Uid: userInfo.Id}, user.UpdatePasswordParams{
			OldPassword: "321321",
			NewPassword: "aaa",
		})

		assert.Equal(t, exception.InvalidPassword.Code(), r.Status)
		assert.Equal(t, exception.InvalidPassword.Error(), r.Message)
	}

	{
		var newPassword = "321321"
		// 2. 更新测试账号的密码
		r := user.UpdatePassword(controller.Context{Uid: userInfo.Id}, user.UpdatePasswordParams{
			OldPassword: "123123",
			NewPassword: newPassword,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.Equal(t, nil, r.Data)

		r2 := auth.SignIn(controller.Context{
			UserAgent: "test",
			Ip:        "0.0.0.0.0",
		}, auth.SignInParams{
			Account:  userInfo.Username,
			Password: newPassword,
		})

		assert.Equal(t, schema.StatusSuccess, r2.Status)
		assert.Equal(t, "", r2.Message)
	}
}

func TestUpdatePasswordByAdmin(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	{
		var newPassword = "321321"
		// 2. 更新测试账号的密码
		r := user.UpdatePasswordByAdmin(controller.Context{Uid: adminInfo.Id}, userInfo.Id, user.UpdatePasswordByAdminParams{
			NewPassword: newPassword,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.Equal(t, nil, r.Data)

		r2 := auth.SignIn(controller.Context{
			UserAgent: "test",
			Ip:        "0.0.0.0.0",
		}, auth.SignInParams{
			Account:  userInfo.Username,
			Password: newPassword,
		})

		assert.Equal(t, schema.StatusSuccess, r2.Status)
		assert.Equal(t, "", r2.Message)
	}
}

func TestUpdatePasswordRouter(t *testing.T) {
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + userInfo.Token,
	}

	// 修改密码
	{

		body, _ := json.Marshal(&user.UpdatePasswordParams{
			OldPassword: "123123",
			NewPassword: "321321",
		})

		r := tester.HttpUser.Put("/v1/user/password", body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		assert.Equal(t, nil, res.Data)

		// 验证密码是否已修改
		user := model.User{Id: userInfo.Id}

		assert.Nil(t, database.Db.First(&user).Error)
		assert.Equal(t, util.GeneratePassword("321321"), user.Password)
	}
}

func TestUpdatePasswordByAdminRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	// 修改密码
	{

		body, _ := json.Marshal(&user.UpdatePasswordByAdminParams{
			NewPassword: "321321",
		})

		r := tester.HttpAdmin.Put("/v1/user/u/"+userInfo.Id+"/password", body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		if !assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)) {
			return
		}

		if !assert.Equal(t, "", res.Message) {
			return
		}

		if !assert.Equal(t, schema.StatusSuccess, res.Status) {
			return
		}

		assert.Equal(t, nil, res.Data)

		// 验证密码是否已修改
		user := model.User{Id: userInfo.Id}

		assert.Nil(t, database.Db.First(&user).Error)
		assert.Equal(t, util.GeneratePassword("321321"), user.Password)
	}
}
