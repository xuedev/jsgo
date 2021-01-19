package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/skip2/go-qrcode"
)

//go build --buildmode=plugin -o test_v1.0.so plugins/test.go
//执行命令函数: 不能感知命令的执行信息，只返回是否执行成功
//commandName 命名名称，如cat，ls，git等
//params 命令参数，如ls -l的-l，git log 的log等
func Shell(in string) (string, error) {
	var data map[string]string
	err := json.Unmarshal([]byte(in), &data)
	if err != nil {
		return "", err
	}
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", data["cmd"])

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err = cmd.Run()

	return out.String(), err
}

func ReadFile(in string) (string, error) {
	var data map[string]string
	err := json.Unmarshal([]byte(in), &data)
	if err != nil {
		return "", err
	}
	f, err := os.Open(data["file"])
	if err != nil {
		return "", err
	}
	dd, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(dd), nil
}

func QR(data string, file string) error {
	return qrcode.WriteFile(data, qrcode.Medium, 512, file)
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GenQr(in string) (string, error) {
	var data map[string]string
	err := json.Unmarshal([]byte(in), &data)
	if err != nil {
		return "", err
	}
	f, err := os.Open(data["file"])
	if err != nil && os.IsNotExist(err) {
		err = QR(data["data"], data["file"])
		if err != nil {
			return "", err
		}
		return data["file"], nil
	}
	//fmt.Printf("file exist!\n")
	defer f.Close()
	return data["file"], nil
}
func main() {
	s := `
	{
		"data": "https://biz.jnfda.gov.cn:9777/checkout/q.jsp?id=2ff1e113dbdb11ea9f83005056bf3be1",
		"file": "./static/JY23701020190003.png"
	  }
	`
	r, e := GenQr(s)
	fmt.Println(r)
	fmt.Println(e)

	/**
	s := `
	{
		"file": "./jsplugin/db.js"
	  }
	`
	r, e := ReadFile(s)
	fmt.Println(r)
	fmt.Println(e)

	s := `
	  {
		  "cmd":"ls -l"
	  }
	`
	r, e := Shell(s)
	fmt.Println(r)
	fmt.Println(e)**/
}
