// Copyright 2019 Axetroy. All rights reserved. MIT license.
package util_test

import (
	"github.com/axetroy/terminal/core/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsNil(t *testing.T) {
	assert.Equal(t, false, util.IsNil("1"))
	assert.Equal(t, false, util.IsNil(1))
	assert.Equal(t, true, util.IsNil(nil))
}
