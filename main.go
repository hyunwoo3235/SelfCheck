package main

import (
	"SelfCheck/eduro"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("template/*")
	r.Static("/assets", "./assets")
	r.StaticFile("/favicon.ico", "./assets/favicon.svg")
	r.NoRoute(func(c *gin.Context) {
		c.Redirect(301, "https://youtu.be/atykIBND1bg?t=35")
		c.Abort()
	})

	v1 := r.Group("/api")
	{
		v1.GET("", eduro.F4p)
		v1.GET("/jaga", eduro.SC)
		v1.GET("/isSurvey", eduro.Elifstatus)
	}

	_ = r.Run(":80")
}
