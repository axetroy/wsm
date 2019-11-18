package tester

import (
	"github.com/axetroy/mocker"
	"github.com/axetroy/terminal/internal/app"
)

var (
	HttpUser = mocker.New(app.UserRouter)
)
