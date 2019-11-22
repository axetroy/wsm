package host_test

import (
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/library/controller"
	"testing"

	"github.com/axetroy/terminal/internal/app/db"
	"github.com/axetroy/terminal/internal/app/host"
	"github.com/stretchr/testify/assert"
)

func TestService_UpdateHost(t *testing.T) {
	userInfo, hostInfo, err := host.TestCreateHost()

	assert.Nil(t, err)

	defer func() {
		_ = db.DeleteRowByTable(new(db.User).TableName(), "id", userInfo.Id)
		_ = db.DeleteRowByTable(new(db.Host).TableName(), "id", hostInfo.Id)
	}()

	newHost := "1.1.1.1"

	r := host.Core.UpdateHost(controller.NewContext(userInfo.Id, "", ""), host.UpdateHostParams{
		Id:   hostInfo.Id,
		Host: &newHost,
	})

	newHostInfo := schema.Host{}

	assert.Nil(t, r.Decode(&newHostInfo))

	assert.Equal(t, newHost, newHostInfo.Host)
	assert.Equal(t, hostInfo.Port, newHostInfo.Port)
}
