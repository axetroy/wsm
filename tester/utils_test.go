package tester_test

import (
	"testing"

	"github.com/axetroy/terminal/internal/library/controller/admin"
	"github.com/axetroy/terminal/internal/library/controller/auth"
	"github.com/axetroy/terminal/tester"
	"github.com/stretchr/testify/assert"
)

func init() {
	admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "admin",
		Password: "admin",
		Name:     "admin",
	}, true)
}

func TestCreateUser(t *testing.T) {
	user, err := tester.CreateUser()

	assert.Nil(t, err)

	defer auth.DeleteUserByUserName(user.Username)

	assert.NotEmpty(t, user.Username)
	assert.NotEmpty(t, user.Token)
	assert.NotEmpty(t, user.Id)
}

func TestLoginAdmin(t *testing.T) {
	admin, err := tester.LoginAdmin()

	assert.Nil(t, err)

	assert.NotEmpty(t, admin.Token)
	assert.NotEmpty(t, admin.Id)
	assert.Equal(t, "admin", admin.Username)
	assert.Equal(t, "admin", admin.Name)
}
