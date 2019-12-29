// Copyright 2019 Axetroy. All rights reserved. MIT license.
package util_test

import (
	"testing"

	"github.com/axetroy/wsm/internal/library/util"
	"github.com/stretchr/testify/assert"
)

func TestIsPoint(t *testing.T) {
	s := "123"
	assert.Equal(t, false, util.IsPoint(""))
	assert.Equal(t, false, util.IsPoint(123))
	assert.Equal(t, true, util.IsPoint(&s))
}
