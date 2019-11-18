package tester

import (
	"github.com/axetroy/mocker"
	"github.com/axetroy/terminal/core/server/user_server"
)

var (
	HttpUser = mocker.New(user_server.UserRouter)
)
