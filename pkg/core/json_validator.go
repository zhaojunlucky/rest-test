package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"reflect"
	"strings"
)

type JSONValidator struct {
	js JSEnvExpander
}

func (j *JSONValidator) Validate(obj any, validators map[string]any) error {
	if len(validators) == 0 {
		return nil
	}
	objType := reflect.TypeOf(obj)
	switch objType.Kind() {
	case reflect.Map:

		if objType.Key().Kind() != reflect.String {
			return fmt.Errorf("unsupported key type: %v", objType.Elem().Kind())
		}
		return j.validate(obj, validators)
	case reflect.Slice, reflect.Array:
		return j.validate(obj, validators)
	default:
		return fmt.Errorf("unsupported type: %v", objType)
	}
}

func (j *JSONValidator) validate(obj any, validators map[string]any) error {

	for k, v := range validators {
		valType := reflect.TypeOf(v)
		vv := v
		if valType.Kind() == reflect.String {
			vStr := strings.TrimSpace(v.(string))
			if strings.HasPrefix(vStr, "$") {
				vvStr, err := j.js.Expand(vStr)
				if err != nil {
					return err
				}
				vv, err = j.parseJSON(vvStr)
				valType = reflect.TypeOf(vv)
			}
		}

		isArray := valType.Kind() == reflect.Array || valType.Kind() == reflect.Slice
		if strings.EqualFold(k, "and") || strings.EqualFold(k, "or") {
			if isArray {
				listData := v.([]any)
				opExecutor := JSONOperator{
					expectCount: len(listData),
					OP:          k,
				}
				for _, child := range listData {
					childValidator, ok := child.(map[string]any)
					if !ok {
						return fmt.Errorf("unsupported validator: %s need map[string]any", k)
					}

					err := opExecutor.Add(j.validate(obj, childValidator))
					if err != nil {
						return err
					}
					if opExecutor.Passed() {
						return nil
					}

				}
				return opExecutor.GetErrors()

			} else {
				return fmt.Errorf("unsupported validator: %s with type []%T", k, valType)
			}
		} else {
			jsonValue, err := jsonpath.Get(k, obj)
			if err != nil {
				return errors.Join(fmt.Errorf("failed to get json path: %s", k), err)
			}
			if !j.compareNum(jsonValue, vv) && !reflect.DeepEqual(jsonValue, vv) {
				return fmt.Errorf("failed to verify json path %s, got %v, want %v", k, jsonValue, valType)
			} else {
				return nil
			}
		}
	}
	return nil
}

func (j *JSONValidator) compareNum(jsonValue, expect any) bool {
	jsonValueType := reflect.ValueOf(jsonValue)
	expectValueType := reflect.ValueOf(expect)

	if expectValueType.CanConvert(jsonValueType.Type()) {
		val := jsonValueType.Convert(expectValueType.Type()).Interface() == expect
		return val
	} else if jsonValueType.CanConvert(expectValueType.Type()) {
		val := expectValueType.Convert(jsonValueType.Type()).Interface() == jsonValue
		return val
	}
	return false
}

func (j *JSONValidator) parseJSON(vv string) (any, error) {
	type JSStruct struct {
		Value any `json:"value"`
	}
	jsonStr := fmt.Sprintf(`{"value":%s}`, vv)
	var js JSStruct
	err := json.Unmarshal([]byte(jsonStr), &js)
	if err != nil {
		return nil, err
	}
	return js.Value, nil

}

func NewJSONValidator(js JSEnvExpander) *JSONValidator {
	return &JSONValidator{
		js: js,
	}
}

type JSONOperator struct {
	expectCount  int
	OP           string
	errors       []error
	successCount int
}

func (j *JSONOperator) Add(err error) error {
	if err != nil {
		j.errors = append(j.errors, err)
	} else {
		j.successCount++
	}

	if len(j.errors)+j.successCount > j.expectCount {
		return fmt.Errorf("failed to verify json path, got %d results, want %d", j.successCount+len(j.errors), j.expectCount)
	}
	return nil
}

func (j *JSONOperator) GetErrors() error {
	return errors.Join(j.errors...)
}

func (j *JSONOperator) Passed() bool {
	if strings.EqualFold(j.OP, "or") {
		return j.successCount > 0
	} else if strings.EqualFold(j.OP, "and") {
		return j.successCount == j.expectCount
	}
	return false
}

func NewJSONOperator() *JSONOperator {
	return &JSONOperator{}
}
