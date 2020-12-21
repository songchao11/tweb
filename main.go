package main

import (
	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"tweb/global"
	"tweb/handler"
	"tweb/model"
)

func main() {

	conf, err := global.LoadConfig()
	if err != nil {
		log.Fatalf("load config failed, %s", err)
	}

	router := gin.Default()

	//init data model
	model.Init(conf.DBType, conf.DBAddr, conf.DebugDB)

	//init admin account
	model.InitAdmin()

	// 应用其它配置项
	if !conf.Apply() {
		log.Fatalf("apply conf failed!")
	}

	handler.Init(router)

	router.Run(conf.Endpoint)
}
