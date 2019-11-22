package oauth2

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/markbates/goth/gothic"
)

func (s *Service) CallbackRouter(c *gin.Context) {
	provider := c.Param("provider")
	c.Request = mux.SetURLVars(c.Request, map[string]string{"provider": provider})
	gothUser, err := gothic.CompleteUserAuth(c.Writer, c.Request)

	if err != nil {
		return
	}

	redirectToClient(c, &gothUser)
}
