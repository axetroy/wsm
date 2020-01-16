// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package util_test

import (
	"github.com/axetroy/wsm/internal/library/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerate(t *testing.T) {
	testPassword := "password"
	s := util.GeneratePassword(testPassword)

	// 生成的密码与预期的一致
	assert.Equal(t, "c52f65639a16da778bd8839424495012", s)

	// 生成两次的密码保持一致
	assert.Equal(t, util.GeneratePassword(testPassword), util.GeneratePassword(testPassword))

	assert.Len(t, util.GeneratePassword(testPassword), 32, "密码必须是32位")
}
