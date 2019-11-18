package host_test

import (
	"testing"

	"github.com/axetroy/terminal/internal/app/host"
	"github.com/axetroy/terminal/internal/app/model"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/app/user"
	"github.com/axetroy/terminal/internal/library/controller"
	"github.com/axetroy/terminal/internal/library/database"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateHost(t *testing.T) {
	profile, err := user.CreateTestUser()

	assert.Nil(t, err)

	defer assert.Nil(t, database.DeleteRowByTable(new(model.User).TableName(), "id", profile.Id))

	remark := "master server"

	r := host.Core.CreateHost(controller.NewContext(profile.Id, "", ""), host.CreateHostParams{
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

	defer assert.Nil(t, database.DeleteRowByTable(new(model.Host).TableName(), "id", hostInfo.Id))

	assert.Equal(t, "192.168.0.1", hostInfo.Host)
	assert.Equal(t, uint(22), hostInfo.Port)
	assert.Equal(t, "root", hostInfo.Username)
	assert.Equal(t, &remark, hostInfo.Remark)
	assert.Equal(t, &remark, hostInfo.Remark)
}
