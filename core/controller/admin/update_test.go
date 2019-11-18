// Copyright 2019 Axetroy. All rights reserved. MIT license.
package admin_test

import (
	"encoding/json"
	"github.com/axetroy/mocker"
	"github.com/axetroy/terminal/core/controller"
	"github.com/axetroy/terminal/core/controller/admin"
	"github.com/axetroy/terminal/core/model"
	"github.com/axetroy/terminal/core/schema"
	"github.com/axetroy/terminal/core/service/database"
	"github.com/axetroy/terminal/core/service/token"
	"github.com/axetroy/terminal/tester"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUpdate(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	context := controller.Context{
		Uid: adminInfo.Id,
	}

	r := admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "123123",
		Password: "123123",
		Name:     "123123",
	}, false)

	assert.Equal(t, "", r.Message)
	assert.Equal(t, schema.StatusSuccess, r.Status)

	defer admin.DeleteAdminByAccount("123123")

	testAdminInfo := model.Admin{}

	assert.Nil(t, tester.Decode(r.Data, &testAdminInfo))

	status1 := model.AdminStatusInactivated

	res := admin.Update(context, testAdminInfo.Id, admin.UpdateParams{Status: &status1})

	assert.Equal(t, "", res.Message)
	assert.Equal(t, schema.StatusSuccess, res.Status)

	m := model.Admin{Id: testAdminInfo.Id}

	assert.Nil(t, database.Db.First(&m).Error)
	assert.Equal(t, status1, m.Status)
}

func TestUpdateRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()

	r1 := admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "123123",
		Password: "123123",
		Name:     "123123",
	}, false)

	assert.Equal(t, "", r1.Message)
	assert.Equal(t, schema.StatusSuccess, r1.Status)

	defer admin.DeleteAdminByAccount("123123")

	testAdminInfo := model.Admin{}

	assert.Nil(t, tester.Decode(r1.Data, &testAdminInfo))

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	newStatus := model.AdminStatusBanned

	body, _ := json.Marshal(&admin.UpdateParams{
		Status: &newStatus,
	})

	r := tester.HttpAdmin.Put("/v1/admin/a/"+testAdminInfo.Id, body, &header)

	res := schema.Response{}
	testProfile := schema.AdminProfile{}

	assert.Equal(t, http.StatusOK, r.Code)
	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, "", res.Message)
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Nil(t, tester.Decode(res.Data, &testProfile))
	assert.Equal(t, newStatus, testProfile.Status)

}
