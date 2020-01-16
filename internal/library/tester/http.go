// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package tester

import (
	"github.com/axetroy/mocker"
	"github.com/axetroy/wsm/internal/app"
)

var (
	Http = mocker.New(app.UserRouter)
)
