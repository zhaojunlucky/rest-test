package core

import (
	"errors"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"reflect"
	"strings"
)

type JSONValidator struct {
}

func (j *JSONValidator) Validate(obj any, validators map[string]any) error {
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
		valType := reflect.ValueOf(v)
		isArray := valType.Kind() == reflect.Array || valType.Kind() == reflect.Slice
		if strings.EqualFold(k, "and") || strings.EqualFold(k, "or") {
			if isArray && valType.Type().Key().Kind() == reflect.Map {
				listValidators := v.([]map[string]any)
				opExecutor := JSONOperator{
					expectCount: len(listValidators),
					OP:          k,
				}
				for i := 0; i < valType.Len(); i++ {
					eleVal := valType.Index(i)
					if eleVal.Type().Key().Kind() != reflect.String {
						return fmt.Errorf("unsupported validator: %s with [%d]map[%d]%T", k, i, eleVal.Type().Key(), eleVal.Type().Elem())
					}
					opExecutor.Add(j.validate(obj, listValidators[i]))
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
			valArr := v
			if !isArray {
				valArr = []any{v}
			}
			if !reflect.DeepEqual(jsonValue, valArr) {
				return fmt.Errorf("failed to verify json path %s, got %v, want %v", k, jsonValue, valType)
			} else {
				return nil
			}
		}
	}
	return nil
}

func NewJSONValidator() *JSONValidator {
	return &JSONValidator{}
}

type JSONOperator struct {
	expectCount  int
	OP           string
	errors       []error
	successCount int
}

func (j *JSONOperator) Add(err error) {
	if err != nil {
		j.errors = append(j.errors, err)
	} else {
		j.successCount++
	}
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
