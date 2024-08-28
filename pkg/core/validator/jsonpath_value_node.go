package validator

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"reflect"
	"strings"
)

type JSONPathValueNode struct {
	data map[string]reflect.Value
	lang *JSONPathLanguage
	path string
}

func (j *JSONPathValueNode) init(data map[string]any, path string) error {
	j.path = path
	if j.data == nil {
		j.data = make(map[string]reflect.Value)
	}
	for k, v := range data {
		childPath := fmt.Sprintf("%s.(%s)", path, k)
		if IsLogicOperator(k) {
			return fmt.Errorf("path %s shouldn't be a logic operator", childPath)
		}

		if err := j.lang.Compile(k); err != nil {
			return fmt.Errorf("path %s key %s is not a valid json path expression", childPath, k)
		}

		vType := reflect.ValueOf(v)
		switch vType.Type().Kind() {
		case reflect.Array, reflect.Slice:
			for i := range vType.Len() {
				cPath := fmt.Sprintf("%s[%d]", childPath, i)
				eleType := reflect.TypeOf(vType.Index(i).Interface())
				if eleType.Kind() == reflect.Map || eleType.Kind() == reflect.Array || eleType.Kind() == reflect.Slice {
					return fmt.Errorf("path: %s only scalar values are supported. But found %s", cPath, eleType.Kind().String())
				}
			}
		case reflect.Map:
			return fmt.Errorf("path %s only scalar, array and slice are supported. But found Map", childPath)
		default:
			break
		}

		j.data[k] = vType
	}

	return nil

}

func (j *JSONPathValueNode) Validate(js core.JSEnvExpander, v any) error {
	for expr, expect := range j.data {
		childPath := fmt.Sprintf("%s.(%s)", j.path, expr)
		actualVal, err := j.lang.Evaluate(expr, v)
		if err != nil {
			log.Errorf("evaluate error %s at path %s", err.Error(), childPath)
			return err
		}
		expect, err = j.handleString(expect, js)
		if err != nil {
			log.Errorf("interpret error %s at path %s", err.Error(), childPath)
			return err
		}
		actual := reflect.ValueOf(actualVal)

		err = j.compare(actual, expect)
		if err != nil {
			log.Errorf("compare error %s at path %s", err.Error(), childPath)
			return err
		}

	}
	return nil
}

func (j *JSONPathValueNode) parseJSON(vv string) (any, error) {
	type JSStruct struct {
		Value any `json:"value"`
	}
	jsonStr := fmt.Sprintf(`{"value":%s}`, vv)
	var js JSStruct
	jsonDecoder := json.NewDecoder(strings.NewReader(jsonStr))
	jsonDecoder.UseNumber()

	err := jsonDecoder.Decode(&js)
	if err != nil {
		return nil, err
	}
	switch js.Value.(type) {
	case json.Number:
		if n, err := js.Value.(json.Number).Int64(); err == nil {
			return n, nil
		} else if f, err := js.Value.(json.Number).Float64(); err == nil {
			return f, nil
		} else {
			return js.Value, nil
		}
	default:
		return js.Value, nil
	}
}

func (j *JSONPathValueNode) handleString(expect reflect.Value, js core.JSEnvExpander) (reflect.Value, error) {
	if expect.Type().Kind() == reflect.String {
		expectStr := strings.TrimSpace(expect.String())
		if strings.HasPrefix(expectStr, "$") {
			vvStr, err := js.Expand(expectStr)
			if err != nil {
				return expect, err
			}
			vv, err := j.parseJSON(vvStr)
			expect = reflect.ValueOf(vv)
			if err != nil {
				return expect, err
			}
		}
	}
	return expect, nil
}

func (j *JSONPathValueNode) compare(actual reflect.Value, expect reflect.Value) error {

	if expect.CanInt() && actual.CanInt() {
		if expect.Int() != actual.Int() {
			return fmt.Errorf("expect int %d, but got %d", expect.Int(), actual.Int())
		}
		return nil
	} else if expect.CanFloat() && actual.CanFloat() {
		if expect.Float() != actual.Float() {
			return fmt.Errorf("expect float %f, but got %f", expect.Float(), actual.Float())
		}
		return nil
	}

	if expect.Kind() != actual.Kind() {
		return fmt.Errorf("expect type %v, but got type %v", expect.Kind().String(), actual.Kind().String())
	}
	switch expect.Type().Kind() {
	case reflect.Func, reflect.Map:
		return fmt.Errorf("func and Map are not supported in value node")
	case reflect.Slice:
		if expect.Len() != actual.Len() {
			return fmt.Errorf("expect %d elements, but got %d elements", expect.Len(), actual.Len())
		}
		for i := 0; i < expect.Len(); i++ {

			if err := j.compare(actual.Index(i), expect.Index(i)); err != nil {
				return err
			}
		}
		return nil
	default:
		if !expect.Equal(actual) {
			return fmt.Errorf("expect %v, but got %v", expect, actual)
		}
	}
	return nil
}
