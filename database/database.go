package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
)

type School struct {
	name    string
	org     string
	geoCode string
	url     string
}

var sqdb, _ = sql.Open("sqlite3", "database.db")

var tokens = buildData()

func buildData() map[string]map[string]string {
	jsonFile, _ := os.Open("authToken.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var authToken map[string]map[string]map[string]string
	_ = json.Unmarshal(byteValue, &authToken)
	return authToken["authToken"]
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
			"url":       "https://" + tokens[cityNm]["areaCode"] + "hcs.eduro.go.kr/",
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
			"url":     "https://" + tokens[cityNm]["areaCode"] + "hcs.eduro.go.kr/",
			"schulNm": schulNm,
		}
	}
	if maps["url"] == "" {
		return "", "", "", errors.New("학교 검색중 에러가 발생했습니다.")
	}
	return maps["url"], maps["cityNm"], maps["schulNm"], nil
}

func GT(geo string) string {
	return tokens[geo]["token"]
}
