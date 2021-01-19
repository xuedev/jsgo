package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Req struct {
	Url     string `json:"url"`
	Method  string `json:"method"`
	Header  string `json:"header"`
	Timeout string `json:"timeout"`
	Body    string `json:"body"`
}

//go build --buildmode=plugin -o test_v1.0.so plugins/test.go
func Download(in string) (string, error) {
	var p map[string]string
	err := json.Unmarshal([]byte(in), &p)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// Get the data
	resp, err := http.Get(p["url"])
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 创建一个文件用于保存
	out, err := os.Create(p["file"])
	if err != nil {
		return "", err
	}
	defer out.Close()

	// 然后将响应流和文件流对接起来
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return "ok", nil
}
func Request(in string) (string, error) {
	p := Req{}
	err := json.Unmarshal([]byte(in), &p)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var header map[string]interface{}
	err = json.Unmarshal([]byte(p.Header), &header)
	if err != nil {
		return "", err
	}
	var req *http.Request
	var client *http.Client
	if len(p.Timeout) > 0 {
		i, _ := strconv.Atoi(p.Timeout)
		t := time.Duration(i)
		client = &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					deadline := time.Now().Add(t * time.Second)
					c, err := net.DialTimeout(netw, addr, time.Second*t)
					if err != nil {
						return nil, err
					}
					c.SetDeadline(deadline)
					return c, nil
				},
			},
		}
	} else {
		client = &http.Client{}
	}

	if p.Method == "post" {
		req, err = http.NewRequest("POST", p.Url, strings.NewReader(p.Body))
	} else {
		req, err = http.NewRequest("GET", p.Url, nil)
	}

	if err != nil {
		return "", err
	}
	println(p.Method)
	req.Header.Set("Accept", "*/*")
	for k, v := range header {
		println(k + ":" + fmt.Sprintf("%s", v))
		req.Header.Set(k, fmt.Sprintf("%s", v))

	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)

	return string(result), nil
}

func Upload(in string) (string, error) {
	var p map[string]string
	err := json.Unmarshal([]byte(in), &p)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	url := p["url"]
	file := p["file"]
	params := map[string]string{}
	req, err := NewFileUploadRequest(url, file, params)
	if err != nil {
		e := fmt.Sprintf("error to new upload file request:%s\n", err.Error())
		return "", errors.New(e)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)

	return string(result), nil
}

// NewFileUploadRequest ...
func NewFileUploadRequest(url, path string, params map[string]string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	// 文件写入 body
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("upload-key", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	// 其他参数列表写入 body
	for k, v := range params {
		if err := writer.WriteField(k, v); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, err
}

func Hmac256(in string) (string, error) {
	var data map[string]string
	err := json.Unmarshal([]byte(in), &data)
	if err != nil {
		return "", err
	}
	return ComputeHmac256(data["message"], data["secret"]), nil
}

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func main() {
	start := time.Now().Unix()
	str := "12345678"
	t := 1605003532771
	m := "f1b21ebfb60846e7ba16bc7c5b14f910" + fmt.Sprintf("%v", t) + str
	k := "e74c2e360d8e4f96b4fba96f964c9ffa"
	s := ComputeHmac256(m, k)
	println(s)

	p := `{
		"url": "http://data.jinan.gov.cn/gateway/api/1/cydwlhdjgsxx?XKZH=JY23701020190003",
		"header": "{\"X-Client-Id\": \"f1b21ebfb60846e7ba16bc7c5b14f910\",\"X-Timestamp\": \"` + fmt.Sprintf("%v", t) + `\",\"X-Nonce\": \"` + str + `\",\"X-Signature\": \"` + s + `\"}",
		"method":"post",
		"timeout":"1",
		"body": ""
	}`
	r, e := Request(p)
	println("resp:" + r)

	fmt.Println(e)
	end := time.Now().Unix()
	fmt.Println(end - start)
	/**
	p := `{
		"url": "http://81.70.218.13/api/v1/xraw/public.v8.test",
		"file": "a.json"
	}`
	r, e := Download(p)
	println(r)
	fmt.Println(e)
	*/
	/**
	p := `{
		"url": "http://api.weixin.qq.com/cgi-bin/media/voice/addvoicetorecofortext?access_token=ACCESS_TOKEN&format=&voice_id=xxxxxx&lang=zh_CN",
		"file": "a.json"
	}`
	r, e := Upload(p)
	println(r)
	fmt.Println(e)
	**/
}
