package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_ "net/http/pprof"

	"github.com/xuedev/jsgo"
)

// go env-w GOPROXY=https://goproxy.cn,direct
func call(path string) (string, error) {
	//fmt.Println("jsgo")
	// init 10 vm
	f, err := os.Open("service/" + path + ".js")
	defer f.Close()
	if err != nil {
		return "", err
	}
	data, err2 := ioutil.ReadAll(f)
	if err2 != nil {
		return "", err2
	}
	js := string(data)

	param := `
		{
			"data":"jsgo"
		}
	`
	vm, ret, err := jsgo.DoInVm(js, param)
	param = ""
	fmt.Println(vm)
	fmt.Println(ret)
	fmt.Println(err)
	defer func() {
		data = nil
		js = ""
		ret = ""
	}()
	return ret, err
}
func main() {
	jsgo.Init(5)

	//处理路由为 / 的方法
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		ret, err := call(path)
		if err == nil {
			fmt.Fprintln(w, "Hello World", ret)
		} else {
			fmt.Fprintln(w, "Error", err)
		}

	})
	fmt.Println("service in 8000")
	//监听3000端口
	http.ListenAndServe(":8000", nil)
}
