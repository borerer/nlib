package app

import (
	"io"
	"net/http"
	"path"

	"github.com/borerer/nlib/logs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (app *App) getFileHandler(c *gin.Context) {
	appID := c.Query("app")
	file := c.Query("file")
	filename := path.Join("apps", appID, file)
	fileReader, err := app.minioClient.GetFile(filename)
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
	logs.Info("get file", zap.String("filename", filename), zap.Int64("size", n))
}

func (app *App) putFileHandler(c *gin.Context) {
	appID := c.Query("app")
	file := c.Query("file")
	filename := path.Join("apps", appID, file)
	defer c.Request.Body.Close()
	n, err := app.minioClient.PutFile(filename, true, c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	logs.Info("put file", zap.String("file", filename), zap.Int64("size", n))
	c.JSON(http.StatusOK, ResponseGeneralOK)
}

func (app *App) deleteFileHandler(c *gin.Context) {
	appID := c.Query("app")
	file := c.Query("file")
	filename := path.Join("apps", appID, file)
	err := app.minioClient.DeleteFile(filename)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	logs.Info("delete file", zap.String("file", filename))
	c.JSON(http.StatusOK, ResponseGeneralOK)
}

func (app *App) fileStatsHandler(c *gin.Context) {
	appID := c.Query("app")
	file := c.Query("file")
	filename := path.Join("apps", appID, file)
	stats, err := app.minioClient.HeadFile(filename)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	logs.Info("file stats", zap.String("file", filename))
	c.JSON(http.StatusOK, stats)
}

func (app *App) listFolderHandler(c *gin.Context) {
	appID := c.Query("app")
	folder := c.Query("folder")
	prefix := path.Join("apps", appID, folder)
	res, err := app.minioClient.ListFolder(prefix)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	logs.Info("list folder", zap.String("prefix", prefix))
	c.JSON(http.StatusOK, res)
}
