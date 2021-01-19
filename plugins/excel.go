package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	excelize "github.com/360EntSecGroup-Skylar/excelize"
)

type Req struct {
	File   string
	Header string
	Data   []map[string]string
}

var cols = [26]string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

//go build --buildmode=plugin -o test_v1.0.so plugins/test.go
func Json2Excel(in string) (string, error) {
	var p Req
	err := json.Unmarshal([]byte(in), &p)
	if err != nil {
		//fmt.Println(err)
		return "", err
	}
	f := excelize.NewFile()
	// 创建一个工作表
	index := f.NewSheet("Sheet1")
	//创建表头
	header := strings.Split(p.Header, "#")
	for col, t := range header {
		//fmt.Println(t)
		kv := strings.Split(t, ",")
		pos := cols[col] + "1"
		//fmt.Println(pos)
		f.SetCellValue("Sheet1", pos, kv[1])
	}

	for row, d := range p.Data {
		//fmt.Println(d["id"])
		for col, t := range header {
			//fmt.Println(t)
			kv := strings.Split(t, ",")
			pos := cols[col] + strconv.Itoa(row+2)
			//fmt.Println(pos)
			f.SetCellValue("Sheet1", pos, d[kv[0]])
		}

	}

	// 设置工作簿的默认工作表
	f.SetActiveSheet(index)
	// 根据指定路径保存文件
	if err := f.SaveAs(p.File); err != nil {
		return "", err
	}

	return p.File, nil
}
func main() {
	data := `[{
		"addr": "",
		"code": "QZX_C1_027085e98c8580dca4ce306046f58c50",
		"create_by": "",
		"created_at": "2020-12-09 17:22:51",
		"deleted_at": "",
		"dtlhdj": "3",
		"gps": "",
		"height": "720",
		"id": "588",
		"isclean": "",
		"jgewm": "https://biz.jnfda.gov.cn:9777/checkout/q.jsp?id=08826844279011e8af920050568f29b8",
		"name": "",
		"remark": "",
		"rt": "0.2833",
		"runt": "运行中",
		"shop": "济南历下奥体中心三屿亭火锅店",
		"sndzhdj": "未评定",
		"type": "device",
		"update_by": "",
		"updated_at": "2020-12-17 14:21:13",
		"width": "1280",
		"x": "0",
		"xkzh": "JY23701020073753",
		"xydm": "",
		"xzqh": "",
		"y": "0"
	}, {
		"addr": "",
		"code": "CM201-1-CW_2a03f66a71805ab620d404ce093314af",
		"create_by": "",
		"created_at": "2020-12-07 04:01:50",
		"deleted_at": "",
		"dtlhdj": "2",
		"gps": "",
		"height": "340",
		"id": "537",
		"isclean": "",
		"jgewm": "https://biz.jnfda.gov.cn:9777/checkout/q.jsp?id=04b8ada798a511e7bec30050568f29b8",
		"name": "",
		"remark": "",
		"rt": "0.2833",
		"runt": "运行中",
		"shop": "济南历下奥体中心锅说火锅店",
		"sndzhdj": "B",
		"type": "device",
		"update_by": "",
		"updated_at": "2020-12-17 14:21:13",
		"width": "600",
		"x": "20",
		"xkzh": "JY23701020055671",
		"xydm": "",
		"xzqh": "",
		"y": "20"
	}]`
	in := `
		{
			"file":"test.xls",
			"header":"shop,店名#xkzh,许可证号",
			"data":` + data + `
		}
	`
	r, e := Json2Excel(in)
	fmt.Println(r)
	fmt.Println(e)
}
