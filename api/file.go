package api

import (
	"io"
	"net/http"
	"path"

	"github.com/borerer/nlib/logs"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (api *API) getFileHandler(c *gin.Context) {
	appID := c.Query("app")
	file := c.Query("file")
	filename := path.Join("apps", appID, file)
	fileReader, err := api.minioClient.GetFile(filename)
	if err != nil {
		abort500(c, err)
		return
	}
	defer fileReader.Close()
	n, err := io.Copy(c.Writer, fileReader)
	if err != nil {
		abort500(c, err)
		return
	}
	logs.Info("get file", zap.String("filename", filename), zap.Int64("size", n))
}

func (api *API) putFileHandler(c *gin.Context) {
	appID := c.Query("app")
	file := c.Query("file")
	filename := path.Join("apps", appID, file)
	defer c.Request.Body.Close()
	n, err := api.minioClient.PutFile(filename, true, c.Request.Body)
	if err != nil {
		abort500(c, err)
		return
	}
	logs.Info("put file", zap.String("file", filename), zap.Int64("size", n))
	c.JSON(http.StatusOK, ResponseGeneralOK)
}

func (api *API) deleteFileHandler(c *gin.Context) {
	appID := c.Query("app")
	file := c.Query("file")
	filename := path.Join("apps", appID, file)
	err := api.minioClient.DeleteFile(filename)
	if err != nil {
		abort500(c, err)
		return
	}
	logs.Info("delete file", zap.String("file", filename))
	c.JSON(http.StatusOK, ResponseGeneralOK)
}

func (api *API) fileStatsHandler(c *gin.Context) {
	appID := c.Query("app")
	file := c.Query("file")
	filename := path.Join("apps", appID, file)
	stats, err := api.minioClient.HeadFile(filename)
	if err != nil {
		abort500(c, err)
		return
	}
	logs.Info("file stats", zap.String("file", filename))
	c.JSON(http.StatusOK, stats)
}

func (api *API) listFolderHandler(c *gin.Context) {
	appID := c.Query("app")
	folder := c.Query("folder")
	prefix := path.Join("apps", appID, folder)
	res, err := api.minioClient.ListFolder(prefix)
	if err != nil {
		abort500(c, err)
		return
	}
	logs.Info("list folder", zap.String("prefix", prefix))
	c.JSON(http.StatusOK, res)
}
