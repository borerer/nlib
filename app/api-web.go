package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *App) getHomepageHandler(c *gin.Context) {
	c.String(http.StatusOK, "good")
}
