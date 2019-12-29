package host_test

import (
	"testing"

	"github.com/axetroy/wsm/internal/app/db"
	"github.com/axetroy/wsm/internal/app/exception"
	"github.com/axetroy/wsm/internal/app/host"
	"github.com/axetroy/wsm/internal/app/schema"
	"github.com/axetroy/wsm/internal/library/controller"
	"github.com/stretchr/testify/assert"
)

func TestService_DeleteHostByID(t *testing.T) {
	userInfo, hostInfo, err := host.TestCreateHost()

	assert.Nil(t, err)

	defer func() {
		_ = db.DeleteRowByTable(new(db.User).TableName(), "id", userInfo.Id)
		_ = db.DeleteRowByTable(new(db.Host).TableName(), "id", hostInfo.Id)
		_ = db.DeleteRowByTable(new(db.HostRecord).TableName(), "host_id", hostInfo.Id)
	}()

	ctx := controller.NewContext(userInfo.Id, "", "")

	r := host.Core.DeleteHostByID(ctx, hostInfo.Id)

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)
	assert.Nil(t, r.Data)

	r1 := host.Core.QueryMyHostByID(ctx, hostInfo.Id)
	assert.Equal(t, exception.NoData.Code(), r1.Status)
	assert.Equal(t, exception.NoData.Error(), r1.Message)
	assert.Nil(t, r1.Data)
}
