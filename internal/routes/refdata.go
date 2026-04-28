package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rm-hull/next-departures-api/internal/models"
)

func StopTypes(c *gin.Context) {
	c.JSON(200, models.StopTypes)
}
