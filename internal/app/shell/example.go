package shell

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Service) ExampleRouter(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
