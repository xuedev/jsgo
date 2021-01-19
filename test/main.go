package main

import (
	"fmt"
	"jsgo"
)

// go env-w GOPROXY=https://goproxy.cn,direct
func main() {
	fmt.Println("jsgo")
	// init 10 vm
	jsgo.Init(3)
	js := `//vm[random]  run in a random vm
		   console.log(param.data)
		`
	param := `
		{
			"data":"jsgo"
		}
	`
	vm, ret, err := jsgo.DoInVm(js, param)
	fmt.Println(vm)
	fmt.Println(ret)
	fmt.Println(err)
}
