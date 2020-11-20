package core

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Join(grade, class string) []map[string]string {
	val := map[string]string{
		"orgCode":   "B100000441",
		"grade":     grade,
		"classCode": class,
	}
	jsonValue, _ := json.Marshal(val)
	req, _ := http.NewRequest("POST", "https://senhcs.eduro.go.kr/join", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiIsImlhdCI6MTYwNTY4Njc0OTc2MH0.eyJubyI6IjIwMjAwMDA0NDAiLCJvcmciOiJCMTAwMDAwNDQxIiwiZXhwIjoxNjM3MjIyNzQ5LCJpYXQiOjE2MDU2ODY3NDk3NjB9.2QGYY1wFo11StgcZijuUVRdic-y0D04dsEcL8EextuQ")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string][]map[string]string

	_ = json.Unmarshal(body, &data)
	joinList := data["joinList"]
	return joinList
}
