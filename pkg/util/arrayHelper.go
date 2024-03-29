package util

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
)

// ArrayColumn 相当于php中的array_column
func ArrayColumn(array interface{},key string)(result map[string]interface{},err error){
	result = make(map[string]interface{})
	t := reflect.TypeOf(array)
	v := reflect.ValueOf(array)
	if t.Kind() != reflect.Slice {
		return nil, errors.New("array type not slice")
	}
	if v.Len() == 0{
		return nil, errors.New("array len is zero")
	}

	for i := 0; i < v.Len(); i++{
		indexv := v.Index(i)
		fmt.Printf("v1 type:%T\n", indexv)
		fmt.Println(indexv)
		if indexv.Type().Kind() != reflect.Struct{
			return nil, errors.New("element type not struct")
		}
		mapKeyInterface := indexv.FieldByName(key)
		if mapKeyInterface.Kind() == reflect.Invalid {
			return nil,errors.New("key not exist")
		}
		mapKeyString, err := interfaceToString(mapKeyInterface.Interface())
		if err != nil{
			return nil, err
		}
		result[mapKeyString] = indexv.Interface()
	}
	return result,err
}


func interfaceToString(v interface{})(result string, err error){
	switch reflect.TypeOf(v).Kind(){
	case reflect.Int64,reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32:
		result = fmt.Sprintf("%v",v)
	case reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
		result = fmt.Sprintf("%v",v)
	case reflect.String:
		result =  v.(string)
	default:
		err =  errors.New("can't transition to string")
	}
	return result,err
}

// ArrayDiff 对比两个数组的差异
func ArrayDiff(array1 []interface{}, othersParams ...[]interface{}) ([]interface{}, error) {
	if len(array1) == 0 {
		return []interface{}{}, nil
	}
	if len(array1) > 0 && len(othersParams) == 0 {
		return array1, nil
	}
	var tmp = make(map[interface{}]int, len(array1))
	count := 1
	for _, v := range array1 {
		tmp[v] = count
		count++
	}
	var res = make([]interface{}, 0, len(tmp))
	for _, param := range othersParams {
		for _, arg := range param {
			if tmp[arg] == 0 {
				res = append(res, arg)
			}
		}
	}
	return res, nil
}

func RemoveDuplicate(v []int) []int {
	// 排序
	sort.Slice(v, func(i, j int) bool { return v[i] <= v[j] })
	//为了性能需要尽可能的减小拷贝，最悲观的情况每个元素只移动一次。
	toIndex := 0
	p := 0

	for i, _ := range v {
		// 为了实际去重结构时减小内存拷贝
		c := &v[i]

		if p == *c && i != 0 {
			// 重复内容，跳过
			continue
		}

		if i != toIndex {
			// 需要移动当前元素
			v[toIndex] = *c
		}

		toIndex++
		p = *c
	}

	return v[:toIndex]
}

// InArrayHelper 目前只支持int int64 string
func InArrayHelper(needle interface{}, hyStack interface{}) bool {
	switch key := needle.(type) {
	case string:
		for _, item := range hyStack.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range hyStack.([]int) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range hyStack.([]int64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}
	return false
}

func InArray(needle interface{}, hystack interface{}) bool {
	switch key := needle.(type) {
	case string:
		for _, item := range hystack.([]string) {
			if key == item {
				return true
			}
		}
	case int:
		for _, item := range hystack.([]int) {
			if key == item {
				return true
			}
		}
	case int64:
		for _, item := range hystack.([]int64) {
			if key == item {
				return true
			}
		}
	default:
		return false
	}
	return false
}