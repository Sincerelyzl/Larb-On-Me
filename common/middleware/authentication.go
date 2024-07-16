package middleware

import (
	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"github.com/gin-gonic/gin"
)

func AuthenticationLOM(c *gin.Context) {
	authHeader := c.GetHeader(LOMCookieAuthPrefix)

	// prepare error response unauthorized
	if authHeader == "" {
		errResponse := utils.NewErrorResponse(401, "unauthenticated")
		c.JSON(errResponse.StatusCode, errResponse)
		c.Abort()
		return
	}

	var lomUser models.LOMUser
	err := ClaimsLOM(authHeader, lomUser)
	if err != nil {
		errResponse := utils.NewErrorResponse(401, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		c.Abort()
		return
	}

	c.Set(LOMUserPrefix, lomUser)
	c.Next()
}
