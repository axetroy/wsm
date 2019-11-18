// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role_test

import (
	"encoding/json"
	"github.com/axetroy/mocker"
	"github.com/axetroy/terminal/core/controller"
	"github.com/axetroy/terminal/core/controller/auth"
	"github.com/axetroy/terminal/core/controller/role"
	"github.com/axetroy/terminal/core/exception"
	"github.com/axetroy/terminal/core/model"
	"github.com/axetroy/terminal/core/rbac/accession"
	"github.com/axetroy/terminal/core/schema"
	"github.com/axetroy/terminal/core/service/database"
	"github.com/axetroy/terminal/core/service/token"
	"github.com/axetroy/terminal/tester"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUpdate(t *testing.T) {
	var (
		roleInfo    = schema.Role{}
		name        = "vip"
		description = "VIP 用户"
		accessions  = accession.Stringify(accession.ProfileUpdate)
	)

	adminInfo, _ := tester.LoginAdmin()

	context := controller.Context{
		Uid: adminInfo.Id,
	}

	{
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

		assert.Nil(t, tester.Decode(r.Data, &roleInfo))

		n := &roleInfo

		assert.Equal(t, name, n.Name)
		assert.Equal(t, description, n.Description)
		assert.Equal(t, accessions, n.Accession)
		assert.Equal(t, false, n.BuildIn)
	}

	{

		var (
			newDescription = "new description"
		)

		r := role.Update(context, roleInfo.Name, role.UpdateParams{
			Description: &newDescription,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &roleInfo))

		assert.Equal(t, newDescription, roleInfo.Description)
	}

	// 修改一个内建的角色，应该报错
	{
		var (
			newDescription = "new description"
		)

		r := role.Update(context, model.DefaultUser.Name, role.UpdateParams{
			Description: &newDescription,
		})

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.RoleCannotUpdate.Error(), r.Message)
	}
}

func TestUpdateRouter(t *testing.T) {
	var (
		roleInfo    = schema.Role{}
		name        = "vip"
		description = "VIP 用户"
		accessions  = accession.Stringify(accession.ProfileUpdate)
	)

	adminInfo, _ := tester.LoginAdmin()

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	{
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

		assert.Nil(t, tester.Decode(r.Data, &roleInfo))

		n := &roleInfo

		assert.Equal(t, name, n.Name)
		assert.Equal(t, description, n.Description)
		assert.Equal(t, accessions, n.Accession)
		assert.Equal(t, false, n.BuildIn)
	}

	{

		var (
			newDescription = "new description"
		)

		body, _ := json.Marshal(&role.UpdateParams{
			Description: &newDescription,
		})

		r := tester.HttpAdmin.Put("/v1/role/r/"+roleInfo.Name, body, &header)

		if !assert.Equal(t, http.StatusOK, r.Code) {
			return
		}

		res := schema.Response{}

		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
		assert.Equal(t, "", res.Message)
		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Nil(t, tester.Decode(res.Data, &roleInfo))
		assert.Equal(t, newDescription, roleInfo.Description)

	}

}

func TestUpdateUserRole(t *testing.T) {
	var (
		roleInfo    = schema.Role{}
		name        = "vip"
		description = "VIP 用户"
		accessions  = accession.Stringify(accession.ProfileUpdate)
	)

	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	context := controller.Context{
		Uid: adminInfo.Id,
	}

	{
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

		assert.Nil(t, tester.Decode(r.Data, &roleInfo))

		n := &roleInfo

		assert.Equal(t, name, n.Name)
		assert.Equal(t, description, n.Description)
		assert.Equal(t, accessions, n.Accession)
		assert.Equal(t, false, n.BuildIn)
	}

	// 更改用户的角色
	{
		r := role.UpdateUserRole(context, userInfo.Id, role.UpdateUserRoleParams{
			Roles: []string{roleInfo.Name},
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		profile := schema.Profile{}

		assert.Nil(t, tester.Decode(r.Data, &profile))

		assert.Equal(t, []string{name}, profile.Role)
	}

	// 查看用户的角色是否正确
	{
		profile := model.User{
			Id: userInfo.Id,
		}

		assert.Nil(t, database.Db.First(&profile).Error)

		assert.Equal(t, userInfo.Username, profile.Username)
		assert.Equal(t, pq.StringArray{roleInfo.Name}, profile.Role)
	}
}

func TestUpdateUserRoleRouter(t *testing.T) {
	var (
		roleInfo    = schema.Role{}
		name        = "vip"
		description = "VIP 用户"
		accessions  = accession.Stringify(accession.ProfileUpdate)
	)

	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	{
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

		assert.Nil(t, tester.Decode(r.Data, &roleInfo))

		n := &roleInfo

		assert.Equal(t, name, n.Name)
		assert.Equal(t, description, n.Description)
		assert.Equal(t, accessions, n.Accession)
		assert.Equal(t, false, n.BuildIn)
	}

	{
		body, _ := json.Marshal(&role.UpdateUserRoleParams{
			Roles: []string{roleInfo.Name},
		})

		r := tester.HttpAdmin.Put("/v1/role/u/"+userInfo.Id, body, &header)

		assert.Equal(t, http.StatusOK, r.Code)
		res := schema.Response{}
		assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
		assert.Equal(t, "", res.Message)
		assert.Equal(t, schema.StatusSuccess, res.Status)
		assert.Nil(t, tester.Decode(res.Data, &roleInfo))
	}

	// 查看用户的角色是否正确
	{
		profile := model.User{
			Id: userInfo.Id,
		}

		assert.Nil(t, database.Db.First(&profile).Error)

		assert.Equal(t, userInfo.Username, profile.Username)
		assert.Equal(t, pq.StringArray{roleInfo.Name}, profile.Role)
	}
}
