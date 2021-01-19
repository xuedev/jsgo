package main

import (
	"encoding/json"
	"database/sql"
	"fmt"
	"errors"
	_ "github.com/go-sql-driver/mysql"
)
 
/**
{
	"js": " var d = callp('mysql','1.0','Init','root:123456a?@tcp(144.34.157.112:3306)/xassistants?charset=utf8'); var ret = callp('mysql','1.0','Query','show tables;') ; return ret;"
}
**/
var DB *sql.DB = nil
 
func ping() error{
	if err := DB.Ping(); err != nil {
		return err
	}
	return nil
}
//path = "root:password@tcp(127.0.0.1:3306)/mydb?charset=utf8"
func Init(path string) (string,error) { //连接到MySQL
    //root = 用户名
    //password = 密码
    //mydb = 数据库名称
 
	fmt.Println(path)
	DB, _ = sql.Open("mysql", path)
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		return "",err
	}
	return "",nil
}
 
func Exec(SQL string) (string,error){
	if ping() == nil {
		ret, err := DB.Exec(SQL) //增、删、改就靠这一条命令就够了，很简单
		//insID, _ := ret.LastInsertId()
		//fmt.Println(insID)
		dataType , _ := json.Marshal(ret)
		str := string(dataType)
		return str,err
	}
	return "",errors.New("exe sql error");
}
 
func Query(SQL string) (string, error){ //通用查询
	if ping() != nil { //连接数据库
		return "", errors.New("db not connected")
	}
	rows, err := DB.Query(SQL) //执行SQL语句，比如select * from users
	if err != nil {
		panic(err)
	}
	columns, _ := rows.Columns()            //获取列的信息
	count := len(columns)                   //列的数量
	var values = make([]interface{}, count) //创建一个与列的数量相当的空接口
	for i, _ := range values {
		var ii interface{} //为空接口分配内存
		values[i] = &ii    //取得这些内存的指针，因后继的Scan函数只接受指针
	}
	ret := make([]map[string]string, 0) //创建返回值：不定长的map类型切片
	for rows.Next() {
		err := rows.Scan(values...)  //开始读行，Scan函数只接受指针变量
		m := make(map[string]string) //用于存放1列的 [键/值] 对
		if err != nil {
			panic(err)
		}
		for i, colName := range columns {
			var raw_value = *(values[i].(*interface{})) //读出raw数据，类型为byte
			b, _ := raw_value.([]byte)
			v := string(b) //将raw数据转换成字符串
			m[colName] = v //colName是键，v是值
		}
		ret = append(ret, m) //将单行所有列的键值对附加在总的返回值上（以行为单位）
	}
    
    defer rows.Close()
 
	if len(ret) != 0 {

		dataType , _ := json.Marshal(ret)
		dataString := string(dataType)
		return dataString, nil
	}
	return "", errors.New("query error")
}
