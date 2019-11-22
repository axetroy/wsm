package tester

import (
	"github.com/axetroy/mocker"
	"github.com/axetroy/terminal/internal/app"
)

var (
	Http = mocker.New(app.UserRouter)
)
