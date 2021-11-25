package util

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
)

// GetBranchInsertSql 获取批量添加数据sql语句
func GetBranchInsertSql(objs []interface{}, tableName string) string {
	if len(objs) == 0 {
		return ""
	}
	fieldName := ""
	var valueTypeList []string
	fieldNum := reflect.TypeOf(objs[0]).NumField()
	fieldT := reflect.TypeOf(objs[0])
	for a := 0; a < fieldNum; a++ {
		name := fieldT.Field(a).Tag.Get("json")
		// 添加字段名
		if a == fieldNum-1 {
			fieldName += fmt.Sprintf("`%s`", name)
		} else {
			fieldName += fmt.Sprintf("`%s`,", name)
		}
		// 获取字段类型
		if fieldT.Field(a).Type.Name() == "string" {
			valueTypeList = append(valueTypeList, "string")
		} else if strings.Index(fieldT.Field(a).Type.Name(), "uint") != -1 {
			valueTypeList = append(valueTypeList, "uint")
		} else if strings.Index(fieldT.Field(a).Type.Name(), "int") != -1 {
			valueTypeList = append(valueTypeList, "int")
		}
	}
	var valueList []string
	for _, obj := range objs {
		objV := reflect.ValueOf(obj)
		v := "("
		for index, i := range valueTypeList {
			if index == fieldNum-1 {
				v += GetFormatField(objV, index, i, "")
			} else {
				v += GetFormatField(objV, index, i, ",")
			}
		}
		v += ")"
		valueList = append(valueList, v)
	}
	insertSql := fmt.Sprintf("insert into `%s` (%s) values %s", tableName, fieldName, strings.Join(valueList, ",")+";")
	return insertSql
}

// GetFormatField 获取字段类型值转为字符串
func GetFormatField(objV reflect.Value, index int, t string, sep string) string {
	v := ""
	if t == "string" {
		v += fmt.Sprintf("'%s'%s", objV.Field(index).String(), sep)
	} else if t == "uint" {
		v += fmt.Sprintf("%d%s", objV.Field(index).Uint(), sep)
	} else if t == "int" {
		v += fmt.Sprintf("%d%s", objV.Field(index).Int(), sep)
	}
	return v

}

// GetColumnName 获取字段名
func GetColumnName(jsonName string) string {
	for _, name := range strings.Split(jsonName, ";") {
		if strings.Index(name, "column") == -1 {
			continue
		}
		return strings.Replace(name, "column:", "", 1)
	}
	return ""
}


// BatchCreateModelsByPage 分页批量插入
func BatchCreateModelsByPage(tx *gorm.DB, dataList []interface{}, tableName string) (err error) {
	if len(dataList) == 0 {
		return
	}
	// 如果超过一百条, 则分批插入
	size := 100
	page := len(dataList) / size
	if len(dataList)%size != 0 {
		page += 1
	}
	for a := 1; a <= page; a++ {
		var bills = make([]interface{}, 0)
		if a == page {
			bills = dataList[(a-1)*size:]
		} else {
			bills = dataList[(a-1)*size : a*size]
		}
		sql := GetBranchInsertSql(bills, tableName)
		if err = tx.Exec(sql).Error; err != nil {
			fmt.Println(fmt.Sprintf("batch create data error: %v, sql: %s, tableName: %s", err, sql, tableName))
			return
		}
	}
	return
}