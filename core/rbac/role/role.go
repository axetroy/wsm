// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role

import (
	"github.com/axetroy/terminal/core/rbac/accession"
)

type Role struct {
	Name        string                `json:"name"`        // 角色名
	Description string                `json:"description"` // 角色描述
	Accession   []accession.Accession `json:"accession"`   // 角色拥有的权限
}

func New(name string, description string, accessions []accession.Accession) *Role {
	return &Role{
		Name:        name,
		Description: description,
		Accession:   accessions,
	}
}

func (r *Role) AccessionArray() (list []string) {
	for _, v := range r.Accession {
		list = append(list, v.Name)
	}
	return
}
