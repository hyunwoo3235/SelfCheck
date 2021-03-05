package eduro

import (
	"SelfCheck/database"
	"github.com/gin-gonic/gin"
)

func F4p(c *gin.Context) {
	pName := c.DefaultQuery("pName", "")
	frnoRidno := c.DefaultQuery("frnoRidno", "")
	schulNm := c.DefaultQuery("schulNm", "")
	geoNm := c.DefaultQuery("geoNm", "")
	al := c.DefaultQuery("al", "")
	if schulNm != "" && geoNm != "" && pName != "" && frnoRidno != "" {
		maps := database.SearchSchul(geoNm, schulNm)
		if len(maps) == 1 {
			maps[0]["pName"] = rsaEncrypt(pName)
			maps[0]["frnoRidno"] = rsaEncrypt(frnoRidno)
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
			"pName":     rsaEncrypt(pName),
			"frnoRidno": rsaEncrypt(frnoRidno),
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

func Schulis(c *gin.Context) {
	pName := c.DefaultQuery("pName", "")
	frnoRidno := c.DefaultQuery("frnoRidno", "")
	schulNm := c.DefaultQuery("schulNm", "")
	geoNm := c.DefaultQuery("geoNm", "")
	if schulNm != "" && geoNm != "" && pName != "" && frnoRidno != "" {
		maps := database.SearchSchul(geoNm, schulNm)
		if len(maps) == 1 {
			maps[0]["pName"] = rsaEncrypt(pName)
			maps[0]["frnoRidno"] = rsaEncrypt(frnoRidno)
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

func Elifstatus(c *gin.Context) {
	org := c.DefaultQuery("org", "")
	grade := c.DefaultQuery("grade", "")
	class := c.DefaultQuery("class", "")
	isJson := c.DefaultQuery("json", "false")

	if org == "" {
		c.JSON(200, gin.H{
			"error": "헤으응",
		})
		return
	}

	url, geo, _, _ := database.SearchURL(org)

	r, err := Join(url, org, grade, class, database.GT(geo))
	if err != nil {
		c.JSON(200, gin.H{
			"message": "서버 오류가 발생했습니다",
		})
		return
	}
	var res []map[string]string
	for _, i := range r {
		rspns00 := i["rspns00"]
		surveyYn := i["surveyYn"]

		isStmptom := "N"
		if rspns00 == "N" && surveyYn == "Y" {
			isStmptom = "Y"
		}
		res = append(res, map[string]string{
			"attNumber":   i["stdntCnEncpt"],
			"isSurvey":    i["surveyYn"],
			"isSymptom":   isStmptom,
			"registerDtm": i["registerDtm"],
		})
	}

	c.JSON(200, gin.H{
		"orgCode":    org,
		"grade":      grade,
		"class":      class,
		"isJson":     isJson,
		"surveyList": res,
	})
}
