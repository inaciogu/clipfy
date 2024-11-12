package handler

import (
	"clipfy/internal/api/command"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UploadHandler struct {
	uploadFile *command.UploadFileCommand
}

func NewUploadHandler(uploadFile *command.UploadFileCommand) *UploadHandler {
	return &UploadHandler{
		uploadFile: uploadFile,
	}
}

func (u *UploadHandler) Handle(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		http.Error(c.Writer, "Error getting file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	input := &command.UploadFileCommandInput{
		FileName:      header.Filename,
		File:          file,
		ContentType:   c.Request.Header.Get("Content-Type"),
		ContentLength: header.Size,
	}

	u.uploadFile.Execute(input)

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
	})
}
