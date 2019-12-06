package host_test

import (
	"testing"

	"github.com/axetroy/terminal/internal/app/db"
	"github.com/axetroy/terminal/internal/app/host"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/app/user"
	"github.com/axetroy/terminal/internal/library/controller"
	"github.com/stretchr/testify/assert"
)

func TestService_TransferHost(t *testing.T) {
	userInfo, hostInfo, err := host.TestCreateHost()

	assert.Nil(t, err)

	targetUserInfo, err := user.TestCreateUser()

	assert.Nil(t, err)

	defer func() {
		_ = db.DeleteRowByTable(new(db.User).TableName(), "id", targetUserInfo.Id)
		_ = db.DeleteRowByTable(new(db.User).TableName(), "id", userInfo.Id)
		_ = db.DeleteRowByTable(new(db.Host).TableName(), "id", hostInfo.Id)
		_ = db.DeleteRowByTable(new(db.HostRecord).TableName(), "host_id", hostInfo.Id)
	}()

	r := host.Core.TransferHost(controller.NewContext(userInfo.Id, "", ""), hostInfo.Id, targetUserInfo.Id)

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)
	assert.Nil(t, r.Data)

	// 再查询是否已转让成功
	r1 := host.Core.QueryMyHostByID(controller.NewContext(targetUserInfo.Id, "", ""), hostInfo.Id)

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)

	newHostInfo := schema.Host{}
	assert.Nil(t, r1.Decode(&newHostInfo))

	assert.Equal(t, hostInfo.Host, newHostInfo.Host)
	assert.Equal(t, hostInfo.Port, newHostInfo.Port)
	assert.Equal(t, targetUserInfo.Id, newHostInfo.OwnerID)
}
