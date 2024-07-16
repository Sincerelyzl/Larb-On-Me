package middleware

import (
	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"github.com/gin-gonic/gin"
)

func Authorization(c *gin.Context, roles map[string]*bool) {
	value, exist := c.Get(LOMUserPrefix)
	if !exist {
		errResponse := utils.NewErrorResponse(401, "unauthorized")
		c.JSON(errResponse.StatusCode, errResponse)
		c.Abort()
		return
	}

	lomUser, ok := value.(models.LOMUser)
	if !ok {
		errResponse := utils.NewErrorResponse(401, "unauthorized user")
		c.JSON(errResponse.StatusCode, errResponse)
		c.Abort()
		return
	}

	if roles[lomUser.Role] == nil || !*roles[lomUser.Role] {
		errResponse := utils.NewErrorResponse(401, "unauthorized role permission denied")
		c.JSON(errResponse.StatusCode, errResponse)
		c.Abort()
		return
	}

	c.Next()
}
