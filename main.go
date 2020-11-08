package main

import (
	"SelfCheck/core"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("template/*")
	r.Static("/assets", "./assets")
	r.GET("/regist", core.Regist)
	r.GET("/", core.DoSC)
	_ = r.Run(":80")
}
