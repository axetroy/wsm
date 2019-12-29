package tester

import (
	"github.com/axetroy/mocker"
	"github.com/axetroy/wsm/internal/app"
)

var (
	Http = mocker.New(app.UserRouter)
)
