package core

import (
	"SelfCheck/database"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Regist(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	birth := c.DefaultQuery("birth", "")
	geo := c.DefaultQuery("geo", "")
	org := c.DefaultQuery("orgName", "")
	c.Header("Content-Type", "text/html")
	if name != "" {
		schul := database.SearchSchul(geo, org)[0]
		_, err := DoLogin(name, birth, schul["schulCode"], schul["url"])
		if err != nil {
			c.HTML(200, "register-fail.html", gin.H{
				"msg": "에러가 발생했습니다",
			})
			return
		}
		c.SetCookie("name", name, 5184000, "/", "127.0.0.1", false, false)
		c.SetCookie("birth", birth, 5184000, "/", "127.0.0.1", false, false)
		c.SetCookie("org", schul["schulCode"], 5184000, "/", "127.0.0.1", false, false)
		c.SetCookie("url", schul["url"], 5184000, "/", "127.0.0.1", false, false)
		c.HTML(200, "register-success.html", gin.H{})
		return
	}
	c.HTML(200, "register.html", gin.H{})
}

func DoSC(c *gin.Context) {
	/*name, err := c.Request.Cookie("name")
	if err != nil {
		c.Redirect(200, "/regist")
	}
	birth, _ := c.Request.Cookie("birth")
	org, _ := c.Request.Cookie("org")
	url, _ := c.Request.Cookie("url")*/
	cookies := c.Request.Cookies()
	for i := range cookies {
		fmt.Println(string(rune(i)))
	}

}
