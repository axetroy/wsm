// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

func New(text string, code int) Error {
	return Error{
		message: text,
		code:    code,
	}
}

type Error struct {
	message string
	code    int
}

func (e Error) Error() string {
	return e.message
}

func (e Error) Code() int {
	return e.code
}
