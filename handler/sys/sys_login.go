package sys

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"net/http"
	"tweb/handler/global"
)

func LoginHandler(c *gin.Context){
	type loginParam struct {
		Account string `json:"account"`
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

	resp.Data = param.Account
	resp.Code = global.ErrCodeSuccess
	resp.Msg = "success"
	return
}
