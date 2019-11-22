// Copyright 2019 Axetroy. All rights reserved. MIT license.
package util_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsPoint(t *testing.T) {
	s := "123"
	assert.Equal(t, false, IsPoint(""))
	assert.Equal(t, false, IsPoint(123))
	assert.Equal(t, true, IsPoint(&s))
}
