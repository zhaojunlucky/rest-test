package core

import (
	"encoding/json"
	"reflect"
)

func ConvertObj(obj map[string]any) map[string]any {
	copyObj := make(map[string]any)

	for k, v := range obj {
		if v == nil {
			copyObj[k] = v
			continue
		}
		valType := reflect.TypeOf(v)
		switch valType.Kind() {
		case reflect.Array, reflect.Slice:
			copyObj[k] = ConvertArr(v.([]any))
			continue
		case reflect.Map:
			copyObj[k] = ConvertObj(v.(map[string]any))
			continue
		default:
			break
		}
		if valType.String() == "json.Number" {
			jn := v.(json.Number)
			if n, err := jn.Int64(); err == nil {
				copyObj[k] = n
			} else if f, err := jn.Float64(); err == nil {
				copyObj[k] = f
			} else {
				copyObj[k] = jn.String()
			}
		} else {
			copyObj[k] = v
		}
	}
	return copyObj
}

func ConvertArr(arr []any) []any {
	var copyObj []any

	for _, v := range arr {
		if v == nil {
			copyObj = append(copyObj, v)
			continue
		}
		valType := reflect.TypeOf(v)
		switch valType.Kind() {
		case reflect.Array, reflect.Slice:
			copyObj = append(copyObj, ConvertArr(v.([]any)))
			continue
		case reflect.Map:
			copyObj = append(copyObj, ConvertObj(v.(map[string]any)))
			continue
		default:
			break
		}
		if valType.String() == "json.Number" {
			jn := v.(json.Number)
			if n, err := jn.Int64(); err == nil {
				copyObj = append(copyObj, n)
			} else if f, err := jn.Float64(); err == nil {
				copyObj = append(copyObj, f)
			} else {
				copyObj = append(copyObj, jn.String())
			}
		} else {
			copyObj = append(copyObj, v)
		}
	}
	return copyObj
}
