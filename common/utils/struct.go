package utils

import (
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

// a 数据源 b 新的结构体 c 指针
func StAtoB(a, b, c interface{}) error {
	//第一步,先将结构体转化为map方便后续遍历
	amap := Struct2Map(a)
	bmap := Struct2Map(b)
	for k1, v1 := range amap {
		if _, ok := bmap[k1]; ok {
			typea := reflect.TypeOf(v1)
			typeb := reflect.TypeOf(bmap[k1])
			if typea == typeb {
				bmap[k1] = v1
			}
		}
	}

	if err := mapstructure.Decode(bmap, c); err != nil {
		fmt.Println("the err when unmarshal mapstructure is:", err)
		return err
	}
	return nil
}

func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func MapToStatus(a interface{}, b map[string]interface{}, c interface{}) error {
	amap := Struct2Map(a)
	for k, v := range amap {
		typeofA := reflect.TypeOf(v)
		sfield, _ := reflect.TypeOf(a).FieldByName(k)
		tag := sfield.Tag.Get("mapkey")
		if len(tag) > 0 {
			if _, ok := b[tag]; ok {
				typeofB := reflect.TypeOf(b[tag])
				if typeofA == typeofB {
					amap[k] = b[tag]
				}
			}
		}
	}
	if err := mapstructure.Decode(amap, c); err != nil {
		fmt.Println("the err when unmarshal mapstructure is:", err)
		return err
	}
	return nil
}
