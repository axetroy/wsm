// Copyright 2019 Axetroy. All rights reserved. Apache License 2.0.
package util_test

import (
	"testing"

	"github.com/axetroy/wsm/internal/library/util"
	"github.com/stretchr/testify/assert"
)

func TestIsNil(t *testing.T) {
	assert.Equal(t, false, util.IsNil("1"))
	assert.Equal(t, false, util.IsNil(1))
	assert.Equal(t, true, util.IsNil(nil))
}
