package middleware

import (
	"net/url"
	"strings"

	"github.com/Sincerelyzl/larb-on-me/common/models"
	"github.com/Sincerelyzl/larb-on-me/common/utils"
	"github.com/gin-gonic/gin"
)

func AuthenticationLOM(c *gin.Context) {
	authHeader := c.GetHeader(LOMCookieAuthPrefix)
	decodedAuthHeader, err := url.QueryUnescape(authHeader)
	decodedAuthHeader = strings.ReplaceAll(decodedAuthHeader, " ", "+")

	if err != nil {
		errResponse := utils.NewErrorResponse(401, "unauthenticated")
		c.JSON(errResponse.StatusCode, errResponse)
		c.SetCookie(LOMCookieAuthPrefix, "", -1, "/", "", false, true)
		c.Abort()
	}

	// prepare error response unauthorized
	if decodedAuthHeader == "" {
		errResponse := utils.NewErrorResponse(401, "unauthenticated")
		c.JSON(errResponse.StatusCode, errResponse)
		c.SetCookie(LOMCookieAuthPrefix, "", -1, "/", "", false, true)
		c.Abort()
		return
	}

	lomUser := models.UserAuthenticationLOM{}
	err = ClaimsLOM(decodedAuthHeader, &lomUser)
	if err != nil {
		errResponse := utils.NewErrorResponse(401, err.Error())
		c.JSON(errResponse.StatusCode, errResponse)
		c.SetCookie(LOMCookieAuthPrefix, "", -1, "/", "", false, true)
		c.Abort()
		return
	}

	c.Set(LOMUserPrefix, lomUser)
	c.Next()
}
