package host_test

import (
	"testing"

	"github.com/axetroy/terminal/internal/app/db"
	"github.com/axetroy/terminal/internal/app/host"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/library/controller"
	"github.com/stretchr/testify/assert"
)

func TestService_QueryHost(t *testing.T) {
	userInfo, hostInfo, err := host.TestCreateHost()

	assert.Nil(t, err)

	defer func() {
		_ = db.DeleteRowByTable(new(db.User).TableName(), "id", userInfo.Id)
		_ = db.DeleteRowByTable(new(db.Host).TableName(), "id", hostInfo.Id)
		_ = db.DeleteRowByTable(new(db.HostRecord).TableName(), "host_id", hostInfo.Id)
	}()

	r := host.Core.QueryHost(controller.NewContext(userInfo.Id, "", ""), hostInfo.Id)

	assert.Equal(t, "", r.Message)
	assert.Equal(t, schema.StatusSuccess, r.Status)

	data := schema.Host{}

	assert.Nil(t, r.Decode(&data))
	assert.Equal(t, hostInfo.Host, data.Host)
	assert.Equal(t, hostInfo.Port, data.Port)
	assert.Equal(t, hostInfo.Username, data.Username)
	assert.Equal(t, hostInfo.Name, data.Name)
	assert.Equal(t, hostInfo.OwnerID, data.OwnerID)
	assert.Equal(t, hostInfo.Remark, data.Remark)
}

func TestService_QueryOperationalServer(t *testing.T) {
	userInfo, hostInfo, err := host.TestCreateHost()

	assert.Nil(t, err)

	defer func() {
		_ = db.DeleteRowByTable(new(db.User).TableName(), "id", userInfo.Id)
		_ = db.DeleteRowByTable(new(db.Host).TableName(), "id", hostInfo.Id)
		_ = db.DeleteRowByTable(new(db.HostRecord).TableName(), "host_id", hostInfo.Id)
	}()

	r := host.Core.QueryOperationalServer(controller.NewContext(userInfo.Id, "", ""), host.QueryList{})

	assert.Equal(t, "", r.Message)
	assert.Equal(t, schema.StatusSuccess, r.Status)

	list := make([]schema.Host, 0)

	assert.Nil(t, r.Decode(&list))
	assert.Len(t, list, 1)

	for _, data := range list {
		assert.Equal(t, hostInfo.Host, data.Host)
		assert.Equal(t, hostInfo.Port, data.Port)
		assert.Equal(t, hostInfo.Username, data.Username)
		assert.Equal(t, hostInfo.Name, data.Name)
		assert.Equal(t, hostInfo.OwnerID, data.OwnerID)
		assert.Equal(t, hostInfo.Remark, data.Remark)
	}
}
