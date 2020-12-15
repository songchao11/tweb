package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"tweb/handler/auth"
	"tweb/handler/sys"
	"tweb/handler/test"
)

func Init(r *gin.Engine) {
	log.Info("初始化HTTP请求句柄")
	v1 := r.Group("/api/v1")

	v1.POST("/sys/login", sys.LoginHandler)

	auth.Init()
	jwtSecured := v1.Group("")
	jwtSecured.Use(auth.JwtWrapper.MiddlewareFunc())
	jwtSecured.Use(auth.JwtRefresher())
	{
		jwtSecured.POST("/test", test.TestHandler)
	}
}
