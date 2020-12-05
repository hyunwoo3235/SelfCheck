package core

import (
	"SelfCheck/database"
	"github.com/gin-gonic/gin"
)

func Schulis(c *gin.Context) {
	pName := c.DefaultQuery("pName", "")
	frnoRidno := c.DefaultQuery("frnoRidno", "")
	schulNm := c.DefaultQuery("schulNm", "")
	geoNm := c.DefaultQuery("geoNm", "")
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
	}
}
