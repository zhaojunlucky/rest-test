package model

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
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
	Script              string
}

func (d *RestTestResponseJSONBody) Validate(ctx *core.RestTestContext, resp *http.Response, js core.JSEnvExpander) (any, error) {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyStr := string(data)
	log.Infof("body: %s", bodyStr)
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
		return d.validate(arr, js)
	} else {
		var obj map[string]any
		err = json.Unmarshal(data, &obj)
		if err != nil {
			return nil, err
		}
		return d.validate(obj, js)
	}
}

func (d *RestTestResponseJSONBody) Parse(mapWrapper *collection.MapWrapper) error {

	if mapWrapper.Has("array") {
		err := mapWrapper.Get("array", &d.Array)
		if err != nil {
			return err
		}
	}

	if mapWrapper.Has("script") {
		err := mapWrapper.Get("script", &d.Script)
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
		var validators []map[string]any
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

func (d *RestTestResponseJSONBody) checkValidators(validators map[string]any, parent string) error {
	if len(validators) == 0 {
		return nil
	}
	keys := maps.Keys(validators)
	var path string
	if len(parent) > 0 {
		path = fmt.Sprintf("%s -> %s", parent, keys[0])

	} else {
		path = keys[0]
	}

	if len(validators) > 1 {

		for i := 0; i < len(keys); i++ {
			if keys[i] == AND || keys[i] == OR {
				return fmt.Errorf("path %s: can't mix operator %s or %s with value checker", path, AND, OR)
			}
		}
	}

	if keys[0] == AND || keys[0] == OR {
		value := validators[keys[0]]

		var listEle []any
		valType := reflect.ValueOf(value)

		if valType.Type().Kind() != reflect.Slice && valType.Type().Kind() != reflect.Array {
			return fmt.Errorf("path %s: only slice and array are supported", path)
		}

		for i := 0; i < valType.Len(); i++ {
			listEle = append(listEle, valType.Index(i).Interface())
		}

		return d.checkOPValidators(listEle, path)
	} else {
		return d.checkValueValidators(validators, parent)
	}
}

func (d *RestTestResponseJSONBody) validate(obj any, js core.JSEnvExpander) (any, error) {
	if len(d.Script) > 0 {
		log.Info("validate body with script")
		_, err := js.RunScriptWithBody(d.Script, obj)
		if err != nil {
			log.Errorf("script error: %s", err.Error())
			return nil, err
		}
	}
	jsonValidator := core.NewJSONValidator(js)
	if err := jsonValidator.Validate(obj, d.Validators); err != nil {
		return nil, err
	}
	return obj, nil
}

func (d *RestTestResponseJSONBody) checkOPValidators(listEle []any, path string) error {

	for i, child := range listEle {
		childPath := fmt.Sprintf("%s[%d]", path, i)
		cType := reflect.TypeOf(child)
		if cType.Kind() != reflect.Map {
			return fmt.Errorf("%s is not a map", childPath)
		}
		if cType.Key().Kind() != reflect.String {
			return fmt.Errorf("sub map %s key is not a string", childPath)
		}
		if err := d.checkValidators(child.(map[string]any), childPath); err != nil {
			return err
		}
	}
	return nil
}

func (d *RestTestResponseJSONBody) checkValueValidators(validators map[string]any, path string) error {
	for k, v := range validators {
		childPath := fmt.Sprintf("%s -> %s", path, k)
		valType := reflect.TypeOf(v)
		isArray := valType.Kind() == reflect.Array || valType.Kind() == reflect.Slice
		if isArray {

			childList, ok := v.([]any)
			if !ok {
				return fmt.Errorf("path %s: only slice and array are supported. But found %T", childPath, valType)
			}

			for i, child := range childList {
				ccPath := fmt.Sprintf("%s[%d]", childPath, i)
				ccType := reflect.TypeOf(child)
				if ccType.Kind() == reflect.Map || ccType.Kind() == reflect.Slice || ccType.Kind() == reflect.Array {
					return fmt.Errorf("path %s: only scalar values are supported. But found %T", ccPath, ccType)
				}
			}

		} else if valType.Kind() == reflect.Map {
			return fmt.Errorf("path %s: only slice, array and scaler values are supported. But found %T", childPath, valType)

		}
	}
	return nil
}
