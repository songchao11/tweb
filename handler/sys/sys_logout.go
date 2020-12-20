package sys

import (
	"github.com/gin-gonic/gin"
	"tweb/handler/auth"
)

func LogoutHandler(c *gin.Context) {
	auth.JwtWrapper.LogoutHandler(c)
}
