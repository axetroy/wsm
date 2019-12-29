package host_test

import (
	"testing"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/host"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/app/user"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateHost(t *testing.T) {
	profile, err := user.TestCreateUser()

	assert.Nil(t, err)

	defer assert.Nil(t, db.DeleteRowByTable(new(db.User).TableName(), "id", profile.Id))

	remark := "master server"

	r := host.Core.CreateHostByUser(controller.NewContext(profile.Id, "", ""), host.CreateHostCommonParams{
		Name:     "test server",
		Host:     "192.168.0.1",
		Port:     22,
		Username: "root",
		Password: "password",
		Remark:   &remark,
	})

	hostInfo := schema.Host{}

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)
	assert.Nil(t, r.Decode(&hostInfo))

	defer assert.Nil(t, db.DeleteRowByTable(new(db.Host).TableName(), "id", hostInfo.Id))

	assert.Equal(t, "192.168.0.1", hostInfo.Host)
	assert.Equal(t, uint(22), hostInfo.Port)
	assert.Equal(t, "root", hostInfo.Username)
	assert.Equal(t, &remark, hostInfo.Remark)
	assert.Equal(t, "test server", hostInfo.Name)
}
