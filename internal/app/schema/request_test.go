package schema_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery_formatSort(t *testing.T) {
	q := NewQuery()

	assert.Equal(t, []Sort{
		{
			Field: "created_at",
			Order: OrderDesc,
		},
	}, q.FormatSort())

	q.Sort = q.Sort + ",-name,balance"

	assert.Equal(t, []Sort{
		{
			Field: "created_at",
			Order: OrderDesc,
		},
		{
			Field: "name",
			Order: OrderDesc,
		},
		{
			Field: "balance",
			Order: OrderAsc,
		},
	}, q.FormatSort())

	assert.Equal(t, DefaultLimit, q.Limit)
	assert.Equal(t, DefaultPage, q.Page)
}
