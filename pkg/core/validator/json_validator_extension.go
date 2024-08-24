package validator

import (
	"fmt"
	"github.com/PaesslerAG/gval"
	"reflect"
	"strconv"
	"strings"
)

var jsonValidatorExtension = map[string]func(args ...any) (any, error){
	"string":  stringFunc,
	"int":     intFunc,
	"float":   floatFunc,
	"bool":    boolFunc,
	"len":     lenFunc,
	"contain": containFunc,
}

func containFunc(args ...any) (any, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("contain function only support 2 parameter")
	}
	arg1 := reflect.ValueOf(args[0])
	arg2 := reflect.ValueOf(args[1])

	if arg1.Type().Kind() == reflect.String || arg2.Type().Kind() == reflect.String {
		return strings.Contains(arg2.String(), arg1.String()), nil
	}

	switch arg1.Type().Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		return nil, fmt.Errorf("array, slice and map are not supported in first arg of contain function")
	default:
		break
	}
	switch arg2.Type().Kind() {
	case reflect.Array, reflect.Slice:
		for i := range arg2.Len() {
			if arg1.Equal(arg2.Index(i)) {
				return true, nil
			}
		}
		return false, nil
	case reflect.Map:
		val := arg2.MapIndex(arg1)
		return val.IsValid(), nil
	default:
		return nil, fmt.Errorf("only array, slice and map are supported in second arg of contain function")
	}

}

func lenFunc(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("len function only support 1 parameter")
	}
	argType := reflect.ValueOf(args[0])

	switch argType.Type().Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return argType.Len(), nil
	case reflect.String:
		return len(args[0].(string)), nil
	default:
		return nil, fmt.Errorf("unsupported type: %v", argType)
	}
}

func intFunc(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("int function only support 1 parameter")
	}
	argType := reflect.ValueOf(args[0])

	if argType.CanInt() {
		return argType.Int(), nil
	} else if argType.CanFloat() {
		return int(argType.Float()), nil
	} else if argType.Kind() == reflect.String {
		return strconv.Atoi(argType.String())
	}
	return nil, fmt.Errorf("unsupported type: %v", argType)
}

func floatFunc(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("float function only support 1 parameter")
	}
	argType := reflect.ValueOf(args[0])

	if argType.CanFloat() {
		return argType.Float(), nil
	} else if argType.CanInt() {
		return float64(argType.Int()), nil
	} else if argType.Kind() == reflect.String {
		return strconv.ParseFloat(argType.String(), 64)
	}
	return nil, fmt.Errorf("unsupported type: %v", argType)
}

func boolFunc(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("bool function only support 1 parameter")
	}
	argType := reflect.ValueOf(args[0])

	if argType.CanInt() {
		return argType.Int() > 0, nil
	} else if argType.CanFloat() {
		return argType.Float() > 0, nil
	} else if argType.CanUint() {
		return argType.Uint() > 0, nil
	}

	switch argType.Type().Kind() {
	case reflect.Bool:
		return argType.Bool(), nil
	case reflect.Array, reflect.Slice, reflect.Map:
		return argType.Len() > 0, nil
	case reflect.String:
		return argType.String() == "true", nil
	default:

		return nil, fmt.Errorf("unsupported type: %v", argType)
	}
}

func stringFunc(args ...any) (any, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("string function only support 1 parameter")
	}
	argType := reflect.ValueOf(args[0])

	if argType.CanInt() {
		return fmt.Sprintf("%d", argType.Int()), nil
	} else if argType.CanFloat() {
		return fmt.Sprintf("%f", argType.Float()), nil
	} else if argType.CanUint() {
		return fmt.Sprintf("%d", argType.Uint()), nil
	}

	switch argType.Type().Kind() {
	case reflect.Bool:
		return fmt.Sprintf("%t", argType.Bool()), nil
	case reflect.String:
		return argType.String(), nil
	default:
		return nil, fmt.Errorf("unsupported type: %v", argType)
	}
}

func GetExtensionLanguage() []gval.Language {
	var langs []gval.Language

	for k, v := range jsonValidatorExtension {
		langs = append(langs, gval.Function(k, v))
	}
	return langs
}
