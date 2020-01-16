// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
package team

var Core *Service

type Service struct {
}

func New() *Service {
	return &Service{}
}

func init() {
	Core = New()
}
