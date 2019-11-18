package tester_test

import (
	"github.com/axetroy/terminal/internal/app/model"
	"github.com/axetroy/terminal/internal/app/schema"
	"github.com/axetroy/terminal/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecode(t *testing.T) {
	source := model.News{
		Title:   "title",
		Content: "content",
	}
	dest := schema.News{}
	assert.Nil(t, tester.Decode(source, &dest))

	assert.Equal(t, "title", dest.Title)
	assert.Equal(t, "content", dest.Content)

	assert.NotNil(t, tester.Decode(source, dest), "decode: dest expect a point")
}
