package handler

import (
	"clipfy/internal/api/command"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type EditionJobsHandler struct {
	createEditionJob *command.CreateEditionJob
	listEditionJobs  *command.ListEditionJobs
}

func NewEditionJobsHandler(createEditionJob *command.CreateEditionJob, listEditionJobs *command.ListEditionJobs) *EditionJobsHandler {
	return &EditionJobsHandler{
		createEditionJob: createEditionJob,
		listEditionJobs:  listEditionJobs,
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

	output, err := e.createEditionJob.Execute(&input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

func (e *EditionJobsHandler) ListEditionJobs(c *gin.Context) {
	user := c.MustGet("user").(jwt.MapClaims)
	userID := user["sub"].(string)

	output := e.listEditionJobs.Execute(&command.ListEditionJobsInput{
		UserID: userID,
	})

	c.JSON(http.StatusOK, output)
}
