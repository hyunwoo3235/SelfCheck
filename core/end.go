package core

import (
	"SelfCheck/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
)

func Regist(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	birth := c.DefaultQuery("birth", "")
	geo := c.DefaultQuery("geo", "")
	org := c.DefaultQuery("orgName", "")
	c.Header("Content-Type", "text/html")
	if name != "" {
		lis := database.SearchSchul(geo, org)
		if len(lis) != 1 {
			c.HTML(200, "register-fail.html", gin.H{
				"msg": "학교가 여러개 검색되거나 찾지 못했습니다.",
			})
			return
		}
		schul := lis[0]
		_, err := DoLogin(name, birth, schul["schulCode"], schul["url"])
		if err != nil {
			c.HTML(200, "register-fail.html", gin.H{
				"msg": "이름이나 학교, 생년월일을 한번 더 확인해 주세요",
			})
			return
		}

		c.SetCookie("name", name, 5184000, "/", "193.123.246.37", false, false)
		c.SetCookie("birth", birth, 5184000, "/", "193.123.246.37", false, false)
		c.SetCookie("org", schul["schulCode"], 5184000, "/", "193.123.246.37", false, false)
		c.SetCookie("url", schul["url"], 5184000, "/", "193.123.246.37", false, false)
		c.HTML(200, "register-success.html", gin.H{})
		return
	}
	c.HTML(200, "register.html", gin.H{})
}

func DoSC(c *gin.Context) {
	_, err := c.Request.Cookie("name")
	if err != nil {
		c.Redirect(301, "/regist")
		c.Abort()
		return
	}
	name, _ := c.Request.Cookie("name")
	birth, _ := c.Request.Cookie("birth")
	org, _ := c.Request.Cookie("org")
	ur, _ := c.Request.Cookie("url")

	names := name.Value
	births := birth.Value
	orgs := org.Value
	urls := ur.Value

	fmt.Println(Selfcheck(dc(names), births, orgs, dc(urls)))
	c.HTML(200, "selfcheck-done.html", gin.H{})
}

func F4p(c *gin.Context) {
	pName := c.DefaultQuery("pName", "")
	frnoRidno := c.DefaultQuery("frnoRidno", "")
	schulNm := c.DefaultQuery("schulNm", "")
	geoNm := c.DefaultQuery("geoNm", "")
	al := c.DefaultQuery("al", "")
	if schulNm != "" && geoNm != "" && pName != "" && frnoRidno != "" {
		maps := database.SearchSchul(geoNm, schulNm)
		if len(maps) == 1 {
			maps[0]["pName"] = RsaEncrypt(pName)
			maps[0]["frnoRidno"] = RsaEncrypt(frnoRidno)
			c.JSON(200, maps[0])
		} else {
			c.JSON(200, gin.H{
				"url":       "",
				"cityNm":    "",
				"schulCode": "",
				"schulNm":   "학교가 여러개가 검색되었습니다. 더 정확한 이름을 입력하거나, 학교 코드로 진행해 주세요.",
			})
		}
	} else if pName != "" && frnoRidno != "" {
		c.JSON(200, gin.H{
			"pName":     RsaEncrypt(pName),
			"frnoRidno": RsaEncrypt(frnoRidno),
		})
	} else if schulNm != "" && geoNm != "" {
		maps := database.SearchSchul(geoNm, schulNm)
		if len(maps) == 1 {
			c.JSON(200, maps[0])
		} else if al != "" {
			c.JSON(200, maps)
		} else {
			c.JSON(200, gin.H{
				"url":       "",
				"cityNm":    "",
				"schulCode": "",
				"schulNm":   "학교가 여러개가 검색되었습니다. 더 정확한 이름을 입력하거나, 학교 코드로 진행해 주세요.",
			})
		}
	} else {
		c.JSON(503, gin.H{"error": "인자좀 제대로 주세요요"})
	}
}

func SC(c *gin.Context) {
	name := c.DefaultQuery("name", "")
	birth := c.DefaultQuery("birth", "")
	org := c.DefaultQuery("org", "")
	pf := c.DefaultQuery("pf", "")
	if name == "" || birth == "" || org == "" {
		c.JSON(503, gin.H{"error": "인자좀 제대로 주세요요"})
		return
	}
	res, schulNm, fname, err := Selfcheck2(name, birth, org, pf)
	if err != nil {
		c.JSON(503, gin.H{"error": "인자에 오류있는데 인자좀 제대로 주세요요"})
		return
	}
	c.JSON(200, gin.H{
		"notice":      "코로나 의심 증상이 있을시 자가진단을 수동으로 하고, 주변 병원에 연락하세요",
		"status":      "success",
		"name":        name,
		"pfname":      fname,
		"school_name": schulNm,
		"submit_time": res,
	})
}

func JoinList(c *gin.Context) {
	c.JSON(200, Join("2", "09"))
}

func dc(q string) string {
	r, _ := url.PathUnescape(q)
	return r
}
