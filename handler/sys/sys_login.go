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
		Account  string `json:"account"`
		Password string `json:"password"`
	}

	param := &loginParam{}

	err := c.ShouldBindJSON(param)
	if err != nil {
		log.Error("illegal body:", err)
		resp := &global.Response{}
		resp.Code = global.ErrCodeParamInvalid
		resp.Msg = "invalid body"
		c.JSON(http.StatusOK, resp)
		return
	}

	// 转由jwt中间件处理登录流程
	login := &auth.Login{}
	login.Account = param.Account
	login.Password = param.Password

	c.Set(auth.LoginKey, login)
	auth.JwtWrapper.LoginHandler(c)
}
