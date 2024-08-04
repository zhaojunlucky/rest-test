package core

import (
	"fmt"
	"reflect"
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
		return j.validateMap(obj.(map[string]any), validators)
	case reflect.Slice, reflect.Array:
		return j.validateArray(obj.([]any), validators)
	default:
		return fmt.Errorf("unsupported type: %v", objType)
	}
	return nil
}

func NewJSONValidator() *JSONValidator {
	return &JSONValidator{}
}

func (j *JSONValidator) validateMap(obj map[string]any, validators map[string]any) error {
	for k, v := range obj {
		if validators[k] == nil {
			continue
		}
		err := j.validate(v, validators[k])
		if err != nil {
			return err
		}
	}
	return nil
}

func (j *JSONValidator) validateArray(obj []any, validators map[string]any) error {
	for _, v := range obj {
		err := j.validate(v, validators)
		if err != nil {
			return err
		}
	}
	return nil
}
