package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (api *API) getHomepageHandler(c *gin.Context) {
	c.String(http.StatusOK, "good")
}
