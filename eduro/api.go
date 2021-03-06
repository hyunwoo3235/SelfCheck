package eduro

import (
	"SelfCheck/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
	pass := c.DefaultQuery("pass", "")
	pf := c.DefaultQuery("pf", "")
	if name == "" || birth == "" || org == "" {
		c.JSON(503, gin.H{"error": "인자좀 제대로 주세요요"})
		return
	}
	res, schulNm, fname, err := Selfcheck(name, birth, org, pass, pf)
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

	if isJson == "false" {
		var res []string
		for _, i := range r {
			rspns00 := i["rspns00"]
			surveyYn := i["surveyYn"]

			var isStmptom, resDtm string
			switch {
			case surveyYn == "N":
				isStmptom = "미참여"
			case rspns00 == "N" && surveyYn == "Y":
				isStmptom = "유증상"
			default:
				isStmptom = "ㅤ정상"
			}
			if i["registerDtm"] == "" {
				resDtm = ""
			} else {
				resDtm = "ㅤ" + i["registerDtm"][:19] + " |"
			}
			res = append(res, fmt.Sprintf("|ㅤ%2s번ㅤ|ㅤ%sㅤ|%s", i["stdntCnEncpt"], isStmptom, resDtm))
		}
		header := fmt.Sprintf("학교 코드:%s\n%s학년 %s반\n\n\n", org, grade, class)
		c.String(http.StatusOK, header+strings.Join(res, "\n"))
		return
	}

	var res []map[string]string
	var stat = [3]int{0, 0, 0}
	for _, i := range r {
		rspns00 := i["rspns00"]
		surveyYn := i["surveyYn"]

		var isStmptom, registerDtm string
		switch {
		case surveyYn == "N":
			isStmptom = "N"
			registerDtm = ""
			stat[0] += 1
		case rspns00 == "N" && surveyYn == "Y":
			isStmptom = "Y"
			registerDtm = i["registerDtm"][:19]
			stat[2] += 1
		default:
			isStmptom = "N"
			registerDtm = i["registerDtm"][:19]
			stat[1] += 1
		}
		res = append(res, map[string]string{
			"attNumber":   i["stdntCnEncpt"],
			"isSurvey":    i["surveyYn"],
			"isSymptom":   isStmptom,
			"registerDtm": registerDtm,
		})
	}

	c.JSON(200, gin.H{
		"orgCode": org,
		"grade":   grade,
		"class":   class,
		"isJson":  isJson,
		"detail": map[string]int{
			"nonpart": stat[0],
			"normal":  stat[1],
			"symptom": stat[2],
		},
		"surveyList": res,
	})
}
