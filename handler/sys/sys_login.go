package sys

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"net/http"
	"tweb/global"
	"tweb/handler/auth"
)

func LoginHandler(c *gin.Context) {
	type loginParam struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	param := &loginParam{}
	resp := &global.Response{}
	defer c.JSON(http.StatusOK, resp)

	err := c.ShouldBindJSON(param)
	if err != nil {
		log.Error("illegal body:", err)
		resp.Code = global.ErrCodeParamInvalid
		resp.Msg = "invalid body"
		return
	}

	// 转由jwt中间件处理登录流程
	login := &auth.Login{}
	login.Username = param.Username
	login.Password = param.Password

	c.Set(auth.LoginKey, login)
	auth.JwtWrapper.LoginHandler(c)
}
