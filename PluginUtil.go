package gxapi

import (
	"plugin"
)

var plugins map[string]*plugin.Plugin = make(map[string]*plugin.Plugin)

func callPlugin(pluginName string, version string, funcName string, param string) (string, error) {
	pkey := pluginName + "_v" + version + ".so"
	var pdll *plugin.Plugin
	if odll, ok := plugins[pkey]; ok {
		pdll = odll
		//fmt.Println("use cached:" + pkey)
	} else {
		//打开动态库
		ndll, err := plugin.Open(pkey)
		if err != nil {
			return "Can not find plugin", err
		}
		plugins[pkey] = ndll
		pdll = ndll
	}

	//获取动态库方法
	funcPrint, err := pdll.Lookup(funcName)
	if err != nil {
		//...
		return "Can not find func", err
	}
	//动态库方法调用
	ret, err := funcPrint.(func(string) (string, error))(param)

	return ret, err
}
