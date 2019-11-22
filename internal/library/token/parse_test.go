// Copyright 2019 Axetroy. All rights reserved. MIT license.
package token_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse(t *testing.T) {
	uid := "123123"
	tokenStr, err := Generate(uid, false)

	assert.Nil(t, err)
	assert.IsType(t, "123", tokenStr)

	c, err1 := Parse(Prefix+" "+tokenStr, false)

	assert.Nil(t, err1)

	assert.Equal(t, uid, c.Id)
}
