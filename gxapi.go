package gxapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-admin/app/admin/models"
	mycasbin "go-admin/pkg/casbin"
	"go-admin/pkg/jwtauth"
	_ "go-admin/pkg/jwtauth"
	"go-admin/tools"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var ii = V8Init()
var vm = CreateV8VM()
var vmcount = 1
var vms []VM
var gid, _ = snowflake.NewNode(1)
var init_vm = 0
var init_top = 0

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
func Init() {
	if init_vm == 0 {
		if vmc := viper.GetString("settings.application.vms"); len(vmc) > 0 {
			vmcount, _ = strconv.Atoi(vmc)
			fmt.Println("-----------------set vmcount-------------------")
			fmt.Println("-----------------" + vmc + "-------------------")
		}
		vms = make([]VM, vmcount)

		rand.Seed(time.Now().UnixNano())
		for i := 0; i < vmcount; i++ {
			vms[i] = CreateV8VM()
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

func resp(c *gin.Context, code int, data string) {
	/**c.JSON(code, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})**/
	println("data:" + data)
	dmap := make(map[string]string)
	defer func() {
		dmap = nil
	}()
	err := json.Unmarshal([]byte(data), &dmap)
	if err != nil {
		c.JSON(code, gin.H{
			"code": code,
			"msg":  "success",
			"data": data + "",
		})
		return
	}

	c.JSON(code, gin.H{
		"code": code,
		"msg":  "success",
		"data": dmap,
	})
}

func checkPermission(system string, module string, api string, c *gin.Context) error {
	//check permission
	jwd, _ := c.Get("JWT_DATA")
	if jwd == nil || jwd == "" {
		return errors.New("no jwd permission")
	}
	v := jwd.(jwtauth.MapClaims)
	e := mycasbin.Casbin()
	///tools.HasError(err, "", 500)

	key0 := system + ".*"
	key := system + "." + module + "." + api
	fmt.Printf("%s [INFO] %s %s %s \r\n",
		tools.GetCurrentTimeStr(),
		c.Request.Method,
		key,
		v["rolekey"],
	)

	res, _ := e.Enforce(v["rolekey"], key0, c.Request.Method)
	if res {
		return nil
	} else {
		res, _ = e.Enforce(v["rolekey"], key, c.Request.Method)
		if res {
			return nil
		} else {
			return errors.New("no permission")
		}

	}
}

// @Summary GXApi
// @Description call gxapi
// @Tags 系统信息
// @Success 200 {object} app.Response "{"code": 200, "data": [...]}"
// @Router /api/v1/x/* [post]
/**
{
        "system":"default",
		"module":"test",
		"code":"test",
	    "param": "{\"aa\":\"123\"}"
}
**/
func DoHandle(c *gin.Context) {
	startTime := tools.GetCurrentTime()
	Init()
	body, err := ioutil.ReadAll(c.Request.Body)
	if body == nil {
		tools.HasError(err, "Parameter must not null", 400)
		return
	}
	//println("--------------")
	//println(fmt.Sprintf("%s", body))
	//println("--------------")
	dmap := make(map[string]string)
	defer func() {
		dmap = nil
	}()
	err = json.Unmarshal([]byte(body), &dmap)
	if err != nil {
		ee := fmt.Sprintf("%s", err)
		resp(c, 400, "Parameter format error:"+ee)
		return
	}

	if dmap["reloadjs"] != "" {
		init_top = 0
		Init()
		//fmt.Println(msg)
		resp(c, 200, "ok")
		return
	}

	if dmap["code"] == "" {
		resp(c, 404, "Api code must not null")
		return
	}

	var data models.XApi
	data.Code = dmap["code"]
	if dmap["system"] != "" {
		data.System = dmap["system"]
	}
	if dmap["module"] != "" {
		data.Module = dmap["module"]
	}

	//checkPermission
	if data.System != "public" {
		err = checkPermission(data.System, data.Module, data.Code, c)
		if err != nil {
			tools.HasError(err, "对不起，您没有该接口访问权限，请联系管理员", 403)
		}
	}

	result, err := data.GetEqual()
	if err != nil {
		tools.HasError(err, "Api not found", 404)
	}

	jp := dmap["param"]
	//query api script
	js := result.Script
	x, str, err := DoInVm(js, jp)
	if err != nil {
		ee := fmt.Sprintf("%s", err)
		resp(c, 500, "Script error:"+ee)
		return
	}
	endTime := tools.GetCurrentTime()
	//记录日志
	log(startTime, endTime, "vm["+strconv.Itoa(x)+"]."+data.System+"."+data.Module+"."+data.Code, string(jp), str, 200, c)

	resp(c, 200, str)
}
func DoInVm(js string, param string) (int, string, error) {
	//随机vm
	x := GetVmNum(js, vmcount)
	id := gid.Generate().String()
	str, err := doExeJsIsolateInVm(vms[x], id, js, param)
	return x, str, err
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

func DoXrawHandle(c *gin.Context) {
	startTime := tools.GetCurrentTime()
	Init()
	body, err := ioutil.ReadAll(c.Request.Body)
	if body == nil {
		tools.HasError(err, "Parameter must not null", 400)
		return
	}
	//println("--------------")
	//println(fmt.Sprintf("%s", body))
	//println("--------------")

	api := c.Param("api")
	var aa []string = make([]string,3)
	defer func(){
		aa = nil
	}()
	if len(api) < 1{
		aa[0] = c.Param("system")
		aa[1] = c.Param("module")
		aa[2] = c.Param("code")

	}else{
		aa = strings.Split(api, ".")
	}
	
	if len(aa) != 3 {
		c.String(404, "ERROR NOT FOUND")
		return
	}
	if aa[2] == "" {
		resp(c, 404, "Api code must not null")
		return
	}

	var data models.XApi
	data.Code = aa[2]
	if aa[0] != "" {
		data.System = aa[0]
	}
	if aa[1] != "" {
		data.Module = aa[1]
	}

	//checkPermission
	if data.System != "public" {
		err = checkPermission(data.System, data.Module, data.Code, c)
		if err != nil {
			tools.HasError(err, "对不起，您没有该接口访问权限，请联系管理员", 403)
		}
	}

	result, err := data.GetEqual()
	if err != nil {
		tools.HasError(err, "Api not found", 404)
	}

	pbody := string(body)
	pp := c.Request.URL.Query()
	pj, err := json.Marshal(pp)
	if err != nil {
		fmt.Println("生成json字符串错误")
	}

	hh := c.Request.Header
	hj, err := json.Marshal(hh)
	if err != nil {
		fmt.Println("生成header字符串错误")
	}

	jp := `
		{
			'query':'` + string(pj) + `',
			'header':'` + string(hj) + `',
			'body':` + "`" + pbody + "`" + `
		}
	`
	//query api script
	js := result.Script
	x, str, err := DoInVm(js, jp)
	if err != nil {
		ee := fmt.Sprintf("%s", err)
		resp(c, 500, "Script error:"+ee)
		return
	}

	endTime := tools.GetCurrentTime()
	//记录日志
	log(startTime, endTime, "vm["+strconv.Itoa(x)+"]."+data.System+"."+data.Module+"."+data.Code, string(jp), str, 200, c)

	rr := strings.Split(str, "@@@")
	if len(rr) == 1 {
		c.String(200, str)
		return
	}
	if len(rr) > 1 {
		re := regexp.MustCompile(rr[0] + "@@@")
		c.Data(http.StatusOK, rr[0], []byte(re.ReplaceAllString(str, "")))
	}

	//c.String(200, str)
}

func DoHandleNoPermission(param map[string]string) (string, error) {
	startTime := tools.GetCurrentTime()
	Init()
	api := param["api"]
	pp := param["param"]

	aa := strings.Split(api, ".")
	if len(aa) != 3 {
		return "", errors.New("api not found")
	}
	if aa[2] == "" {
		return "", errors.New("api code must not null")
	}

	var data models.XApi
	data.Code = aa[2]
	if aa[0] != "" {
		data.System = aa[0]
	}
	if aa[1] != "" {
		data.Module = aa[1]
	}

	result, err := data.GetEqual()
	if err != nil {
		return "", errors.New("api code must not null")
	}

	//query api script
	js := result.Script
	x, str, err := DoInVm(js, pp)

	if err != nil {
		ee := fmt.Sprintf("%s", err)
		return "", errors.New("script error:" + ee)
	}

	endTime := tools.GetCurrentTime()
	//记录日志
	log(startTime, endTime, "vm["+strconv.Itoa(x)+"]."+data.System+"."+data.Module+"."+data.Code, pp, str, 200, nil)

	return str, nil
}

func log(startTime time.Time, endTime time.Time, apiCode string, params string, result string, statusCode int, c *gin.Context) {
	// 执行时间
	latencyTime := endTime.Sub(startTime)
	// 请求方式
	reqMethod := "POST"
	// 请求IP
	clientIP := "localhost"
	if c != nil {
		clientIP = c.ClientIP()
	}

	sysOperLog := models.SysOperLog{}
	sysOperLog.Title = "动态接口"
	sysOperLog.BusinessType = "99"
	sysOperLog.OperatorType = "99"
	sysOperLog.OperParam = params
	sysOperLog.JsonResult = result
	sysOperLog.Params = params

	sysOperLog.OperIp = clientIP
	sysOperLog.Status = tools.IntToString(statusCode)
	if c != nil {
		sysOperLog.OperLocation = tools.GetLocation(clientIP)
		sysOperLog.OperName = tools.GetUserName(c)
		sysOperLog.CreateBy = tools.GetUserName(c)
		sysOperLog.UserAgent = c.Request.UserAgent()
	} else {
		sysOperLog.OperLocation = "self"
		sysOperLog.OperName = "self"
		sysOperLog.CreateBy = "self"
		sysOperLog.UserAgent = "gxapi"
	}

	sysOperLog.Method = apiCode
	sysOperLog.RequestMethod = reqMethod
	sysOperLog.OperUrl = apiCode

	sysOperLog.OperTime = tools.GetCurrentTime()
	sysOperLog.LatencyTime = (latencyTime).String()
	//go sysOperLog.Create()
	go func() {
		_, err := sysOperLog.Create()
		if err != nil {
			println(err)
		}
	}()

}

var sqlt = ``
