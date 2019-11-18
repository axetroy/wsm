// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role

import (
	"errors"
	"github.com/axetroy/terminal/core/exception"
	"github.com/axetroy/terminal/core/helper"
	"github.com/axetroy/terminal/core/rbac/accession"
	"github.com/axetroy/terminal/core/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAccession() (res schema.Response) {
	var (
		err  error
		data []*accession.Accession
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		helper.Response(&res, data, err)
	}()

	data = accession.List

	return
}

func GetAccessionRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	res = GetAccession()
}
