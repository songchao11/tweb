package main

import (
	"github.com/gin-gonic/gin"
	"tweb/handler"
)

func main(){

	router := gin.Default()

	handler.Init(router)

	router.Run(":9090")
}
