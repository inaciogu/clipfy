package handler

import (
	"clipfy/internal/api/command"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SegmentsHandler struct {
	listSegments *command.ListSegmentsCommand
}

func NewSegmentsHandler(listSegments *command.ListSegmentsCommand) *SegmentsHandler {
	return &SegmentsHandler{
		listSegments: listSegments,
	}
}

func (h *SegmentsHandler) ListSegments(c *gin.Context) {
	parentID := c.Param("parent_id")
	fmt.Println(parentID)

	output, err := h.listSegments.Execute(parentID)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, output)
}
