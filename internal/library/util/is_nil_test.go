// Copyright 2019 Axetroy. All rights reserved. MIT license.
package util_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsNil(t *testing.T) {
	assert.Equal(t, false, IsNil("1"))
	assert.Equal(t, false, IsNil(1))
	assert.Equal(t, true, IsNil(nil))
}
