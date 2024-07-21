package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		res := NewSuccessResponse(http.StatusOK, "ok", nil)
		c.JSON(res.StatusCode, res)
	}
}
