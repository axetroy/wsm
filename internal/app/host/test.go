package host

import (
	"errors"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/internal/app/user"
	"github.com/axetroy/terminal/internal/library/controller"
)

func CreateTestHost() (profile schema.ProfileWithToken, hostInfo schema.Host, err error) {
	if profile, err = user.CreateTestUser(); err != nil {
		return
	}

	remark := "remark"

	r := Core.CreateHost(controller.NewContext(profile.Id, "", ""), CreateHostParams{
		Host:     "192.168.0.1",
		Port:     22,
		Username: "root",
		Password: "password",
		Remark:   &remark,
	})

	if r.Status != schema.StatusSuccess {
		err = errors.New(r.Message)
		return
	}

	err = r.Decode(&hostInfo)

	return
}
