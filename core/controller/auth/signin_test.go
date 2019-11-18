// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth_test

import (
	"encoding/json"
	"github.com/axetroy/terminal/core/controller"
	"github.com/axetroy/terminal/core/controller/auth"
	"github.com/axetroy/terminal/core/exception"
	"github.com/axetroy/terminal/core/schema"
	"github.com/axetroy/terminal/core/service/token"
	"github.com/axetroy/terminal/tester"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"testing"
)

func TestSignInWithEmptyBody(t *testing.T) {
	// empty body
	r := tester.HttpUser.Post("/v1/auth/signin", []byte(nil), nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

	assert.Equal(t, schema.StatusFail, res.Status)
	assert.Equal(t, exception.InvalidParams.Error(), res.Message)
	assert.Nil(t, res.Data)
}

func TestSignInWithInvalidPassword(t *testing.T) {
	body, _ := json.Marshal(&auth.SignInParams{
		Account:  "TestSignInWithInvalidPassword",
		Password: "abc", // 输入错误的密码
	})

	// empty body
	r := tester.HttpUser.Post("/v1/auth/signin", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

	assert.Equal(t, exception.InvalidAccountOrPassword.Code(), res.Status)
	assert.Equal(t, exception.InvalidAccountOrPassword.Error(), res.Message)
	assert.Nil(t, res.Data)
}

func TestSignInSuccess(t *testing.T) {
	rand.Seed(111)
	// 先注册一个账号
	username := "test-TestSignInSuccess"
	password := "123123"

	if r := auth.SignUpWithUsername(auth.SignUpWithUsernameParams{
		Username: username,
		Password: password,
	}); r.Status != schema.StatusSuccess {
		t.Error(r.Message)
		return
	} else {
		defer auth.DeleteUserByUserName(username)
	}

	res := auth.SignIn(controller.Context{
		UserAgent: "test-user-agent",
		Ip:        "0.0.0.0.0",
	}, auth.SignInParams{
		Account:  username,
		Password: password,
	})

	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := schema.ProfileWithToken{}

	if err := mapstructure.Decode(res.Data, &profile); err != nil {
		assert.Error(t, err, err.Error())
	}

	assert.NotEmpty(t, profile.Token)

	if c, err := token.Parse("Bearer "+profile.Token, false); err != nil {
		t.Error(err)
		return
	} else {
		assert.IsType(t, "", c.Uid, "UID必须是字符串")
	}
}

func TestSignInWithWechat(t *testing.T) {
	res := auth.SignInWithWechat(controller.Context{}, auth.SignInWithWechatParams{
		Code: "test_code",
	})

	profile := schema.ProfileWithToken{}

	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)
	assert.Nil(t, mapstructure.Decode(res.Data, &profile))
	assert.NotEmpty(t, profile)

	defer auth.DeleteUserByUid(profile.Id)
}
