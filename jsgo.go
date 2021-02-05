package jsgo

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
)

var ii = V8Init()
var vm = CreateV8VM()
var vmcount = 1
var vms []VM
var gid, _ = snowflake.NewNode(1)
var init_vm = 0
var init_top = 0
var resetting []int

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

func LoadAllJsFile(cvm VM, pathname string) string {
	rmsg := ""
	//load first v8_top.js
	jf := pathname + "/v8_top.js"
	_, msg := cvm.Load(jf)
	rmsg = rmsg + jf + ":" + msg + ";"

	rd, _ := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			fmt.Printf("[%s]\n", pathname+"\\"+fi.Name())
			LoadAllJsFile(cvm, pathname+fi.Name()+"\\")
		} else {
			fmt.Println(fi.Name())
			if fi.Name() == "v8_top.js" {
				continue
			}
			if strings.Contains(fi.Name(), ".js") {
				jf := pathname + "/" + fi.Name()
				_, msg := cvm.Load(jf)
				rmsg = rmsg + jf + ":" + msg + ";"

			}
		}
	}
	return rmsg
}

/**
load js files from ./jsplugin
**/
func Init(initcount int) {
	if init_vm == 0 {
		if initcount < 1 {
			initcount = 1
		}
		vmcount = initcount
		vms = make([]VM, vmcount)
		resetting = make([]int, vmcount)
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < vmcount; i++ {
			vms[i] = CreateV8VM()
			resetting[i] = 0
		}
		init_vm = 1
	}
	if init_top == 0 {
		fmt.Println("------Start init js plugin------")
		for v := range vms {
			msg := LoadAllJsFile(vms[v], "./jsplugin")
			fmt.Println(msg)
		}
		msg := LoadAllJsFile(vm, "./jsplugin")
		fmt.Println(msg)
		init_top = 1
		fmt.Println("------End init js plugin------")

	}

}

func DoInVm(js string, param string) (int, string, error) {
	//随机vm
	x := GetVmNum(js, vmcount)
	if resetting[x] == 1 {
		return x, "", nil
	}
	id := gid.Generate().String()
	str, err := doExeJsIsolateInVm(vms[x], id, js, param)
	return x, str, err
}
func ResetVm(id int) error {
	resetting[id] = 1
	vms[id].Reset()
	resetting[id] = 0
	return nil
}

func GetVmNum(js string, count int) int {
	x := 0
	if strings.HasPrefix(js, "//vm[") {
		s := strings.Split(js, "[")
		if len(s) > 1 {
			se := strings.Split(s[1], "]")
			if len(se) > 1 {
				sn := se[0]
				if sn == "random" {
					return rand.Intn(count)
				}
				x, _ = strconv.Atoi(sn)
				if x > count-1 {
					x = count - 1
				}
			}
		}
	}
	return x
}
func DoHandleNoPermission(param map[string]string) (string, error) {
	return "to do", nil
}
