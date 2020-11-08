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

	v1 := r.Group("/api")
	{
		v1.GET("", core.F4p)
		v1.GET("/jaga", core.SC)
	}

	_ = r.Run(":80")
}
