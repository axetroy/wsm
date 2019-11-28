package helper

import (
	"github.com/axetroy/terminal/internal/app/exception"
	"github.com/axetroy/terminal/internal/app/schema"
)

func Response(res *schema.Response, data interface{}, meta *schema.Meta, err error) {
	if err != nil {
		res.Data = nil
		res.Message = err.Error()
		res.Status = exception.GetCodeFromError(err)
		res.Meta = nil
	} else {
		res.Data = data
		res.Status = schema.StatusSuccess
	}
}
