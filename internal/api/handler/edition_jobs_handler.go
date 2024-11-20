package handler

import (
	"clipfy/internal/api/command"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type EditionJobsHandler struct {
	createEditionJob *command.CreateEditionJob
}

func NewEditionJobsHandler(createEditionJob *command.CreateEditionJob) *EditionJobsHandler {
	return &EditionJobsHandler{
		createEditionJob: createEditionJob,
	}
}

func (e *EditionJobsHandler) CreateEditionJob(c *gin.Context) {
	var input command.CreateEditionJobInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.MustGet("user").(jwt.MapClaims)
	input.UserID = user["sub"].(string)

	output := e.createEditionJob.Execute(&input)

	c.JSON(http.StatusOK, output)
}
