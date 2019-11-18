// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role_test

import (
	"encoding/json"
	"github.com/axetroy/mocker"
	"github.com/axetroy/terminal/core/controller"
	"github.com/axetroy/terminal/core/controller/role"
	"github.com/axetroy/terminal/core/rbac/accession"
	"github.com/axetroy/terminal/core/schema"
	"github.com/axetroy/terminal/core/service/token"
	"github.com/axetroy/terminal/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetList(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	context := controller.Context{
		Uid: adminInfo.Id,
	}

	{
		var (
			name        = "vip"
			description = "VIP 用户"
			accessions  = accession.Stringify(accession.ProfileUpdate)
			n           = schema.Role{}
		)

		r := role.Create(controller.Context{
			Uid: adminInfo.Id,
		}, role.CreateParams{
			Name:        name,
			Description: description,
			Accession:   accessions,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		defer role.DeleteRoleByName(name)

		assert.Nil(t, tester.Decode(r.Data, &n))

		assert.Equal(t, name, n.Name)
		assert.Equal(t, description, n.Description)
		assert.Equal(t, accessions, n.Accession)
		assert.Equal(t, false, n.BuildIn)
	}

	// 获取列表
	{
		r := role.GetList(context, role.Query{})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		list := make([]schema.Role, 0)

		assert.Nil(t, tester.Decode(r.Data, &list))

		assert.Equal(t, schema.DefaultLimit, r.Meta.Limit)
		assert.Equal(t, schema.DefaultPage, r.Meta.Page)
		assert.IsType(t, 1, r.Meta.Num)
		assert.IsType(t, int64(1), r.Meta.Total)

		if !assert.True(t, len(list) >= 1) {
			return
		}

		for _, n := range list {
			assert.IsType(t, "string", n.Name)
			assert.IsType(t, "string", n.Description)
			assert.IsType(t, []string{}, n.Accession)
			assert.IsType(t, true, n.BuildIn)
		}
	}
}

func TestGetListRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	{
		var (
			name        = "vip"
			description = "VIP 用户"
			accessions  = accession.Stringify(accession.ProfileUpdate)
			n           = schema.Role{}
		)

		r := role.Create(controller.Context{
			Uid: adminInfo.Id,
		}, role.CreateParams{
			Name:        name,
			Description: description,
			Accession:   accessions,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		defer role.DeleteRoleByName(name)

		assert.Nil(t, tester.Decode(r.Data, &n))

		assert.Equal(t, name, n.Name)
		assert.Equal(t, description, n.Description)
		assert.Equal(t, accessions, n.Accession)
		assert.Equal(t, false, n.BuildIn)
	}

	{
		r := tester.HttpAdmin.Get("/v1/role", nil, &header)

		res := schema.Response{}

		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Equal(t, "", res.Message)

		list := make([]schema.Role, 0)

		assert.Nil(t, tester.Decode(res.Data, &list))

		for _, n := range list {
			assert.IsType(t, "string", n.Name)
			assert.IsType(t, "string", n.Description)
			assert.IsType(t, []string{}, n.Accession)
			assert.IsType(t, true, n.BuildIn)
		}
	}
}
