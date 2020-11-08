package core

import (
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
	"selfcheck-go/database"
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

func RsaEncrypt(origData string) string {
	block, _ := pem.Decode(publicKey)
	pubInterface, _ := x509.ParsePKIXPublicKey(block.Bytes)
	pub := pubInterface.(*rsa.PublicKey)
	enc, _ := rsa.EncryptPKCS1v15(rand.Reader, pub, []byte(origData))
	return base64.StdEncoding.EncodeToString(enc)
}

func doLogin(name, birth, school, url string) (string, error) {
	val := map[string]string{
		"orgCode":   school,
		"loginType": "school",
		"stdntPNo":  "",
		"name":      RsaEncrypt(name),
		"birthday":  RsaEncrypt(birth),
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
	return token, nil
}

func doSumit(name, fname, url, token string) (string, error) {
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

func Selfcheck(name, birth, school, url string) (string, error) {
	token, err := doLogin(name, birth, school, url)
	if err != nil {
		return "", err
	}
	res, err := doSumit(name, "", url, token)
	if err != nil {
		return "", err
	}
	return res, nil
}

func Selfcheck2(name, birth, org, prefix, ip string) (string, string, string, error) {
	url, city, schulNm, err := database.SearchURL(org)
	if err != nil {
		return "", "", "", err
	}
	token, err := doLogin(name, birth, org, url)
	if err != nil {
		return "", "", "", err
	}
	fname := ""
	if prefix != "" {
		fname = GenerateResult(name, city)
	}
	res, err := doSumit(name, fname, url, token)
	if err != nil {
		return "", "", "", err
	}
	return res, schulNm, fname, nil
}
