package test

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"net/http"
	"tweb/global"
)

func TestHandler(c *gin.Context) {
	type postParam struct {
		Text string `json:"text"`
	}

	param := &postParam{}
	resp := &global.Response{}
	defer c.JSON(http.StatusOK, resp)

	err := c.ShouldBindJSON(param)
	if err != nil {
		log.Error("illegal body:", err)
		resp.Code = global.ErrCodeParamInvalid
		resp.Msg = "invalid body"
		return
	}

	resp.Data = param.Text
	resp.Code = global.ErrCodeSuccess
	resp.Msg = "success"
	return
}
