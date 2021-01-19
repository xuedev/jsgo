package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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

//获取微信accesstoken
func AccessToken(in string) (string, error) {
	token := token{}
	var data map[string]string
	err := json.Unmarshal([]byte(in), &data)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", data["appid"], data["secret"])
	resp, err := http.Get(url)
	if err != nil {
		e := fmt.Sprintf("获取微信token失败", err)
		return "", errors.New(e)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e := fmt.Sprintf("获取微信token失败", err)
		return "", errors.New(e)
	}

	err = json.Unmarshal(body, &token)
	if err != nil {
		e := fmt.Sprintf("微信token解析json失败", err)
		return "", errors.New(e)
	}
	if token.AccessToken == "" {
		return "", errors.New("获取失败，请检查appid和secret")
	}
	return token.AccessToken, nil
}

func GetFollowList(in string) (string, error) {
	var data map[string]string
	err := json.Unmarshal([]byte(in), &data)
	if err != nil {
		return "", err
	}
	url := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + data["token"] + "&next_openid="
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取关注列表失败", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取内容失败", err)
		return "", err
	}
	//flist := gjson.Get(string(body), "data.openid").Array()
	return string(body), nil
}

//发送模板消息
func SendTemplateMsg(in string) (string, error) {
	var data map[string]string
	err := json.Unmarshal([]byte(in), &data)
	if err != nil {
		return "", err
	}
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + data["token"]

	reqbody := "{\"touser\":\"" + data["touser"] + "\", \"template_id\":\"" + data["template_id"] + "\", \"url\":\"" + data["url"] + "\", \"data\": " + data["data"] + "}"

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
func main() {
	//ret, err := sendTodo("test", "remark")
	//fmt.Println(ret)
	//fmt.Println(err)

	/**
	j := `{
		"appid":"wxc5c96be3cb9b4652",
		"secret":"0926be326cb0ff53d39c52114831c1bc"
	}
	`
	token, err := AccessToken(j)
	fmt.Println(token)
	fmt.Println(err)

	li, e := GetFollowList(token)
	fmt.Println(li)
	fmt.Println(e)

	//now := time.Now()
	//reqdata := "{\"content\":{\"value\":\"send by golang\", \"color\":\"#0000CD\"}, \"note\":{\"value\":\"ok\"}, \"time\":{\"value\":\"" + now.Format("2006-01-02 15:04:05") + "\"}}"

	j = `{
		"token":"` + token + `",
		"touser":"o79o16zvw02Xvlwn_bx6D91bFS_M",
		"template_id":"DVmeujVmcP5R4_a_aoreiRdfh5RMfELfiQKUWXQRemA",
		"url":"http://81.70.218.13/",
		"data":"{\"content\":{\"value\":\"send by golang\", \"color\":\"#0000CD\"}, \"note\":{\"value\":\"ok\"}, \"time\":{\"value\":\" 11:50:00 \"}}"
	}`
	fmt.Println(j)
	ret, e1 := SendTemplateMsg(j)
	fmt.Println(ret)
	fmt.Println(e1)
	**/
	p := `
		{"token":"39_r7wwxzSmtf0pfkTDxKVxM-X655nRZqJhuOE0_550GaEjebElXhKRP6R0Jmjzo_cWHyGZ7j2NpGXwScsJHW1clkyK4DkPY2ZTv03gTMDhFL_a0b4aShgr6868VYAvPC596f6YwdzFwo_3x52CADYiADAGJE"}
	`
	s, e := GetFollowList(p)
	fmt.Println(s)
	fmt.Println(e)
}
