// Copyright 2019 Axetroy. All rights reserved. MIT license.
package util_test

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestRandomString(t *testing.T) {
	assert.Len(t, RandomString(1), 1)
	assert.Len(t, RandomString(2), 2)
	assert.Len(t, RandomString(3), 3)
	assert.Len(t, RandomString(4), 4)
	assert.Len(t, RandomString(8), 8)
	assert.Len(t, RandomString(16), 16)
	assert.IsType(t, "string", RandomString(16))
}

func TestRandomNumeric(t *testing.T) {
	assert.Len(t, RandomNumeric(1), 1)
	assert.Len(t, RandomNumeric(2), 2)
	assert.Len(t, RandomNumeric(3), 3)
	assert.Len(t, RandomNumeric(4), 4)
	assert.Len(t, RandomNumeric(8), 8)
	assert.Len(t, RandomNumeric(16), 16)
	assert.IsType(t, "string", RandomNumeric(16))
	assert.True(t, regexp.MustCompile("^\\d+$").MatchString(RandomNumeric(32)))
}
