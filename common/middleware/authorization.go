package middleware

import (
	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"github.com/gin-gonic/gin"
)

func Authorization(roles map[string]bool) func(c *gin.Context) {
	return func(c *gin.Context) {
		value, exist := c.Get(LOMUserPrefix)
		if !exist {
			errResponse := utils.NewErrorResponse(401, "unauthorized")
			c.JSON(errResponse.StatusCode, errResponse)
			c.Abort()
			return
		}

		user, ok := value.(models.UserAuthenticationLOM)
		if !ok {
			errResponse := utils.NewErrorResponse(401, "unauthorized user")
			c.JSON(errResponse.StatusCode, errResponse)
			c.Abort()
			return
		}

		// @TODO: implement role-based authorization
		if !roles[user.PermissionGroup] {
			errResponse := utils.NewErrorResponse(401, "unauthorized role permission denied")
			c.JSON(errResponse.StatusCode, errResponse)
			c.Abort()
			return
		}
	}
}
