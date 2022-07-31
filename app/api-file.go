package app

import (
	"io"
	"net/http"

	"github.com/borerer/nlib/engines"
	"github.com/borerer/nlib/logs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (app *App) getObjectHandler(c *gin.Context) {
	appEngine := c.MustGet("app-engine").(*engines.AppEngine)
	file := c.Query("file")
	fileReader, err := appEngine.FileEngine.GetFile(file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	n, err := io.Copy(c.Writer, fileReader)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	logs.Info("get object", zap.String("file", file), zap.Int64("size", n))
}
