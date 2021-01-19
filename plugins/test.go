package main

import (
	"encoding/json"
	"fmt"
)

//go build --buildmode=plugin -o test_v1.0.so plugins/test.go
func Call(in string) (string, error) {
	fmt.Println("call test.call v2 ,param:", in)
	var p map[string]string
	err := json.Unmarshal([]byte(in), &p)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return in, nil
}
func main() {
	println("test")
}
