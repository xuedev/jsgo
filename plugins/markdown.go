package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Depado/bfchroma"
	bf "github.com/russross/blackfriday/v2"
)

func MarkdownToHtml(input string) string {
	//output := blackfriday.Run([]byte(input), blackfriday.WithNoExtensions())
	//unsafe := blackfriday.Run([]byte(input))
	//html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	//p := bluemonday.UGCPolicy()
	//p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	//html := p.SanitizeBytes(unsafe)
	html := bf.Run([]byte(input), bf.WithRenderer(bfchroma.NewRenderer()))
	return string(html)
}

//go build --buildmode=plugin -o test_v1.0.so plugins/test.go
func FromFile(in string) (string, error) {
	var p map[string]string
	err := json.Unmarshal([]byte(in), &p)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	f, err := os.Open(p["file"])
	if err != nil {
		return "", err
	}
	dd, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return MarkdownToHtml(string(dd)), nil
}

func FromString(in string) (string, error) {
	var p map[string]string
	err := json.Unmarshal([]byte(in), &p)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return MarkdownToHtml(p["data"]), nil
}
func main() {
	p := `
	{"file":"readme.md"}
	`
	r, e := FromFile(p)
	fmt.Println(r)
	fmt.Println(e)
}
