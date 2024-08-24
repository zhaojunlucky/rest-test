package core

func ConvertObj(obj map[string]any) map[string]any {
	copyObj := make(map[string]any)

	for k, v := range obj {
		//valType := reflect.ValueOf(v)
		//switch valType.(type) {
		//case reflect.Value:
		//
		//}
		copyObj[k] = v
	}
	return copyObj
}

func ConvertArr(arr []any) []any {
	return nil
}
