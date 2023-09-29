package files

import (
	"context"
	"fmt"
	"helloGo/jwt-api/model"
	"helloGo/jwt-api/orm"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	projectID  = "fleet-parity-400408" // FILL IN WITH YOURS
	bucketName = "navy-12"             // FILL IN WITH YOURS
)

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	bucketName string
	uploadPath string
}

var uploader *ClientUploader

func init() {

	// Creates a client.
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	uploader = &ClientUploader{
		cl:         client,
		bucketName: bucketName,
		projectID:  projectID,
		uploadPath: "files/",
	}

}
func UploadFile(c *gin.Context) {
	f, header, err := c.Request.FormFile("file")
	if err != nil {
		logrus.Errorf("find error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer f.Close()
	originalFileName := strings.TrimSuffix(filepath.Base(header.Filename), filepath.Ext(header.Filename))
	if err != nil {
		logrus.Errorf("find error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = uploader.UploadFile(f, originalFileName)
	if err != nil {
		logrus.Errorf("find error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	url := `https://storage.cloud.google.com/` + uploader.bucketName + "/" + uploader.uploadPath + originalFileName
	file := orm.Files{Filename: originalFileName, Bucket: uploader.bucketName, Type: filepath.Ext(header.Filename), Path: uploader.uploadPath, Url: url}

	err = orm.DB.Create(&file).Error // pass pointer of data to Create
	if err != nil {
		logrus.Errorf("find error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    &file,
	})

}

// UploadFile     godoc
// @Summary		Update Users
// @Tags       File
// @Produce		json
// @Router			/upload [post]
// @param Body formData file true "body"
// @Success 200 {object} model.Response "OK"
// @Failure 400 {object} model.Response "Bad Request"
// @Failure 500 {object} model.Response "Internal Server Error"
func (c *ClientUploader) UploadFile(file multipart.File, object string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Upload an object with storage.Writer.
	wc := c.cl.Bucket(c.bucketName).Object(c.uploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}
