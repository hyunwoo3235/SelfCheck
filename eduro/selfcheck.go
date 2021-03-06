package eduro

import (
	"SelfCheck/database"
	"SelfCheck/util"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/http"
)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA81dCnCKt0NVH7j5Oh2+S
GgEU0aqi5u6sYXemouJWXOlZO3jqDsHYM1qfEjVvCOmeoMNFXYSXdNhflU7mjWP8
jWUmkYIQ8o3FGqMzsMTNxr+bAp0cULWu9eYmycjJwWIxxB7vUwvpEUNicgW7v5nC
wmF5HS33Hmn7yDzcfjfBs99K5xJEppHG0qc+q3YXxxPpwZNIRFn0Wtxt0Muh1U8a
vvWyw03uQ/wMBnzhwUC8T4G5NclLEWzOQExbQ4oDlZBv8BM/WxxuOyu0I8bDUDdu
tJOfREYRZBlazFHvRKNNQQD2qDfjRz484uFs7b5nykjaMB9k/EJAuHjJzGs9MMMW
tQIDAQAB
-----END PUBLIC KEY-----
`)

func rsaEncrypt(origData string) string {
	block, _ := pem.Decode(publicKey)
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	pub := pubInterface.(*rsa.PublicKey)
	enc, _ := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(origData))
	return base64.StdEncoding.EncodeToString(enc)
}

func DoLogin(name, birth, school, url string) (string, error) {
	val := map[string]string{
		"orgCode":   school,
		"loginType": "school",
		"stdntPNo":  "",
		"name":      rsaEncrypt(name),
		"birthday":  rsaEncrypt(birth),
	}
	jsonValue, _ := json.Marshal(val)
	req, _ := http.NewRequest("POST", url+"v2/findUser", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", errors.New("서버 오류가 발생했습니다. 이름이나 생년월일 또는 학교 정보를 다시 한번 확인해 주세요")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]string
	_ = json.Unmarshal(body, &data)
	token := data["token"]
	PNo := getPNo(url, token)
	if PNo == "" {
		return "", errors.New(" 학생 정보를 불러오는데 에러가 발생했습니다.")
	}
	token2 := getToken2(PNo, school, url, token)
	return token2, nil
}

func DoLoginPno(name, pno, school, url string) (string, error) {
	val := map[string]string{
		"orgCode":   school,
		"loginType": "school",
		"name":      rsaEncrypt(name),
		"userPNo":   pno,
	}
	jsonValue, _ := json.Marshal(val)
	req, _ := http.NewRequest("POST", url+"v2/findUser", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]string
	_ = json.Unmarshal(body, &data)
	return data["token"], nil
}

func validatePassword(url, token, pass string) (string, error) {
	val := map[string]string{
		"password": rsaEncrypt(pass),
	}
	jsonvalue, _ := json.Marshal(val)
	req, _ := http.NewRequest("POST", url+"v2/validatePassword", bytes.NewBuffer(jsonvalue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	client := &http.Client{}
	respon, _ := client.Do(req)
	defer respon.Body.Close()
	var resdata map[string]string
	body, _ := ioutil.ReadAll(respon.Body)
	_ = json.Unmarshal(body, &resdata)
	return resdata["token"], nil
}

func getToken2(pno, org, url, token string) string {
	val := map[string]string{
		"orgCode": org,
		"userPNo": pno,
	}
	jsonvalue, _ := json.Marshal(val)
	req, _ := http.NewRequest("POST", url+"v2/getUserInfo", bytes.NewBuffer(jsonvalue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	client := &http.Client{}
	respon, _ := client.Do(req)
	defer respon.Body.Close()
	var resdata map[string]string
	body, _ := ioutil.ReadAll(respon.Body)
	_ = json.Unmarshal(body, &resdata)
	return resdata["token"]
}

func getPNo(url, token string) string {
	jsonValue, _ := json.Marshal(map[string]string{})
	req, _ := http.NewRequest("POST", url+"v2/selectUserGroup", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var data []map[string]string
	_ = json.Unmarshal(body, &data)
	return data[0]["userPNo"]
}

func DoSumit(name, fname, url, token string) (string, error) {
	val := map[string]string{
		"rspns00":            "y",
		"rspns01":            "1",
		"rspns02":            "1",
		"rspns09":            "0",
		"deviceuuid":         "",
		"upperToken":         token,
		"upperUserNameEncpt": name,
	}
	if fname != "" {
		val["upperUserNameEncpt"] = fname
	}
	jsonvalue, _ := json.Marshal(val)
	req, _ := http.NewRequest("POST", url+"registerServey", bytes.NewBuffer(jsonvalue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	client := &http.Client{}
	respon, _ := client.Do(req)
	defer respon.Body.Close()
	var resdata map[string]string
	body, _ := ioutil.ReadAll(respon.Body)
	_ = json.Unmarshal(body, &resdata)
	return resdata["inveYmd"], nil
}

func Join(url, org, grade, class, token string) ([]map[string]string, error) {
	val := map[string]string{
		"orgCode":   org,
		"grade":     grade,
		"classCode": class,
	}
	jsonvalue, _ := json.Marshal(val)
	req, _ := http.NewRequest("POST", url+"join", bytes.NewBuffer(jsonvalue))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	client := &http.Client{}
	respon, _ := client.Do(req)
	defer respon.Body.Close()
	var resdata map[string][]map[string]string
	body, _ := ioutil.ReadAll(respon.Body)
	err := json.Unmarshal(body, &resdata)
	if err != nil {
		return resdata["joinList"], err
	}
	return resdata["joinList"], nil
}

func Selfcheck(name, birth, school, url string) (string, error) {
	token, err := DoLogin(name, birth, school, url)
	if err != nil {
		return "", err
	}
	res, err := DoSumit(name, "", url, token)
	if err != nil {
		return "", err
	}
	return res, nil
}

func Selfcheck2(name, birth, org, pass, prefix string) (string, string, string, error) {
	url, city, schulNm, err := database.SearchURL(org)
	if err != nil {
		return "", "", "", err
	}
	token, err := DoLogin(name, birth, org, url)
	if err != nil {
		return "", "", "", err
	}
	fname := ""
	if prefix != "" {
		fname = util.GenerateResult(name, city)
	}
	res, err := DoSumit(name, fname, url, token)
	if err != nil {
		return "", "", "", err
	}
	return res, schulNm, fname, nil
}
