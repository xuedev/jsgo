package main

import (
	"strings"

	xj "github.com/basgys/goxml2json"
)

//go build --buildmode=plugin -o test_v1.0.so plugins/test.go
func Xml2Json(in string) (string, error) {
	// xml is an io.Reader
	xml := strings.NewReader(in)
	json, err := xj.Convert(xml)
	if err != nil {
		return "", err
	}

	return json.String(), nil
	// {"hello": "world"}
}
func main() {
	println("conv")
}
