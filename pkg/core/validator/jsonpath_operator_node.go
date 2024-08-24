package validator

import (
	"fmt"
	"github.com/zhaojunlucky/rest-test/pkg/core"
	"golang.org/x/exp/maps"
	"reflect"
)

type JSONPathOperatorNode struct {
	children []JSONPathNode
	operator string
	lang     *JSONPathLanguage
}

func (j *JSONPathOperatorNode) init(data map[string]any, path string) error {
	keys := maps.Keys(data)
	if len(keys) != 1 {
		return fmt.Errorf("only one operator is supported for operator node at path: %s", path)
	}
	operator := keys[0]

	if !IsLogicOperator(operator) {
		return fmt.Errorf("unsupported operator %s in opereator node at path: %s", operator, path)
	}
	j.operator = keys[0]

	var childPath string
	if path == "" {
		childPath = j.operator
	} else {
		childPath = fmt.Sprintf("%s.%s", path, j.operator)
	}

	return j.initChildren(data[keys[0]], childPath)
}

func (j *JSONPathOperatorNode) initChildren(data any, path string) error {

	cType := reflect.ValueOf(data)
	if cType.Type().Kind() != reflect.Slice && cType.Type().Kind() != reflect.Array {
		return fmt.Errorf("%s is not array or slice", path)
	}

	for i := 0; i < cType.Len(); i++ {
		childPath := fmt.Sprintf("%s[%d]", path, i)
		childValType := reflect.TypeOf(cType.Index(i).Interface())
		if childValType.Kind() != reflect.Map {
			return fmt.Errorf("path %s is not map", childPath)
		} else if childValType.Key().Kind() != reflect.String {
			return fmt.Errorf("path %s map key is not string", childPath)
		}

		cData := cType.Index(i).Interface().(map[string]any)
		cKeys := maps.Keys(cData)
		if len(cKeys) == 1 && IsLogicOperator(cKeys[0]) {
			cNode := &JSONPathOperatorNode{
				lang: j.lang,
			}
			if err := cNode.init(cData, childPath); err != nil {
				return err
			}
			j.children = append(j.children, cNode)
		} else {
			cNode := &JSONPathValueNode{
				lang: j.lang,
			}
			if err := cNode.init(cData, childPath); err != nil {
				return err
			}
			j.children = append(j.children, cNode)
		}
	}
	return nil
}

func (j *JSONPathOperatorNode) Validate(js core.JSEnvExpander, v any) error {
	jsonOP := JSONOperator{
		expectCount: len(j.children),
		OP:          j.operator,
	}
	for _, child := range j.children {
		if err := jsonOP.Add(child.Validate(js, v)); err != nil {
			return err
		}
		if jsonOP.Passed() {
			return nil
		}
	}
	if jsonOP.Passed() {
		return nil
	}
	return jsonOP.GetErrors()
}

func NewJSONPathRootValidator(def any) (*JSONPathOperatorNode, error) {
	jType := reflect.TypeOf(def)
	var data map[string]any
	if jType.Kind() == reflect.Map {
		if jType.Key().Kind() != reflect.String {
			return nil, fmt.Errorf("unsupported map key type: %s", jType.Elem().Kind().String())
		}
		data = def.(map[string]any)
	} else if jType.Kind() == reflect.Slice || jType.Kind() == reflect.Array {
		data = make(map[string]any)
		data[And] = def
	} else {
		return nil, fmt.Errorf("unsupported data type: %s", jType.Kind().String())
	}

	node := &JSONPathOperatorNode{
		lang: NewJSONPathLanguage(),
	}

	err := node.init(data, "")
	if err != nil {
		return nil, err
	}
	return node, nil
}
