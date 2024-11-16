package handler

import (
	"clipfy/internal/api/command"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UploadHandler struct {
	uploadFile *command.GenUploadFileURLCommand
}

func NewUploadHandler(uploadFile *command.GenUploadFileURLCommand) *UploadHandler {
	return &UploadHandler{
		uploadFile: uploadFile,
	}
}

func (u *UploadHandler) Handle(c *gin.Context) {
	var input command.GenUploadFileURLCommandInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output := u.uploadFile.Execute(&input)

	c.JSON(http.StatusOK, output)
}
