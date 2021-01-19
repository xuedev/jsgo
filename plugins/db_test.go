package main

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var dbs = make(map[string]*gorm.DB)

type DbParam struct {
	Db  string `json:"db"`
	Url string `json:"url"`
	Sql string `json:"sql"`
}

// 正则过滤sql注入的方法
// 参数 : 要匹配的语句
func FilteredSQLInject(to_match_str string) bool {
	//过滤 ‘
	//ORACLE 注解 --  /**/
	//关键字过滤 update ,delete
	// 正则的字符串, 不能用 " " 因为" "里面的内容会转义
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|insert|char|chr|into|substr|ascii|declare|exec|count|master|into|execute)\b)`
	re, err := regexp.Compile(str)
	if err != nil {
		panic(err.Error())
		return false
	}
	return re.MatchString(to_match_str)
}
func DbQuery(dbparam string) (string, error) {
	p := DbParam{}
	err := json.Unmarshal([]byte(dbparam), &p)
	if err != nil {
		return "", err
	}
	result := make(map[int]map[string]string)
	db := dbs[p.Db]
	if db == nil {
		return "", errors.New("db[" + p.Db + "] not inited")
	}
	if !FilteredSQLInject(p.Sql) {
		return "", errors.New("SQL is illegal")
	}
	rows, err := db.Raw(p.Sql).Rows() // (*sql.Rows, error)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	i := 0
	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err = rows.Scan(columnPointers...); err != nil {
			return "", err
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]string)
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = Strval(*val)
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		//fmt.Println(m)
		result[i] = m
		i = i + 1
	}

	out := make([]map[string]string, i)
	for j := 0; j < i; j++ {
		out[j] = result[j]
	}

	str, err := json.Marshal(out)
	if err != nil {
		return "", nil
	}
	return string(str), nil

}

func DbInit(dbp string) (string, error) {
	p := DbParam{}
	err := json.Unmarshal([]byte(dbp), &p)
	if err != nil {
		return "error parse param", err
	}

	db := dbs[p.Db]
	if db == nil {
		ndb, err := gorm.Open("mysql", p.Url)
		if err != nil {
			return "error connect to db", err
		}
		dbs[p.Db] = ndb
	}
	//err := db.DB().Ping()
	//if err != nil {
	//	return DbInit(dbkey, dburl)
	//}
	return "ok", nil
}

func TestDB(t *testing.T) {
	dbi := `
		{
			"db": "test1",
			"url": "root:123456a?@tcp(127.0.0.1:3306)/jsgo?charset=utf8&parseTime=True&loc=Local&timeout=1000s"
		}
	`
	DbInit(dbi)
	param := `{
				"db":"test1",
				"sql": "select * from sys_user"
			}`
	data, err := DbQuery(param)
	if err != nil {
		t.Log(err)
	} else {
		t.Log(data)
	}
}

// Strval 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
