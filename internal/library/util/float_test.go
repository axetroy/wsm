// Copyright 2019 Axetroy. All rights reserved. MIT license.
package util_test

import (
	"testing"

	"github.com/axetroy/wsm/internal/library/util"
	"github.com/stretchr/testify/assert"
)

func TestFloatToStr(t *testing.T) {
	assert.Equal(t, "0.10000000", util.FloatToStr(0.1))
	assert.Equal(t, "12.20000000", util.FloatToStr(12.2))
	assert.Equal(t, "5.00000000", util.FloatToStr(5))
}
