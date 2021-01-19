package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
)

type Req struct {
	File string `json:"file"`
	Data string `json:"data"`
}

//go build --buildmode=plugin -o test_v1.0.so plugins/test.go
func Execute(in string) (string, error) {
	p := Req{}
	err := json.Unmarshal([]byte(in), &p)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var mapResult map[string]interface{}
	err = json.Unmarshal([]byte(in), &mapResult)
	if err != nil {
		return "", err
	}
	fmt.Println(mapResult)

	// 解析指定文件生成模板对象
	tem, err := template.ParseFiles(p.File)
	if err != nil {
		return "", err
	}
	// 利用给定数据渲染模板，并将结果写入w
	var buf bytes.Buffer
	if err := tem.Execute(&buf, mapResult); err != nil {
		return "", err
	}
	fmt.Println(buf.String()) // 渲染后的字符串 // <p> hello Tom </p>

	return buf.String(), nil
}
func main() {
	str := `
	{
		"file": "app/admin/apis/jsgo/plugins/tmpl.html",
		"data": "123456"
	}
	`
	out, err := Execute(str)
	println(out)
	fmt.Println(err)

}
