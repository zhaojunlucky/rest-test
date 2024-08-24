package validator

import (
	"fmt"
	"github.com/dop251/goja"
	"testing"
)

type TestStruct struct {
}

func (ts *TestStruct) Expand(s string) (string, error) {
	return s, nil
}
func (ts *TestStruct) ExpandMap(m map[string]string) (map[string]string, error) {
	return m, nil
}
func (ts *TestStruct) ExpandScript(s string) (string, error) {
	return s, nil
}
func (ts *TestStruct) ExpandScriptWithBody(script string, body string) (string, error) {
	return script, nil
}
func (ts *TestStruct) RunScriptWithBody(script string, body any) (goja.Value, error) {
	return nil, nil
}

func TestJSONPathOperatorNode_Simple(t *testing.T) {
	js := &TestStruct{}
	validators := map[string]any{
		"or": []map[string]any{
			{"$.a": 1},
		},
	}

	jsonValidator, err := NewJSONPathRootValidator(validators)
	if err != nil {
		t.Fatal(err)
	}
	jsonData := map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	err = jsonValidator.Validate(js, jsonData)
	if err != nil {
		t.Fatal(err)
	}

	validators = map[string]any{
		"and": []map[string]any{
			{"$.a": 1},
			{"$.c": 2},
			{"or": []map[string]any{
				{"$.c.f": []int{5, 6}},
				{"$.c.e": "4", "$.c.g": []int{5, 6}},
			}},
		},
	}

	jsonValidator, err = NewJSONPathRootValidator(validators)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(jsonValidator)

	validators = map[string]any{
		"or": []map[string]any{
			{"$.a": 1},
			{"$.a": 2},
		},
	}

	jsonValidator, err = NewJSONPathRootValidator(validators)
	if err != nil {
		t.Fatal(err)
	}
	err = jsonValidator.Validate(js, jsonData)
	if err != nil {
		t.Fatal(err)
	}

	validators = map[string]any{
		"or": []map[string]any{
			{"$.c.f": []int{5, 6}},
		},
	}

	jsonValidator, err = NewJSONPathRootValidator(validators)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(jsonValidator)

	validators = map[string]any{
		"or": []map[string]any{
			{"$.a": 2},
			{"$.c.f": []int{5, 6}},
		},
	}

	jsonValidator, err = NewJSONPathRootValidator(validators)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(jsonValidator)

	validators2 := []map[string]any{
		map[string]any{"$.a": 1},
		map[string]any{"$.c.f": []int{5, 6}},
		map[string]any{"$.c.e": "4"},
	}

	jsonValidator, err = NewJSONPathRootValidator(validators2)
	if err != nil {
		t.Fatal(err)
	}

	if jsonValidator.operator != "and" {
		t.Fatal("expect 'and', but got", jsonValidator.operator)
	}
	fmt.Println(jsonValidator)
}

func TestJSONPathOperatorNode_SimpleFail(t *testing.T) {
	js := &TestStruct{}
	jsonData := map[string]any{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	validators := map[string]any{
		"or": []map[string]any{
			{"$.a": 2},
		},
	}
	jsonValidator, err := NewJSONPathRootValidator(validators)
	if err != nil {
		t.Fatal(err)
	}

	err = jsonValidator.Validate(js, jsonData)
	if err == nil {
		t.Fatal("expect failed")
	}

	validators = map[string]any{
		"and": []map[string]any{
			{"$.a": 1},
			{"$.c": 2},
		},
	}
	jsonValidator, err = NewJSONPathRootValidator(validators)
	if err != nil {
		t.Fatal(err)
	}
	err = jsonValidator.Validate(js, jsonData)
	if err == nil {
		t.Fatal("expect failed")
	}

	validators = map[string]any{
		"or": []map[string]any{
			{"$.a": 3},
			{"$.a": 2},
		},
	}
	jsonValidator, err = NewJSONPathRootValidator(validators)
	if err != nil {
		t.Fatal(err)
	}
	err = jsonValidator.Validate(js, jsonData)
	if err == nil {
		t.Fatal("expect failed")
	}

}

func TestJSONPathOperatorNode_ComplexPass(t *testing.T) {
	js := &TestStruct{}
	jsonData := map[string]any{
		"a": 1,
		"b": 2,
		"c": map[string]any{
			"d": 3,
			"e": "4",
			"f": []int{5, 6},
		},
	}

	validators := map[string]any{
		"or": []map[string]any{
			{"$.c.f": []int{5, 6}},
		},
	}
	jsonValidator, err := NewJSONPathRootValidator(validators)
	if err != nil {
		t.Fatal(err)
	}

	err = jsonValidator.Validate(js, jsonData)
	if err != nil {
		t.Fatal(err)
	}

	validators = map[string]any{
		"or": []map[string]any{
			{"$.a": 2},
			{"$.c.f": []int{5, 6}},
		},
	}
	jsonValidator, err = NewJSONPathRootValidator(validators)
	if err != nil {
		t.Fatal(err)
	}

	err = jsonValidator.Validate(js, jsonData)
	if err != nil {
		t.Fatal(err)
	}
	validators = map[string]any{
		"and": []map[string]any{
			{"$.a": 1},
			{"$.c.f": []int{5, 6}},
			{"$.c.e": "4"},
		},
	}
	jsonValidator, err = NewJSONPathRootValidator(validators)
	if err != nil {
		t.Fatal(err)
	}

	err = jsonValidator.Validate(js, jsonData)
	if err != nil {
		t.Fatal(err)
	}
}
