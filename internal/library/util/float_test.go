// Copyright 2019 Axetroy. All rights reserved. MIT license.
package util_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloatToStr(t *testing.T) {
	assert.Equal(t, "0.10000000", FloatToStr(0.1))
	assert.Equal(t, "12.20000000", FloatToStr(12.2))
	assert.Equal(t, "5.00000000", FloatToStr(5))
}
