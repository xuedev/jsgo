package jsgo

import (
	"fmt"
	"testing"
)

func TestVm(t *testing.T) {
	Init()
	vm := CreateV8VM()
	//r := vm.Load("/run/media/xuegx/hook/go/workspace/v8go/t.js")
	//r, _ := vm.Eval("var ret = '';try{var top = callp()}catch(e){ret = e+''};ret")

	b := vm.Load("top.js")
	println(b)

	js := `var sql = squel.select()
			.from("table", "t1")
			.field("t1.id")
			.field("t2.name")
			.left_join("table2", "t2", "t1.id = t2.id")
			.group("t1.id")
			.where("t2.name <> 'Mark'")
			.where("t2.name <> 'John'")
			.toString();sql`

	r, err := vm.Eval(js)
	//r, err := vm.Eval("go.v8();")
	println(r)
	println(fmt.Sprintf("%s", err))

	//defer func() {
	//	C.free(unsafe.Pointer(r))
	//}()
}
