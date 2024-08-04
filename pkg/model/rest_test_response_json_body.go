package model

import (
	"encoding/json"
	"fmt"
	"github.com/zhaojunlucky/golib/pkg/collection"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"golang.org/x/exp/maps"
	"io"
	"math"
	"net/http"
	"reflect"
	"unicode"
)

const OR = "or"
const AND = "and"

type RestTestResponseJSONBody struct {
	RestTestRequest     *RestTestRequestDef
	Array               bool
	Length              int
	ContainsRequestJSON bool
	Validators          map[string]any
}

func (d RestTestResponseJSONBody) Validate(ctx *core.RestTestContext, resp *http.Response) (any, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyStr := string(data)
	for _, r := range bodyStr {
		if unicode.IsSpace(r) {
			continue
		}

		if r == '{' {
			if d.Array {
				return nil, fmt.Errorf("invalid json body, must be array, but got object")
			}
			break
		} else if r == '[' {
			if !d.Array {
				return nil, fmt.Errorf("invalid json body, must be object, but got array")
			}
		} else {
			return nil, fmt.Errorf("invalid json body, must start with [ or {, but got %s", string(r))
		}
	}

	if d.Array {
		var arr []any
		err = json.Unmarshal(data, &arr)
		if err != nil {
			return nil, err
		}

		if d.Length != math.MinInt && d.Length != len(arr) {
			return nil, fmt.Errorf("invalid JSON Array length: %d, expect %d", len(arr), d.Length)
		}
		return d.validate(arr)
	} else {
		var obj map[string]any
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return nil, err
		}
		return d.validate(obj)
	}
}

func (d RestTestResponseJSONBody) Parse(mapWrapper *collection.MapWrapper) error {

	if mapWrapper.Has("array") {
		err := mapWrapper.Get("array", &d.Array)
		if err != nil {
			return err
		}
	}

	if mapWrapper.Has("length") {
		err := mapWrapper.Get("length", &d.Length)
		if err != nil {
			return err
		}

	} else {
		d.Length = math.MinInt
	}

	if mapWrapper.Has("containsRequestJSON") {
		err := mapWrapper.Get("containsRequestJSON", &d.ContainsRequestJSON)
		if err != nil {
			return err
		}
	}

	if !mapWrapper.Has("validators") {
		return nil
	}

	valType, err := mapWrapper.GetType("validators")
	if err != nil {
		return err
	}

	if valType.Type().Kind() == reflect.Map {
		err = mapWrapper.Get("validators", &d.Validators)
		if err != nil {
			return err
		}
	} else if valType.Type().Kind() == reflect.Slice {
		var validators []string
		err = mapWrapper.Get("validators", &validators)
		if err != nil {
			return err
		}

		d.Validators = map[string]any{
			AND: validators,
		}
	}

	return d.checkValidators(d.Validators, "")
}

func (d RestTestResponseJSONBody) checkValidators(validators map[string]any, parent string) error {
	if len(validators) == 0 {
		return nil
	}

	keys := maps.Keys(validators)

	if len(keys) > 1 {
		return fmt.Errorf("only one validator is supported. But found multiple: %v", keys)
	}

	if keys[0] != AND && keys[0] != OR {
		return fmt.Errorf("only %s and %s are supported. But found %s", AND, OR, keys[0])
	}

	path := fmt.Sprintf("%s -> %s", parent, keys[0])
	value := validators[keys[0]]

	valType := reflect.TypeOf(value)

	switch valType.Kind() {
	case reflect.Slice, reflect.Array:
		if valType.Elem().Kind() == reflect.String {
			return nil
		} else if valType.Elem().Kind() == reflect.Map {
			mapList, ok := value.([]map[string]any)
			if !ok {
				return fmt.Errorf("list value of key %s[*] are not all maps", path)
			}
			for i, child := range mapList {
				cType := reflect.TypeOf(child)
				if cType.Key().Kind() != reflect.String {
					return fmt.Errorf("SubMap key is not a string for key %s[%d]", path, i)
				}
				if err := d.checkValidators(child, fmt.Sprintf("%s[%d]", path, i)); err != nil {
					return err
				}
			}
		}
		return fmt.Errorf("list value of key %s[*] are not all maps/string", path)

	case reflect.Map:
		if valType.Key().Kind() != reflect.String {
			return fmt.Errorf("SubMap key is not a string for key %s", path)
		}
		return d.checkValidators(value.(map[string]any), path)
	default:
		return fmt.Errorf("value is not a slice/array/map")
	}
}

func (d RestTestResponseJSONBody) validate(obj any) (any, error) {
	jsonValidator := core.NewJSONValidator()
	if err := jsonValidator.Validate(obj, d.Validators); err != nil {
		return nil, err
	}
	return obj, nil
}
