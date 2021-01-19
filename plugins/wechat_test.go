package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

var (
	APPID          = "wxc5c96be3cb9b4652"
	APPSECRET      = "0926be326cb0ff53d39c52114831c1bc"
	SentTemplateID = "DVmeujVmcP5R4_a_aoreiRdfh5RMfELfiQKUWXQRemA" //模板ID，替换成自己的
)

type token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type sentence struct {
	Content     string `json:"content"`
	Note        string `json:"note"`
	Translation string `json:"translation"`
}

//发送todo
func sendTodo(content string, note string) (string, error) {
	fxurl := "https://mp.weixin.qq.com/debug/cgi-bin/sandboxinfo?action=showinfo&t=sandbox/index"
	access_token := getaccesstoken()
	if access_token == "" {
		return "", errors.New("get accesstocken failed")
	}

	flist := getflist(access_token)
	if flist == nil {
		return "", errors.New("cust has not follow")
	}
	now := time.Now()
	reqdata := "{\"content\":{\"value\":\"" + content + "\", \"color\":\"#0000CD\"}, \"note\":{\"value\":\"" + note + "\"}, \"time\":{\"value\":\"" + now.Format("2006-01-02 15:04:05") + "\"}}"
	for _, v := range flist {
		templatepost(access_token, reqdata, fxurl, SentTemplateID, v.Str)
	}
	return "ok", nil
}

//获取微信accesstoken
func getaccesstoken() string {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", APPID, APPSECRET)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取微信token失败", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("微信token读取失败", err)
		return ""
	}

	token := token{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println("微信token解析json失败", err)
		return ""
	}

	return token.AccessToken
}

//获取关注者列表
func getflist(access_token string) []gjson.Result {
	url := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + access_token + "&next_openid="
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取关注列表失败", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取内容失败", err)
		return nil
	}
	flist := gjson.Get(string(body), "data.openid").Array()
	return flist
}

//发送模板消息
func templatepost(access_token string, reqdata string, fxurl string, templateid string, openid string) (string, error) {
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + access_token

	reqbody := "{\"touser\":\"" + openid + "\", \"template_id\":\"" + templateid + "\", \"url\":\"" + fxurl + "\", \"data\": " + reqdata + "}"

	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(string(reqbody)))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(body), nil
}

//go build --buildmode=plugin -o test_v1.0.so plugins/test.go
func Call(in string) (string, error) {
	fmt.Println("call test.call v2 ,param:", in)
	return in, nil
}

func main() {
	ret, err := sendTodo("test", "remark")
	fmt.Println(ret)
	fmt.Println(err)
}
