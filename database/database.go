package database

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var sqdb, _ = sql.Open("sqlite3", "database.db")

var urlList = map[string]string{
	"서울": "https://senhcs.eduro.go.kr/",
	"경기": "https://goehcs.eduro.go.kr/",
	"대전": "https://djehcs.eduro.go.kr/",
	"대구": "https://dgehcs.eduro.go.kr/",
	"부산": "https://penhcs.eduro.go.kr/",
	"인천": "https://icehcs.eduro.go.kr/",
	"광주": "https://genhcs.eduro.go.kr/",
	"울산": "https://usehcs.eduro.go.kr/",
	"세종": "https://sjehcs.eduro.go.kr/",
	"충북": "https://cbehcs.eduro.go.kr/",
	"충남": "https://cnehcs.eduro.go.kr/",
	"경북": "https://gbehcs.eduro.go.kr/",
	"경남": "https://gnehcs.eduro.go.kr/",
	"강원": "https://kwehcs.eduro.go.kr/",
	"전북": "https://jbehcs.eduro.go.kr/",
	"전남": "https://jnehcs.eduro.go.kr/",
	"제주": "https://jjehcs.eduro.go.kr/",
}

func SearchSchul(cityNm, schulNm string) []map[string]string {
	query := "SELECT * FROM school_data WHERE schulNM LIKE '%" + schulNm + "%'"
	if cityNm != "" {
		query += " AND cityNm = '" + cityNm + "'"
	}
	row, err := sqdb.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var maps []map[string]string
	for row.Next() {
		var cityNm string
		var schulNm string
		var schulCode string
		_ = row.Scan(&cityNm, &schulNm, &schulCode)
		maps = append(maps, map[string]string{
			"cityNm":    cityNm,
			"schulNm":   schulNm,
			"schulCode": schulCode,
			"url":       urlList[cityNm],
		})
	}
	return maps
}

func SearchURL(schulCode string) (string, string, string, error) {
	query := "SELECT * FROM school_data WHERE schulCode = '" + schulCode + "'"
	row, err := sqdb.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	maps := map[string]string{
		"cityNm":  "",
		"url":     "",
		"schulNm": "",
	}
	for row.Next() {
		var cityNm string
		var schulNm string
		var schulCode string
		_ = row.Scan(&cityNm, &schulNm, &schulCode)
		maps = map[string]string{
			"cityNm":  cityNm,
			"url":     urlList[cityNm],
			"schulNm": schulNm,
		}
	}
	if maps["url"] == "" {
		return "", "", "", errors.New("학교 검색중 에러가 발생했습니다.")
	}
	return maps["url"], maps["cityNm"], maps["schulNm"], nil
}
