package gxapi

func getTopVar(key string) (string, error) {
	//id := gid.Generate().String()
	value, err := vm.Eval(key)
	if err != nil {
		return "", err
	}
	return value, nil
}
func doExeTopJs(js string) (string, error) {
	//id := gid.Generate().String()
	value, err := vm.Eval(js)
	if err != nil {
		return "", err
	}
	println("eval top ret:" + value)
	return value, nil
}
func doExeJsIsolate(id string, js string, param string) (string, error) {
	jsFunc := `
		var p_` + id + ` = eval(` + param + `);
		function func_` + id + `(param){
			//console.log('param:'+param);
			try{
				` + js + `;
			}catch(e){
				return e+'';
			}
		};
		var ret_` + id + ` = func_` + id + `(p_` + id + `);
	`
	//println(jsFunc)
	_, err := vm.Eval(jsFunc)
	if err != nil {
		return "", err
	}

	value, err := vm.Eval("ret_" + id)
	if err != nil {
		return "", err
	}

	clearJs := `
		func_` + id + ` = null;
		ret_` + id + ` = null;
	`
	//println(clearJs)
	_, err = vm.Eval(clearJs)
	if err != nil {
		return "", err
	}

	return value, err
}

func doExeJsIsolateInVm(cvm VM, id string, js string, param string) (string, error) {
	jsFunc := `
		var p_` + id + ` = eval(` + param + `);
		function func_` + id + `(param){
			//console.log('param:'+param);
			try{
				` + js + `;
			}catch(e){
				return e+'';
			}
		};
		var ret_` + id + ` = func_` + id + `(p_` + id + `);
	`
	//println(jsFunc)
	_, err := cvm.Eval(jsFunc)
	if err != nil {
		return "", err
	}

	value, err := cvm.Eval("ret_" + id)
	if err != nil {
		return "", err
	}

	clearJs := `
		func_` + id + ` = null;
		ret_` + id + ` = null;
	`
	//println(clearJs)
	_, err = cvm.Eval(clearJs)
	if err != nil {
		return "", err
	}

	return value, err
}
