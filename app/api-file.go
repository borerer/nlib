package app

import (
	"io"
	"net/http"
	"path"

	"github.com/borerer/nlib/file"
	"github.com/borerer/nlib/logs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (app *App) getObjectHandler(c *gin.Context) {
	fileHelper := c.MustGet("file-helper").(file.FileHelper)
	appID := c.Query("app")
	file := c.Query("file")
	filename := path.Join("apps", appID, file)
	fileReader, err := fileHelper.GetFile(filename)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer fileReader.Close()
	n, err := io.Copy(c.Writer, fileReader)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	logs.Info("get object", zap.String("filename", filename), zap.Int64("size", n))
}

func (app *App) putObjectHandler(c *gin.Context) {
	fileHelper := c.MustGet("file-helper").(file.FileHelper)
	appID := c.Query("app")
	file := c.Query("file")
	filename := path.Join("apps", appID, file)
	defer c.Request.Body.Close()
	n, err := fileHelper.PutFile(filename, true, c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	logs.Info("put object", zap.String("file", filename), zap.Int64("size", n))
}
