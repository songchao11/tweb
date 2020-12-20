package main

import (
	"github.com/gin-gonic/gin"
	"tweb/handler"
	"tweb/model"
)

func main() {

	router := gin.Default()

	//init data model
	model.Init()

	//init admin account
	model.InitAdmin()

	handler.Init(router)

	router.Run(":9090")
}
